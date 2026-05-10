package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPromo handles fetching active promos
func GetPromo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berhasil mengambil data promo",
		"data": []gin.H{
			{
				"id_diskon":   1,
				"nama_diskon": "Promo Awal Tahun",
				"banner_url":  "https://api.mantra.com/storage/banner/promo-1.jpg",
				"tgl_selesai": "2026-12-31T23:59:59Z",
			},
		},
	})
}

// GetKategori handles fetching product categories
func GetKategori(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berhasil mengambil daftar kategori",
		"data": []gin.H{
			{
				"id_kategori":   1,
				"nama_kategori": "Gadget",
				"icon_kategori": "https://api.mantra.com/storage/icons/gadget.png",
			},
		},
	})
}

// GetDaftarBarang handles fetching catalog products
func GetDaftarBarang(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berhasil mengambil daftar barang",
		"data": []gin.H{
			{
				"id_barang":       1,
				"nama_barang":     "Laptop Gaming X",
				"harga_terendah":  13500000,
				"harga_tertinggi": 15000000,
				"harga_diskon":    12150000,
				"punya_diskon":    true,
				"gambar_barang":   "https://api.mantra.com/storage/barang/laptop-x.jpg",
			},
		},
		"meta": gin.H{
			"page":        1,
			"limit":       10,
			"total":       50,
			"total_pages": 5,
		},
	})
}

// GetDetailBarangByScan handles fetching product by barcode
func GetDetailBarangByScan(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data barang ditemukan",
		"data": gin.H{
			"id_barang":     10,
			"nama_barang":   "Laptop Gaming X",
			"kode_barcode":  "89912345678",
			"gambar_barang": "https://api.mantra.com/storage/barang/laptop-x.jpg",
			"kategori":      "Elektronik",
			"satuan":        "Unit",
			"diskon": gin.H{
				"nama_diskon":  "Promo Awal Tahun",
				"besar_diskon": 10,
				"tgl_selesai":  "2026-12-31T23:59:59Z",
			},
			"varian": []gin.H{
				{
					"id_spesifikasi_barang": 5,
					"nama_spesifikasi":      "RAM",
					"nama_detail":           "16GB",
					"harga_barang":          15000000,
					"harga_diskon":          13500000,
					"stok":                  5,
				},
			},
		},
	})
}

// TambahKeKeranjang handles adding items to cart
func TambahKeKeranjang(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Barang berhasil ditambahkan ke keranjang",
	})
}

// UpdateKeranjang handles updating cart item quantity
func UpdateKeranjang(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Keranjang berhasil diperbarui",
	})
}

// HapusItemKeranjang handles deleting item from cart
func HapusItemKeranjang(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Item keranjang berhasil dihapus",
	})
}

// GetNotifikasiCustomer handles fetching customer notifications
func GetNotifikasiCustomer(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Notifikasi berhasil diambil",
		"data": []gin.H{
			{
				"id_notifikasi": 1,
				"judul":         "Diskon Menanti!",
				"pesan":         "Ada diskon 10% untuk barang favoritmu hari ini.",
				"status":        "unread",
				"created_at":    "2026-05-05T10:00:00Z",
			},
		},
	})
}

// GetDaftarPesananCustomer handles fetching customer order history
func GetDaftarPesananCustomer(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Daftar pesanan berhasil diambil",
		"data": []gin.H{
			{
				"id_pesanan":    "12345678",
				"status":        "diproses",
				"tanggal_pesan": "2026-04-10T22:39:00Z",
				"total_bayar":   50000,
				"items": []gin.H{
					{
						"id_barang":       101,
						"nama_barang":     "Novel Ancika 1995",
						"jumlah":          1,
						"harga_saat_beli": 50000,
						"gambar":          "https://api.mantra.com/storage/barang/ancika.jpg",
					},
				},
			},
		},
		"meta": gin.H{
			"page":        1,
			"limit":       10,
			"total":       5,
			"total_pages": 1,
		},
	})
}

