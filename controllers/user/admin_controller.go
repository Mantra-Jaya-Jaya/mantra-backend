package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetProfilAdmin mengambil data profil admin yang sedang login.
// Dipakai oleh: admin (GET /admin/profil)
// Auth: Wajib login, role admin
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

// UpdateProfilAdmin memperbarui data profil admin.
// Dipakai oleh: admin (PUT /admin/profil)
// Auth: Wajib login, role admin
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

// GetDashboardAdmin mengambil data ringkasan dashboard admin.
// Dipakai oleh: admin (GET /admin/dashboard)
// Auth: Wajib login, role admin
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
