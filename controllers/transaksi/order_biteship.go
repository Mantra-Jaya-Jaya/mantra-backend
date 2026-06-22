package transaksi

import (
	"net/http"

	"backend-mantra/config"
	"backend-mantra/controllers/auth"
	"backend-mantra/models"
	"backend-mantra/services"
	"backend-mantra/utils"

	"github.com/gin-gonic/gin"
)

// GetBiteshipOrderStatus mengambil status order Biteship untuk pesanan customer.
// Route: GET /api/v1/customer/pesanan/:public_id/status-biteship
func GetBiteshipOrderStatus(c *gin.Context) {
	publicID := c.Param("public_id")
	role := auth.NormalizeRoleName(c.GetString("role"))

	if role != "admin" {
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "User ID tidak ditemukan di session/token",
			})
			return
		}

		userIDInt64, ok := userIDInterface.(int64)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Terjadi kesalahan sistem: tipe data User ID tidak valid",
			})
			return
		}

		owned, err := isCustomerPesananOwner(publicID, uint(userIDInt64))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Gagal memverifikasi kepemilikan pesanan",
			})
			return
		}

		if !owned {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "Anda tidak memiliki akses ke resource ini",
				"error": gin.H{
					"code":   "AUTH_002",
					"detail": "Pesanan ini bukan milik Anda",
				},
			})
			return
		}
	}

	var pesanan models.Pesanan
	if err := config.DB.Preload("Ekspedisi").Where("public_id = ?", publicID).First(&pesanan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Pesanan tidak ditemukan"})
		return
	}

	// Hanya berlaku untuk pesanan eksternal
	if pesanan.BiteshipOrderID == nil || *pesanan.BiteshipOrderID == "" {
		if pesanan.TipeKurirID != utils.GetTipeKurirID("external") {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Pesanan ini bukan ekspedisi eksternal"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Order Biteship belum dibuat atau waybill belum terbit",
			"data": gin.H{
				"order_id":   nil,
				"waybill_id": pesanan.NomorResi,
				"status":     "pending",
				"ekspedisi":  "",
			},
		})
		return
	}

	adapter := services.NewBiteshipAdapter()
	order, err := adapter.GetOrderStatus(*pesanan.BiteshipOrderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil status order Biteship: " + err.Error(),
		})
		return
	}

	namaEkspedisi := ""
	if pesanan.Ekspedisi != nil {
		namaEkspedisi = pesanan.Ekspedisi.NamaEkspedisi
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Status order berhasil diambil",
		"data": gin.H{
			"order_id":   order.ID,
			"waybill_id": order.WaybillID,
			"status":     order.Status,
			"ekspedisi":  namaEkspedisi,
			"courier":    order.Courier,
			"nomor_resi": pesanan.NomorResi,
		},
	})
}

// CancelBiteshipOrder membatalkan order Biteship (hanya untuk admin).
// Route: POST /api/v1/admin/pesanan/:public_id/cancel-shipment
func CancelBiteshipOrder(c *gin.Context) {
	publicID := c.Param("public_id")

	var pesanan models.Pesanan
	if err := config.DB.Preload("Ekspedisi").Where("public_id = ?", publicID).First(&pesanan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Pesanan tidak ditemukan"})
		return
	}

	// Validasi tipe pesanan
	if pesanan.TipeKurirID != utils.GetTipeKurirID("external") {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Pesanan ini bukan ekspedisi eksternal"})
		return
	}

	// Validasi status: hanya bisa cancel kalau belum dikirim
	statusPesanan := ""
	if pesanan.StatusPesanan != nil {
		statusPesanan = pesanan.StatusPesanan.NamaStatus
	}
	if statusPesanan == "Dikirim" || statusPesanan == "Selesai" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Pesanan sudah dalam perjalanan, tidak bisa dibatalkan"})
		return
	}

	// Cek apakah ada biteship_order_id
	if pesanan.BiteshipOrderID == nil || *pesanan.BiteshipOrderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Order Biteship belum dibuat"})
		return
	}

	// Cancel di Biteship
	adapter := services.NewBiteshipAdapter()
	if err := adapter.CancelOrder(*pesanan.BiteshipOrderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal membatalkan order Biteship: " + err.Error(),
		})
		return
	}

	// Update status pesanan jadi Dibatalkan
	dibatalkanID := utils.GetStatusPesananIDSafe("Dibatalkan")
	if dibatalkanID > 0 {
		config.DB.Model(&pesanan).Updates(map[string]any{
			"id_status_pesanan": dibatalkanID,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Pesanan berhasil dibatalkan",
	})
}
