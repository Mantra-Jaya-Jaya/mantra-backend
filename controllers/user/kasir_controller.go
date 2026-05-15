package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetProfilKasir mengambil data profil kasir yang sedang login.
// Dipakai oleh: kasir (GET /kasir/profil)
// Auth: Wajib login, role kasir

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

// GetDaftarKasir mengambil daftar semua kasir.
// Dipakai oleh: admin (GET /admin/user/kasir)
// Auth: Wajib login, role admin
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

// TambahKasir membuat akun kasir baru.
// Dipakai oleh: admin (POST /admin/user/kasir)
// Auth: Wajib login, role admin
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

// GetDetailKasir mengambil detail profil satu kasir.
// Dipakai oleh: admin (GET /admin/user/kasir/:id_kasir)
// Auth: Wajib login, role admin
func GetDetailKasir(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Detail kasir berhasil diambil",
		"data": gin.H{
			"id_kasir":      1,
			"nama_lengkap":  "Budi Santoso",
			"email":         "budi@mantra.com",
			"no_telp":       "081234567890",
			"alamat":        "Jl. Mawar No. 5",
			"tanggal_lahir": "1995-03-15",
			"shift":         "Pagi",
			"foto":          "https://api.mantra.com/storage/kasir/budi.jpg",
		},
	})
}

// UpdateKasir memperbarui data kasir.
// Dipakai oleh: admin (PUT /admin/user/kasir/:id_kasir)
// Auth: Wajib login, role admin
func UpdateKasir(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data kasir berhasil diperbarui",
	})
}

// HapusKasir menghapus akun kasir.
// Dipakai oleh: admin (DELETE /admin/user/kasir/:id_kasir)
// Auth: Wajib login, role admin
func HapusKasir(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Kasir berhasil dihapus",
	})
}

// GetDaftarKaryawan mengambil semua karyawan (kasir + kurir).
// Dipakai oleh: admin (GET /admin/user/karyawan)
// Auth: Wajib login, role admin
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
