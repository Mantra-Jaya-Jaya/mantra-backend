package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetDashboardKasir handles fetching cashier dashboard data
func GetDashboardKasir(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data dashboard berhasil diambil",
		"data": gin.H{
			"user": gin.H{
				"nama_kasir":        "Budi Santoso",
				"status_notifikasi": true,
			},
			"statistik_hari_ini": gin.H{
				"total_pendapatan":   1500000,
				"jumlah_transaksi":   12,
				"total_item_terjual": 35,
			},
			"aktivitas_terkini": []gin.H{
				{
					"id_transaksi":      101,
					"nomor_invoice":     "INV-20260509-001",
					"metode_pembayaran": "cash",
					"waktu":             "09:30",
					"total_bayar":       125000,
				},
			},
		},
	})
}

// GetLaporanRingkasan handles fetching sales report summary
func GetLaporanRingkasan(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data laporan berhasil diambil",
		"data": gin.H{
			"header_statistik": gin.H{
				"total_pendapatan":               5000000,
				"persentase_kenaikan_pendapatan": 12.5,
				"total_transaksi":                40,
				"persentase_kenaikan_transaksi":  5.0,
				"rata_rata_pesanan":              125000,
				"status_rata_rata":               "naik",
			},
			"grafik_pendapatan": []gin.H{
				{
					"label":        "08:00",
					"nilai":        250000,
					"is_highlight": false,
				},
			},
			"produk_terlaris": []gin.H{
				{
					"id_produk":      1,
					"nama_produk":    "Laptop Gaming X",
					"deskripsi":      "16GB RAM, 1TB SSD",
					"jumlah_terjual": 15,
					"gambar":         "https://api.mantra.com/storage/barang/laptop-x.jpg",
				},
			},
		},
	})
}

// GetDetailLaporanProduk handles fetching transaction details per product
func GetDetailLaporanProduk(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Detail produk berhasil diambil",
		"data": gin.H{
			"produk": gin.H{
				"id_produk":   1,
				"nama_produk": "Laptop Gaming X",
				"kategori":    "Elektronik",
				"gambar":      "https://api.mantra.com/storage/barang/laptop-x.jpg",
			},
			"statistik_produk": gin.H{
				"total_terjual":       50,
				"terjual_periode_ini": 15,
				"label_periode":       "mingguan",
			},
			"riwayat_transaksi": []gin.H{
				{
					"id_transaksi":  101,
					"nomor_invoice": "INV-20260509-001",
					"tanggal_waktu": "2026-05-09T09:30:00Z",
					"subtotal":      13500000,
					"quantity":      1,
				},
			},
		},
	})
}

// GetDetailPesananDariLaporan handles fetching single order details from report
func GetDetailPesananDariLaporan(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Detail pesanan berhasil diambil",
		"data": gin.H{
			"order_info": gin.H{
				"nomor_order":   "ORD-20260509-001",
				"tanggal_waktu": "2026-05-09T09:30:00Z",
				"status_order": gin.H{
					"kode":  "selesai",
					"label": "Selesai",
				},
			},
			"pelanggan": gin.H{
				"nama":   "Aarav Lysander",
				"alamat": "Jl. Cempaka Putih No. 12",
			},
			"daftar_item": []gin.H{
				{
					"id_produk":        1,
					"nama_produk":      "Laptop Gaming X",
					"qty":              1,
					"total_harga_item": 13500000,
				},
			},
			"rincian_pembayaran": gin.H{
				"metode":        "qris",
				"subtotal":      13500000,
				"pajak_nominal": 1485000,
				"total_akhir":   14985000,
			},
		},
	})
}

