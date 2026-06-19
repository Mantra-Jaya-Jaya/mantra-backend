package transaksi

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"backend-mantra/config"
	"backend-mantra/models"
	"backend-mantra/utils"

	"github.com/gin-gonic/gin"
)

type midtransNotification struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionID     string `json:"transaction_id"`
	TransactionStatus string `json:"transaction_status"`
	TransactionType   string `json:"transaction_type"`
	FraudStatus       string `json:"fraud_status"`
	OrderID           string `json:"order_id"`
	StatusCode        string `json:"status_code"`
	GrossAmount       string `json:"gross_amount"`
	PaymentType       string `json:"payment_type"`
	SignatureKey      string `json:"signature_key"`
	VaNumbers         []struct {
		Bank     string `json:"bank"`
		VaNumber string `json:"va_number"`
	} `json:"va_numbers"`
	BillKey         string `json:"bill_key"`
	BillCode        string `json:"bill_code"`
	Store           string `json:"store"`
	PaymentCode     string `json:"payment_code"`
	PermataVaNumber string `json:"permata_va_number"`
	QrCodeUrl       string `json:"qr_code_url"`
	MerchantID      string `json:"merchant_id"`
	Acquirer        string `json:"acquirer"`
	Currency        string `json:"currency"`
}

func MidtransNotificationHandler(c *gin.Context) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println("❌ Webhook Error: Gagal membaca body")
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Gagal membaca body"})
		return
	}

	var notif midtransNotification
	if err := json.Unmarshal(bodyBytes, &notif); err != nil {
		fmt.Println("❌ Webhook Error: Format notifikasi tidak valid")
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format notifikasi tidak valid"})
		return
	}

	fmt.Printf("ℹ️ Webhook Received: OrderID=%s, Status=%s\n", notif.OrderID, notif.TransactionStatus)

	// 1. Verifikasi Keamanan
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	expectedSignature := computeSignature(notif.OrderID, notif.StatusCode, notif.GrossAmount, serverKey)

	if !hmac.Equal([]byte(notif.SignatureKey), []byte(expectedSignature)) {
		if os.Getenv("MIDTRANS_ENVIRONMENT") == "production" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Signature tidak valid"})
			return
		}
	}

	// 2. Cari Data Pembayaran
	var pembayaran models.Pembayaran
	if err := config.DB.Where("order_id_midtrans = ?", notif.OrderID).First(&pembayaran).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Pembayaran tidak ditemukan"})
		return
	}

	// 3. Siapkan Update Pembayaran
	updatePembayaran := map[string]interface{}{
		"transaksi_midtrans_id": notif.TransactionID,
	}

	tipePembayaranDB := notif.PaymentType
	if notif.PaymentType == "gopay" {
		tipePembayaranDB = "qris"
	}
	if id := utils.GetTipePembayaranIDSafe(tipePembayaranDB); id != 0 {
		updatePembayaran["id_tipe_pembayaran"] = id
	}
	if id := utils.GetStatusTransaksiIDSafe(notif.TransactionStatus); id != 0 {
		updatePembayaran["id_status_transaksi"] = id
	}
	if id := utils.GetFraudStatusIDSafe(notif.FraudStatus); id != 0 {
		updatePembayaran["id_fraud_status"] = id
	}

	grossAmount := 0
	if err := parseGrossAmount(notif.GrossAmount, &grossAmount); err == nil {
		updatePembayaran["total_dibayar"] = grossAmount
	}

	isLunas := notif.TransactionStatus == "settlement" || notif.TransactionStatus == "capture"
	if isLunas {
		updatePembayaran["waktu_pembayaran"] = time.Now()
	}

	// 4. Update Database Pembayaran
	config.DB.Model(&pembayaran).Updates(updatePembayaran)

	if notif.TransactionStatus == "settlement" || notif.TransactionStatus == "capture" {
		dikemasID := utils.GetStatusPesananIDSafe("Dikemas")
		if dikemasID != 0 {
			err := config.DB.Model(&models.Pesanan{}).Where("id_pesanan = ?", pembayaran.PesananID).
				Update("id_status_pesanan", dikemasID).Error
			if err != nil {
				fmt.Printf("❌ Webhook Error: Gagal update status pesanan: %s\n", err.Error())
			} else {
				fmt.Printf("✅ Webhook Success: Pesanan ID %d status berubah jadi 'Dikemas'\n", pembayaran.PesananID)
			}
		}
		processExternalShipment(pembayaran.PesananID)
	}

	// 🚀 6. PANGGIL FUNGSI SAVE DETAIL BUAT INVOICE LUUU!!!
	savePaymentDetails(&pembayaran, &notif)

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Notifikasi diproses"})
}

// =========================================================================
// 2. FUNGSI SAVE DETAIL (SUDAH DIPASANG HELM ANTI-DUPLIKAT)
// =========================================================================
func savePaymentDetails(pembayaran *models.Pembayaran, notif *midtransNotification) {
	// 🚀 HELM PENGAMAN: Cek dulu apakah detailnya udah dibikin pas kasir ngeklik "Bayar"?
	var count int64
	config.DB.Model(&models.DetailPembayaran{}).Where("id_pembayaran = ?", pembayaran.IdPembayaran).Count(&count)

	// Kalau udah ada (count > 0), STOP! Gak usah Create lagi biar tabel invoice lu gak beranak-pinak
	if count > 0 {
		return
	}

	// Kalau ternyata belum ada (misal transaksinya digenerate dari luar aplikasi kasir), baru kita Create
	detail := models.DetailPembayaran{
		PembayaranID: pembayaran.IdPembayaran,
	}

	// Isi data sesuai balikan Midtrans
	if len(notif.VaNumbers) > 0 {
		detail.KanalPembayaran = "va_" + strings.ToLower(notif.VaNumbers[0].Bank)
		detail.NomorVA = notif.VaNumbers[0].VaNumber
		detail.NamaBank = notif.VaNumbers[0].Bank
	} else if notif.PermataVaNumber != "" {
		detail.KanalPembayaran = "va_permata"
		detail.NomorVA = notif.PermataVaNumber
		detail.NamaBank = "permata"
	} else if notif.QrCodeUrl != "" {
		detail.KanalPembayaran = "qris"
		detail.QrCode = notif.QrCodeUrl
	} else if notif.PaymentCode != "" {
		detail.KanalPembayaran = "payment_code"
		detail.MerchantID = notif.MerchantID
		detail.NamaBank = notif.Store
	}

	// Simpan 1 kali aja ke database!
	config.DB.Create(&detail)
}

func computeSignature(orderID, statusCode, grossAmount, serverKey string) string {
	payload := orderID + statusCode + grossAmount + serverKey
	hash := sha512.New()
	hash.Write([]byte(payload))
	return hex.EncodeToString(hash.Sum(nil))
}

func parseGrossAmount(amount string, result *int) error {
	val, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return err
	}
	*result = int(math.Round(val))
	return nil
}
