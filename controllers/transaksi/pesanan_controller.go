package transaksi

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"backend-mantra/config"
	"backend-mantra/models"
	"backend-mantra/services"
	"backend-mantra/utils"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
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
	statusFilter := c.Query("status") // Ambil filter status dari query param

	// 1. Ambil data dari context dengan aman
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User ID tidak ditemukan di session/token",
		})
		return
	}

	// 2. Lakukan type assertion ke int64 (sesuai data dari middleware)
	userIDInt64, ok := userIDInterface.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Terjadi kesalahan sistem: tipe data User ID tidak valid",
		})
		return
	}

	// 3. Konversi ke uint agar variabel `userID` di bawahnya tetap berfungsi tanpa mengubah query
	userID := uint(userIDInt64)

	var pesanan []models.Pesanan
	var total int64
	query := config.DB.Model(&models.Pesanan{}).Preload("StatusPesanan")

	// Terapkan filter status jika ada — konversi snake_case query ke nama_status di DB
	if statusFilter != "" && statusFilter != "Semua" {
		statusNameMap := map[string]string{
			"menunggu_pembayaran": "Menunggu Pembayaran",
			"diproses":            "Diproses",
			"dikemas":             "Dikemas",
			"dikirim":             "Dikirim",
			"selesai":             "Selesai",
			"dibatalkan":          "Dibatalkan",
			"draft":               "Draft",
		}
		if namaStatus, ok := statusNameMap[statusFilter]; ok {
			statusID := utils.GetStatusPesananIDSafe(namaStatus)
			if statusID > 0 {
				query = query.Where("id_status_pesanan = ?", statusID)
			}
		}
	}

	switch role {
	case "customer":
		query = query.Where("id_customer = (SELECT id_customer FROM customer WHERE id_user = ?)", userID)
	case "kasir":
		query = query.Where("id_kasir = (SELECT id_kasir FROM kasir k JOIN karyawan kw ON kw.id_karyawan = k.id_karyawan WHERE kw.id_user = ?)", userID)
	case "admin":
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
	// 1. Ambil data dari context dengan aman
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User ID tidak ditemukan di session/token",
		})
		return
	}

	// 2. Lakukan type assertion ke int64 (sesuai data dari middleware)
	userIDInt64, ok := userIDInterface.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Terjadi kesalahan sistem: tipe data User ID tidak valid",
		})
		return
	}

	// 3. Konversi ke uint agar variabel `userID` di bawahnya tetap berfungsi tanpa mengubah query
	userID := uint(userIDInt64)

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
	if err := config.DB.Preload("SpesifikasiBarang.Barang").Preload("SpesifikasiBarang.DetailSpesifikasi").Where("id_pesanan = ?", pesanan.IdPesanan).Find(&details).Error; err != nil {
		details = []models.DetailPesanan{}
	}

	var pengantaran models.Pengantaran
	if err := config.DB.Preload("Kurir.Karyawan.User").Preload("Ekspedisi").Where("id_pesanan = ?", pesanan.IdPesanan).First(&pengantaran).Error; err != nil {
		pengantaran = models.Pengantaran{}
	}

	var items []gin.H
	subtotalItems := 0
	for _, d := range details {
		subtotalItems += d.Subtotal
		if d.SpesifikasiBarang.IdSpesifikasiBarang == 0 {
			continue
		}
		idBarang := uint(0)
		namaBarang := ""
		gambar := ""
		if d.SpesifikasiBarang.Barang.IdBarang != 0 {
			idBarang = d.SpesifikasiBarang.BarangID
			namaBarang = d.SpesifikasiBarang.Barang.NamaBarang
			gambar = d.SpesifikasiBarang.Barang.GambarBarang
		}
		varian := ""
		if d.SpesifikasiBarang.DetailSpesifikasi.IdDetailSpesifikasi != 0 {
			varian = d.SpesifikasiBarang.DetailSpesifikasi.NamaDetailSpesifikasi
		}
		items = append(items, gin.H{
			"id_barang":    idBarang,
			"nama_barang":  namaBarang,
			"varian":       varian,
			"jumlah":       d.Jumlah,
			"harga_satuan": d.HargaSatuan,
			"gambar":       gambar,
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

	// --- DATA PEMBAYARAN ---
	var pembayaran models.Pembayaran
	var payDetail models.DetailPembayaran
	config.DB.Preload("TipePembayaranRel").Where("id_pesanan = ?", pesanan.IdPesanan).Order("id_pembayaran DESC").First(&pembayaran)
	if pembayaran.IdPembayaran != 0 {
		config.DB.Where("id_pembayaran = ?", pembayaran.IdPembayaran).First(&payDetail)
	}

	// Ambil nama tipe pembayaran dari relasi
	metodePembayaran := ""
	if pembayaran.TipePembayaranRel != nil {
		metodePembayaran = pembayaran.TipePembayaranRel.NamaTipe
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Detail pesanan berhasil diambil",
		"data": gin.H{
			"no_pesanan":          pesanan.PublicId,
			"id_status_pesanan":   pesanan.StatusPesananID,
			"nama_status_pesanan": pesanan.StatusPesanan.NamaStatus,
			"tanggal_pesan":       pesanan.TanggalPesanan,
			"items":               items,
			"tujuan_pengantaran":  tujuanPengantaran,
			"kurir":               kurirData,
			"rincian_pembayaran": gin.H{
				"subtotal_items": subtotalItems,
				"ongkir":         pesanan.OngkosKirim,
				"biaya_proteksi": 0,
				"total":          pesanan.TotalPembayaran,
				"metode":         metodePembayaran,
				"va_number":      payDetail.NomorVA,
				"qr_url":         payDetail.QrCode,
				"bill_key":       payDetail.BillKey,
				"bill_code":      payDetail.BillCode,
				"order_id":       pembayaran.OrderIdMidtrans,
			},
		},
	})
}

// CheckoutPesanan membuat pesanan baru dari data yang dikirim frontend.
// Bisa dari keranjang belanja atau langsung klik "Beli Sekarang".
// Dipakai oleh: customer (POST /customer/pesanan/checkout)
func CheckoutPesanan(c *gin.Context) {

	type CheckoutItem struct {
		IdSpesifikasiBarang uint `json:"id_spesifikasi_barang" binding:"required"`
		Quantity            int  `json:"qty" binding:"required"`
	}

	type CheckoutInput struct {
		IdAlamat           string         `json:"id_alamat" binding:"required"`
		MetodePembayaran   string         `json:"metode_pembayaran" binding:"required"`
		Items              []CheckoutItem `json:"items" binding:"required"`
		Catatan            string         `json:"catatan"`
		IdEkspedisi        *uint          `json:"id_ekspedisi"`
		IdLayananEkspedisi *uint          `json:"id_layanan_ekspedisi"`
		OngkosKirim        int            `json:"ongkos_kirim"`
		IdTipeKurir        *uint          `json:"id_tipe_kurir"`
	}

	var input CheckoutInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Input tidak valid",
			"detail":  err.Error(),
		})
		return
	}

	// 1. Ambil data dari context dengan aman
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User ID tidak ditemukan di session/token",
		})
		return
	}

	// 2. Lakukan type assertion ke int64 (sesuai data dari middleware)
	userIDInt64, ok := userIDInterface.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Terjadi kesalahan sistem: tipe data User ID tidak valid",
		})
		return
	}

	// 3. Konversi ke uint agar variabel `userID` di bawahnya tetap berfungsi tanpa mengubah query
	userID := uint(userIDInt64)

	// Cari id_customer
	var customerIDResult struct{ IdCustomer uint }
	if err := config.DB.Raw("SELECT id_customer FROM customer WHERE id_user = ?", userID).Scan(&customerIDResult).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengidentifikasi customer",
		})
		return
	}

	customerID := customerIDResult.IdCustomer

	// Cari ID alamat asli (uint) berdasarkan PublicId (UUID string) yang dikirim
	var alamat models.Alamat
	if err := config.DB.Where("public_id = ? AND id_customer = ?", input.IdAlamat, customerID).First(&alamat).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Alamat pengiriman tidak ditemukan",
		})
		return
	}

	tx := config.DB.Begin()

	now := time.Now()
	totalBayar := 0
	var detailsToInsert []models.DetailPesanan

	for _, item := range input.Items {

		var spec models.SpesifikasiBarang
		if err := tx.Preload("Barang.Diskon").Where("id_spesifikasi_barang = ?", item.IdSpesifikasiBarang).First(&spec).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("Produk dengan ID varian %d tidak ditemukan", item.IdSpesifikasiBarang),
			})
			return
		}

		// Cek Stok
		if spec.Jumlah < item.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("Stok %s tidak mencukupi", spec.Barang.NamaBarang),
			})
			return
		}

		// Hitung Harga (Gunakan harga dari DB agar aman dari manipulasi frontend)
		hargaSatuan := spec.HargaBarang

		if spec.Barang.DiskonID != nil && spec.Barang.Diskon != nil && spec.Barang.Diskon.IdDiskon != 0 {
			d := spec.Barang.Diskon
			if d.TglMulai.Before(now) && d.TglSelesai.After(now) {
				if d.TipeDiskon == "persen" {
					hargaSatuan = spec.HargaBarang - (spec.HargaBarang * d.BesarDiskon / 100)
				} else if d.TipeDiskon == "nominal" {
					hargaSatuan = spec.HargaBarang - d.BesarDiskon
					if hargaSatuan < 0 {
						hargaSatuan = 0
					}
				}
			}
		}

		subtotal := item.Quantity * hargaSatuan
		totalBayar += subtotal

		detailsToInsert = append(detailsToInsert, models.DetailPesanan{
			SpesifikasiBarangID: item.IdSpesifikasiBarang,
			Jumlah:              item.Quantity,
			HargaSatuan:         hargaSatuan,
			Subtotal:            subtotal,
		})

		// Update Stok
		if err := tx.Model(&spec).Update("jumlah", spec.Jumlah-item.Quantity).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Gagal memperbarui stok",
			})
			return
		}
	}

	// Hitung grand total (subtotal items + ongkir)
	ongkir := input.OngkosKirim
	grandTotal := totalBayar + ongkir

	// Tentukan status awal pesanan berdasarkan metode pembayaran
	metodeInput := strings.ToLower(input.MetodePembayaran)
	var initialStatusID uint
	if metodeInput == "tunai" || metodeInput == "cash" {
		initialStatusID = utils.GetStatusPesananID("Dikemas")
	} else {
		initialStatusID = utils.GetStatusPesananID("Menunggu Pembayaran")
	}

	// Tentukan tipe kurir (default internal)
	tipeKurirID := utils.GetTipeKurirID("internal")
	if input.IdTipeKurir != nil && *input.IdTipeKurir == utils.GetTipeKurirID("external") {
		tipeKurirID = utils.GetTipeKurirID("external")
	}

	// Buat Pesanan
	pesanan := models.Pesanan{
		CustomerID:         customerID,
		TotalPembayaran:    grandTotal,
		TanggalPesanan:     now,
		TipePesananID:      utils.GetTipePesananID("Online"),
		StatusPesananID:    initialStatusID,
		TipeKurirID:        tipeKurirID,
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal membuat pesanan",
		})
		return
	}

	for _, item := range detailsToInsert {
		detail := models.DetailPesanan{
			PesananID:           pesanan.IdPesanan,
			SpesifikasiBarangID: item.SpesifikasiBarangID,
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
		// Logika: Hapus barang dari keranjang jika barang tersebut dibeli
		tx.Where("id_customer = ? AND id_spesifikasi_barang = ?", customerID, detail.SpesifikasiBarangID).Delete(&models.Keranjang{})
	}

	// --- INTEGRASI MIDTRANS CORE API ---
	var qrUrl, vaNumber, billKey, billCode string
	var orderIDMidtrans string

	// Normalisasi metodeInput (Flutter kirim va_bni, va_bri, dll)
	cleanMetode := metodeInput
	if strings.HasPrefix(metodeInput, "va_") {
		cleanMetode = strings.Replace(metodeInput, "va_", "", 1)
	}

	if metodeInput != "tunai" && metodeInput != "cash" {
		// Cari ID Metode Pembayaran dari database
		kodeMetodeDb := "qris" // Default ke QRIS
		if cleanMetode == "bca" || cleanMetode == "bni" || cleanMetode == "bri" || cleanMetode == "mandiri" || cleanMetode == "permata" {
			kodeMetodeDb = "va"
		}

		var metodeDb models.MetodePembayaran
		if err := tx.Where("kode_metode = ?", kodeMetodeDb).First(&metodeDb).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Metode pembayaran belum disetup"})
			return
		}

		// Setup Midtrans Core API
		var coreClient coreapi.Client
		coreClient.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)
		orderIDMidtrans = "MNT-CUST-" + strconv.Itoa(int(pesanan.IdPesanan)) + "-" + strconv.FormatInt(time.Now().Unix(), 10)

		req := &coreapi.ChargeReq{
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  orderIDMidtrans,
				GrossAmt: int64(pesanan.TotalPembayaran),
			},
		}

		// Atur Request Tipe Pembayaran
		if cleanMetode == "bca" {
			req.PaymentType = coreapi.PaymentTypeBankTransfer
			req.BankTransfer = &coreapi.BankTransferDetails{Bank: midtrans.BankBca}
		} else if cleanMetode == "bni" {
			req.PaymentType = coreapi.PaymentTypeBankTransfer
			req.BankTransfer = &coreapi.BankTransferDetails{Bank: midtrans.BankBni}
		} else if cleanMetode == "bri" {
			req.PaymentType = coreapi.PaymentTypeBankTransfer
			req.BankTransfer = &coreapi.BankTransferDetails{Bank: midtrans.BankBri}
		} else if cleanMetode == "mandiri" {
			req.PaymentType = coreapi.PaymentTypeEChannel
			req.EChannel = &coreapi.EChannelDetail{
				BillInfo1: "Pembayaran Mantra",
				BillInfo2: "Order ID: " + orderIDMidtrans,
			}
		} else if cleanMetode == "permata" {
			req.PaymentType = coreapi.PaymentTypeBankTransfer
			req.BankTransfer = &coreapi.BankTransferDetails{Bank: midtrans.BankPermata}
		} else if cleanMetode == "qris" {
			req.PaymentType = coreapi.PaymentTypeQris
		} else {
			req.PaymentType = coreapi.PaymentTypeGopay // Default QRIS via Gopay
		}

		// Tembak Midtrans
		coreResp, err := coreClient.ChargeTransaction(req)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal terhubung ke Midtrans"})
			return
		}

		// Ekstrak Data Balikan
		if cleanMetode == "bca" || cleanMetode == "bni" || cleanMetode == "bri" {
			if len(coreResp.VaNumbers) > 0 {
				vaNumber = coreResp.VaNumbers[0].VANumber
			}
		} else if cleanMetode == "mandiri" {
			billKey = coreResp.BillKey
			billCode = coreResp.BillerCode
		} else if cleanMetode == "permata" {
			vaNumber = coreResp.PermataVaNumber
		} else {
			// Cek di Actions untuk QRIS/Gopay
			for _, action := range coreResp.Actions {
				if action.Name == "generate-qr-code" {
					qrUrl = action.URL
					break
				}
			}
		}

		// Simpan Record Pembayaran
		// Petakan nama metode ke kode tipe pembayaran di DB
		tipePembayaranDB := metodeInput
		if cleanMetode == "bca" || cleanMetode == "bni" || cleanMetode == "bri" || cleanMetode == "mandiri" || cleanMetode == "permata" {
			tipePembayaranDB = "bank_transfer"
		}

		pembayaran := models.Pembayaran{
			PesananID:          pesanan.IdPesanan,
			OrderIdMidtrans:    orderIDMidtrans,
			TipePembayaranID:   utils.GetTipePembayaranID(tipePembayaranDB),
			StatusTransaksiID:  utils.GetStatusTransaksiID("pending"),
			FraudStatusID:      utils.GetFraudStatusID("accept"),
			MetodePembayaranID: &metodeDb.IdMetodePembayaran,
			TotalDibayar:       pesanan.TotalPembayaran,
		}
		if err := tx.Create(&pembayaran).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal membuat record pembayaran"})
			return
		}

		// Simpan Detail Pembayaran
		detailPembayaran := models.DetailPembayaran{
			PembayaranID:    pembayaran.IdPembayaran,
			KanalPembayaran: metodeInput,
		}
		if qrUrl != "" {
			detailPembayaran.QrCode = qrUrl
		}
		if vaNumber != "" {
			detailPembayaran.NomorVA = vaNumber
			detailPembayaran.NamaBank = strings.ToUpper(cleanMetode)
		}
		if billKey != "" && billCode != "" {
			detailPembayaran.BillKey = billKey
			detailPembayaran.BillCode = billCode
		}

		if err := tx.Create(&detailPembayaran).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menyimpan detail pembayaran"})
			return
		}
	}

	tx.Commit()

	// Notifikasi ke semua kasir bahwa ada pesanan online baru
	go func() {
		var kasirUsers []models.User
		config.DB.Joins("JOIN role ON role.id_role = users.id_role").
			Where("role.nama_role = ?", "Kasir").
			Find(&kasirUsers)

		for _, u := range kasirUsers {
			notif := models.Notifikasi{
				UserID:             u.IdUser,
				Judul:              "Transaksi Baru",
				Pesan:              "Ada pesanan online baru yang perlu dikonfirmasi.",
				StatusNotifikasiID: utils.GetStatusNotifikasiID("unread"),
			}
			config.DB.Create(&notif)
		}
	}()

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Pesanan berhasil dibuat",
		"data": gin.H{
			"id_pesanan": pesanan.PublicId,
			"metode":     metodeInput,
			"qr_url":     qrUrl,
			"va_number":  vaNumber,
			"bill_key":   billKey,
			"bill_code":  billCode,
			"order_id":   orderIDMidtrans,
		},
	})
}

