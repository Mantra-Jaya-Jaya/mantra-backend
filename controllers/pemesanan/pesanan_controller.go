package pemesanan

import (
	"fmt"
	"net/http"
	"time"

	"backend-mantra/config"
	"backend-mantra/models"
	"backend-mantra/utils"

	"github.com/gin-gonic/gin"
)

func GetPesananTerbaru(c *gin.Context) {
	// 🚀 1. ERROR HANDLING: CEK AUTH (401 Unauthorized)
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

	// Pastikan user ini kurir (Biar makin ketat)
	var count int64
	config.DB.Table("kurir").Joins("JOIN karyawan ON kurir.id_karyawan = karyawan.id_karyawan").Where("karyawan.id_user = ?", userID).Count(&count)
	if count == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Auth salah: Anda bukan kurir"})
		return
	}

	// 🚀 2. QUERY DATABASE (kecualikan yang udah punya pengantaran)
	var pesanan models.Pesanan
	dikemasID := utils.GetStatusPesananID("Dikemas")
	internalID := utils.GetTipeKurirID("internal")
	err := config.DB.
		Preload("StatusPesanan").
		Preload("Customer.User").
		Preload("Alamat").
		Preload("DetailPesanan.SpesifikasiBarang.Barang").
		Preload("DetailPesanan.SpesifikasiBarang.DetailSpesifikasi").
		Preload("DetailPesanan.SpesifikasiBarang.DetailSpesifikasi.Spesifikasi").
		Where("id_tipe_pesanan = ?", utils.GetTipePesananID("Online")).
		Where("NOT EXISTS (SELECT 1 FROM pengantaran WHERE pengantaran.id_pesanan = pesanan.id_pesanan)").
		Where("id_status_pesanan = ?", dikemasID).
		Where("id_tipe_kurir = ?", internalID).
		Order("tanggal_pesanan DESC").
		Limit(1).
		Find(&pesanan).Error

	if err != nil || pesanan.IdPesanan == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Belum ada pesanan baru",
			"data":    nil,
		})
		return
	}

	// 🚀 4. REFACTOR JSON (Bikin Custom Struct Lokal Biar Rapi!)
	type ItemBarang struct {
		NamaBarang string `json:"nama_barang"`
		Varian     string `json:"varian"`
		Jumlah     int    `json:"jumlah"`
	}

	type PesananRingkas struct {
		PublicID        string       `json:"public_id"`
		TotalPembayaran int          `json:"total_pembayaran"`
		TanggalPesanan  time.Time    `json:"tanggal_pesanan"`
		IDStatusPesanan uint         `json:"id_status_pesanan"`
		NamaCustomer    string       `json:"nama_customer"`
		NoTelp          string       `json:"no_telp"`
		AlamatLengkap   string       `json:"alamat_lengkap"`
		CatatanLokasi   string       `json:"catatan_lokasi"`
		DaftarBarang    []ItemBarang `json:"daftar_barang"`
	}

	// Siapin variabel buat alamat (jaga-jaga kalau null)
	namaCust := "Customer Offline"
	noTelp := "-"
	alamatLengkap := "Ambil di Toko"
	catatan := "-"

	if pesanan.Alamat != nil {
		namaCust = pesanan.Alamat.NamaPenerima
		noTelp = pesanan.Alamat.NoTelpPenerima
		alamatLengkap = pesanan.Alamat.AlamatLengkap
		catatan = pesanan.Alamat.CatatanLokasi
	} else if pesanan.Customer.User.NamaLengkap != "" {
		// Fallback ke nama akun kalau alamat null
		namaCust = pesanan.Customer.User.NamaLengkap
	}

	// Looping isi keranjang biar rapi
	var listBarang []ItemBarang
	for _, detail := range pesanan.DetailPesanan {
		varian := "Default"
		if detail.SpesifikasiBarang.DetailSpesifikasi.NamaDetailSpesifikasi != "" {
			namaSpesifikasi := detail.SpesifikasiBarang.DetailSpesifikasi.Spesifikasi.NamaSpesifikasi
			nilaiSpesifikasi := detail.SpesifikasiBarang.DetailSpesifikasi.NamaDetailSpesifikasi

			// Kalau master namanya ada, gabungin jadi "Warna: Hitam"
			if namaSpesifikasi != "" {
				varian = namaSpesifikasi + ": " + nilaiSpesifikasi
			} else {
				// Jaga-jaga kalau master namanya kosong
				varian = nilaiSpesifikasi
			}
		}

		listBarang = append(listBarang, ItemBarang{
			NamaBarang: detail.SpesifikasiBarang.Barang.NamaBarang,
			Varian:     varian,
			Jumlah:     detail.Jumlah,
		})
	}

	// Bungkus ke DTO
	var statusID uint
	if pesanan.StatusPesanan != nil {
		statusID = pesanan.StatusPesanan.IdStatusPesanan
	}
	dataBungkus := PesananRingkas{
		PublicID:        pesanan.PublicId.String(),
		TotalPembayaran: pesanan.TotalPembayaran,
		TanggalPesanan:  pesanan.TanggalPesanan,
		IDStatusPesanan: statusID,
		NamaCustomer:    namaCust,
		NoTelp:          noTelp,
		AlamatLengkap:   alamatLengkap,
		CatatanLokasi:   catatan,
		DaftarBarang:    listBarang,
	}

	// 🚀 5. SUKSES (200 OK)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Pesanan terbaru berhasil diambil",
		"data":    dataBungkus,
	})
}

