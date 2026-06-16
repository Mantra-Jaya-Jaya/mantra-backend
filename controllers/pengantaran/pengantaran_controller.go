package pengantaran

import (
	"net/http"
	"time"

	"backend-mantra/config"
	"backend-mantra/models"
	"backend-mantra/utils"

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
			Status:          namaStatus,    // Pakai variabel yang udah aman
			Ekspedisi:       namaEkspedisi, // Pakai variabel yang udah aman
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
// Dipakai oleh: kurir (PATCH /kurir/pengantaran/:public_id/lokasi)
// Auth: Wajib login, role kurir
// Ownership: kurir hanya bisa update lokasi pengantaran yang ditugaskan kepadanya
func UpdateLokasiKurir(c *gin.Context) {
	idPengantaran := c.Param("public_id")
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
	selesaiPengantaranID := utils.GetStatusPengantaranID("Selesai") // lookup by name, bukan hardcode
	config.DB.Model(&models.Pengantaran{}).
		Where("id_kurir = ?", idKurir).
		Where("id_status_pengantaran = ?", selesaiPengantaranID).
		Where("waktu_sampai >= ? AND waktu_sampai < ?", startOfDay, endOfDay).
		Count(&selesaiCount)

	// 4. QUERY 2: Hitung Pesanan yang "BELUM SELESAI" (Masih dipegang kurir ini)
	config.DB.Model(&models.Pengantaran{}).
		Where("id_kurir = ?", idKurir).
		Where("id_status_pengantaran != ?", selesaiPengantaranID).
		Count(&belumSelesaiCount)

	// 🚀 5. QUERY 3 (BARU!): Hitung Pesanan Online yang NGANGGUR / Siap Direbut
	// Kita hitung dari tabel Pesanan langsung yang statusnya siap antar
	config.DB.Model(&models.Pesanan{}).
		Where("id_tipe_pesanan = ?", utils.GetTipePesananID("Online")).
		Where("id_status_pesanan = ?", utils.GetStatusPesananID("Dikemas")).
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
// Dipakai oleh: kurir (GET /kurir/pengantaran/:public_id/detail)
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
		Preload("Pesanan.DetailPesanan.SpesifikasiBarang.Barang").
		Preload("Pesanan.DetailPesanan.SpesifikasiBarang.DetailSpesifikasi.Spesifikasi").
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

	// 🚀 Build daftar barang dengan nil-guard
	var totalPembayaran int
	var listBarang []gin.H
	if pengantaran.Pesanan != nil {
		totalPembayaran = pengantaran.Pesanan.TotalPembayaran
		for _, detail := range pengantaran.Pesanan.DetailPesanan {
			varian := "Default"
			if detail.SpesifikasiBarang.DetailSpesifikasi.IdDetailSpesifikasi != 0 {
				varian = detail.SpesifikasiBarang.DetailSpesifikasi.Spesifikasi.NamaSpesifikasi + ": " +
					detail.SpesifikasiBarang.DetailSpesifikasi.NamaDetailSpesifikasi
			}
			listBarang = append(listBarang, gin.H{
				"nama_barang":   detail.SpesifikasiBarang.Barang.NamaBarang,
				"variasi":       varian,
				"jumlah_beli":   detail.Jumlah,
				"harga_satuan":  detail.HargaSatuan,
				"subtotal_item": detail.Subtotal,
			})
		}
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
			"foto_bukti":         pengantaran.FotoBuktiPengiriman,
			"total_pembayaran":   totalPembayaran,
			"daftar_barang":      listBarang,

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

func UploadBuktiPengiriman(c *gin.Context) {
	// 🚀 1. TANGKAP PUBLIC ID & CEK AUTH
	idPengantaran := c.Param("public_id")

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

	// Cari ID Kurir buat proteksi
	var result struct{ IdKurir uint }
	if err := config.DB.Raw("SELECT id_kurir FROM kurir JOIN karyawan ON kurir.id_karyawan = karyawan.id_karyawan WHERE karyawan.id_user = ?", userID).Scan(&result).Error; err != nil || result.IdKurir == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengidentifikasi kurir"})
		return
	}
	kurirID := result.IdKurir

	// 🚀 2. CARI DATA PENGANTARAN DARI DATABASE
	var pengantaran models.Pengantaran
	err := config.DB.
		Preload("Pesanan.Alamat").
		Preload("Pesanan.Customer.User").
		Preload("Pesanan.DetailPesanan.SpesifikasiBarang.Barang").
		Preload("Pesanan.DetailPesanan.SpesifikasiBarang.DetailSpesifikasi.Spesifikasi").
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

	// Proteksi: Cuma kurir yang bawa paket ini yang boleh upload!
	if pengantaran.KurirID == nil || *pengantaran.KurirID != kurirID {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": "Akses Ditolak: Bukan paket Anda!"})
		return
	}

	// 🚀 3. UPLOAD GAMBAR KE MINIO
	// Kita set form data key-nya "foto_bukti" dan folder-nya "bukti_pengantaran"
	fileUrl, err := utils.UploadFileToMinio(c, "foto_bukti", "bukti_pengantaran")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Gagal mengunggah foto bukti: " + err.Error(),
		})
		return
	}

	// 🚀 4. UPDATE DATA DI DATABASE (TABEL PENGANTARAN)
	waktuSekarang := time.Now()

	pengantaran.FotoBuktiPengiriman = fileUrl // Simpan URL dari MinIO
	pengantaran.StatusPengantaranID = utils.GetStatusPengantaranID("Selesai")
	pengantaran.WaktuSampai = &waktuSekarang // Catat waktu selesai realtime

	if err := config.DB.Save(&pengantaran).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gambar berhasil diupload, tapi gagal update database"})
		return
	}

	// 🚀 5. UPDATE STATUS DI TABEL PESANAN (Biar sinkron!)
	if pengantaran.Pesanan != nil {
		config.DB.Model(&models.Pesanan{}).
			Where("id_pesanan = ?", pengantaran.PesananID).
			Update("id_status_pesanan", utils.GetStatusPesananID("Selesai"))
	}

	var listBarang []gin.H
	if pengantaran.Pesanan != nil {
		for _, detail := range pengantaran.Pesanan.DetailPesanan {
			varian := "Default"
			if detail.SpesifikasiBarang.DetailSpesifikasi.IdDetailSpesifikasi != 0 {
				varian = detail.SpesifikasiBarang.DetailSpesifikasi.Spesifikasi.NamaSpesifikasi + ": " +
					detail.SpesifikasiBarang.DetailSpesifikasi.NamaDetailSpesifikasi
			}
			listBarang = append(listBarang, gin.H{
				"nama_barang":   detail.SpesifikasiBarang.Barang.NamaBarang,
				"variasi":       varian,
				"jumlah_beli":   detail.Jumlah,
				"harga_satuan":  detail.HargaSatuan,
				"subtotal_item": detail.Subtotal,
			})
		}
	}

	// 🚀 6. KEMBALIKAN RESPON SUKSES
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Mantap! Bukti pengiriman berhasil diupload dan pesanan diselesaikan.",
		"data": gin.H{
			"url_bukti":      fileUrl,
			"waktu_sampai":   waktuSekarang,
			"id_pengantaran": pengantaran.PublicId,
			"waktu_pickup":   pengantaran.WaktuPickup,
			"daftar_barang":  listBarang,
		},
	})
}

