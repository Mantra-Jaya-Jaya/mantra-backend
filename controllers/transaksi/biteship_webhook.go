package transaksi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"backend-mantra/config"
	"backend-mantra/models"
	"backend-mantra/utils"

	"github.com/gin-gonic/gin"
)

// BiteshipWebhookPayload merepresentasikan struktur notifikasi webhook Biteship
type BiteshipWebhookPayload struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

// biteshipOrderData detail order Biteship
type biteshipOrderData struct {
	ID        string `json:"id"`
	WaybillID string `json:"waybill_id"`
	Status    string `json:"status"`
}

// BiteshipWebhookHandler menerima notifikasi webhook dari Biteship.
// Route: POST /api/v1/webhook/biteship
// Header: X-Biteship-Signature (HMAC-SHA256)
// Auth: Tidak perlu (public endpoint dengan signature verification)
func BiteshipWebhookHandler(c *gin.Context) {
	// 1. Baca payload mentah
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Gagal membaca body"})
		return
	}

	// 2. Verifikasi signature (opsional, tapi disarankan)
	signature := c.GetHeader("X-Biteship-Signature")
	if signature != "" && os.Getenv("BITESHIP_WEBHOOK_SECRET") != "" {
		expectedSig := computeBiteshipSignature(bodyBytes)
		if !hmac.Equal([]byte(signature), []byte(expectedSig)) {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Signature tidak valid"})
			return
		}
	}

	var payload BiteshipWebhookPayload
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format payload tidak valid"})
		return
	}

	fmt.Printf("📦 Biteship Webhook: event=%s\n", payload.Event)

	// 3. Tangani event sesuai tipe
	switch payload.Event {
	case "order.status":
		handleOrderStatus(payload.Data)
	case "order.waybill_id":
		handleWaybillUpdate(payload.Data)
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Webhook diproses"})
}

func handleOrderStatus(data json.RawMessage) {
	var order biteshipOrderData
	if err := json.Unmarshal(data, &order); err != nil {
		fmt.Printf("❌ Webhook Error: parse order.status: %s\n", err.Error())
		return
	}

	// Cari pesanan berdasarkan biteship_order_id
	var pesanan models.Pesanan
	if err := config.DB.Where("biteship_order_id = ?", order.ID).First(&pesanan).Error; err != nil {
		fmt.Printf("❌ Webhook Error: pesanan tidak ditemukan untuk order_id=%s\n", order.ID)
		return
	}

	// Update status pesanan
	statusMap := map[string]string{
		"pending":  "Dikemas",
		"active":   "Dikirim",
		"delivered": "Selesai",
		"returned": "Dikirim", // atau Dibatalkan, tergantung logika bisnis
	}

	if newStatus, ok := statusMap[order.Status]; ok {
		statusID := utils.GetStatusPesananIDSafe(newStatus)
		if statusID > 0 {
			config.DB.Model(&pesanan).Update("id_status_pesanan", statusID)
			fmt.Printf("✅ Webhook Success: Pesanan %s status → %s\n", pesanan.PublicId, newStatus)
		}
	}
}

func handleWaybillUpdate(data json.RawMessage) {
	var order biteshipOrderData
	if err := json.Unmarshal(data, &order); err != nil {
		fmt.Printf("❌ Webhook Error: parse order.waybill_id: %s\n", err.Error())
		return
	}

	// Cari pesanan berdasarkan biteship_order_id
	var pesanan models.Pesanan
	if err := config.DB.Where("biteship_order_id = ?", order.ID).First(&pesanan).Error; err != nil {
		// fallback: cari via nomor_resi yang sudah ada
		if err := config.DB.Where("nomor_resi = ?", order.WaybillID).First(&pesanan).Error; err != nil {
			fmt.Printf("❌ Webhook Error: waybill tidak dikenali: %s\n", order.WaybillID)
			return
		}
	}

	// Update nomor resi kalau belum ada
	if pesanan.NomorResi == nil || *pesanan.NomorResi == "" {
		config.DB.Model(&pesanan).Update("nomor_resi", order.WaybillID)
		fmt.Printf("✅ Webhook Success: Pesanan %s waybill → %s\n", pesanan.PublicId, order.WaybillID)
	}
}

func computeBiteshipSignature(body []byte) string {
	secret := os.Getenv("BITESHIP_WEBHOOK_SECRET")
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(body)
	return hex.EncodeToString(h.Sum(nil))
}