func GetAllPesananOnline(c *gin.Context) {
	// 🚀 1. ERROR HANDLING: CEK AUTH (401 Unauthorized)
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

	// Pastikan user ini kurir
	var count int64
	config.DB.Table("kurir").Joins("JOIN karyawan ON kurir.id_karyawan = karyawan.id_karyawan").Where("karyawan.id_user = ?", userID).Count(&count)
	if count == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Auth salah: Anda bukan kurir"})
		return
	}

	// 🚀 2. QUERY DATABASE (Tarik Semua Pesanan)
	var daftarPesanan []models.Pesanan
	dikemasID := utils.GetStatusPesananID("Dikemas")
	internalID := utils.GetTipeKurirID("internal")
	err := config.DB.
		Preload("StatusPesanan").
		Preload("Customer.User").
		Preload("Alamat").
		Preload("DetailPesanan.SpesifikasiBarang.Barang").
		Preload("DetailPesanan.SpesifikasiBarang.DetailSpesifikasi").
		Preload("DetailPesanan.SpesifikasiBarang.DetailSpesifikasi.Spesifikasi").
		Where("id_tipe_pesanan = ?", utils.GetTipePesananID("Online")).
		Where("id_status_pesanan = ?", dikemasID).
		Where("id_tipe_kurir = ?", internalID).
		Where("NOT EXISTS (SELECT 1 FROM pengantaran WHERE pengantaran.id_pesanan = pesanan.id_pesanan)").
		Order("tanggal_pesanan DESC").
		Find(&daftarPesanan).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Terjadi kesalahan pada server saat mengambil data",
			"data":    nil,
		})
		return
	}

	if len(daftarPesanan) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Tidak ada pesanan online yang tersedia saat ini",
			"data":    []interface{}{},
		})
		return
	}

	// 🚀 4. REFACTOR JSON (Bikin Struktur DTO)
	type ItemBarang struct {
		NamaBarang string `json:"nama_barang"`
		Varian     string `json:"varian"`
		Jumlah     int    `json:"jumlah"`
	}

	type PesananRingkas struct {
		PublicID        string       `json:"public_id"`
		TotalPembayaran int          `json:"total_pembayaran"`
		TanggalPesanan  time.Time    `json:"tanggal_pesanan"`
		IDStatusPesanan uint         `json:"id_status_pesanan"`
		NamaCustomer    string       `json:"nama_customer"`
		NoTelp          string       `json:"no_telp"`
		AlamatLengkap   string       `json:"alamat_lengkap"`
		CatatanLokasi   string       `json:"catatan_lokasi"`
		DaftarBarang    []ItemBarang `json:"daftar_barang"`
	}

	// Siapin "keranjang" buat nampung semua DTO
	var hasilAkhir []PesananRingkas

	// 🚀 5. LOOPING DAN MAPPING DATA
	for _, pesanan := range daftarPesanan {
		// Default value kalau alamat null
		namaCust := "Customer Offline"
		noTelp := "-"
		alamatLengkap := "Ambil di Toko"
		catatan := "-"

		// Timpa pakai data asli kalau ada alamatnya
		if pesanan.Alamat != nil {
			namaCust = pesanan.Alamat.NamaPenerima
			noTelp = pesanan.Alamat.NoTelpPenerima
			alamatLengkap = pesanan.Alamat.AlamatLengkap
			catatan = pesanan.Alamat.CatatanLokasi
		} else if pesanan.Customer.User.NamaLengkap != "" {
			namaCust = pesanan.Customer.User.NamaLengkap
		}

		// Mapping barangnya
		var listBarang []ItemBarang
		for _, detail := range pesanan.DetailPesanan {
			varian := "Default"
			if detail.SpesifikasiBarang.DetailSpesifikasi.NamaDetailSpesifikasi != "" {
				namaSpesifikasi := detail.SpesifikasiBarang.DetailSpesifikasi.Spesifikasi.NamaSpesifikasi
				nilaiSpesifikasi := detail.SpesifikasiBarang.DetailSpesifikasi.NamaDetailSpesifikasi

				// Kalau master namanya ada, gabungin jadi "Warna: Hitam"
				if namaSpesifikasi != "" {
					varian = namaSpesifikasi + ": " + nilaiSpesifikasi
				} else {
					// Jaga-jaga kalau master namanya kosong
					varian = nilaiSpesifikasi
				}
			}

			listBarang = append(listBarang, ItemBarang{
				NamaBarang: detail.SpesifikasiBarang.Barang.NamaBarang,
				Varian:     varian,
				Jumlah:     detail.Jumlah,
			})
		}

		// Masukin pesanan yang udah langsing ini ke "keranjang" hasil akhir
		var statusID uint
		if pesanan.StatusPesanan != nil {
			statusID = pesanan.StatusPesanan.IdStatusPesanan
		}
		hasilAkhir = append(hasilAkhir, PesananRingkas{
			PublicID:        pesanan.PublicId.String(),
			TotalPembayaran: pesanan.TotalPembayaran,
			TanggalPesanan:  pesanan.TanggalPesanan,
			IDStatusPesanan: statusID,
			NamaCustomer:    namaCust,
			NoTelp:          noTelp,
			AlamatLengkap:   alamatLengkap,
			CatatanLokasi:   catatan,
			DaftarBarang:    listBarang,
		})
	}

	// 🚀 6. SUKSES (200 OK)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Daftar pesanan online berhasil diambil",
		"data":    hasilAkhir,
	})
}

