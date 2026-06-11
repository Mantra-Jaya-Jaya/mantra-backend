package pengantaran

import (
	"net/http"
	"time"

	"backend-mantra/config"
	"backend-mantra/models"

	"github.com/gin-gonic/gin"
)

// GetDaftarPengantaran mengambil daftar pengantaran yang sedang aktif.
// Dipakai oleh: admin (GET /admin/pengantaran), kurir (GET /kurir/pengantaran)
// Auth: Wajib login, role admin atau kurir (dikontrol di route)
func GetDaftarPengantaran(c *gin.Context) {
	role := c.GetString("role")

	// 🚀 1. PENJINAK TOKEN (Gantiin c.GetInt64 yang rawan bug)
	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Auth salah: User belum login"})
		return
	}

	var userID int64
	switch v := val.(type) {
	case float64:
		userID = int64(v)
	case int64:
		userID = v
	case int:
		userID = int64(v)
	case uint:
		userID = int64(v)
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Auth salah: Format token tidak valid"})
		return
	}

	var pengantarans []models.Pengantaran
	
	// Preload standar (gak perlu sampai ke spesifikasi barang karena ini halaman daftar tugas)
	query := config.DB.Preload("Pesanan").
		Preload("Pesanan.Alamat").
		Preload("Pesanan.Customer.User").
		Preload("StatusPengantaran").
		Preload("Ekspedisi")

	// 🚀 2. FILTER UTAMA: Khusus untuk Kurir
	if role == "kurir" {
		var result struct{ IdKurir uint }
		if err := config.DB.Raw("SELECT id_kurir FROM kurir JOIN karyawan ON kurir.id_karyawan = karyawan.id_karyawan WHERE karyawan.id_user = ?", userID).Scan(&result).Error; err != nil || result.IdKurir == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Data kurir tidak ditemukan. Pastikan Anda login sebagai kurir.",
			})
			return
		}
		// Kurir cuma bisa liat tugasnya sendiri
		query = query.Where("id_kurir = ?", result.IdKurir)
	}

	// 🚀 3. EKSEKUSI QUERY
	if err := query.Find(&pengantarans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil daftar pengantaran",
			"data":    nil,
		})
		return
	}

	// 🚀 4. REFACTOR JSON (Bikin Cetakan DTO biar Langsing)
	type PengantaranRingkas struct {
		PublicID        string      `json:"public_id"`
		Status          string      `json:"status"`
		Ekspedisi       string      `json:"ekspedisi"`
		WaktuPickup     interface{} `json:"waktu_pickup"` // Pakai interface{} biar aman nampung null / time
		WaktuSampai     interface{} `json:"waktu_sampai"`
		NamaCustomer    string      `json:"nama_customer"`
		NoTelp          string      `json:"no_telp"`
		AlamatLengkap   string      `json:"alamat_lengkap"`
		TotalPendapatan int         `json:"total_pendapatan"` 
	}

	var hasilAkhir []PengantaranRingkas

	// 🚀 5. LOOPING DAN MAPPING DATA
	for _, p := range pengantarans {
		// Fallback data alamat
		namaCust := "Customer Offline"
		noTelp := "-"
		alamatLengkap := "Ambil di Toko"

		// Amanin Alamat
		if p.Pesanan.Alamat != nil {
			namaCust = p.Pesanan.Alamat.NamaPenerima
			noTelp = p.Pesanan.Alamat.NoTelpPenerima
			alamatLengkap = p.Pesanan.Alamat.AlamatLengkap
		} else if p.Pesanan.Customer.User.NamaLengkap != "" {
			namaCust = p.Pesanan.Customer.User.NamaLengkap
		}

		// 🚀 AMANIN EKSPEDISI BIAR GAK PANIC
		namaEkspedisi := "Internal / Belum Ada"
		if p.Ekspedisi != nil {
			namaEkspedisi = p.Ekspedisi.NamaEkspedisi
		}

		// 🚀 AMANIN STATUS BIAR GAK PANIC
		namaStatus := "Menunggu"
		if p.StatusPengantaran != nil {
			namaStatus = p.StatusPengantaran.NamaStatus
		}

		hasilAkhir = append(hasilAkhir, PengantaranRingkas{
			PublicID:        p.PublicId.String(),
			Status:          namaStatus,     // Pakai variabel yang udah aman
			Ekspedisi:       namaEkspedisi,  // Pakai variabel yang udah aman
			WaktuPickup:     p.WaktuPickup,
			WaktuSampai:     p.WaktuSampai,
			NamaCustomer:    namaCust,
			NoTelp:          noTelp,
			AlamatLengkap:   alamatLengkap,
			TotalPendapatan: 35000, 
		})
	}

	// Biar array-nya tetap [] bukan null kalau lagi gak ada tugas
	if len(hasilAkhir) == 0 {
		hasilAkhir = []PengantaranRingkas{}
	}

	// 🚀 6. RESPONSE SUKSES
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Daftar pengantaran berhasil diambil",
		"data":    hasilAkhir,
	})
}

