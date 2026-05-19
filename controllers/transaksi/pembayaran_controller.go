package transaksi

import (
	"net/http"
	"strconv"

	"backend-mantra/config"
	"backend-mantra/models"

	"github.com/gin-gonic/gin"
)

// GetRingkasanCheckout mengambil ringkasan belanja sebelum pembayaran di POS kasir.
// Dipakai oleh: kasir (GET /kasir/transaksi/:id_transaksi/checkout)
// Auth: Wajib login, role kasir
func GetRingkasanCheckout(c *gin.Context) {
	idPesananStr := c.Query("id_pesanan")
	if idPesananStr == "" {
		idPesananStr = c.Query("id_transaksi")
	}

	var pesanan models.Pesanan
	var err error
	if idPesananStr != "" {
		err = config.DB.First(&pesanan, "id_pesanan = ?", idPesananStr).Error
	} else {
		// Fallback ke pesanan paling akhir
		err = config.DB.Order("id_pesanan DESC").First(&pesanan).Error
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Transaksi/Pesanan tidak ditemukan",
		})
		return
	}

	var details []models.DetailPesanan
	config.DB.Preload("SpesifikasiBarang.Barang").Preload("SpesifikasiBarang.DetailSpesifikasi").Where("id_pesanan = ?", pesanan.IdPesanan).Find(&details)

	itemCheckout := []gin.H{}
	subtotal := 0
	for _, d := range details {
		subtotal += d.Subtotal
		varian := "Default"
		if d.SpesifikasiBarang.DetailSpesifikasi.NamaDetailSpesifikasi != "" {
			varian = d.SpesifikasiBarang.DetailSpesifikasi.NamaDetailSpesifikasi
		}
		itemCheckout = append(itemCheckout, gin.H{
			"nama_produk":    d.SpesifikasiBarang.Barang.NamaBarang,
			"varian":         varian,
			"qty":            d.Jumlah,
			"total_per_item": d.Subtotal,
		})
	}

	pajakNominal := int(float64(subtotal) * 0.11) // 11% PPN
	totalAkhir := subtotal + pajakNominal

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data checkout berhasil diambil",
		"data": gin.H{
			"order_info": gin.H{
				"id_order":    pesanan.IdPesanan,
				"nomor_order": "ORD-" + pesanan.TanggalPesanan.Format("20060102") + "-" + strconv.Itoa(int(pesanan.IdPesanan)),
			},
			"item_checkout": itemCheckout,
			"ringkasan_biaya": gin.H{
				"subtotal":      subtotal,
				"pajak_nominal": pajakNominal,
				"total_akhir":   totalAkhir,
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
	var input struct {
		IdPesanan           uint `json:"id_pesanan"`
		IdSpesifikasiBarang uint `json:"id_spesifikasi_barang"`
		Jumlah              int  `json:"jumlah"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Format inputan salah",
		})
		return
	}

	var detail models.DetailPesanan
	if err := config.DB.Where("id_pesanan = ? AND id_spesifikasi_barang = ?", input.IdPesanan, input.IdSpesifikasiBarang).First(&detail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Item pesanan tidak ditemukan",
		})
		return
	}

	tx := config.DB.Begin()

	if input.Jumlah <= 0 {
		if err := tx.Delete(&detail).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Gagal menghapus item dari transaksi",
			})
			return
		}
	} else {
		detail.Jumlah = input.Jumlah
		detail.Subtotal = detail.Jumlah * detail.HargaSatuan
		if err := tx.Save(&detail).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Gagal memperbarui jumlah item",
			})
			return
		}
	}

	// Hitung ulang total_pembayaran
	var totalBayar int64
	tx.Model(&models.DetailPesanan{}).Where("id_pesanan = ?", input.IdPesanan).Select("COALESCE(SUM(subtotal), 0)").Scan(&totalBayar)

	if err := tx.Model(&models.Pesanan{}).Where("id_pesanan = ?", input.IdPesanan).Update("total_pembayaran", totalBayar).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memperbarui total pembayaran pesanan",
		})
		return
	}

	tx.Commit()

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
	var input struct {
		IdPesanan uint `json:"id_pesanan"`
		Bayar     int  `json:"bayar"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Format inputan salah",
		})
		return
	}

	var pesanan models.Pesanan
	if err := config.DB.First(&pesanan, "id_pesanan = ?", input.IdPesanan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Pesanan tidak ditemukan",
		})
		return
	}

	pajakNominal := int(float64(pesanan.TotalPembayaran) * 0.11) // PPN 11%
	totalAkhir := pesanan.TotalPembayaran + pajakNominal

	if input.Bayar < totalAkhir {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Uang pembayaran kurang untuk melunasi transaksi",
		})
		return
	}

	kembalian := input.Bayar - totalAkhir

	tx := config.DB.Begin()

	pesanan.StatusPesanan = "Selesai"
	if err := tx.Save(&pesanan).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menyelesaikan status pesanan",
		})
		return
	}

	pembayaran := models.Pembayaran{
		PesananID:       pesanan.IdPesanan,
		PaymentType:     "cash",
		StatusTransaksi: "settlement",
		FraudStatus:     "accept",
	}
	if err := tx.Create(&pembayaran).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menyimpan data pembayaran",
		})
		return
	}

	tx.Commit()

	invoiceNum := "INV-" + pesanan.TanggalPesanan.Format("20060102") + "-" + strconv.Itoa(int(pesanan.IdPesanan))

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Pembayaran tunai berhasil",
		"data": gin.H{
			"kembalian": kembalian,
			"invoice": gin.H{
				"nomor_invoice":   invoiceNum,
				"url_print_struk": "https://api.mantra.com/struk/" + invoiceNum,
			},
		},
	})
}

// BayarNonTunai memproses pembayaran non-tunai via Midtrans (QRIS, transfer, dll).
// Dipakai oleh: kasir (POST /kasir/transaksi/:id_transaksi/bayar/non-tunai)
// Auth: Wajib login, role kasir
func BayarNonTunai(c *gin.Context) {
	var input struct {
		IdPesanan uint `json:"id_pesanan"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Format inputan salah",
		})
		return
	}

	var pesanan models.Pesanan
	if err := config.DB.First(&pesanan, "id_pesanan = ?", input.IdPesanan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Pesanan tidak ditemukan",
		})
		return
	}

	tx := config.DB.Begin()

	pesanan.StatusPesanan = "Diproses"
	if err := tx.Save(&pesanan).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memperbarui status pesanan",
		})
		return
	}

	pembayaran := models.Pembayaran{
		PesananID:       pesanan.IdPesanan,
		OrderIdMidtrans: "MID-" + pesanan.TanggalPesanan.Format("20060102") + "-" + strconv.Itoa(int(pesanan.IdPesanan)),
		PaymentType:     "qris",
		StatusTransaksi: "pending",
		FraudStatus:     "challenge",
	}
	if err := tx.Create(&pembayaran).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mencatat data pembayaran non-tunai",
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Pembayaran non-tunai diproses",
		"data": gin.H{
			"midtrans_data": gin.H{
				"token":        "snap-token-" + strconv.Itoa(int(pesanan.IdPesanan)),
				"redirect_url": "https://app.sandbox.midtrans.com/snap/v2/vtweb/" + strconv.Itoa(int(pesanan.IdPesanan)),
			},
		},
	})
}
