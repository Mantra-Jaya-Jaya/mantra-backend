package transaksi

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"backend-mantra/config"
	"backend-mantra/models"

	"github.com/gin-gonic/gin"
)

// GetDaftarPesanan mengambil daftar pesanan.
// Customer: hanya pesanan milik sendiri.
// Kasir & Admin: semua pesanan dengan filter opsional.
// Dipakai oleh: customer (GET /customer/pesanan), kasir (GET /kasir/pesanan), admin (GET /admin/pesanan)
// Auth: Wajib login, semua role boleh akses (dikontrol di route)
func GetDaftarPesanan(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	offset := (page - 1) * limit

	role := c.GetString("role")
	userID := c.GetInt64("user_id")

	var pesanan []models.Pesanan
	var total int64
	query := config.DB.Model(&models.Pesanan{})

	switch role {
	case "customer":
		// Customer hanya bisa lihat pesanan milik sendiri
		query = query.Where("id_customer = (SELECT id_customer FROM customer WHERE id_user = ?)", userID)
	case "kasir", "admin":
		// Kasir dan admin bisa lihat semua pesanan
	}

	query.Count(&total)

	if err := query.Limit(limit).Offset(offset).Order("tanggal_pesanan desc").Find(&pesanan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil daftar pesanan",
		})
		return
	}

	if len(pesanan) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Daftar pesanan berhasil diambil",
			"data":    []gin.H{},
			"meta": gin.H{
				"page":        page,
				"limit":       limit,
				"total":       total,
				"total_pages": 0,
			},
		})
		return
	}

	var pesananIDs []uint
	for _, p := range pesanan {
		pesananIDs = append(pesananIDs, p.IdPesanan)
	}

	var details []models.DetailPesanan
	config.DB.Preload("SpesifikasiBarang.Barang").Where("id_pesanan IN ?", pesananIDs).Find(&details)

	itemsMap := make(map[uint][]gin.H)
	for _, d := range details {
		itemsMap[d.PesananId] = append(itemsMap[d.PesananId], gin.H{
			"id_barang":       d.SpesifikasiBarang.BarangID,
			"nama_barang":     d.SpesifikasiBarang.Barang.NamaBarang,
			"jumlah":          d.Jumlah,
			"harga_saat_beli": d.HargaSatuan,
			"gambar":          d.SpesifikasiBarang.Barang.GambarBarang,
		})
	}

	var responseData []gin.H
	for _, p := range pesanan {
		items := itemsMap[p.IdPesanan]
		if items == nil {
			items = []gin.H{}
		}
		responseData = append(responseData, gin.H{
			"id_pesanan":    p.PublicId,
			"status":        p.StatusPesanan,
			"tanggal_pesan": p.TanggalPesanan,
			"total_bayar":   p.TotalPembayaran,
			"items":         items,
		})
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Daftar pesanan berhasil diambil",
		"data":    responseData,
		"meta": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetDetailPesanan mengambil detail satu pesanan.
// Ownership: customer hanya boleh akses pesanan miliknya sendiri (selalu 403 bukan 404 jika bukan miliknya).
// Kasir dan admin boleh akses semua pesanan.
// Dipakai oleh: customer (GET /customer/pesanan/:id_pesanan), kasir (GET /kasir/pesanan/:id_pesanan), admin (GET /admin/pesanan/:id_pesanan)
// Auth: Wajib login
func GetDetailPesanan(c *gin.Context) {
	idPesanan := c.Param("id_pesanan")
	role := c.GetString("role")
	userID := c.GetInt64("user_id")

	var pesanan models.Pesanan
	query := config.DB.Preload("Alamat").Where("public_id = ?", idPesanan)

	if role == "customer" {
		// Ownership check: pesanan harus milik customer yang login
		var count int64
		config.DB.Raw(`
			SELECT COUNT(*) FROM pesanan p
			JOIN customer c ON c.id_customer = p.id_customer
			WHERE p.public_id = ? AND c.id_user = ?
		`, idPesanan, userID).Scan(&count)

		if count == 0 {
			// Selalu 403, bukan 404 — jangan bocorkan bahwa ID ada
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "Anda tidak memiliki akses ke resource ini",
				"error": gin.H{
					"code":   "AUTH_002",
					"detail": "Pesanan ini bukan milik Anda",
				},
			})
			return
		}
	}

	if err := query.First(&pesanan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Pesanan tidak ditemukan",
		})
		return
	}

	var details []models.DetailPesanan
	config.DB.Preload("SpesifikasiBarang.Barang").Preload("SpesifikasiBarang.DetailSpesifikasi").Where("id_pesanan = ?", pesanan.IdPesanan).Find(&details)

	var pengantaran models.Pengantaran
	config.DB.Preload("Kurir.User").Preload("Ekspedisi").Where("id_pesanan = ?", pesanan.IdPesanan).First(&pengantaran)

	var items []gin.H
	subtotalItems := 0
	for _, d := range details {
		subtotalItems += d.Subtotal
		items = append(items, gin.H{
			"id_barang":    d.SpesifikasiBarang.BarangID,
			"nama_barang":  d.SpesifikasiBarang.Barang.NamaBarang,
			"varian":       d.SpesifikasiBarang.DetailSpesifikasi.NamaDetailSpesifikasi,
			"jumlah":       d.Jumlah,
			"harga_satuan": d.HargaSatuan,
			"gambar":       d.SpesifikasiBarang.Barang.GambarBarang,
		})
	}

	var tujuanPengantaran interface{} = nil
	if pesanan.Alamat != nil {
		tujuanPengantaran = gin.H{
			"nama_penerima":  pesanan.Alamat.NamaPenerima,
			"alamat_lengkap": pesanan.Alamat.AlamatLengkap,
		}
	}

	var kurirData interface{} = nil
	if pengantaran.IdPengantaran != 0 {
		kurirData = gin.H{
			"nama_kurir": pengantaran.Kurir.User.NamaLengkap,
			"plat_nomor": "H 6582 TH", // Mock karena tidak ada di DB
			"ekspedisi":  pengantaran.Ekspedisi.NamaEkspedisi,
			"foto_kurir": "https://api.mantra.com/storage/kurir/ricardo.jpg", // Mock
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Detail pesanan berhasil diambil",
		"data": gin.H{
			"no_pesanan":         pesanan.PublicId,
			"status":             pesanan.StatusPesanan,
			"tanggal_pesan":      pesanan.TanggalPesanan,
			"items":              items,
			"tujuan_pengantaran": tujuanPengantaran,
			"kurir":              kurirData,
			"rincian_pembayaran": gin.H{
				"subtotal_items": subtotalItems,
				"ongkir":         20000, // Mock karena belum ada kalkulasi ongkir
				"biaya_proteksi": 2000,  // Mock
				"total":          subtotalItems + 20000 + 2000,
			},
		},
	})
}