// UpdateLokasiKurir memperbarui koordinat lokasi kurir yang sedang bertugas.
// Dipakai oleh: kurir (PATCH /kurir/pengantaran/:id_pengantaran/lokasi)
// Auth: Wajib login, role kurir
// Ownership: kurir hanya bisa update lokasi pengantaran yang ditugaskan kepadanya
func UpdateLokasiKurir(c *gin.Context) {
	idPengantaran := c.Param("id_pengantaran")
	userID := c.GetInt64("user_id")

	type UpdateLokasiInput struct {
		Latitude  float64 `json:"latitude" binding:"required"`
		Longitude float64 `json:"longitude" binding:"required"`
	}

	var input UpdateLokasiInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Format inputan salah, pastikan latitude dan longitude diisi dengan benar",
		})
		return
	}

	// 1. Cari id_kurir berdasarkan user_id dari JWT
	var result struct{ IdKurir uint }
	if err := config.DB.Raw("SELECT id_kurir FROM kurir JOIN karyawan ON kurir.id_karyawan = karyawan.id_karyawan WHERE karyawan.id_user = ?", userID).Scan(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengidentifikasi kurir",
		})
		return
	}
	kurirID := result.IdKurir

	// 2. Cari data pengantaran dan periksa kepemilikan (ownership check)
	var pengantaran models.Pengantaran
	if err := config.DB.Where("public_id = ?", idPengantaran).First(&pengantaran).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data pengantaran tidak ditemukan",
		})
		return
	}

	if pengantaran.KurirID == nil || *pengantaran.KurirID != kurirID {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "Anda tidak memiliki akses ke resource ini",
			"error": gin.H{
				"code":   "AUTH_002",
				"detail": "Pengantaran ini tidak ditugaskan kepada Anda",
			},
		})
		return
	}

	// 3. Update koordinat
	pengantaran.LastLatitude = input.Latitude
	pengantaran.LastLongitude = input.Longitude

	if err := config.DB.Save(&pengantaran).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memperbarui lokasi kurir",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Lokasi kurir berhasil diperbarui",
	})
}

func GetLaporanHariIni(c *gin.Context) {
	// 🚀 PENJINAK TOKEN (Tetap sama)
	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "User belum login"})
		return
	}

	var userID int64
	switch v := val.(type) {
	case float64:
		userID = int64(v)
	case int64:
		userID = v
	case int:
		userID = int64(v)
	case uint:
		userID = int64(v)
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Format token tidak valid"})
		return
	}

	// 1. Cari ID Kurir (Tetap sama)
	var result struct{ IdKurir uint }
	if err := config.DB.Raw("SELECT id_kurir FROM kurir JOIN karyawan ON kurir.id_karyawan = karyawan.id_karyawan WHERE karyawan.id_user = ?", userID).Scan(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengidentifikasi kurir"})
		return
	}

	if result.IdKurir == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Data kurir tidak ditemukan."})
		return
	}

	idKurir := result.IdKurir

	// 2. Set rentang waktu HARI INI (Tetap sama)
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	// 🚀 SIAPIN 3 WADAH VARIABEL SEKARANG
	var selesaiCount int64
	var belumSelesaiCount int64
	var pesananBaruCount int64 

	// 3. QUERY 1: Hitung Pesanan yang "SELESAI" hari ini
	config.DB.Model(&models.Pengantaran{}).
		Where("id_kurir = ?", idKurir).
		Where("id_status_pengantaran = ?", 4).
		Where("waktu_sampai >= ? AND waktu_sampai < ?", startOfDay, endOfDay).
		Count(&selesaiCount)

	// 4. QUERY 2: Hitung Pesanan yang "BELUM SELESAI" (Masih dipegang kurir ini)
	config.DB.Model(&models.Pengantaran{}).
		Where("id_kurir = ?", idKurir).
		Where("id_status_pengantaran != ?", 4).
		Count(&belumSelesaiCount)

	// 🚀 5. QUERY 3 (BARU!): Hitung Pesanan Online yang NGANGGUR / Siap Direbut
	// Kita hitung dari tabel Pesanan langsung yang statusnya siap antar
	config.DB.Model(&models.Pesanan{}).
		Where("tipe_pesanan = ?", "Online").
		Where("status_pesanan = ?", "Dikemas"). // NOTE: Sesuaikan nama status lu kalau beda!
		Count(&pesananBaruCount)

	// 6. Kembalikan 3 data tersebut ke Flutter
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Laporan hari ini berhasil diambil",
		"data": gin.H{
			"pesanan_selesai":  selesaiCount,
			"pesanan_proses":   belumSelesaiCount,
			"pesanan_tersedia": pesananBaruCount, // 🚀 Data incaran lu nangkring di sini!
		},
	})
}