// GetDaftarPesananKasir handles fetching all orders for cashier
func GetDaftarPesananKasir(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Daftar pesanan berhasil diambil",
		"data": gin.H{
			"daftar_pesanan": []gin.H{
				{
					"id_order":    101,
					"nomor_order": "ORD-20260509-001",
					"sumber_pesanan": gin.H{
						"kode":      "online",
						"icon_type": "globe",
					},
					"status": gin.H{
						"kode":  "diproses",
						"label": "Diproses",
					},
					"ringkasan_item": "2x Kopi, 1x Roti",
					"waktu_relatif":  "5 menit lalu",
					"total_harga":    45000,
				},
			},
		},
		"meta": gin.H{
			"page":        1,
			"limit":       20,
			"total":       45,
			"total_pages": 3,
		},
	})
}

// GetDetailPesananKasir handles fetching single order details for cashier
func GetDetailPesananKasir(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Detail pesanan berhasil diambil",
		"data": gin.H{
			"order_header": gin.H{
				"id_order":      101,
				"nomor_order":   "ORD-20260509-001",
				"tanggal_waktu": "2026-05-09T09:30:00Z",
				"status": gin.H{
					"kode":      "selesai",
					"label":     "Selesai",
					"warna_hex": "#22C55E",
				},
			},
			"informasi_pelanggan": gin.H{
				"nama":              "Aarav Lysander",
				"alamat_pengiriman": "Jl. Cempaka Putih No. 12",
				"map_url":           nil,
			},
			"item_pesanan": []gin.H{
				{
					"id_produk":        1,
					"nama_produk":      "Laptop Gaming X",
					"varian":           "16GB RAM",
					"gambar":           "https://api.mantra.com/storage/barang/laptop-x.jpg",
					"qty":              1,
					"harga_satuan":     13500000,
					"total_harga_item": 13500000,
				},
			},
			"data_kasir": gin.H{
				"nama_kasir": "Budi Santoso",
				"shift_info": "Shift Pagi",
			},
			"informasi_pembayaran": gin.H{
				"metode_pembayaran": "cash",
				"status_pembayaran": "lunas",
				"rincian_kalkulasi": gin.H{
					"subtotal":      13500000,
					"pajak_persen":  11,
					"pajak_nominal": 1485000,
					"total_akhir":   14985000,
				},
			},
		},
	})
}

// CariProdukTransaksi handles product lookup by barcode or name for transaction
func CariProdukTransaksi(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Produk berhasil ditemukan",
		"data": gin.H{
			"id_produk":    1,
			"nama_produk":  "Laptop Gaming X",
			"harga_satuan": 13500000,
			"gambar":       "https://api.mantra.com/storage/barang/laptop-x.jpg",
			"varian": []gin.H{
				{
					"id_spesifikasi_barang": 5,
					"label":                 "16GB RAM",
					"stok":                  5,
				},
			},
		},
	})
}

// UpdateQuantityItem handles updating quantity in active transaction
func UpdateQuantityItem(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Quantity berhasil diperbarui",
		"data":    nil,
	})
}

// GetRingkasanCheckout handles fetching pre-payment summary
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

// BayarTunai handles cash payment
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

// BayarNonTunai handles non-cash payment (Midtrans)
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

// GetProfilKasir handles fetching cashier profile
func GetProfilKasir(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data profil berhasil diambil",
		"data": gin.H{
			"nama_kasir":  "Budi Santoso",
			"email":       "budi@mantra.com",
			"role":        "kasir",
			"shift":       "Pagi",
			"status_akun": "aktif",
		},
	})
}

// GetNotifikasiKasir handles fetching cashier notifications
func GetNotifikasiKasir(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Daftar notifikasi berhasil diambil",
		"data": gin.H{
			"notifikasi_grup": []gin.H{
				{
					"label_waktu": "Hari ini",
					"daftar_notifikasi": []gin.H{
						{
							"id":         1,
							"judul":      "Stok Menipis",
							"pesan":      "Stok Laptop Gaming X tersisa 3 unit",
							"tipe_icon":  "warning",
							"created_at": "2026-05-09T08:00:00Z",
						},
					},
				},
			},
		},
	})
}