// CheckoutPesanan membuat pesanan baru dari isi keranjang customer.
// Dipakai oleh: customer (POST /customer/pesanan/checkout)
// Auth: Wajib login, role customer
// Ownership: id_customer diambil dari JWT (user_id), bukan dari body request
func CheckoutPesanan(c *gin.Context) {
	userID := c.GetInt64("user_id")

	// Cari id_customer berdasarkan user_id dari JWT
	var customerIDResult struct{ IdCustomer uint }
	if err := config.DB.Raw("SELECT id_customer FROM customer WHERE id_user = ?", userID).Scan(&customerIDResult).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengidentifikasi customer",
		})
		return
	}
	customerID := customerIDResult.IdCustomer

	var cartItems []models.Keranjang
	if err := config.DB.Preload("SpesifikasiBarang.Barang.Diskon").Where("id_customer = ?", customerID).Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data keranjang",
		})
		return
	}

	if len(cartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Keranjang belanja masih kosong",
		})
		return
	}

	totalBayar := 0
	now := time.Now()

	type ItemDetail struct {
		SpesifikasiBarangID uint
		Jumlah              int
		HargaSatuan         int
		Subtotal            int
	}

	var itemsToInsert []ItemDetail

	for _, item := range cartItems {
		hargaSatuan := item.SpesifikasiBarang.HargaBarang
		b := item.SpesifikasiBarang.Barang

		if b.DiskonId != nil && b.Diskon.IdDiskon != 0 {
			if b.Diskon.TglMulai.Before(now) && b.Diskon.TglSelesai.After(now) {
				hargaSatuan = item.SpesifikasiBarang.HargaBarang - (item.SpesifikasiBarang.HargaBarang * b.Diskon.BesarDiskon / 100)
			}
		}

		subtotal := item.Quantity * hargaSatuan
		totalBayar += subtotal

		itemsToInsert = append(itemsToInsert, ItemDetail{
			SpesifikasiBarangID: item.SpesifikasiBarangID,
			Jumlah:              item.Quantity,
			HargaSatuan:         hargaSatuan,
			Subtotal:            subtotal,
		})
	}

	tx := config.DB.Begin()

	pesanan := models.Pesanan{
		CustomerId:      customerID,
		TotalPembayaran: totalBayar,
		TanggalPesanan:  now,
		TipePesanan:     "Online",
		StatusPesanan:   "Diproses",
	}

	if err := tx.Create(&pesanan).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal membuat pesanan",
		})
		return
	}

	for _, item := range itemsToInsert {
		detail := models.DetailPesanan{
			PesananId:           pesanan.IdPesanan,
			SpesifikasiBarangId: item.SpesifikasiBarangID,
			Jumlah:              item.Jumlah,
			HargaSatuan:         item.HargaSatuan,
			Subtotal:            item.Subtotal,
		}
		if err := tx.Create(&detail).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Gagal menyimpan detail pesanan",
			})
			return
		}
	}

	if err := tx.Where("id_customer = ?", customerID).Delete(&models.Keranjang{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengosongkan keranjang",
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Pesanan berhasil dibuat",
		"data": gin.H{
			"id_pesanan":     pesanan.PublicId,
			"midtrans_token": "token-untuk-sdk-flutter",
			"redirect_url":   "https://app.sandbox.midtrans.com/snap/v2/vtweb/...",
		},
	})
}

