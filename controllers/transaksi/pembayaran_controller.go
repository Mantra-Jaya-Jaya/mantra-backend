package transaksi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetRingkasanCheckout mengambil ringkasan belanja sebelum pembayaran di POS kasir.
// Dipakai oleh: kasir (GET /kasir/transaksi/:id_transaksi/checkout)
// Auth: Wajib login, role kasir
func GetRingkasanCheckout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data checkout berhasil diambil",
		"data": gin.H{
			"order_info": gin.H{
				"id_order":    101,
				"nomor_order": "ORD-20260509-001",
			},
			"item_checkout": []gin.H{
				{
					"nama_produk":    "Laptop Gaming X",
					"varian":         "16GB RAM",
					"qty":            1,
					"total_per_item": 13500000,
				},
			},
			"ringkasan_biaya": gin.H{
				"subtotal":      13500000,
				"pajak_nominal": 1485000,
				"total_akhir":   14985000,
			},
			"pilihan_pembayaran": []gin.H{
				{
					"id_metode": 1,
					"label":     "Cash",
					"tipe":      "cash",
				},
				{
					"id_metode": 2,
					"label":     "QRIS",
					"tipe":      "non-cash",
				},
			},
		},
	})
}

// UpdateQuantityItem memperbarui quantity item dalam transaksi POS yang sedang berjalan.
// Dipakai oleh: kasir (PATCH /kasir/transaksi/:id_transaksi/item/:id_item)
// Auth: Wajib login, role kasir
func UpdateQuantityItem(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Quantity berhasil diperbarui",
		"data":    nil,
	})
}

// BayarTunai memproses pembayaran tunai (cash) di POS kasir.
// Dipakai oleh: kasir (POST /kasir/transaksi/:id_transaksi/bayar/tunai)
// Auth: Wajib login, role kasir
func BayarTunai(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Pembayaran tunai berhasil",
		"data": gin.H{
			"kembalian": 15000,
			"invoice": gin.H{
				"nomor_invoice":   "INV-20260509-001",
				"url_print_struk": "https://api.mantra.com/struk/INV-20260509-001",
			},
		},
	})
}

// BayarNonTunai memproses pembayaran non-tunai via Midtrans (QRIS, transfer, dll).
// Dipakai oleh: kasir (POST /kasir/transaksi/:id_transaksi/bayar/non-tunai)
// Auth: Wajib login, role kasir
func BayarNonTunai(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Pembayaran non-tunai diproses",
		"data": gin.H{
			"midtrans_data": gin.H{
				"token":        "snap-token-dari-midtrans",
				"redirect_url": "https://app.sandbox.midtrans.com/snap/v2/vtweb/...",
			},
		},
	})
}