// GetDetailPengantaran mengambil detail lengkap untuk halaman Peta Kurir
func GetDetailPengantaran(c *gin.Context) {
  // 🚀 1. PERBAIKAN PARAMETER: Harus sama kayak di Route!
  idPengantaran := c.Param("public_id") 
  
  // ==========================================
  // PENJINAK TOKEN
  // ==========================================
  val, exists := c.Get("user_id")
  if !exists {
    c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "User belum login"})
    return
  }

  var userID int64
  switch v := val.(type) {
  case float64: userID = int64(v)
  case int64: userID = v
  case int: userID = int64(v)
  case uint: userID = int64(v)
  default:
    c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Format token tidak valid"})
    return
  }

  // Cari ID Kurir
  var result struct{ IdKurir uint }
  if err := config.DB.Raw("SELECT id_kurir FROM kurir JOIN karyawan ON kurir.id_karyawan = karyawan.id_karyawan WHERE karyawan.id_user = ?", userID).Scan(&result).Error; err != nil || result.IdKurir == 0 {
    c.JSON(http.StatusInternalServerError, gin.H{ 
      "status":  "error",
      "code":    500, // Cek juga result.IdKurir == 0 buat jaga-jaga
      "message": "Gagal mengidentifikasi data kurir di server",
    })
    return
  }
  kurirID := result.IdKurir

  // ==========================================
  // TARIK DATA DATABASE
  // ==========================================
  var pengantaran models.Pengantaran
  err := config.DB.
    Preload("Pesanan.Alamat").       
    Preload("Pesanan.Customer.User"). 
    Preload("StatusPengantaran").    
    Where("public_id = ?", idPengantaran).
    First(&pengantaran).Error

  if err != nil {
    c.JSON(http.StatusNotFound, gin.H{ 
      "status":  "error",
      "code":    404,
      "message": "Data pengantaran tidak ditemukan",
    })
    return
  }

  if pengantaran.KurirID == nil || *pengantaran.KurirID != kurirID {
    c.JSON(http.StatusForbidden, gin.H{ 
      "status":  "error",
      "code":    403,
      "message": "Akses Ditolak: Pengantaran ini tidak ditugaskan kepada Anda",
    })
    return
  }

  // ==========================================
  // 🚀 2. SUSUN DATA DENGAN PELINDUNG NIL POINTER
  // ==========================================
  namaPenerima := "Customer"
  noTelpPenerima := "-"
  alamatLengkap := "Alamat tidak tersedia"
  var lat, lng float64

  // Pelindung Pesanan & Customer
  if pengantaran.Pesanan != nil {
      if pengantaran.Pesanan.Customer.User.NamaLengkap != "" {
          namaPenerima = pengantaran.Pesanan.Customer.User.NamaLengkap
      }
      // Asumsi ada field NoTelp di Customer lu (sesuaikan kalau beda)
      if pengantaran.Pesanan.Customer.NoTelp != "" {
          noTelpPenerima = pengantaran.Pesanan.Customer.NoTelp
      }

      // Pelindung Alamat
      if pengantaran.Pesanan.Alamat != nil {
          alamatLengkap = pengantaran.Pesanan.Alamat.AlamatLengkap
          lat = pengantaran.Pesanan.Alamat.Latitude
          lng = pengantaran.Pesanan.Alamat.Longitude
      }
  }

  // Pelindung Status Pengantaran
  statusNama := "Diproses"
  if pengantaran.StatusPengantaran != nil && pengantaran.StatusPengantaran.NamaStatus != "" {
    statusNama = pengantaran.StatusPengantaran.NamaStatus
  }

  // ✅ [SUCCESS 200: OK]
  c.JSON(http.StatusOK, gin.H{ 
    "status":  "success",
    "code":    200,
    "message": "Detail pengantaran berhasil ditarik",
    "data": gin.H{
      "id_pengantaran":     pengantaran.PublicId,
      "status_pengantaran": statusNama,
      "waktu_pickup":       pengantaran.WaktuPickup,
      "waktu_sampai":       pengantaran.WaktuSampai,
      
      "penerima": gin.H{
        "nama":    namaPenerima,
        "no_telp": noTelpPenerima,
      },
      "tujuan": gin.H{
        "alamat_lengkap": alamatLengkap,
        "latitude":       lat,
        "longitude":      lng,
      },
    },
  })
}