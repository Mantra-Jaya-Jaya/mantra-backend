package transaksi

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"backend-mantra/config"
	"backend-mantra/models"
	"backend-mantra/services"
	"backend-mantra/utils"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
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
	query := config.DB.Model(&models.Pesanan{}).Preload("StatusPesanan")

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
		itemsMap[d.PesananID] = append(itemsMap[d.PesananID], gin.H{
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

		namaStatus := ""
		if p.StatusPesanan != nil {
			namaStatus = p.StatusPesanan.NamaStatus
		}

		responseData = append(responseData, gin.H{
			"id_pesanan":          p.PublicId,
			"id_status_pesanan":   p.StatusPesananID,
			"nama_status_pesanan": namaStatus, // nama string dari relasi, bukan hardcode
			"tanggal_pesan":       p.TanggalPesanan,
			"total_bayar":         p.TotalPembayaran,
			"items":               items,
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
	idPesanan := c.Param("public_id")
	role := c.GetString("role")
	userID := c.GetInt64("user_id")

	var pesanan models.Pesanan
	query := config.DB.Preload("StatusPesanan").Preload("Alamat").Where("public_id = ?", idPesanan)

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
	config.DB.Preload("Kurir.Karyawan.User").Preload("Ekspedisi").Where("id_pesanan = ?", pesanan.IdPesanan).First(&pengantaran)

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
	if pengantaran.IdPengantaran != 0 && pengantaran.Kurir != nil {
		ekspedisi := ""
		if pengantaran.Ekspedisi != nil {
			ekspedisi = pengantaran.Ekspedisi.NamaEkspedisi
		}
		kurirData = gin.H{
			"nama_kurir": pengantaran.Kurir.Karyawan.User.NamaLengkap,
			"plat_nomor": "",
			"ekspedisi":  ekspedisi,
			"foto_kurir": pengantaran.Kurir.Karyawan.User.FotoProfil,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Detail pesanan berhasil diambil",
		"data": gin.H{
			"no_pesanan":         pesanan.PublicId,
			"id_status_pesanan":  pesanan.StatusPesananID,
			"tanggal_pesan":      pesanan.TanggalPesanan,
			"items":              items,
			"tujuan_pengantaran": tujuanPengantaran,
			"kurir":              kurirData,
			"rincian_pembayaran": gin.H{
				"subtotal_items": subtotalItems,
				"ongkir":         pesanan.OngkosKirim,
				"biaya_proteksi": 0,
				"total":          pesanan.TotalPembayaran,
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

	var input struct {
		IdAlamat           string `json:"id_alamat"`
		IdEkspedisi        *uint  `json:"id_ekspedisi"`
		IdLayananEkspedisi *uint  `json:"id_layanan_ekspedisi"`
		OngkosKirim        int    `json:"ongkos_kirim"`
		Catatan            string `json:"catatan"`
		IdMetodePembayaran *uint  `json:"id_metode_pembayaran"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Input tidak valid"})
		return
	}

	var customerIDResult struct{ IdCustomer uint }
	if err := config.DB.Raw("SELECT id_customer FROM customer WHERE id_user = ?", userID).Scan(&customerIDResult).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengidentifikasi customer"})
		return
	}
	customerID := customerIDResult.IdCustomer

	var cartItems []models.Keranjang
	if err := config.DB.Preload("SpesifikasiBarang.Barang.Diskon").Where("id_customer = ?", customerID).Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengambil data keranjang"})
		return
	}

	if len(cartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Keranjang belanja masih kosong"})
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

		if b.DiskonID != nil && b.Diskon.IdDiskon != 0 {
			if b.Diskon.TglMulai.Before(now) && b.Diskon.TglSelesai.After(now) {
				if b.Diskon.TipeDiskon == "persen" {
					hargaSatuan = item.SpesifikasiBarang.HargaBarang - (item.SpesifikasiBarang.HargaBarang * b.Diskon.BesarDiskon / 100)
				} else if b.Diskon.TipeDiskon == "nominal" {
					hargaSatuan = item.SpesifikasiBarang.HargaBarang - b.Diskon.BesarDiskon
					if hargaSatuan < 0 {
						hargaSatuan = 0
					}
				}
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

	// Cari alamat jika ada
	var alamat models.Alamat
	if input.IdAlamat != "" {
		config.DB.Where("public_id = ?", input.IdAlamat).First(&alamat)
	}

	// Hitung grand total: subtotal + ongkir + pajak (11%)
	ongkir := input.OngkosKirim
	pajak := int(float64(totalBayar+ongkir) * 0.11)
	grandTotal := totalBayar + ongkir + pajak

	tx := config.DB.Begin()

	pesanan := models.Pesanan{
		CustomerID:         customerID,
		TotalPembayaran:    grandTotal,
		TanggalPesanan:     now,
		TipePesananID:      utils.GetTipePesananID("Online"),
		StatusPesananID:    utils.GetStatusPesananID("Diproses"),
		OngkosKirim:        ongkir,
		Catatan:            input.Catatan,
		EkspedisiID:        input.IdEkspedisi,
		LayananEkspedisiID: input.IdLayananEkspedisi,
	}
	if alamat.IdAlamat != 0 {
		pesanan.AlamatID = &alamat.IdAlamat
	}

	if err := tx.Create(&pesanan).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal membuat pesanan"})
		return
	}

	for _, item := range itemsToInsert {
		detail := models.DetailPesanan{
			PesananID:           pesanan.IdPesanan,
			SpesifikasiBarangID: item.SpesifikasiBarangID,
			Jumlah:              item.Jumlah,
			HargaSatuan:         item.HargaSatuan,
			Subtotal:            item.Subtotal,
		}
		if err := tx.Create(&detail).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menyimpan detail pesanan"})
			return
		}
	}

	if err := tx.Where("id_customer = ?", customerID).Delete(&models.Keranjang{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengosongkan keranjang"})
		return
	}

	// Buat record pembayaran
	statusTransaksi := "pending"
	fraudStatus := "accept"
	if input.IdMetodePembayaran != nil {
		var metode models.MetodePembayaran
		config.DB.First(&metode, *input.IdMetodePembayaran)

		if metode.KodeMetode == "cod" || metode.KodeMetode == "cash" {
			statusTransaksi = "settlement"
		}
	}

	pembayaran := models.Pembayaran{
		PesananID:          pesanan.IdPesanan,
		TipePembayaranID:   utils.GetTipePembayaranID("non-cash"),
		StatusTransaksiID:  utils.GetStatusTransaksiID(statusTransaksi),
		FraudStatusID:      utils.GetFraudStatusID(fraudStatus),
		MetodePembayaranID: input.IdMetodePembayaran,
	}

	if statusTransaksi == "settlement" {
		pembayaran.TotalDibayar = grandTotal
		nowTime := time.Now()
		pembayaran.WaktuPembayaran = &nowTime
		pesanan.StatusPesananID = utils.GetStatusPesananID("Selesai")
		tx.Save(&pesanan)
	}

	if err := tx.Create(&pembayaran).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mencatat pembayaran"})
		return
	}

	tx.Commit()

	// Generate Midtrans Snap jika metode dari Midtrans
	var midtransToken string
	var redirectURL string

	if input.IdMetodePembayaran != nil {
		var metode models.MetodePembayaran
		if err := config.DB.First(&metode, *input.IdMetodePembayaran).Error; err == nil {
			if metode.Penyedia == "midtrans" && statusTransaksi != "settlement" {
				orderID := "MID-" + pesanan.TanggalPesanan.Format("20060102") + "-" + strconv.Itoa(int(pesanan.IdPesanan))
				pembayaran.OrderIdMidtrans = orderID
				config.DB.Save(&pembayaran)

				var s snap.Client
				s.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

				req := &snap.Request{
					TransactionDetails: midtrans.TransactionDetails{
						OrderID:  orderID,
						GrossAmt: int64(grandTotal),
					},
				}

				snapResp, err := s.CreateTransaction(req)
				if err == nil {
					midtransToken = snapResp.Token
					redirectURL = snapResp.RedirectURL
				}
			}
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Pesanan berhasil dibuat",
		"data": gin.H{
			"id_pesanan":     pesanan.PublicId,
			"total_bayar":    grandTotal,
			"ongkos_kirim":   ongkir,
			"pajak":          pajak,
			"midtrans_token": midtransToken,
			"redirect_url":   redirectURL,
		},
	})
}

// BatalkanPesanan membatalkan pesanan yang masih berstatus "Diproses".
// Dipakai oleh: customer (PATCH /customer/pesanan/:id_pesanan/batal)
// Auth: Wajib login, role customer
// Ownership: pesanan harus milik customer yang login (id_customer dari JWT)
func BatalkanPesanan(c *gin.Context) {
	idPesanan := c.Param("public_id")
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
	if err := config.DB.Preload("StatusPesanan").Where("public_id = ?", idPesanan).First(&pesanan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Pesanan tidak ditemukan",
		})
		return
	}

	statusName := ""
	if pesanan.StatusPesanan != nil {
		statusName = pesanan.StatusPesanan.NamaStatus
	}
	if statusName != "Diproses" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Pesanan tidak bisa dibatalkan karena status saat ini: " + statusName,
		})
		return
	}

	pesanan.StatusPesananID = utils.GetStatusPesananID("Dibatalkan")
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
	idPesanan := c.Param("public_id")
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
	if err := config.DB.Preload("StatusPesanan").Preload("Ekspedisi").Where("public_id = ?", idPesanan).First(&pesanan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Pesanan tidak ditemukan",
		})
		return
	}

	// 1. Ekspedisi Eksternal (Biteship) - memiliki NomorResi dan Ekspedisi
	if pesanan.NomorResi != nil && *pesanan.NomorResi != "" && pesanan.Ekspedisi != nil && pesanan.Ekspedisi.KodeApi != "" {
		biteship := services.NewBiteshipAdapter()
		history, err := biteship.TrackShipment(*pesanan.NomorResi, pesanan.Ekspedisi.KodeApi)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  "success",
				"message": "Data lacak pesanan berhasil diambil (Biteship offline/pending)",
				"data": gin.H{
					"id_pesanan":     idPesanan,
					"nomor_resi":     *pesanan.NomorResi,
					"ekspedisi":      pesanan.Ekspedisi.NamaEkspedisi,
					"tipe_ekspedisi": "eksternal",
					"history":        []interface{}{},
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Data lacak pesanan berhasil diambil",
			"data": gin.H{
				"id_pesanan":     idPesanan,
				"nomor_resi":     *pesanan.NomorResi,
				"ekspedisi":      pesanan.Ekspedisi.NamaEkspedisi,
				"tipe_ekspedisi": "eksternal",
				"history":        history,
			},
		})
		return
	}

	// 2. Ekspedisi Internal (Kurir Toko)
	var pengantaran models.Pengantaran
	if err := config.DB.Preload("Kurir.Karyawan.User").
		Joins("JOIN pesanan ON pesanan.id_pesanan = pengantaran.id_pesanan").
		Where("pesanan.public_id = ?", idPesanan).
		First(&pengantaran).Error; err != nil {
		// Jika tidak ada record pengantaran, return status default
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Pesanan sedang dipersiapkan di toko",
			"data": gin.H{
				"id_pesanan":        idPesanan,
				"tipe_ekspedisi":    "internal",
				"id_status_pesanan": pesanan.StatusPesananID,
				"kurir":             nil,
			},
		})
		return
	}

	fotoKurir := ""
	var namaKurir string
	if pengantaran.Kurir != nil {
		namaKurir = pengantaran.Kurir.Karyawan.User.NamaLengkap
		if pengantaran.Kurir.Karyawan.User.FotoProfil != "" {
			fotoKurir = pengantaran.Kurir.Karyawan.User.FotoProfil
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data lacak pesanan berhasil diambil",
		"data": gin.H{
			"id_pesanan":     idPesanan,
			"tipe_ekspedisi": "internal",
			"kurir": gin.H{
				"nama":       namaKurir,
				"plat_nomor": "",
				"foto":       fotoKurir,
			},
			"lokasi_kurir": gin.H{
				"latitude":  pengantaran.LastLatitude,
				"longitude": pengantaran.LastLongitude,
			},
			"estimasi_tiba": "",
			"jarak_meter":   0,
		},
	})
}

// GetDashboardKasir mengambil data ringkasan dashboard kasir.
// Dipakai oleh: kasir (GET /kasir/dashboard)
// Auth: Wajib login, role kasir
func GetDashboardKasir(c *gin.Context) {
	// Ambil ID User dari context secara aman (karena di middleware diset sebagai int64)
	var userID uint
	if id, exists := c.Get("user_id"); exists {
		if idInt64, ok := id.(int64); ok {
			userID = uint(idInt64)
		} else if idUint, ok := id.(uint); ok {
			userID = idUint
		}
	}

	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User tidak terautentikasi",
		})
		return
	}

	// Ambil nama kasir dari data user yang login
	var kasir models.Kasir
	if err := config.DB.Joins("JOIN karyawan ON karyawan.id_karyawan = kasir.id_karyawan").Preload("Karyawan").Preload("Karyawan.User").Where("karyawan.id_user = ?", userID).First(&kasir).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data kasir tidak ditemukan",
		})
		return
	}

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Hitung total pendapatan hari ini (Gross)
	var totalPendapatan struct{ Total int }
	config.DB.Model(&models.Pesanan{}).
		Select("COALESCE(SUM(total_pembayaran), 0) as total").
		Where("tanggal_pesanan >= ? AND tanggal_pesanan < ?", startOfDay, endOfDay).
		Scan(&totalPendapatan)

	// Hitung total pendapatan bersih hari ini (Net - Hanya Selesai)
	var totalPendapatanBersih struct{ Total int }
	statusSelesaiID := utils.GetStatusPesananID("Selesai")
	config.DB.Model(&models.Pesanan{}).
		Select("COALESCE(SUM(total_pembayaran), 0) as total").
		Where("tanggal_pesanan >= ? AND tanggal_pesanan < ? AND id_status_pesanan = ?", startOfDay, endOfDay, statusSelesaiID).
		Scan(&totalPendapatanBersih)

	// Hitung jumlah transaksi hari ini (Gross)
	var jumlahTransaksi int64
	config.DB.Model(&models.Pesanan{}).
		Where("tanggal_pesanan >= ? AND tanggal_pesanan < ?", startOfDay, endOfDay).
		Count(&jumlahTransaksi)

	// Hitung jumlah transaksi selesai hari ini (Net)
	var jumlahTransaksiSelesai int64
	config.DB.Model(&models.Pesanan{}).
		Where("tanggal_pesanan >= ? AND tanggal_pesanan < ? AND id_status_pesanan = ?", startOfDay, endOfDay, statusSelesaiID).
		Count(&jumlahTransaksiSelesai)

	// Hitung total item terjual hari ini
	var totalItemTerjual struct{ Total int }
	config.DB.Model(&models.DetailPesanan{}).
		Select("COALESCE(SUM(detail_pesanan.jumlah), 0) as total").
		Joins("JOIN pesanan ON pesanan.id_pesanan = detail_pesanan.id_pesanan").
		Where("pesanan.tanggal_pesanan >= ? AND tanggal_pesanan < ?", startOfDay, endOfDay).
		Scan(&totalItemTerjual)

	// Ambil 5 aktivitas terkini hari ini
	type AktivitasResult struct {
		IdPesanan       uint
		TanggalPesanan  time.Time
		TotalPembayaran int
		RawPaymentType  string
	}

	var aktivitasRaw []AktivitasResult
	config.DB.Table("pesanan").
		Select("pesanan.id_pesanan, pesanan.tanggal_pesanan, pesanan.total_pembayaran, COALESCE(tipe_pembayaran.nama_tipe, 'tunai') as raw_payment_type").
		Joins("LEFT JOIN pembayaran ON pembayaran.id_pesanan = pesanan.id_pesanan").
		Joins("LEFT JOIN tipe_pembayaran ON tipe_pembayaran.id = pembayaran.id_tipe_pembayaran").
		Where("pesanan.tanggal_pesanan >= ? AND tanggal_pesanan < ?", startOfDay, endOfDay).
		Order("pesanan.tanggal_pesanan DESC").
		Limit(5).
		Scan(&aktivitasRaw)

	var aktivitasTerkini []gin.H
	for _, a := range aktivitasRaw {
		aktivitasTerkini = append(aktivitasTerkini, gin.H{
			"id_transaksi":      a.IdPesanan,
			"nomor_invoice":     "INV-" + a.TanggalPesanan.Format("20060102") + "-" + strconv.Itoa(int(a.IdPesanan)),
			"metode_pembayaran": a.RawPaymentType,
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
				"nama_kasir":        kasir.Karyawan.User.NamaLengkap,
				"status_notifikasi": true,
			},
			"statistik_hari_ini": gin.H{
				"total_pendapatan":         totalPendapatan.Total,
				"total_pendapatan_bersih":  totalPendapatanBersih.Total,
				"jumlah_transaksi":         jumlahTransaksi,
				"jumlah_transaksi_selesai": jumlahTransaksiSelesai,
				"total_item_terjual":       totalItemTerjual.Total,
			},
			"aktivitas_terkini": aktivitasTerkini,
		},
	})
}

// GetSemuaAktivitasHariIni mengambil semua aktivitas transaksi pada hari ini.
// Dipakai oleh: kasir (GET /kasir/aktivitas-hari-ini)
// Auth: Wajib login, role kasir
func GetSemuaAktivitasHariIni(c *gin.Context) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	type AktivitasResult struct {
		IdPesanan       uint
		TanggalPesanan  time.Time
		TotalPembayaran int
		RawPaymentType  string
	}

	var aktivitasRaw []AktivitasResult
	config.DB.Table("pesanan").
		Select("pesanan.id_pesanan, pesanan.tanggal_pesanan, pesanan.total_pembayaran, COALESCE(tipe_pembayaran.nama_tipe, 'tunai') as raw_payment_type").
		Joins("LEFT JOIN pembayaran ON pembayaran.id_pesanan = pesanan.id_pesanan").
		Joins("LEFT JOIN tipe_pembayaran ON tipe_pembayaran.id = pembayaran.id_tipe_pembayaran").
		Where("pesanan.tanggal_pesanan >= ? AND tanggal_pesanan < ?", startOfDay, endOfDay).
		Order("pesanan.tanggal_pesanan DESC").
		Scan(&aktivitasRaw)

	var responseData []gin.H
	for _, a := range aktivitasRaw {
		responseData = append(responseData, gin.H{
			"id_transaksi":      a.IdPesanan,
			"nomor_invoice":     "INV-" + a.TanggalPesanan.Format("20060102") + "-" + strconv.Itoa(int(a.IdPesanan)),
			"metode_pembayaran": a.RawPaymentType,
			"waktu":             a.TanggalPesanan.Format("15:04"),
			"total_bayar":       a.TotalPembayaran,
		})
	}
	if responseData == nil {
		responseData = []gin.H{}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Semua aktivitas hari ini berhasil diambil",
		"data":    responseData,
	})
}

// GetLaporanRingkasan mengambil ringkasan laporan penjualan.
// Dipakai oleh: kasir (GET /kasir/laporan)
// Auth: Wajib login, role kasir
func GetLaporanRingkasan(c *gin.Context) {
	var totalPendapatan int64
	statusSelesaiID := utils.GetStatusPesananID("Selesai")
	config.DB.Model(&models.Pesanan{}).Where("id_status_pesanan = ?", statusSelesaiID).Select("COALESCE(SUM(total_pembayaran), 0)").Scan(&totalPendapatan)

	var totalTransaksi int64
	config.DB.Model(&models.Pesanan{}).Where("id_status_pesanan = ?", statusSelesaiID).Count(&totalTransaksi)

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
		Where("id_status_pesanan = ? AND tanggal_pesanan >= ? AND tanggal_pesanan < ?", statusSelesaiID, startOfDay, endOfDay).
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
		PublicId      string
		NamaBarang    string
		Deskripsi     string
		GambarBarang  string
		JumlahTerjual int64
	}

	var topProducts []TopProduct
	config.DB.Table("detail_pesanan").
		Select("barang.id_barang, barang.public_id, barang.nama_barang, barang.deskripsi, barang.gambar_barang, SUM(detail_pesanan.jumlah) as jumlah_terjual").
		Joins("JOIN spesifikasi_barang ON spesifikasi_barang.id_spesifikasi_barang = detail_pesanan.id_spesifikasi_barang").
		Joins("JOIN barang ON barang.id_barang = spesifikasi_barang.id_barang").
		Joins("JOIN pesanan ON pesanan.id_pesanan = detail_pesanan.id_pesanan").
		Where("pesanan.id_status_pesanan = ?", statusSelesaiID).
		Group("barang.id_barang, barang.public_id, barang.nama_barang, barang.deskripsi, barang.gambar_barang").
		Order("jumlah_terjual DESC").
		Limit(5).
		Scan(&topProducts)

	produkTerlaris := []gin.H{}
	for _, tp := range topProducts {
		produkTerlaris = append(produkTerlaris, gin.H{
			"id_produk":      tp.IdBarang,
			"public_id":      tp.PublicId,
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
	idProdukStr := c.Param("public_id")
	var barang models.Barang
	if err := config.DB.Preload("Kategori").First(&barang, "public_id = ?", idProdukStr).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Produk tidak ditemukan",
		})
		return
	}

	var totalTerjual int64
	statusSelesaiID := utils.GetStatusPesananID("Selesai")
	config.DB.Model(&models.DetailPesanan{}).
		Joins("JOIN spesifikasi_barang ON spesifikasi_barang.id_spesifikasi_barang = detail_pesanan.id_spesifikasi_barang").
		Joins("JOIN pesanan ON pesanan.id_pesanan = detail_pesanan.id_pesanan").
		Where("pesanan.id_status_pesanan = ? AND spesifikasi_barang.id_barang = ?", statusSelesaiID, barang.IdBarang).
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
		Where("pesanan.id_status_pesanan = ? AND spesifikasi_barang.id_barang = ?", statusSelesaiID, barang.IdBarang).
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
	idPesananStr := c.Param("public_id")

	var pesanan models.Pesanan
	if err := config.DB.Preload("StatusPesanan").Preload("Customer.User").Preload("Alamat").First(&pesanan, "public_id = ?", idPesananStr).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Pesanan tidak ditemukan",
		})
		return
	}

	// --- LOGIKA SYNC MIDTRANS (Hanya jika status belum Selesai dan Tipe Offline/Online non-tunai) ---
	statusName := ""
	if pesanan.StatusPesanan != nil {
		statusName = pesanan.StatusPesanan.NamaStatus
	}
	if statusName == "Draft" || statusName == "Menunggu Pembayaran" {
		var pembayaran models.Pembayaran
		if err := config.DB.Where("id_pesanan = ? AND order_id_midtrans != ?", pesanan.IdPesanan, "").Order("id_pembayaran DESC").First(&pembayaran).Error; err == nil {

			// Panggil API Midtrans untuk cek status asli
			var s snap.Client
			s.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

			// Catatan: midtrans-go snap client tidak punya GetStatus, harus pakai coreapi atau manual
			// Karena kita hanya punya snap client, kita asumsikan jika lunas di dashboard, user bisa klik simulasi
			// ATAU kita gunakan manual http call ke Midtrans API

			// Namun agar simpel dan tidak menambah dependency, kita biarkan GetDetailPesanan tetap ori
			// dan pastikan polling frontend menggunakan ID yang benar.
		}
	}
	// --------------------------------------------------------------------------------------------

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
	config.DB.Where("id_pesanan = ?", pesanan.IdPesanan).Preload("TipePembayaranRel").First(&pembayaran)
	metodePembayaran := "tunai"
	if pembayaran.TipePembayaranRel != nil {
		metodePembayaran = pembayaran.TipePembayaranRel.NamaTipe
	}

	customerNama := "Walk-in Customer"
	customerAlamat := "-"
	if pesanan.CustomerID != 0 && pesanan.Customer.IdCustomer != 0 {
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
				"nomor_order":       "ORD-" + pesanan.TanggalPesanan.Format("20060102") + "-" + strconv.Itoa(int(pesanan.IdPesanan)),
				"tanggal_waktu":     pesanan.TanggalPesanan,
				"id_status_pesanan": pesanan.StatusPesananID,
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