func GetDetailPesanan(c *gin.Context) {
	// 🚀 1. ERROR HANDLING: CEK AUTH (Sama kayak referensi lu)
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

	var count int64
	config.DB.Table("kurir").Joins("JOIN karyawan ON kurir.id_karyawan = karyawan.id_karyawan").Where("karyawan.id_user = ?", userID).Count(&count)
	if count == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Auth salah: Anda bukan kurir"})
		return
	}

	// 🚀 2. TANGKAP PUBLIC ID DARI URL
	publicID := c.Param("public_id")

	// Variabel penampung ID asli (Primary Key)
	var idPesananAsli uint

	// 🕵️ INTEL 1: Cek apakah ini UUID milik tabel Pesanan?
	var cekPesanan models.Pesanan
	if err := config.DB.Select("id_pesanan").Where("public_id = ?", publicID).First(&cekPesanan).Error; err == nil {
		idPesananAsli = cekPesanan.IdPesanan // Dapet! Ini dari halaman Home
	} else {
		// 🕵️ INTEL 2: Kalau bukan, cek apakah ini UUID milik tabel Pengantaran?
		var cekPengantaran models.Pengantaran

		// ⚠️ PERHATIAN: Pastikan nama kolom UUID di tabel pengantaran lu bener (biasanya 'public_id' atau 'id_pengantaran')
		if err := config.DB.Select("id_pesanan").Where("public_id = ?", publicID).First(&cekPengantaran).Error; err == nil {
			idPesananAsli = cekPengantaran.PesananID // Dapet! Ini dari halaman Peta, kita ambil FK pesanannya
		}
	}

	// 🚨 Kalau dua-duanya gagal total (Zonk)
	if idPesananAsli == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data tidak ditemukan (Bukan ID Pesanan maupun ID Pengantaran yang valid)",
			"data":    nil,
		})
		return
	}

	// 🚀 3. TARIK DATA FULL (Karena ID Aslinya Udah Ketemu!)
	var pesanan models.Pesanan
	err := config.DB.
		Preload("Customer.User").
		Preload("Alamat").
		Preload("Pembayaran.MetodePembayaran").
		Preload("DetailPesanan.SpesifikasiBarang.Barang").
		Preload("DetailPesanan.SpesifikasiBarang.DetailSpesifikasi").
		Preload("DetailPesanan.SpesifikasiBarang.DetailSpesifikasi.Spesifikasi").
		Where("id_pesanan = ?", idPesananAsli). // 👈 Tarik pakai ID asli (Primary Key) biar kenceng!
		First(&pesanan).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal memuat detail pesanan"})
		return
	}

	// 🚀 4. REFACTOR JSON (DTO KHUSUS DETAIL PESANAN)
	type ItemBarang struct {
		NamaBarang   string `json:"nama_barang"`
		GambarBarang string `json:"gambar_barang"`
		Varian       string `json:"varian"`
		Jumlah       int    `json:"jumlah_beli"`
		HargaSatuan  int    `json:"harga_satuan"`
		Subtotal     int    `json:"subtotal_item"`
	}

	type MetodeBayarDTO struct {
		IDMetode   string `json:"id_metode_bayar"`
		NamaMetode string `json:"nama_metode"`
	}

	type DetailPesananBungkus struct {
		PublicID        string         `json:"public_id"`
		NamaCustomer    string         `json:"nama_customer"`
		NoTelp          string         `json:"no_telp"`
		AlamatLengkap   string         `json:"alamat_lengkap"`
		TotalPembayaran int            `json:"total_pembayaran"`
		TanggalPesanan  time.Time      `json:"tanggal_pesanan"`
		IDStatusPesanan uint           `json:"id_status_pesanan"`
		MetodeBayar     MetodeBayarDTO `json:"metode_bayar"`
		DaftarBarang    []ItemBarang   `json:"daftar_barang"`
	}

	// 🚀 5. MAPPING BARANG (Looping keranjang)
	var listBarang []ItemBarang
	for _, detail := range pesanan.DetailPesanan {
		varian := "Default"
		if detail.SpesifikasiBarang.DetailSpesifikasi.NamaDetailSpesifikasi != "" {
			namaSpesifikasi := detail.SpesifikasiBarang.DetailSpesifikasi.Spesifikasi.NamaSpesifikasi
			nilaiSpesifikasi := detail.SpesifikasiBarang.DetailSpesifikasi.NamaDetailSpesifikasi

			if namaSpesifikasi != "" {
				varian = namaSpesifikasi + ": " + nilaiSpesifikasi
			} else {
				varian = nilaiSpesifikasi
			}
		}

		// Asumsi: field harga di tabel DetailPesanan lu namanya 'Harga' atau 'HargaSatuan'
		hargaSatuan := detail.HargaSatuan           // Ubah 'Harga' jadi nama field lu yang bener kalau beda
		subtotalItem := hargaSatuan * detail.Jumlah // Kalau lu udah punya field Subtotal, tinggal panggil detail.Subtotal

		listBarang = append(listBarang, ItemBarang{
			NamaBarang:   detail.SpesifikasiBarang.Barang.NamaBarang,
			GambarBarang: detail.SpesifikasiBarang.Barang.GambarBarang,
			Varian:       varian,
			Jumlah:       detail.Jumlah,
			HargaSatuan:  hargaSatuan,
			Subtotal:     subtotalItem,
		})
	}

	// 🚀 6. MAPPING METODE BAYAR (Aman dari Nil Pointer)
	metodeBayar := MetodeBayarDTO{
		IDMetode:   "-",
		NamaMetode: "Belum ada pembayaran",
	}

	// 🚀 PERBAIKAN DI SINI: Cek ganda! Pastikan Pembayaran ADA dan MetodePembayaran ADA!
	if pesanan.Pembayaran != nil && pesanan.Pembayaran.MetodePembayaran != nil {
		metodeBayar = MetodeBayarDTO{
			IDMetode:   fmt.Sprintf("%d", pesanan.Pembayaran.MetodePembayaran.IdMetodePembayaran),
			NamaMetode: pesanan.Pembayaran.MetodePembayaran.NamaMetode,
		}
	}

	namaCust := "Customer Offline"
	noTelp := "-"
	alamatLengkap := "Ambil di Toko"

	if pesanan.Alamat != nil {
		namaCust = pesanan.Alamat.NamaPenerima
		noTelp = pesanan.Alamat.NoTelpPenerima
		alamatLengkap = pesanan.Alamat.AlamatLengkap
	} else if pesanan.Customer.User.NamaLengkap != "" {
		namaCust = pesanan.Customer.User.NamaLengkap
	}

	var statusID uint
	if pesanan.StatusPesanan != nil {
		statusID = pesanan.StatusPesanan.IdStatusPesanan
	}

	// 🚀 7. BUNGKUS KE DTO FINAL
	hasilAkhir := DetailPesananBungkus{
		PublicID:        pesanan.PublicId.String(),
		NamaCustomer:    namaCust,
		NoTelp:          noTelp,
		AlamatLengkap:   alamatLengkap,
		TotalPembayaran: pesanan.TotalPembayaran,
		TanggalPesanan:  pesanan.TanggalPesanan,
		IDStatusPesanan: statusID,
		MetodeBayar:     metodeBayar,
		DaftarBarang:    listBarang,
	}

	// 🚀 8. SUKSES (200 OK)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Detail pesanan berhasil diambil",
		"data":    hasilAkhir,
	})
}

