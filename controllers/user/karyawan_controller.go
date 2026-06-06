package user

import (
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"backend-mantra/config"
	"backend-mantra/models"
	"backend-mantra/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// GetDaftarKaryawan mengambil semua karyawan (baik kasir maupun kurir).
func GetDaftarKaryawan(c *gin.Context) {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 { page = 1 }
	if limit < 1 { limit = 10 }
	offset := (page - 1) * limit

	var karyawans []models.Karyawan
	var total int64

	query := config.DB.Model(&models.Karyawan{}).Joins("User").Joins("User.Role")
	if search != "" {
		query = query.Where("User.nama_lengkap ILIKE ? OR User.username ILIKE ? OR User.email ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	
	role := c.Query("role")
	if role != "" && role != "Semua Role" {
		query = query.Where("\"User__Role\".nama_role = ?", role)
	}

	status := c.Query("status")
	if status != "" && status != "Semua Status" {
		query = query.Where("karyawan.status = ?", status)
	}

	query.Count(&total)

	if err := query.Preload("User").Preload("User.Role").Offset(offset).Limit(limit).Find(&karyawans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data karyawan",
		})
		return
	}

	var response []gin.H
	baseURL := os.Getenv("BASE_URL")

	for _, k := range karyawans {
		fotoProfil := k.User.FotoProfil
		if fotoProfil != "" && !strings.HasPrefix(fotoProfil, "http") && baseURL != "" {
			fotoProfil = strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(fotoProfil, "/")
		}

		response = append(response, gin.H{
			"id_karyawan":   k.IdKaryawan,
			"public_id":     k.PublicId,
			"id_user":       k.User.IdUser,
			"nama_lengkap":  k.User.NamaLengkap,
			"username":      k.User.Username,
			"email":         k.User.Email,
			"role":          k.User.Role.NamaRole,
			"no_telp":       k.NoTelp,
			"status":        k.Status,
			"foto_profil":   fotoProfil,
			"terakhir_login": "Belum pernah", // TODO: Implement using RefreshToken table if needed
			"inisial":       getInisial(k.User.NamaLengkap),
		})
	}
	
	if response == nil { response = []gin.H{} }

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Daftar karyawan berhasil diambil",
		"data":    response,
		"meta": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": math.Ceil(float64(total) / float64(limit)),
		},
	})
}