// BatalkanPesanan membatalkan pesanan yang masih berstatus "Diproses".
// Dipakai oleh: customer (PATCH /customer/pesanan/:id_pesanan/batal)
// Auth: Wajib login, role customer
// Ownership: pesanan harus milik customer yang login (id_customer dari JWT)
func BatalkanPesanan(c *gin.Context) {
	idPesanan := c.Param("id_pesanan")
	userID := c.GetInt64("user_id")

	// Ownership check
	var count int64
	config.DB.Raw(`
		SELECT COUNT(*) FROM pesanan p
		JOIN customer c ON c.id_customer = p.id_customer
		WHERE p.public_id = ? AND c.id_user = ?
	`, idPesanan, userID).Scan(&count)

	if count == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "Anda tidak memiliki akses ke resource ini",
			"error": gin.H{
				"code":   "AUTH_002",
				"detail": "Pesanan ini bukan milik Anda",
			},
		})
		return
	}

	var pesanan models.Pesanan
	if err := config.DB.Where("public_id = ?", idPesanan).First(&pesanan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Pesanan tidak ditemukan",
		})
		return
	}

	if pesanan.StatusPesanan != "Diproses" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Pesanan tidak bisa dibatalkan karena status saat ini: " + pesanan.StatusPesanan,
		})
		return
	}

	pesanan.StatusPesanan = "Dibatalkan"
	if err := config.DB.Save(&pesanan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal membatalkan pesanan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Pesanan berhasil dibatalkan",
	})
}

// LacakPesanan mengambil data lokasi kurir untuk lacak pesanan.
// Dipakai oleh: customer (GET /customer/pesanan/:id_pesanan/lacak)
// Auth: Wajib login, role customer
func LacakPesanan(c *gin.Context) {
	idPesanan := c.Param("id_pesanan")

	var pengantaran models.Pengantaran
	if err := config.DB.Preload("Kurir.User").
		Joins("JOIN pesanan ON pesanan.id_pesanan = pengantaran.id_pesanan").
		Where("pesanan.public_id = ?", idPesanan).
		First(&pengantaran).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data pelacakan untuk pesanan ini tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data lacak pesanan berhasil diambil",
		"data": gin.H{
			"id_pesanan": idPesanan,
			"kurir": gin.H{
				"nama":       pengantaran.Kurir.User.NamaLengkap,
				"plat_nomor": "H 6582 TH",                                        // Mock karena tidak ada di DB
				"foto":       "https://api.mantra.com/storage/kurir/ricardo.jpg", // Mock
			},
			"lokasi_kurir": gin.H{
				"latitude":  pengantaran.LastLatitude,
				"longitude": pengantaran.LastLongitude,
			},
			"estimasi_tiba": "8 mins", // Mock
			"jarak_meter":   1500,     // Mock
		},
	})
}