// BatalkanPesanan membatalkan pesanan yang masih berstatus "Diproses".
// Dipakai oleh: customer (PATCH /customer/pesanan/:id_pesanan/batal)
// Auth: Wajib login, role customer
// Ownership: pesanan harus milik customer yang login (id_customer dari JWT)
func BatalkanPesanan(c *gin.Context) {
	idPesanan := c.Param("public_id")
	// 1. Ambil data dari context dengan aman
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User ID tidak ditemukan di session/token",
		})
		return
	}

	// 2. Lakukan type assertion ke int64 (sesuai data dari middleware)
	userIDInt64, ok := userIDInterface.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Terjadi kesalahan sistem: tipe data User ID tidak valid",
		})
		return
	}

	// 3. Konversi ke uint agar variabel `userID` di bawahnya tetap berfungsi tanpa mengubah query
	userID := uint(userIDInt64)

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

	if statusName != "Menunggu Pembayaran" {
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
	if err := config.DB.Preload("StatusPesanan").Preload("Ekspedisi").Preload("Alamat").Where("public_id = ?", idPesanan).First(&pesanan).Error; err != nil {
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

	jarakMeter := 0
	var estimasiTiba string
	if pesanan.Alamat != nil && pengantaran.LastLatitude != 0 && pengantaran.LastLongitude != 0 {
		distanceKm := services.Haversine(
			pengantaran.LastLatitude, pengantaran.LastLongitude,
			pesanan.Alamat.Latitude, pesanan.Alamat.Longitude,
		)
		jarakMeter = int(math.Round(distanceKm * 1000))
		avgSpeedKmh := 30.0
		hours := distanceKm / avgSpeedKmh
		minutes := int(math.Round(hours * 60))
		if minutes > 0 {
			estimasiTiba = fmt.Sprintf("%d menit", minutes)
		} else {
			estimasiTiba = "< 1 menit"
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
			"estimasi_tiba": estimasiTiba,
			"jarak_meter":   jarakMeter,
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
		Where("tanggal_pesanan >= ? AND tanggal_pesanan < ? AND id_kasir = ?", startOfDay, endOfDay, kasir.IdKasir).
		Scan(&totalPendapatan)

	// Hitung total pendapatan bersih hari ini (Net - Hanya Selesai)
	var totalPendapatanBersih struct{ Total int }
	statusSelesaiID := utils.GetStatusPesananID("Selesai")
	config.DB.Model(&models.Pesanan{}).
		Select("COALESCE(SUM(total_pembayaran), 0) as total").
		Where("tanggal_pesanan >= ? AND tanggal_pesanan < ? AND id_status_pesanan = ? AND id_kasir = ?", startOfDay, endOfDay, statusSelesaiID, kasir.IdKasir).
		Scan(&totalPendapatanBersih)

	// Hitung jumlah transaksi hari ini (Gross)
	var jumlahTransaksi int64
	config.DB.Model(&models.Pesanan{}).
		Where("tanggal_pesanan >= ? AND tanggal_pesanan < ? AND id_kasir = ?", startOfDay, endOfDay, kasir.IdKasir).
		Count(&jumlahTransaksi)

	// Hitung jumlah transaksi selesai hari ini (Net)
	var jumlahTransaksiSelesai int64
	config.DB.Model(&models.Pesanan{}).
		Where("tanggal_pesanan >= ? AND tanggal_pesanan < ? AND id_status_pesanan = ? AND id_kasir = ?", startOfDay, endOfDay, statusSelesaiID, kasir.IdKasir).
		Count(&jumlahTransaksiSelesai)

	// Hitung total item terjual hari ini
	var totalItemTerjual struct{ Total int }
	config.DB.Model(&models.DetailPesanan{}).
		Select("COALESCE(SUM(detail_pesanan.jumlah), 0) as total").
		Joins("JOIN pesanan ON pesanan.id_pesanan = detail_pesanan.id_pesanan").
		Where("pesanan.tanggal_pesanan >= ? AND tanggal_pesanan < ? AND pesanan.id_kasir = ?", startOfDay, endOfDay, kasir.IdKasir).
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
		Where("pesanan.tanggal_pesanan >= ? AND tanggal_pesanan < ? AND pesanan.id_kasir = ?", startOfDay, endOfDay, kasir.IdKasir).
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

// KirimPesanan digunakan untuk input nomor resi oleh kasir/admin untuk ekspedisi eksternal.
// Dipakai oleh: kasir (PATCH /kasir/pesanan/:public_id/kirim)
// Auth: Wajib login, role kasir
func KirimPesanan(c *gin.Context) {
	publicID := c.Param("public_id")

	var input struct {
		NomorResi string `json:"nomor_resi" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Nomor resi harus diisi"})
		return
	}

	var pesanan models.Pesanan
	if err := config.DB.Where("public_id = ?", publicID).First(&pesanan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Pesanan tidak ditemukan"})
		return
	}

	if pesanan.StatusPesananID != utils.GetStatusPesananID("Dikemas") {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Pesanan harus berstatus Dikemas untuk dikirim"})
		return
	}

	if pesanan.TipeKurirID != utils.GetTipeKurirID("external") {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Pesanan ini bukan ekspedisi eksternal"})
		return
	}

	dikirimID := utils.GetStatusPesananID("Dikirim")
	updates := map[string]interface{}{
		"nomor_resi":        input.NomorResi,
		"id_status_pesanan": dikirimID,
	}
	if err := config.DB.Model(&pesanan).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengirim pesanan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Pesanan berhasil dikirim, resi: " + input.NomorResi,
	})
}

type PengantaranListResponse struct {
	PublicID          string  `json:"public_id"`
	NoPesanan         string  `json:"no_pesanan"`
	CustomerNama      string  `json:"customer_nama"`
	Ekspedisi         string  `json:"ekspedisi"`
	StatusPengantaran string  `json:"status_pengantaran"`
	WaktuPickup       *string `json:"waktu_pickup,omitempty"`
	WaktuSampai       *string `json:"waktu_sampai,omitempty"`
	KurirNama         string  `json:"kurir_nama"`
	AlamatTujuan      string  `json:"alamat_tujuan"`
	LastLatitude      float64 `json:"last_latitude"`
	LastLongitude     float64 `json:"last_longitude"`
	IsExternal        bool    `json:"is_external"`
	NomorResi         string  `json:"nomor_resi,omitempty"`
}

// GetDaftarPengantaranAdmin daftar pengantaran untuk admin monitoring.
// Dipakai oleh: admin (GET /admin/pengantaran)
func GetDaftarPengantaranAdmin(c *gin.Context) {
	var results []PengantaranListResponse

	rows, err := config.DB.Raw(`
		SELECT
			peng.public_id,
			CONCAT('ORD-', DATE_FORMAT(pes.tanggal_pesanan, '%Y%m%d'), '-', pes.id_pesanan) AS no_pesanan,
			COALESCE(cust.nama_lengkap, '-') AS customer_nama,
			COALESCE(eks.nama_ekspedisi, 'Kurir Toko') AS ekspedisi,
			COALESCE(sp.nama_status, 'Menunggu Kurir') AS status_pengantaran,
			DATE_FORMAT(peng.waktu_pickup, '%Y-%m-%d %H:%i:%s') AS waktu_pickup,
			DATE_FORMAT(peng.waktu_sampai, '%Y-%m-%d %H:%i:%s') AS waktu_sampai,
			COALESCE(CONCAT(kary.nama_lengkap), '-') AS kurir_nama,
			COALESCE(alamat.alamat_lengkap, '-') AS alamat_tujuan,
			COALESCE(peng.last_latitude, 0) AS last_latitude,
			COALESCE(peng.last_longitude, 0) AS last_longitude,
			CASE WHEN tk.nama_tipe = 'external' THEN TRUE ELSE FALSE END AS is_external,
			COALESCE(pes.nomor_resi, '') AS nomor_resi
		FROM pengantaran peng
		JOIN pesanan pes ON pes.id_pesanan = peng.id_pesanan
		JOIN customer cust ON cust.id_customer = pes.id_customer
		JOIN tipe_kurir tk ON tk.id_tipe_kurir = pes.id_tipe_kurir
		LEFT JOIN status_pengantaran sp ON sp.id_status_pengantaran = peng.id_status_pengantaran
		LEFT JOIN kurir kr ON kr.id_kurir = peng.id_kurir
		LEFT JOIN karyawan kary ON kary.id_karyawan = kr.id_karyawan
		LEFT JOIN alamat ON alamat.id_alamat = pes.id_alamat
		LEFT JOIN ekspedisi eks ON eks.id_ekspedisi = peng.id_ekspedisi
		ORDER BY peng.id_pengantaran DESC
	`).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengambil data pengantaran"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var r PengantaranListResponse
		var waktuPickup, waktuSampai *string
		err := rows.Scan(
			&r.PublicID,
			&r.NoPesanan,
			&r.CustomerNama,
			&r.Ekspedisi,
			&r.StatusPengantaran,
			&waktuPickup,
			&waktuSampai,
			&r.KurirNama,
			&r.AlamatTujuan,
			&r.LastLatitude,
			&r.LastLongitude,
			&r.IsExternal,
			&r.NomorResi,
		)
		if err != nil {
			continue
		}
		if waktuPickup != nil && *waktuPickup != "" {
			r.WaktuPickup = waktuPickup
		}
		if waktuSampai != nil && *waktuSampai != "" {
			r.WaktuSampai = waktuSampai
		}
		results = append(results, r)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   results,
		"meta": gin.H{
			"total": len(results),
		},
	})
}