func TerimaPesanan(c *gin.Context) {
	// 🚀 1. CEK AUTH (Identifikasi Kurir)
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

	// Cari ID Kurir aslinya berdasarkan user_id
	var kurir struct {
		IdKurir uint `gorm:"column:id_kurir"`
	}
	errKurir := config.DB.Table("kurir").
		Joins("JOIN karyawan ON kurir.id_karyawan = karyawan.id_karyawan").
		Where("karyawan.id_user = ?", userID).
		Select("kurir.id_kurir").
		First(&kurir).Error

	if errKurir != nil || kurir.IdKurir == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Anda tidak terdaftar sebagai kurir"})
		return
	}

	// 🚀 2. TANGKAP PUBLIC ID PESANAN DARI URL
	publicIDPesanan := c.Param("public_id")

	// Cari ID Pesanan aslinya
	var pesanan models.Pesanan
	if err := config.DB.Where("public_id = ?", publicIDPesanan).First(&pesanan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Pesanan tidak ditemukan"})
		return
	}

	// 🚀 3. CEK TIPE KURIR (Hanya internal yang bisa diklaim kurir toko)
	if pesanan.TipeKurirID != utils.GetTipeKurirID("internal") {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Pesanan ini menggunakan ekspedisi eksternal, tidak bisa diklaim kurir toko!"})
		return
	}

	// 🚀 4. SEMUA DALAM SATU TRANSAKSI DB
	tx := config.DB.Begin()

	var count int64
	tx.Model(&models.Pengantaran{}).Where("id_pesanan = ?", pesanan.IdPesanan).Count(&count)
	if count > 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Waduh, pesanan ini sudah diambil kurir lain!"})
		return
	}

	pengantaranBaru := models.Pengantaran{
		PesananID:           pesanan.IdPesanan,
		KurirID:             &kurir.IdKurir,
		StatusPengantaranID: utils.GetStatusPengantaranID("Menunggu Pickup"),
		EkspedisiID:         pesanan.EkspedisiID,
	}

	if err := tx.Create(&pengantaranBaru).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal membuat jadwal pengantaran"})
		return
	}

	if err := tx.Model(&pesanan).Update("id_status_pesanan", utils.GetStatusPesananID("Dikirim")).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal memperbarui status pesanan"})
		return
	}

	tx.Commit()

	config.DB.Where("id_pengantaran = ?", pengantaranBaru.IdPengantaran).First(&pengantaranBaru)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Mantap! Pesanan berhasil diterima",
		"data":    pengantaranBaru.PublicId.String(),
	})
}
