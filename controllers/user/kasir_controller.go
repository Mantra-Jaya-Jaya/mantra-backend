package user

import (
	"net/http"
	"os"
	"strings"

	"backend-mantra/config"
	"backend-mantra/models"

	"github.com/gin-gonic/gin"
)

// GetProfilKasir mengambil data profil kasir yang sedang login.
// Dipakai oleh: kasir (GET /kasir/profil)
// Auth: Wajib login, role kasir
func GetProfilKasir(c *gin.Context) {
	uid, exists := getUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "User belum login"})
		return
	}

	var kasir models.Kasir
	if err := config.DB.Joins("JOIN karyawan ON karyawan.id_karyawan = kasir.id_karyawan").
		Preload("Karyawan").Preload("Karyawan.User").
		Where("karyawan.id_user = ?", uid).First(&kasir).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Data kasir tidak ditemukan"})
		return
	}

	fotoProfil := kasir.Karyawan.User.FotoProfil
	if fotoProfil != "" && !strings.HasPrefix(fotoProfil, "http") {
		baseURL := os.Getenv("BASE_URL")
		if baseURL != "" {
			fotoProfil = strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(fotoProfil, "/")
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data profil berhasil diambil",
		"data": gin.H{
			"id_kasir":            kasir.IdKasir,
			"public_id":           kasir.Karyawan.PublicId,
			"no_telp":             kasir.Karyawan.NoTelp,
			"tempat_lahir":        kasir.Karyawan.TempatLahir,
			"tanggal_lahir":       kasir.Karyawan.TanggalLahir,
			"jenis_kelamin":       kasir.Karyawan.JenisKelamin,
			"alamat":              kasir.Karyawan.Alamat,
			"pendidikan_terakhir": kasir.Karyawan.PendidikanTerakhir,
			"nik":                 kasir.Karyawan.Nik,
			"id_user":             kasir.Karyawan.User.IdUser,
			"user_public_id":      kasir.Karyawan.User.PublicId,
			"username":            kasir.Karyawan.User.Username,
			"email":               kasir.Karyawan.User.Email,
			"nama_lengkap":        kasir.Karyawan.User.NamaLengkap,
			"foto_profil":         fotoProfil,
		},
	})
}
