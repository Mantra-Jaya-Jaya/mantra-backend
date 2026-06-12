package pemesanan

import (
	"fmt"
	"net/http"
	"time"

	"backend-mantra/config"
	"backend-mantra/models"

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

	// 🚀 2. QUERY DATABASE
	var pesanan models.Pesanan
	err := config.DB.
		Preload("Customer.User").
		Preload("Alamat").
		Preload("DetailPesanan.SpesifikasiBarang.Barang").
		Preload("DetailPesanan.SpesifikasiBarang.DetailSpesifikasi").
		Preload("DetailPesanan.SpesifikasiBarang.DetailSpesifikasi.Spesifikasi").
		Where("tipe_pesanan = ?", "Online").
		Where("status_pesanan = ?", "Dikemas"). // Sesuaikan status
		Order("tanggal_pesanan DESC").
		First(&pesanan).Error

	// 🚀 3. ERROR HANDLING: DATA TIDAK DITEMUKAN (404 Not Found)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data pesanan tidak ditemukan (Belum ada pesanan baru)",
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
		StatusPesanan   string       `json:"status_pesanan"`
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
	dataBungkus := PesananRingkas{
		PublicID:        pesanan.PublicId.String(),
		TotalPembayaran: pesanan.TotalPembayaran,
		TanggalPesanan:  pesanan.TanggalPesanan,
		StatusPesanan:   pesanan.StatusPesanan,
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
	err := config.DB.
		Preload("Customer.User").
		Preload("Alamat").
		Preload("DetailPesanan.SpesifikasiBarang.Barang").
		Preload("DetailPesanan.SpesifikasiBarang.DetailSpesifikasi").
		Preload("DetailPesanan.SpesifikasiBarang.DetailSpesifikasi.Spesifikasi").
		Where("tipe_pesanan = ?", "Online").
		Where("status_pesanan = ?", "Dikemas"). // Filter biar kurir cuma liat yang nganggur/siap antar
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

	// 🚀 3. ERROR HANDLING: DATA KOSONG (404 Not Found)
	if len(daftarPesanan) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Tidak ada pesanan online yang tersedia saat ini",
			"data":    []interface{}{}, // Lempar array kosong biar UI Flutter lu gak jebol (null safety)
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
		StatusPesanan   string       `json:"status_pesanan"`
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
		hasilAkhir = append(hasilAkhir, PesananRingkas{
			PublicID:        pesanan.PublicId.String(),
			TotalPembayaran: pesanan.TotalPembayaran,
			TanggalPesanan:  pesanan.TanggalPesanan,
			StatusPesanan:   pesanan.StatusPesanan,
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
	publicID := c.Param("public_id") // Ambil public_id dari URL

	// 1. Definisikan struct response yang bakal dikirim ke Flutter
	type ItemBarangDTO struct {
		NamaBarang   string `json:"nama_barang"`
		Variasi      string `json:"variasi"`
		JumlahBeli   int    `json:"jumlah_beli"`
		HargaSatuan  int    `json:"harga_satuan"`
		SubtotalItem int    `json:"subtotal_item"`
	}

	type MetodeBayarDTO struct {
		ID         string `json:"id_metode_bayar"`
		NamaMetode string `json:"nama_metode"`
	}

	type DetailPesananDTO struct {
		PublicID        string          `json:"public_id"`
		TotalPembayaran int             `json:"total_pembayaran"`
		MetodeBayar     MetodeBayarDTO  `json:"metode_bayar"`
		DaftarBarang    []ItemBarangDTO `json:"daftar_barang"`
	}

	// 2. Query Database dengan Preload Super Lengkap
	var pesanan models.Pesanan
	err := config.DB.
		Preload("Pembayaran.MetodeBayar"). // 🚀 Pastikan relasi ini ada di model!
		Preload("DetailPesanan.SpesifikasiBarang.Barang").
		Preload("DetailPesanan.SpesifikasiBarang.DetailSpesifikasi.Spesifikasi").
		Where("public_id = ?", publicID).
		First(&pesanan).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Pesanan tidak ditemukan"})
		return
	}

	// 3. Mapping Data Barang
	var listBarang []ItemBarangDTO
	for _, detail := range pesanan.DetailPesanan {
		// Logika gabungin variasi (Warna: Hitam, Size: XL)
		varian := "Default"
		if detail.SpesifikasiBarang.DetailSpesifikasi.IdDetailSpesifikasi != 0 {
			varian = detail.SpesifikasiBarang.DetailSpesifikasi.Spesifikasi.NamaSpesifikasi + ": " +
				detail.SpesifikasiBarang.DetailSpesifikasi.NamaDetailSpesifikasi
		}

		listBarang = append(listBarang, ItemBarangDTO{
			NamaBarang:   detail.SpesifikasiBarang.Barang.NamaBarang,
			Variasi:      varian,
			JumlahBeli:   detail.Jumlah,
			HargaSatuan:  detail.HargaSatuan, // Pastikan field ini ada di tabel DetailPesanan
			SubtotalItem: detail.Subtotal,    // Pastikan field ini ada di tabel DetailPesanan
		})
	}

	// 4. Mapping Data Metode Bayar
	metodeBayar := MetodeBayarDTO{
		ID:         "",
		NamaMetode: "Belum Ada Metode",
	}
	if pesanan.Pembayaran.MetodePembayaran.IdMetodePembayaran != 0 {
		metodeBayar = MetodeBayarDTO{
			ID:         fmt.Sprintf("%d", pesanan.Pembayaran.MetodePembayaran.IdMetodePembayaran),
			NamaMetode: pesanan.Pembayaran.MetodePembayaran.NamaMetode,
		}
	}

	// 5. Bungkus Final
	response := DetailPesananDTO{
		PublicID:        pesanan.PublicId.String(),
		TotalPembayaran: pesanan.TotalPembayaran,
		MetodeBayar:     metodeBayar,
		DaftarBarang:    listBarang,
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
}