// CheckoutPesanan handles creating an order
func CheckoutPesanan(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Pesanan berhasil dibuat",
		"data": gin.H{
			"id_pesanan":     "12345678",
			"midtrans_token": "token-untuk-sdk-flutter",
			"redirect_url":   "https://app.sandbox.midtrans.com/snap/v2/vtweb/...",
		},
	})
}

// BatalkanPesanan handles canceling an order
func BatalkanPesanan(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Pesanan berhasil dibatalkan",
	})
}

// GetDetailPesananCustomer handles fetching single order details
func GetDetailPesananCustomer(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Detail pesanan berhasil diambil",
		"data": gin.H{
			"no_pesanan":    "12345678",
			"status":        "diproses",
			"tanggal_pesan": "2026-05-05T22:39:00Z",
			"items": []gin.H{
				{
					"id_barang":    101,
					"nama_barang":  "Kipas Angin Portable",
					"varian":       "Putih",
					"jumlah":       1,
					"harga_satuan": 50000,
					"gambar":       "https://api.mantra.com/storage/barang/kipas.jpg",
				},
			},
			"tujuan_pengantaran": gin.H{
				"nama_penerima":  "Ibu Yunani",
				"alamat_lengkap": "Jl. Melati Merah No. 35, Surakarta 50341",
			},
			"kurir": gin.H{
				"nama_kurir": "Ricardo Holahilo",
				"plat_nomor": "H 6582 TH",
				"ekspedisi":  "SPEX Express",
				"foto_kurir": "https://api.mantra.com/storage/kurir/ricardo.jpg",
			},
			"rincian_pembayaran": gin.H{
				"subtotal_items": 150000,
				"ongkir":         20000,
				"biaya_proteksi": 2000,
				"total":          172000,
			},
		},
	})
}

// LacakPesanan handles fetching courier location
func LacakPesanan(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data lacak pesanan berhasil diambil",
		"data": gin.H{
			"id_pesanan": "12345678",
			"kurir": gin.H{
				"nama":       "Ricardo Holahilo",
				"plat_nomor": "H 6582 TH",
				"foto":       "https://api.mantra.com/storage/kurir/ricardo.jpg",
			},
			"lokasi_kurir": gin.H{
				"latitude":  -7.052,
				"longitude": 110.439,
			},
			"estimasi_tiba": "8 mins",
			"jarak_meter":   1500,
		},
	})
}

// GetProfilCustomer handles fetching customer profile
func GetProfilCustomer(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data profil berhasil diambil",
		"data": gin.H{
			"user": gin.H{
				"nama_lengkap": "Aarav Lysander",
				"no_telp":      "+62 81222222222",
				"email":        "lysander@gmail.com",
				"username":     "aarav_",
			},
			"daftar_alamat": []gin.H{
				{
					"id_alamat":        1,
					"label_alamat":     "Rumah",
					"nama_penerima":    "Aarav",
					"no_telp_penerima": "0812...",
					"alamat_lengkap":   "Jl. Cempaka Putih No. 12...",
					"is_utama":         true,
				},
			},
		},
	})
}

// UpdateAkunCustomer handles updating customer account info
func UpdateAkunCustomer(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Informasi akun berhasil diperbarui",
		"data": gin.H{
			"nama_lengkap": "Aarav Lysander",
			"no_telp":      "+62 81222222222",
			"email":        "lysander@gmail.com",
			"username":     "aarav_new",
		},
	})
}

// TambahAlamat handles adding new address
func TambahAlamat(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Alamat baru berhasil ditambahkan",
		"data": gin.H{
			"id_alamat":    15,
			"label_alamat": "Rumah",
			"is_utama":     false,
		},
	})
}

// UpdateAlamat handles updating address
func UpdateAlamat(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Alamat berhasil diperbarui",
		"data": gin.H{
			"id_alamat":    1,
			"label_alamat": "Kos",
			"is_utama":     true,
		},
	})
}

// HapusAlamat handles deleting address
func HapusAlamat(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Alamat berhasil dihapus",
	})
}
