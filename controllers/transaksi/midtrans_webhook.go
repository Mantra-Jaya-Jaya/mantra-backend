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

	fmt.Printf("ℹ️ Webhook Received: OrderID=%s, Status=%s, PaymentType=%s\n", notif.OrderID, notif.TransactionStatus, notif.PaymentType)

	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	expectedSignature := computeSignature(notif.OrderID, notif.StatusCode, notif.GrossAmount, serverKey)

	if !hmac.Equal([]byte(notif.SignatureKey), []byte(expectedSignature)) {
		fmt.Printf("❌ Webhook Error: Signature tidak valid. Expected: %s, Got: %s\n", expectedSignature, notif.SignatureKey)
		if os.Getenv("MIDTRANS_ENVIRONMENT") == "production" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Signature tidak valid"})
			return
		}
		fmt.Println("⚠️  MIDTRANS_ENVIRONMENT bukan 'production' — signature tetap diverifikasi tapi dilanjutkan (development mode)")
	}

	var pembayaran models.Pembayaran
	if err := config.DB.Where("order_id_midtrans = ?", notif.OrderID).First(&pembayaran).Error; err != nil {
		fmt.Printf("❌ Webhook Error: Pembayaran dengan OrderID %s tidak ditemukan di DB\n", notif.OrderID)
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Pembayaran tidak ditemukan"})
		return
	}

	pembayaran.TransaksiMidtransID = notif.TransactionID

	// Mapping payment_type dari Midtrans ke standar database kita
	tipePembayaranDB := notif.PaymentType
	if notif.PaymentType == "gopay" {
		tipePembayaranDB = "qris"
	}
	tipePembayaranID := utils.GetTipePembayaranIDSafe(tipePembayaranDB)
	if tipePembayaranID != 0 {
		pembayaran.TipePembayaranID = tipePembayaranID
	}

	statusTransaksiID := utils.GetStatusTransaksiIDSafe(notif.TransactionStatus)
	if statusTransaksiID != 0 {
		pembayaran.StatusTransaksiID = statusTransaksiID
	}

	fraudStatusID := utils.GetFraudStatusIDSafe(notif.FraudStatus)
	if fraudStatusID != 0 {
		pembayaran.FraudStatusID = fraudStatusID
	}

	grossAmount := 0
	if err := parseGrossAmount(notif.GrossAmount, &grossAmount); err == nil {
		pembayaran.TotalDibayar = grossAmount
	}

	if notif.TransactionStatus == "settlement" || notif.TransactionStatus == "capture" {
		now := time.Now()
		pembayaran.WaktuPembayaran = &now
		fmt.Println("✅ Webhook Info: Transaksi LUNAS")
	}

	if err := config.DB.Save(&pembayaran).Error; err != nil {
		fmt.Printf("❌ Webhook Error: Gagal simpan status pembayaran: %s\n", err.Error())
	}

	if notif.TransactionStatus == "settlement" || notif.TransactionStatus == "capture" {
		diprosesID := utils.GetStatusPesananIDSafe("Diproses")
		if diprosesID != 0 {
			err := config.DB.Model(&models.Pesanan{}).Where("id_pesanan = ?", pembayaran.PesananID).
				Update("id_status_pesanan", diprosesID).Error
			if err != nil {
				fmt.Printf("❌ Webhook Error: Gagal update status pesanan: %s\n", err.Error())
			} else {
				fmt.Printf("✅ Webhook Success: Pesanan ID %d status berubah jadi 'Diproses'\n", pembayaran.PesananID)
			}
		}
	}

	savePaymentDetails(&pembayaran, &notif)

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Notifikasi diterima"})
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

func savePaymentDetails(pembayaran *models.Pembayaran, notif *midtransNotification) {
	for _, va := range notif.VaNumbers {
		detail := models.DetailPembayaran{
			PembayaranID:    pembayaran.IdPembayaran,
			KanalPembayaran: "va_" + strings.ToLower(va.Bank),
			NomorVA:         va.VaNumber,
			NamaBank:        va.Bank,
		}
		config.DB.Create(&detail)
	}

	if notif.BillKey != "" && notif.BillCode != "" {
		detail := models.DetailPembayaran{
			PembayaranID:    pembayaran.IdPembayaran,
			KanalPembayaran: "retail",
			BillKey:         notif.BillKey,
			BillCode:        notif.BillCode,
			MerchantID:      notif.MerchantID,
		}
		config.DB.Create(&detail)
	}

	if notif.PermataVaNumber != "" {
		detail := models.DetailPembayaran{
			PembayaranID:    pembayaran.IdPembayaran,
			KanalPembayaran: "va_permata",
			NomorVA:         notif.PermataVaNumber,
		}
		config.DB.Create(&detail)
	}

	if notif.QrCodeUrl != "" {
		detail := models.DetailPembayaran{
			PembayaranID:    pembayaran.IdPembayaran,
			KanalPembayaran: "qris",
			QrCode:          notif.QrCodeUrl,
		}
		config.DB.Create(&detail)
	}

	if notif.PaymentCode != "" {
		detail := models.DetailPembayaran{
			PembayaranID:    pembayaran.IdPembayaran,
			KanalPembayaran: "payment_code",
			BillCode:        notif.PaymentCode,
			MerchantID:      notif.MerchantID,
			NamaBank:        notif.Store,
		}
		config.DB.Create(&detail)
	}
}