// AmbilPesanan membuat/mengupdate record pengantaran untuk kurir yang mengambil pesanan.
// Dipakai oleh: kurir (POST /kurir/pengantaran/:public_id/ambil)
// Auth: Wajib login, role kurir
func AmbilPesanan(c *gin.Context) {
	pesananPublicID := c.Param("public_id")
	userID := c.GetInt64("user_id")

	// 1. Cari ID Kurir berdasarkan user_id dari JWT
	var result struct{ IdKurir uint }
	if err := config.DB.Raw("SELECT id_kurir FROM kurir JOIN karyawan ON kurir.id_karyawan = karyawan.id_karyawan WHERE karyawan.id_user = ?", userID).Scan(&result).Error; err != nil || result.IdKurir == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Data kurir tidak ditemukan",
		})
		return
	}
	kurirID := result.IdKurir

	// 2. Cari Pesanan berdasarkan public_id
	var pesanan models.Pesanan
	if err := config.DB.Preload("StatusPesanan").Where("public_id = ?", pesananPublicID).First(&pesanan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Pesanan tidak ditemukan",
		})
		return
	}

	if pesanan.StatusPesanan.NamaStatus != "Dikemas" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Pesanan tidak dapat diambil karena status bukan Dikemas",
		})
		return
	}

	var status models.StatusPengantaran
	status.IdStatusPengantaran = utils.GetStatusPengantaranID("Dalam Perjalanan")

	tx := config.DB.Begin()

	// Update status pesanan ke "Dikirim"
	pesanan.StatusPesananID = utils.GetStatusPesananID("Dikirim")
	if err := tx.Save(&pesanan).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal memperbarui status pesanan"})
		return
	}

	// 4. Buat atau update record Pengantaran
	var pengantaran models.Pengantaran
	now := time.Now()
	if err := tx.Where("id_pesanan = ?", pesanan.IdPesanan).First(&pengantaran).Error; err == nil {
		pengantaran.KurirID = &kurirID
		pengantaran.WaktuPickup = &now
		pengantaran.StatusPengantaranID = status.IdStatusPengantaran
		if err := tx.Save(&pengantaran).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal memperbarui data pengantaran"})
			return
		}
	} else {
		pengantaran = models.Pengantaran{
			WaktuPickup:         &now,
			PesananID:           pesanan.IdPesanan,
			KurirID:             &kurirID,
			StatusPengantaranID: status.IdStatusPengantaran,
		}
		if err := tx.Create(&pengantaran).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal membuat data pengantaran"})
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Pesanan berhasil diambil",
		"data": gin.H{
			"id_pengantaran": pengantaran.PublicId,
			"status":         "Dalam Perjalanan",
		},
	})
}