// GetDashboardKasir mengambil data ringkasan dashboard kasir.
// Dipakai oleh: kasir (GET /kasir/dashboard)
// Auth: Wajib login, role kasir
func GetDashboardKasir(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Ambil nama kasir dari data user yang login
	var kasir models.Kasir
	if err := config.DB.Preload("User").Where("id_user = ?", userID).First(&kasir).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data kasir tidak ditemukan",
		})
		return
	}

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Hitung total pendapatan hari ini
	var totalPendapatan struct{ Total int }
	config.DB.Model(&models.Pesanan{}).
		Select("COALESCE(SUM(total_pembayaran), 0) as total").
		Where("tanggal_pesanan >= ? AND tanggal_pesanan < ?", startOfDay, endOfDay).
		Scan(&totalPendapatan)

	// Hitung jumlah transaksi hari ini
	var jumlahTransaksi int64
	config.DB.Model(&models.Pesanan{}).
		Where("tanggal_pesanan >= ? AND tanggal_pesanan < ?", startOfDay, endOfDay).
		Count(&jumlahTransaksi)

	// Hitung total item terjual hari ini
	var totalItemTerjual struct{ Total int }
	config.DB.Model(&models.DetailPesanan{}).
		Select("COALESCE(SUM(detail_pesanan.jumlah), 0) as total").
		Joins("JOIN pesanan ON pesanan.id_pesanan = detail_pesanan.id_pesanan").
		Where("pesanan.tanggal_pesanan >= ? AND pesanan.tanggal_pesanan < ?", startOfDay, endOfDay).
		Scan(&totalItemTerjual)

	// Ambil 5 aktivitas terkini hari ini
	type AktivitasResult struct {
		IdPesanan       uint
		TanggalPesanan  time.Time
		TotalPembayaran int
		PaymentType     string
	}

	var aktivitasRaw []AktivitasResult
	config.DB.Table("pesanan").
		Select("pesanan.id_pesanan, pesanan.tanggal_pesanan, pesanan.total_pembayaran, COALESCE(pembayaran.payment_type, 'tunai') as payment_type").
		Joins("LEFT JOIN pembayaran ON pembayaran.id_pesanan = pesanan.id_pesanan").
		Where("pesanan.tanggal_pesanan >= ? AND pesanan.tanggal_pesanan < ?", startOfDay, endOfDay).
		Order("pesanan.tanggal_pesanan DESC").
		Limit(5).
		Scan(&aktivitasRaw)

	var aktivitasTerkini []gin.H
	for _, a := range aktivitasRaw {
		aktivitasTerkini = append(aktivitasTerkini, gin.H{
			"id_transaksi":      a.IdPesanan,
			"nomor_invoice":     "INV-" + a.TanggalPesanan.Format("20060102") + "-" + strconv.Itoa(int(a.IdPesanan)),
			"metode_pembayaran": a.PaymentType,
			"waktu":             a.TanggalPesanan.Format("15:04"),
			"total_bayar":       a.TotalPembayaran,
		})
	}
	if aktivitasTerkini == nil {
		aktivitasTerkini = []gin.H{}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data dashboard berhasil diambil",
		"data": gin.H{
			"user": gin.H{
				"nama_kasir":        kasir.User.NamaLengkap,
				"status_notifikasi": true,
			},
			"statistik_hari_ini": gin.H{
				"total_pendapatan":   totalPendapatan.Total,
				"jumlah_transaksi":   jumlahTransaksi,
				"total_item_terjual": totalItemTerjual.Total,
			},
			"aktivitas_terkini": aktivitasTerkini,
		},
	})
}

