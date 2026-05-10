package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetDashboardAdmin handles fetching admin dashboard data
func GetDashboardAdmin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data dashboard berhasil diambil",
		"data": gin.H{
			"penjualan_hari_ini": 5000000,
			"penjualan_mingguan": []gin.H{
				{
					"tanggal": "2026-05-01T00:00:00Z",
					"jumlah":  3500000,
				},
			},
			"stok_menipis": []gin.H{
				{
					"id_barang":   1,
					"nama_barang": "Laptop Gaming X",
					"varian":      "16GB RAM",
					"stok":        3,
				},
			},
		},
	})
}

// GetDaftarBarangAdmin handles fetching all products for admin
func GetDaftarBarangAdmin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Daftar barang berhasil diambil",
		"data": []gin.H{
			{
				"id_barang":    1,
				"nama_barang":  "Laptop Gaming X",
				"kategori":     "Elektronik",
				"satuan":       "Unit",
				"total_stok":   35,
				"gambar":       "https://api.mantra.com/storage/barang/laptop-x.jpg",
				"punya_diskon": true,
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

// TambahBarang handles creating a new product
func TambahBarang(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Barang berhasil ditambahkan",
		"data": gin.H{
			"id_barang":   51,
			"nama_barang": "Laptop Gaming X",
		},
	})
}

// GetDetailBarangAdmin handles fetching single product details for admin
func GetDetailBarangAdmin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Detail barang berhasil diambil",
		"data": gin.H{
			"id_barang":   1,
			"nama_barang": "Laptop Gaming X",
			"kategori":    "Elektronik",
			"satuan":      "Unit",
			"gambar":      "https://api.mantra.com/storage/barang/laptop-x.jpg",
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
					"stok":                  5,
				},
			},
		},
	})
}

// UpdateBarang handles updating product details
func UpdateBarang(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Barang berhasil diperbarui",
	})
}

// HapusBarang handles deleting product
func HapusBarang(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Barang berhasil dihapus",
	})
}

// TambahDiskonBarang handles adding discount to product
func TambahDiskonBarang(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Diskon berhasil ditambahkan",
		"data": gin.H{
			"id_diskon":    3,
			"nama_diskon":  "Promo Lebaran",
			"besar_diskon": 15,
		},
	})
}

// GetDaftarKaryawan handles fetching all employees
func GetDaftarKaryawan(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Daftar karyawan berhasil diambil",
		"data": []gin.H{
			{
				"id_user":      5,
				"nama_lengkap": "Budi Santoso",
				"email":        "budi@mantra.com",
				"role":         "kasir",
			},
		},
	})
}

// GetDaftarKasir handles fetching all cashiers
func GetDaftarKasir(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Daftar kasir berhasil diambil",
		"data": []gin.H{
			{
				"id_kasir":     1,
				"nama_lengkap": "Budi Santoso",
				"foto":         "https://api.mantra.com/storage/kasir/budi.jpg",
			},
		},
	})
}

// TambahKasir handles creating a new cashier
func TambahKasir(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Kasir berhasil ditambahkan",
		"data": gin.H{
			"id_kasir":     5,
			"nama_lengkap": "Budi Santoso",
		},
	})
}

// GetDetailKasir handles fetching single cashier details
func GetDetailKasir(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Detail kasir berhasil diambil",
		"data": gin.H{
			"id_kasir":      1,
			"nama_lengkap":  "Budi Santoso",
			"email":         "budi@mantra.com",
			"no_telp":        "081234567890",
			"alamat":         "Jl. Mawar No. 5",
			"tanggal_lahir": "1995-03-15",
			"shift":         "Pagi",
			"foto":          "https://api.mantra.com/storage/kasir/budi.jpg",
		},
	})
}

// UpdateKasir handles updating cashier details
func UpdateKasir(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data kasir berhasil diperbarui",
	})
}

// HapusKasir handles deleting cashier
func HapusKasir(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Kasir berhasil dihapus",
	})
}

// GetDaftarKurir handles fetching all couriers
func GetDaftarKurir(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Daftar kurir berhasil diambil",
		"data": []gin.H{
			{
				"id_kurir":     1,
				"nama_lengkap": "Udin Kurir",
				"foto":         "https://api.mantra.com/storage/kurir/udin.jpg",
			},
		},
	})
}

// TambahKurir handles creating a new courier
func TambahKurir(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Kurir berhasil ditambahkan",
		"data": gin.H{
			"id_kurir":     5,
			"nama_lengkap": "Udin Kurir",
		},
	})
}

// GetDetailKurir handles fetching single courier details
func GetDetailKurir(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Detail kurir berhasil diambil",
		"data": gin.H{
			"id_kurir":      1,
			"nama_lengkap":  "Udin Kurir",
			"email":         "udin@mantra.com",
			"no_telp":        "081234567891",
			"alamat":         "Jl. Melati No. 6",
			"tanggal_lahir": "1996-04-16",
			"foto":          "https://api.mantra.com/storage/kurir/udin.jpg",
		},
	})
}

// UpdateKurir handles updating courier details
func UpdateKurir(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data kurir berhasil diperbarui",
	})
}

// HapusKurir handles deleting courier
func HapusKurir(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Kurir berhasil dihapus",
	})
}

// GetNotifikasiAdmin handles fetching admin notifications
func GetNotifikasiAdmin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Notifikasi admin berhasil diambil",
		"data": []gin.H{
			{
				"id_notifikasi": 1,
				"id_barang":     1,
				"nama_barang":   "Laptop Gaming X",
				"varian":        "16GB RAM",
				"stok_saat_ini": 3,
				"batas_minimum": 5,
				"pesan":         "Stok Laptop Gaming X (16GB RAM) hampir habis",
				"created_at":    "2026-05-09T08:00:00Z",
			},
		},
	})
}

// GetProfilAdmin handles fetching admin profile
func GetProfilAdmin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data profil berhasil diambil",
		"data": gin.H{
			"nama_lengkap": "Admin Mantra",
			"username":     "admin_mantra",
			"foto":         "https://api.mantra.com/storage/admin/admin.jpg",
		},
	})
}

// UpdateProfilAdmin handles updating admin profile
func UpdateProfilAdmin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Profil admin berhasil diperbarui",
		"data": gin.H{
			"nama_lengkap": "Admin Mantra",
			"username":     "admin_mantra_baru",
		},
	})
}