// TambahKaryawan menambah data karyawan baru.
func TambahKaryawan(c *gin.Context) {
	var input struct {
		Username           string `json:"username" binding:"required"`
		Email              string `json:"email" binding:"required"`
		Password           string `json:"password" binding:"required"`
		NamaLengkap        string `json:"nama_lengkap" binding:"required"`
		NoTelp             string `json:"no_telp"`
		TempatLahir        string `json:"tempat_lahir"`
		TanggalLahir       string `json:"tanggal_lahir"`
		JenisKelamin       string `json:"jenis_kelamin"`
		Alamat             string `json:"alamat"`
		PendidikanTerakhir string `json:"pendidikan_terakhir"`
		Nik                string `json:"nik"`
		Role               string `json:"role" binding:"required"`
		Shift              string `json:"shift"`
		FotoProfil         string `json:"foto_profil"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "error",
			"message": "Validasi gagal",
			"error":   err.Error(),
		})
		return
	}

	// Cek Role Exist
	var role models.Role
	if err := config.DB.Where("nama_role = ?", input.Role).First(&role).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Role tidak valid",
		})
		return
	}

	// Check duplicates
	var existing models.User
	if err := config.DB.Where("username = ?", input.Username).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"status": "error", "message": "Username sudah terdaftar"})
		return
	}
	if err := config.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"status": "error", "message": "Email sudah terdaftar"})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 12)

	tx := config.DB.Begin()

	newUser := models.User{
		Username:    input.Username,
		Email:       input.Email,
		Password:    string(hashed),
		NamaLengkap: input.NamaLengkap,
		RoleID:      role.IdRole,
		FotoProfil:  input.FotoProfil,
	}
	if err := tx.Create(&newUser).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal buat user"})
		return
	}

	var tglLahir time.Time
	if input.TanggalLahir != "" {
		tglLahir, _ = time.Parse("2006-01-02", input.TanggalLahir)
	}

	newKaryawan := models.Karyawan{
		NoTelp:             input.NoTelp,
		TempatLahir:        input.TempatLahir,
		TanggalLahir:       tglLahir,
		JenisKelamin:       input.JenisKelamin,
		Alamat:             input.Alamat,
		PendidikanTerakhir: input.PendidikanTerakhir,
		Nik:                input.Nik,
		UserId:             newUser.IdUser,
		Status:             "Aktif",
	}

	if err := tx.Create(&newKaryawan).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal buat karyawan"})
		return
	}

	if input.Role == "Kasir" {
		kasir := models.Kasir{KaryawanId: newKaryawan.IdKaryawan, Shift: input.Shift}
		tx.Create(&kasir)
	} else if input.Role == "Kurir" {
		kurir := models.Kurir{KaryawanId: newKaryawan.IdKaryawan}
		tx.Create(&kurir)
	}

	tx.Commit()
	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Karyawan ditambahkan"})
}

func HapusKaryawan(c *gin.Context) {
	id := c.Param("public_id")
	var karyawan models.Karyawan
	if err := config.DB.Preload("User").Preload("User.Role").Where("public_id = ?", id).First(&karyawan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Karyawan tidak ditemukan"})
		return
	}

	tx := config.DB.Begin()

	if karyawan.User.Role.NamaRole == "Kasir" {
		tx.Where("id_karyawan = ?", karyawan.IdKaryawan).Delete(&models.Kasir{})
	} else if karyawan.User.Role.NamaRole == "Kurir" {
		tx.Where("id_karyawan = ?", karyawan.IdKaryawan).Delete(&models.Kurir{})
	}

	tx.Where("id_karyawan = ?", karyawan.IdKaryawan).Delete(&models.Karyawan{})
	tx.Where("id_user = ?", karyawan.UserId).Delete(&models.User{})

	now := time.Now()
	tx.Model(&models.RefreshToken{}).Where("id_user = ? AND revoked_at IS NULL", karyawan.UserId).Update("revoked_at", &now)

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Karyawan berhasil dihapus"})
}

func getInisial(name string) string {
	parts := strings.Split(strings.TrimSpace(name), " ")
	if len(parts) == 0 {
		return "U"
	}
	if len(parts) == 1 {
		return strings.ToUpper(string(parts[0][0]))
	}
	return strings.ToUpper(string(parts[0][0]) + string(parts[1][0]))
}

func GetDetailKaryawan(c *gin.Context) {
	id := c.Param("public_id")
	var karyawan models.Karyawan
	if err := config.DB.Preload("User").Preload("User.Role").Where("public_id = ?", id).First(&karyawan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Karyawan tidak ditemukan"})
		return
	}

	var shift string
	if karyawan.User.Role.NamaRole == "Kasir" {
		var kasir models.Kasir
		config.DB.Where("id_karyawan = ?", karyawan.IdKaryawan).First(&kasir)
		shift = kasir.Shift
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"public_id":           karyawan.PublicId,
			"username":            karyawan.User.Username,
			"email":               karyawan.User.Email,
			"nama_lengkap":        karyawan.User.NamaLengkap,
			"no_telp":             karyawan.NoTelp,
			"tempat_lahir":        karyawan.TempatLahir,
			"tanggal_lahir":       karyawan.TanggalLahir.Format("2006-01-02"),
			"jenis_kelamin":       karyawan.JenisKelamin,
			"alamat":              karyawan.Alamat,
			"pendidikan_terakhir": karyawan.PendidikanTerakhir,
			"nik":                 karyawan.Nik,
			"role":                karyawan.User.Role.NamaRole,
			"shift":               shift,
			"status":              karyawan.Status,
			"foto_profil":         karyawan.User.FotoProfil,
			"dibuat_pada":         "Tidak tersedia", // TODO: Tambahkan field created_at di tabel
			"login_terakhir":      "Belum pernah",   // TODO: Ambil dari refresh token atau tracking login
		},
	})
}

func UpdateKaryawan(c *gin.Context) {
	id := c.Param("public_id")
	var input struct {
		Username           string `json:"username"`
		Email              string `json:"email"`
		Password           string `json:"password"`
		NamaLengkap        string `json:"nama_lengkap"`
		NoTelp             string `json:"no_telp"`
		TempatLahir        string `json:"tempat_lahir"`
		TanggalLahir       string `json:"tanggal_lahir"`
		JenisKelamin       string `json:"jenis_kelamin"`
		Alamat             string `json:"alamat"`
		PendidikanTerakhir string `json:"pendidikan_terakhir"`
		Nik                string `json:"nik"`
		Status             string `json:"status"`
		Shift              string `json:"shift"`
		FotoProfil         string `json:"foto_profil"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "error", "message": "Validasi gagal"})
		return
	}

	var karyawan models.Karyawan
	if err := config.DB.Preload("User").Preload("User.Role").Where("public_id = ?", id).First(&karyawan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Karyawan tidak ditemukan"})
		return
	}

	tx := config.DB.Begin()

	// Update User
	if input.Username != "" { karyawan.User.Username = input.Username }
	if input.Email != "" { karyawan.User.Email = input.Email }
	if input.NamaLengkap != "" { karyawan.User.NamaLengkap = input.NamaLengkap }
	if input.FotoProfil != "" { karyawan.User.FotoProfil = input.FotoProfil }
	if input.Password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
		karyawan.User.Password = string(hashed)
	}
	if err := tx.Save(&karyawan.User).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal update user"})
		return
	}

	// Update Karyawan
	if input.NoTelp != "" { karyawan.NoTelp = input.NoTelp }
	if input.TempatLahir != "" { karyawan.TempatLahir = input.TempatLahir }
	if input.TanggalLahir != "" {
		tgl, _ := time.Parse("2006-01-02", input.TanggalLahir)
		karyawan.TanggalLahir = tgl
	}
	if input.JenisKelamin != "" { karyawan.JenisKelamin = input.JenisKelamin }
	if input.Alamat != "" { karyawan.Alamat = input.Alamat }
	if input.PendidikanTerakhir != "" { karyawan.PendidikanTerakhir = input.PendidikanTerakhir }
	if input.Nik != "" { karyawan.Nik = input.Nik }
	if input.Status != "" { karyawan.Status = input.Status }

	if err := tx.Save(&karyawan).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal update karyawan"})
		return
	}

	// Update Shift if Kasir
	if karyawan.User.Role.NamaRole == "Kasir" && input.Shift != "" {
		tx.Model(&models.Kasir{}).Where("id_karyawan = ?", karyawan.IdKaryawan).Update("shift", input.Shift)
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Karyawan berhasil diupdate"})
}

// UploadFotoKaryawan mengunggah gambar foto profil karyawan ke MinIO.
// Dipakai oleh: admin (POST /admin/karyawan/upload)
// Auth: Wajib login, role admin
func UploadFotoKaryawan(c *gin.Context) {
	// Biasanya nama form datanya adalah "file" atau "foto", mari kita coba sesuaikan dengan util yang ada.
	// Jika dari sisi klien mengirim field "foto" atau "file", utils.UploadFileToMinio mungkin pakai field itu.
	// Mari kita asumsi dari parameter fungsi utamanya.
	fileUrl, err := utils.UploadFileToMinio(c, "foto", "karyawan")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Gagal mengunggah foto profil: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Foto berhasil diunggah ke server storage",
		"url":     fileUrl,
	})
}