// GetLaporanRingkasan mengambil ringkasan laporan penjualan.
// Dipakai oleh: kasir (GET /kasir/laporan)
// Auth: Wajib login, role kasir
func GetLaporanRingkasan(c *gin.Context) {
	var totalPendapatan int64
	config.DB.Model(&models.Pesanan{}).Where("status_pesanan = ?", "Selesai").Select("COALESCE(SUM(total_pembayaran), 0)").Scan(&totalPendapatan)

	var totalTransaksi int64
	config.DB.Model(&models.Pesanan{}).Where("status_pesanan = ?", "Selesai").Count(&totalTransaksi)

	var rataRataPesanan int64
	if totalTransaksi > 0 {
		rataRataPesanan = totalPendapatan / totalTransaksi
	}

	// Grafik hourly sales untuk hari ini
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	type HourlySale struct {
		Hour  int
		Total int
	}
	var hourlySales []HourlySale
	config.DB.Model(&models.Pesanan{}).
		Select("EXTRACT(HOUR FROM tanggal_pesanan) as hour, SUM(total_pembayaran) as total").
		Where("status_pesanan = ? AND tanggal_pesanan >= ? AND tanggal_pesanan < ?", "Selesai", startOfDay, endOfDay).
		Group("hour").
		Order("hour ASC").
		Scan(&hourlySales)

	grafikData := []gin.H{}
	for _, hs := range hourlySales {
		label := fmt.Sprintf("%02d:00", hs.Hour)
		grafikData = append(grafikData, gin.H{
			"label":        label,
			"nilai":        hs.Total,
			"is_highlight": false,
		})
	}
	if len(grafikData) == 0 {
		grafikData = append(grafikData, gin.H{
			"label":        "08:00",
			"nilai":        0,
			"is_highlight": false,
		})
	}

	// 5 produk terlaris
	type TopProduct struct {
		IdBarang      uint
		NamaBarang    string
		Deskripsi     string
		GambarBarang  string
		JumlahTerjual int64
	}

	var topProducts []TopProduct
	config.DB.Table("detail_pesanan").
		Select("barang.id_barang, barang.nama_barang, barang.deskripsi, barang.gambar_barang, SUM(detail_pesanan.jumlah) as jumlah_terjual").
		Joins("JOIN spesifikasi_barang ON spesifikasi_barang.id_spesifikasi_barang = detail_pesanan.id_spesifikasi_barang").
		Joins("JOIN barang ON barang.id_barang = spesifikasi_barang.id_barang").
		Joins("JOIN pesanan ON pesanan.id_pesanan = detail_pesanan.id_pesanan").
		Where("pesanan.status_pesanan = ?", "Selesai").
		Group("barang.id_barang, barang.nama_barang, barang.deskripsi, barang.gambar_barang").
		Order("jumlah_terjual DESC").
		Limit(5).
		Scan(&topProducts)

	produkTerlaris := []gin.H{}
	for _, tp := range topProducts {
		produkTerlaris = append(produkTerlaris, gin.H{
			"id_produk":      tp.IdBarang,
			"nama_produk":    tp.NamaBarang,
			"deskripsi":      tp.Deskripsi,
			"jumlah_terjual": tp.JumlahTerjual,
			"gambar":         tp.GambarBarang,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data laporan berhasil diambil",
		"data": gin.H{
			"header_statistik": gin.H{
				"total_pendapatan":               totalPendapatan,
				"persentase_kenaikan_pendapatan": 0.0, // Default/Mock
				"total_transaksi":                totalTransaksi,
				"persentase_kenaikan_transaksi":  0.0,
				"rata_rata_pesanan":              rataRataPesanan,
				"status_rata_rata":               "stabil",
			},
			"grafik_pendapatan": grafikData,
			"produk_terlaris":   produkTerlaris,
		},
	})
}

// GetDetailLaporanProduk mengambil detail transaksi berdasarkan produk tertentu.
// Dipakai oleh: kasir (GET /kasir/laporan/produk/:id_produk)
// Auth: Wajib login, role kasir
func GetDetailLaporanProduk(c *gin.Context) {
	idProdukStr := c.Param("id_produk")
	var barang models.Barang
	if err := config.DB.Preload("Kategori").First(&barang, "public_id = ?", idProdukStr).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Produk tidak ditemukan",
		})
		return
	}

	var totalTerjual int64
	config.DB.Model(&models.DetailPesanan{}).
		Joins("JOIN spesifikasi_barang ON spesifikasi_barang.id_spesifikasi_barang = detail_pesanan.id_spesifikasi_barang").
		Joins("JOIN pesanan ON pesanan.id_pesanan = detail_pesanan.id_pesanan").
		Where("pesanan.status_pesanan = ? AND spesifikasi_barang.id_barang = ?", "Selesai", barang.IdBarang).
		Select("COALESCE(SUM(detail_pesanan.jumlah), 0)").
		Scan(&totalTerjual)

	type TxHistory struct {
		IdPesanan      uint
		TanggalPesanan time.Time
		HargaSatuan    int
		Jumlah         int
		Subtotal       int
	}

	var txHistory []TxHistory
	config.DB.Table("detail_pesanan").
		Select("pesanan.id_pesanan, pesanan.tanggal_pesanan, detail_pesanan.harga_satuan, detail_pesanan.jumlah, detail_pesanan.subtotal").
		Joins("JOIN pesanan ON pesanan.id_pesanan = detail_pesanan.id_pesanan").
		Joins("JOIN spesifikasi_barang ON spesifikasi_barang.id_spesifikasi_barang = detail_pesanan.id_spesifikasi_barang").
		Where("pesanan.status_pesanan = ? AND spesifikasi_barang.id_barang = ?", "Selesai", barang.IdBarang).
		Order("pesanan.tanggal_pesanan DESC").
		Scan(&txHistory)

	riwayatTransaksi := []gin.H{}
	for _, tx := range txHistory {
		riwayatTransaksi = append(riwayatTransaksi, gin.H{
			"id_transaksi":  tx.IdPesanan,
			"nomor_invoice": "INV-" + tx.TanggalPesanan.Format("20060102") + "-" + strconv.Itoa(int(tx.IdPesanan)),
			"tanggal_waktu": tx.TanggalPesanan,
			"subtotal":      tx.Subtotal,
			"quantity":      tx.Jumlah,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Detail produk berhasil diambil",
		"data": gin.H{
			"produk": gin.H{
				"id_produk":   barang.PublicId,
				"nama_produk": barang.NamaBarang,
				"kategori":    barang.Kategori.NamaKategori,
				"gambar":      barang.GambarBarang,
			},
			"statistik_produk": gin.H{
				"total_terjual":       totalTerjual,
				"terjual_periode_ini": totalTerjual,
				"label_periode":       "semua",
			},
			"riwayat_transaksi": riwayatTransaksi,
		},
	})
}

