package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetDaftarKurir mengambil daftar semua kurir.
// Dipakai oleh: admin (GET /admin/user/kurir)
// Auth: Wajib login, role admin
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

// TambahKurir membuat akun kurir baru.
// Dipakai oleh: admin (POST /admin/user/kurir)
// Auth: Wajib login, role admin
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

// GetDetailKurir mengambil detail profil satu kurir.
// Dipakai oleh: admin (GET /admin/user/kurir/:id_kurir)
// Auth: Wajib login, role admin
func GetDetailKurir(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Detail kurir berhasil diambil",
		"data": gin.H{
			"id_kurir":      1,
			"nama_lengkap":  "Udin Kurir",
			"email":         "udin@mantra.com",
			"no_telp":       "081234567891",
			"alamat":        "Jl. Melati No. 6",
			"tanggal_lahir": "1996-04-16",
			"foto":          "https://api.mantra.com/storage/kurir/udin.jpg",
		},
	})
}

// UpdateKurir memperbarui data kurir.
// Dipakai oleh: admin (PUT /admin/user/kurir/:id_kurir)
// Auth: Wajib login, role admin
func UpdateKurir(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data kurir berhasil diperbarui",
	})
}

// HapusKurir menghapus akun kurir.
// Dipakai oleh: admin (DELETE /admin/user/kurir/:id_kurir)
// Auth: Wajib login, role admin
func HapusKurir(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Kurir berhasil dihapus",
	})
}
