package transaksi

import (
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
			"id_pesanan":    strconv.FormatUint(uint64(p.IdPesanan), 10),
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
	query := config.DB.Preload("Alamat").Where("id_pesanan = ?", idPesanan)

	if role == "customer" {
		// Ownership check: pesanan harus milik customer yang login
		var count int64
		config.DB.Raw(`
			SELECT COUNT(*) FROM pesanan p
			JOIN customer c ON c.id_customer = p.id_customer
			WHERE p.id_pesanan = ? AND c.id_user = ?
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
			"no_pesanan":         strconv.FormatUint(uint64(pesanan.IdPesanan), 10),
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

		if b.DiskonId != 0 && b.Diskon.IdDiskon != 0 {
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
			"id_pesanan":     strconv.FormatUint(uint64(pesanan.IdPesanan), 10),
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
		WHERE p.id_pesanan = ? AND c.id_user = ?
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
	if err := config.DB.Where("id_pesanan = ?", idPesanan).First(&pesanan).Error; err != nil {
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
	if err := config.DB.Preload("Kurir.User").Where("id_pesanan = ?", idPesanan).First(&pengantaran).Error; err != nil {
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
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data dashboard berhasil diambil",
		"data": gin.H{
			"user": gin.H{
				"nama_kasir":        "Budi Santoso",
				"status_notifikasi": true,
			},
			"statistik_hari_ini": gin.H{
				"total_pendapatan":   1500000,
				"jumlah_transaksi":   12,
				"total_item_terjual": 35,
			},
			"aktivitas_terkini": []gin.H{
				{
					"id_transaksi":      101,
					"nomor_invoice":     "INV-20260509-001",
					"metode_pembayaran": "cash",
					"waktu":             "09:30",
					"total_bayar":       125000,
				},
			},
		},
	})
}

// GetLaporanRingkasan mengambil ringkasan laporan penjualan.
// Dipakai oleh: kasir (GET /kasir/laporan)
// Auth: Wajib login, role kasir
func GetLaporanRingkasan(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data laporan berhasil diambil",
		"data": gin.H{
			"header_statistik": gin.H{
				"total_pendapatan":               5000000,
				"persentase_kenaikan_pendapatan": 12.5,
				"total_transaksi":                40,
				"persentase_kenaikan_transaksi":  5.0,
				"rata_rata_pesanan":              125000,
				"status_rata_rata":               "naik",
			},
			"grafik_pendapatan": []gin.H{
				{
					"label":        "08:00",
					"nilai":        250000,
					"is_highlight": false,
				},
			},
			"produk_terlaris": []gin.H{
				{
					"id_produk":      1,
					"nama_produk":    "Laptop Gaming X",
					"deskripsi":      "16GB RAM, 1TB SSD",
					"jumlah_terjual": 15,
					"gambar":         "https://api.mantra.com/storage/barang/laptop-x.jpg",
				},
			},
		},
	})
}

// GetDetailLaporanProduk mengambil detail transaksi berdasarkan produk tertentu.
// Dipakai oleh: kasir (GET /kasir/laporan/produk/:id_produk)
// Auth: Wajib login, role kasir
func GetDetailLaporanProduk(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Detail produk berhasil diambil",
		"data": gin.H{
			"produk": gin.H{
				"id_produk":   1,
				"nama_produk": "Laptop Gaming X",
				"kategori":    "Elektronik",
				"gambar":      "https://api.mantra.com/storage/barang/laptop-x.jpg",
			},
			"statistik_produk": gin.H{
				"total_terjual":       50,
				"terjual_periode_ini": 15,
				"label_periode":       "mingguan",
			},
			"riwayat_transaksi": []gin.H{
				{
					"id_transaksi":  101,
					"nomor_invoice": "INV-20260509-001",
					"tanggal_waktu": "2026-05-09T09:30:00Z",
					"subtotal":      13500000,
					"quantity":      1,
				},
			},
		},
	})
}

// GetDetailPesananDariLaporan mengambil detail satu pesanan dari view laporan kasir.
// Dipakai oleh: kasir (GET /kasir/laporan/pesanan/:id_pesanan)
// Auth: Wajib login, role kasir
func GetDetailPesananDariLaporan(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Detail pesanan berhasil diambil",
		"data": gin.H{
			"order_info": gin.H{
				"nomor_order":   "ORD-20260509-001",
				"tanggal_waktu": "2026-05-09T09:30:00Z",
				"status_order": gin.H{
					"kode":  "selesai",
					"label": "Selesai",
				},
			},
			"pelanggan": gin.H{
				"nama":   "Aarav Lysander",
				"alamat": "Jl. Cempaka Putih No. 12",
			},
			"daftar_item": []gin.H{
				{
					"id_produk":        1,
					"nama_produk":      "Laptop Gaming X",
					"qty":              1,
					"total_harga_item": 13500000,
				},
			},
			"rincian_pembayaran": gin.H{
				"metode":        "qris",
				"subtotal":      13500000,
				"pajak_nominal": 1485000,
				"total_akhir":   14985000,
			},
		},
	})
}
