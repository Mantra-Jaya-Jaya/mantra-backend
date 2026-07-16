package user

import (
	"net/http"
	"os"
	"strings"

	"backend-mantra/config"
	"backend-mantra/models"
	"backend-mantra/utils"

	"github.com/gin-gonic/gin"
)

// GetProfilKurir mengambil data profil kurir yang sedang login.
// UploadFotoProfilKurir mengunggah foto profil kurir ke MinIO dan memperbarui database.
func UploadFotoProfilKurir(c *gin.Context) {
	uid, exists := getUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "User belum login"})
		return
	}

	var kurir models.Kurir
	if err := config.DB.Preload("Karyawan").Preload("Karyawan.User").
		Joins("JOIN karyawan ON karyawan.id_karyawan = kurir.id_karyawan").
		Where("karyawan.id_user = ?", uid).
		First(&kurir).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Data kurir tidak ditemukan"})
		return
	}

	fotoURL, err := utils.UploadFileToMinio(c, "foto", "profil")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Gagal mengunggah foto profil: " + err.Error(),
		})
		return
	}

	kurir.Karyawan.User.FotoProfil = fotoURL
	if err := config.DB.Save(&kurir.Karyawan.User).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menyimpan foto profil"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Foto berhasil diunggah",
		"url":     fotoURL,
	})
}

func GetProfilKurir(c *gin.Context) {
	uid, exists := getUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "User belum login"})
		return
	}

	var kurir models.Kurir
	if err := config.DB.Joins("JOIN karyawan ON karyawan.id_karyawan = kurir.id_karyawan").
		Preload("Karyawan").Preload("Karyawan.User").
		Where("karyawan.id_user = ?", uid).First(&kurir).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Data kurir tidak ditemukan"})
		return
	}

	fotoProfil := kurir.Karyawan.User.FotoProfil
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
			"id_kurir":            kurir.IdKurir,
			"public_id":           kurir.Karyawan.PublicId,
			"no_telp":             kurir.Karyawan.NoTelp,
			"tempat_lahir":        kurir.Karyawan.TempatLahir,
			"tanggal_lahir":       kurir.Karyawan.TanggalLahir,
			"jenis_kelamin":       kurir.Karyawan.JenisKelamin,
			"alamat":              kurir.Karyawan.Alamat,
			"pendidikan_terakhir": kurir.Karyawan.PendidikanTerakhir,
			"nik":                 kurir.Karyawan.Nik,
			"id_user":             kurir.Karyawan.User.IdUser,
			"user_public_id":      kurir.Karyawan.User.PublicId,
			"username":            kurir.Karyawan.User.Username,
			"email":               kurir.Karyawan.User.Email,
			"nama_lengkap":        kurir.Karyawan.User.NamaLengkap,
			"foto_profil":         fotoProfil,
		},
	})
}