// GetDetailPesananDariLaporan mengambil detail satu pesanan dari view laporan kasir.
// Dipakai oleh: kasir (GET /kasir/laporan/produk/:id_produk/:id_pesanan)
// Auth: Wajib login, role kasir
func GetDetailPesananDariLaporan(c *gin.Context) {
	idPesananStr := c.Param("id_pesanan")

	var pesanan models.Pesanan
	if err := config.DB.Preload("Customer.User").Preload("Alamat").First(&pesanan, "public_id = ?", idPesananStr).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Pesanan tidak ditemukan",
		})
		return
	}

	var details []models.DetailPesanan
	config.DB.Preload("SpesifikasiBarang.Barang").Where("id_pesanan = ?", pesanan.IdPesanan).Find(&details)

	var daftarItem []gin.H
	for _, d := range details {
		daftarItem = append(daftarItem, gin.H{
			"id_produk":        d.SpesifikasiBarang.BarangID,
			"nama_produk":      d.SpesifikasiBarang.Barang.NamaBarang,
			"qty":              d.Jumlah,
			"total_harga_item": d.Subtotal,
		})
	}

	var pembayaran models.Pembayaran
	config.DB.Where("id_pesanan = ?", pesanan.IdPesanan).First(&pembayaran)
	metodePembayaran := pembayaran.PaymentType
	if metodePembayaran == "" {
		metodePembayaran = "tunai"
	}

	customerNama := "Walk-in Customer"
	customerAlamat := "-"
	if pesanan.CustomerId != 0 && pesanan.Customer.IdCustomer != 0 {
		customerNama = pesanan.Customer.User.NamaLengkap
	}
	if pesanan.Alamat != nil {
		customerAlamat = pesanan.Alamat.AlamatLengkap
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Detail pesanan berhasil diambil",
		"data": gin.H{
			"order_info": gin.H{
				"nomor_order":   "ORD-" + pesanan.TanggalPesanan.Format("20060102") + "-" + strconv.Itoa(int(pesanan.IdPesanan)),
				"tanggal_waktu": pesanan.TanggalPesanan,
				"status_order": gin.H{
					"kode":  pesanan.StatusPesanan,
					"label": pesanan.StatusPesanan,
				},
			},
			"pelanggan": gin.H{
				"nama":   customerNama,
				"alamat": customerAlamat,
			},
			"daftar_item": daftarItem,
			"rincian_pembayaran": gin.H{
				"metode":        metodePembayaran,
				"subtotal":      pesanan.TotalPembayaran,
				"pajak_nominal": 0,
				"total_akhir":   pesanan.TotalPembayaran,
			},
		},
	})
}