// UpdateStatusPengantaran memperbarui status pengantaran dan mengunggah foto bukti jika statusnya Selesai.
// Dipakai oleh: kurir (POST /kurir/pengantaran/:public_id/status)
// Auth: Wajib login, role kurir
func UpdateStatusPengantaran(c *gin.Context) {
	idPengantaran := c.Param("public_id")
	userID := c.GetInt64("user_id")

	// Ambil status dari form-data atau JSON
	statusInput := c.PostForm("status")
	if statusInput == "" {
		var jsonInput struct {
			Status string `json:"status"`
		}
		if err := c.ShouldBindJSON(&jsonInput); err == nil {
			statusInput = jsonInput.Status
		}
	}

	if statusInput == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Status harus diisi"})
		return
	}

	// 1. Cari ID Kurir berdasarkan user_id dari JWT
	var result struct{ IdKurir uint }
	if err := config.DB.Raw("SELECT id_kurir FROM kurir JOIN karyawan ON kurir.id_karyawan = kurir.id_karyawan WHERE karyawan.id_user = ?", userID).Scan(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengidentifikasi kurir"})
		return
	}
	kurirID := result.IdKurir

	// 2. Cari data pengantaran dan periksa kepemilikan
	var pengantaran models.Pengantaran
	if err := config.DB.Preload("Pesanan").Where("public_id = ?", idPengantaran).First(&pengantaran).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Data pengantaran tidak ditemukan"})
		return
	}

	if pengantaran.KurirID == nil || *pengantaran.KurirID != kurirID {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": "Akses ditolak: Tugas ini bukan milik Anda"})
		return
	}

	// 3. Cari status di database berdasarkan nama status
	var status models.StatusPengantaran
	if err := config.DB.Where("nama_status = ?", statusInput).First(&status).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Status pengantaran tidak valid"})
		return
	}

	tx := config.DB.Begin()

	pengantaran.StatusPengantaranID = status.IdStatusPengantaran

	// 4. Logika tambahan jika status adalah Selesai
	if statusInput == "Selesai" {
		now := time.Now()
		pengantaran.WaktuSampai = &now

		// Upload foto ke MinIO jika ada file "foto" terunggah
		fileUrl, err := utils.UploadFileToMinio(c, "foto", "pengantaran")
		if err == nil && fileUrl != "" {
			pengantaran.FotoBuktiPengiriman = fileUrl
		}

		// Update status pesanan ke "Selesai"
		if err := tx.Model(&models.Pesanan{}).Where("id_pesanan = ?", pengantaran.PesananID).Update("id_status_pesanan", utils.GetStatusPesananID("Selesai")).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengupdate status pesanan"})
			return
		}

		// Update status pembayaran ke settlement jika COD/Cash
		var pembayaran models.Pembayaran
		if err := tx.Where("id_pesanan = ?", pengantaran.PesananID).First(&pembayaran).Error; err == nil {
			var metode models.MetodePembayaran
			if err := tx.First(&metode, pembayaran.MetodePembayaranID).Error; err == nil {
				if metode.KodeMetode == "cod" || metode.KodeMetode == "cash" {
					pembayaran.StatusTransaksiID = utils.GetStatusTransaksiID("settlement")
					pembayaran.TotalDibayar = pengantaran.Pesanan.TotalPembayaran
					pembayaran.WaktuPembayaran = &now
					tx.Save(&pembayaran)
				}
			}
		}
	}

	if err := tx.Save(&pengantaran).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal memperbarui status pengantaran"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Status pengantaran berhasil diperbarui",
		"data": gin.H{
			"foto_bukti":   pengantaran.FotoBuktiPengiriman,
			"status":       statusInput,
			"waktu_sampai": pengantaran.WaktuSampai,
		},
	})
}
