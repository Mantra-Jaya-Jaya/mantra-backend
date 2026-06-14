package transaksi

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	"strings"

	"backend-mantra/config"
	"backend-mantra/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

// StartTransaksi membuat pesanan baru dengan status Draft untuk kasir POS.
// Dipakai oleh: kasir (POST /kasir/transaksi/start)
// Auth: Wajib login, role kasir
func StartTransaksi(c *gin.Context) {
	var userID uint
	if id, exists := c.Get("user_id"); exists {
		if idInt64, ok := id.(int64); ok {
			userID = uint(idInt64)
		} else if idUint, ok := id.(uint); ok {
			userID = idUint
		}
	}

	var kasir models.Kasir
	if err := config.DB.Joins("JOIN karyawan ON karyawan.id_karyawan = kasir.id_karyawan").Where("karyawan.id_user = ?", userID).First(&kasir).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Data kasir tidak ditemukan"})
		return
	}

	pesanan := models.Pesanan{
		KasirId:        &kasir.IdKasir,
		TanggalPesanan: time.Now(),
		TipePesanan:    "Offline",
		StatusPesanan:  "Draft",
	}

	if err := config.DB.Create(&pesanan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal memulai transaksi baru"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Transaksi baru berhasil dimulai",
		"data": gin.H{
			"id_pesanan": pesanan.IdPesanan,
		},
	})
}

// GetRingkasanCheckout mengambil ringkasan belanja sebelum pembayaran di POS kasir.
// Dipakai oleh: kasir (GET /kasir/transaksi/checkout)
// Auth: Wajib login, role kasir
// Support dual mode query param:
//   - ?id_pesanan=<UUID>    → lookup by public_id (preferred)
//   - ?id_pesanan=<integer> → lookup by id_pesanan (backward compat)
func GetRingkasanCheckout(c *gin.Context) {
	idPesananStr := c.Query("id_pesanan")
	if idPesananStr == "" {
		idPesananStr = c.Query("id_transaksi")
	}

	var pesanan models.Pesanan
	var err error
	if idPesananStr != "" {
		// Coba parse sebagai UUID (public_id) dulu
		if parsed, parseErr := uuid.Parse(idPesananStr); parseErr == nil {
			err = config.DB.First(&pesanan, "public_id = ?", parsed).Error
		} else {
			// Fallback: anggap sebagai integer id_pesanan (backward compat POS kasir)
			err = config.DB.First(&pesanan, "id_pesanan = ?", idPesananStr).Error
		}
	} else {
		err = config.DB.Order("id_pesanan DESC").First(&pesanan).Error
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Transaksi/Pesanan tidak ditemukan"})
		return
	}

	var details []models.DetailPesanan
	config.DB.Preload("SpesifikasiBarang.Barang").Preload("SpesifikasiBarang.DetailSpesifikasi").Where("id_pesanan = ?", pesanan.IdPesanan).Find(&details)

	itemCheckout := []gin.H{}
	subtotal := 0
	for _, d := range details {
		subtotal += d.Subtotal
		varian := "Default"
		if d.SpesifikasiBarang.DetailSpesifikasi.NamaDetailSpesifikasi != "" {
			varian = d.SpesifikasiBarang.DetailSpesifikasi.NamaDetailSpesifikasi
		}
		itemCheckout = append(itemCheckout, gin.H{
			"nama_produk":    d.SpesifikasiBarang.Barang.NamaBarang,
			"varian":         varian,
			"qty":            d.Jumlah,
			"total_per_item": d.Subtotal,
		})
	}

	// Tanpa Pajak, Total Akhir = Subtotal
	totalAkhir := subtotal

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data checkout berhasil diambil",
		"data": gin.H{
			"order_info": gin.H{
				"id_order":    pesanan.IdPesanan,
				"public_id":   pesanan.PublicId,
				"nomor_order": "ORD-" + pesanan.TanggalPesanan.Format("20060102") + "-" + strconv.Itoa(int(pesanan.IdPesanan)),
			},
			"item_checkout": itemCheckout,
			"ringkasan_biaya": gin.H{
				"subtotal":      subtotal,
				"pajak_nominal": 0, // Pajak 0
				"total_akhir":   totalAkhir,
			},
			"pilihan_pembayaran": []gin.H{
				{"id_metode": 1, "label": "Cash", "tipe": "cash"},
				{"id_metode": 2, "label": "QRIS", "tipe": "non-cash"},
			},
		},
	})
}

// UpdateQuantityItem memperbarui quantity item dalam transaksi POS yang sedang berjalan.
// Dipakai oleh: kasir (PATCH /kasir/transaksi/item/update)
// Auth: Wajib login, role kasir
func UpdateQuantityItem(c *gin.Context) {
	var input struct {
		IdPesanan           uint `json:"id_pesanan"`
		IdSpesifikasiBarang uint `json:"id_spesifikasi_barang"`
		Jumlah              int  `json:"jumlah"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format inputan salah"})
		return
	}

	tx := config.DB.Begin()

	// 1. Jika IdPesanan adalah 0, buat pesanan baru (Draft)
	if input.IdPesanan == 0 {
		var userID uint
		if id, exists := c.Get("user_id"); exists {
			if idInt64, ok := id.(int64); ok {
				userID = uint(idInt64)
			} else if idUint, ok := id.(uint); ok {
				userID = idUint
			}
		}

		// Cari data kasir
		var kasir models.Kasir
		if err := tx.Joins("JOIN karyawan ON karyawan.id_karyawan = kasir.id_karyawan").Where("karyawan.id_user = ?", userID).First(&kasir).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Data kasir tidak ditemukan. Pastikan Anda login sebagai kasir."})
			return
		}

		// Cari customer pertama sebagai placeholder (Walk-in Customer)
		var customer models.Customer
		if err := tx.First(&customer).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusFailedDependency, gin.H{"status": "error", "message": "Data customer tidak ditemukan. Harap buat minimal satu data customer sebagai placeholder (Pelanggan POS)."})
			return
		}

		pesananBaru := models.Pesanan{
			CustomerId:     customer.IdCustomer,
			KasirId:        &kasir.IdKasir,
			TanggalPesanan: time.Now(),
			TipePesanan:    "Offline",
			StatusPesanan:  "Draft",
		}
		if err := tx.Create(&pesananBaru).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal membuat transaksi baru di database."})
			return
		}
		input.IdPesanan = pesananBaru.IdPesanan
	} else {
		// Proteksi: Jangan edit pesanan yang sudah Selesai
		var pesanan models.Pesanan
		if err := tx.First(&pesanan, input.IdPesanan).Error; err == nil {
			if pesanan.StatusPesanan == "Selesai" || pesanan.StatusPesanan == "Dibatalkan" {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Pesanan ini sudah selesai dan tidak dapat diubah lagi."})
				return
			}
		}
	}

	var detail models.DetailPesanan
	// 2. Cek apakah barang udah ada di dalam pesanan
	err := tx.Where("id_pesanan = ? AND id_spesifikasi_barang = ?", input.IdPesanan, input.IdSpesifikasiBarang).First(&detail).Error

	if err != nil {
		// 3. KALAU BELUM ADA (BARU DI-SCAN), BIKIN BARU!
		if input.Jumlah > 0 {
			var spek models.SpesifikasiBarang
			if errSpek := tx.First(&spek, input.IdSpesifikasiBarang).Error; errSpek != nil {
				tx.Rollback()
				c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Spesifikasi barang tidak ditemukan"})
				return
			}

			detailBaru := models.DetailPesanan{
				PesananId:           input.IdPesanan,
				SpesifikasiBarangId: input.IdSpesifikasiBarang,
				Jumlah:              input.Jumlah,
				HargaSatuan:         spek.HargaBarang,
				Subtotal:            input.Jumlah * spek.HargaBarang,
			}
			if errCreate := tx.Create(&detailBaru).Error; errCreate != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menambahkan barang ke pesanan"})
				return
			}
		}
	} else {
		// 4. KALAU BARANGNYA UDAH ADA DI KERANJANG KASIR
		if input.Jumlah <= 0 {
			if errDel := tx.Delete(&detail).Error; errDel != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menghapus item"})
				return
			}
		} else {
			detail.Jumlah = input.Jumlah
			detail.Subtotal = detail.Jumlah * detail.HargaSatuan
			if errUpdate := tx.Save(&detail).Error; errUpdate != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal memperbarui jumlah item"})
				return
			}
		}
	}

	// 5. Hitung ulang Total Pembayaran
	var totalBayar int64
	tx.Model(&models.DetailPesanan{}).Where("id_pesanan = ?", input.IdPesanan).Select("COALESCE(SUM(subtotal), 0)").Scan(&totalBayar)

	if err := tx.Model(&models.Pesanan{}).Where("id_pesanan = ?", input.IdPesanan).Update("total_pembayaran", totalBayar).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal memperbarui total harga"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Item berhasil diperbarui",
		"data": gin.H{
			"id_pesanan": input.IdPesanan,
		},
	})
}

// BayarTunai memproses pembayaran tunai (cash) di POS kasir.
// Dipakai oleh: kasir (POST /kasir/transaksi/bayar/tunai)
// Auth: Wajib login, role kasir
func BayarTunai(c *gin.Context) {
	var input struct {
		IdPesanan uint `json:"id_pesanan"`
		Bayar     int  `json:"bayar"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format inputan salah"})
		return
	}

	var pesanan models.Pesanan
	if err := config.DB.First(&pesanan, "id_pesanan = ?", input.IdPesanan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Pesanan tidak ditemukan"})
		return
	}

	// Langsung gunakan TotalPembayaran tanpa pajak
	var totalAkhir int

	config.DB.
    	Model(&models.DetailPesanan{}).
    	Where("id_pesanan = ?", input.IdPesanan).
    	Select("COALESCE(SUM(subtotal),0)").
    	Scan(&totalAkhir)

	fmt.Println("ID Pesanan =", input.IdPesanan)
	fmt.Println("Bayar =", input.Bayar)
	fmt.Println("TotalAkhir =", totalAkhir)
	fmt.Println("TotalPembayaran =", pesanan.TotalPembayaran)
	
	if input.Bayar < totalAkhir {
    	c.JSON(http.StatusBadRequest, gin.H{
        	"status": "error",
        	"message": "Uang pembayaran kurang. Total yang harus dibayar: Rp " +
            	strconv.Itoa(totalAkhir),
    	})
    	return
	}

	kembalian := input.Bayar - totalAkhir
	tx := config.DB.Begin()

	pesanan.StatusPesanan = "Selesai"
	if err := tx.Save(&pesanan).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal update status"})
		return
	}

	pembayaran := models.Pembayaran{
		PesananID:       pesanan.IdPesanan,
		PaymentType:     "cash",
		StatusTransaksi: "settlement",
		FraudStatus:     "accept",
		TotalDibayar:    totalAkhir,
	}
	now := time.Now()
	pembayaran.WaktuPembayaran = &now

	if err := tx.Create(&pembayaran).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal simpan pembayaran"})
		return
	}

	tx.Commit()

	invoiceNum := "INV-" + pesanan.TanggalPesanan.Format("20060102") + "-" + strconv.Itoa(int(pesanan.IdPesanan))
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   gin.H{"kembalian": kembalian, "invoice": gin.H{"nomor_invoice": invoiceNum}},
	})
}

// --- FUNGSI PEMBANTU (Wajib ada agar hitungan konsisten) ---

func hitungPajak(subtotal int) int {
	return int(float64(subtotal)*0.11 + 0.5)
}

// BayarNonTunai memproses pembayaran non-tunai via Midtrans (QRIS, transfer, dll).
// Dipakai oleh: kasir (POST /kasir/transaksi/bayar/non-tunai)
// Auth: Wajib login, role kasir
func BayarNonTunai(c *gin.Context) {
	var input struct {
		IdPesanan uint   `json:"id_pesanan"`
		Metode    string `json:"metode"` // "qris", "bca", "bni", atau "bri"
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format inputan salah"})
		return
	}

	// 1. Cari Data Pesanan
	var pesanan models.Pesanan
	if err := config.DB.First(&pesanan, "id_pesanan = ?", input.IdPesanan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Pesanan tidak ditemukan"})
		return
	}

	// 🚀 2. CARI ID METODE PEMBAYARAN DARI DATABASE LU
	metodeInput := strings.ToLower(input.Metode)
	kodeMetodeDb := "qris" // Default ke QRIS (id: 2 di seeder lu)
	if metodeInput == "bca" || metodeInput == "bni" || metodeInput == "bri" {
		kodeMetodeDb = "va" // Set ke VA (id: 3 di seeder lu)
	}

	var metodeDb models.MetodePembayaran
	if err := config.DB.Where("kode_metode = ?", kodeMetodeDb).First(&metodeDb).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Metode pembayaran belum disetup di database"})
		return
	}

	// 3. Setup Midtrans Core API
	var coreClient coreapi.Client
	coreClient.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)
	orderID := "MID-" + strconv.Itoa(int(pesanan.IdPesanan)) + "-" + strconv.FormatInt(time.Now().Unix(), 10)

	req := &coreapi.ChargeReq{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(pesanan.TotalPembayaran), 
		},
	}

	// 4. Atur Request Tipe Pembayaran
	if metodeInput == "bca" {
		req.PaymentType = coreapi.PaymentTypeBankTransfer
		req.BankTransfer = &coreapi.BankTransferDetails{Bank: midtrans.BankBca}
	} else if metodeInput == "bni" {
		req.PaymentType = coreapi.PaymentTypeBankTransfer
		req.BankTransfer = &coreapi.BankTransferDetails{Bank: midtrans.BankBni}
	} else if metodeInput == "bri" {
		req.PaymentType = coreapi.PaymentTypeBankTransfer
		req.BankTransfer = &coreapi.BankTransferDetails{Bank: midtrans.BankBri}
	} else {
		req.PaymentType = coreapi.PaymentTypeGopay // Pakai trik Gopay buat narik QRIS
	}

	// 5. Tembak Midtrans
	coreResp, err := coreClient.ChargeTransaction(req)
	if err != nil {
		fmt.Println("❌ ERROR MIDTRANS:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal terhubung ke Midtrans"})
		return
	}

	// 6. Ekstrak Data Balikan (QR atau VA)
	var qrUrl, vaNumber string
	if metodeInput == "bca" || metodeInput == "bni" || metodeInput == "bri" {
		if len(coreResp.VaNumbers) > 0 {
			vaNumber = coreResp.VaNumbers[0].VANumber
		}
	} else {
		for _, action := range coreResp.Actions {
			if action.Name == "generate-qr-code" {
				qrUrl = action.URL
				break
			}
		}
	}

	// 🚀 7. SIMPAN KE DATABASE SECARA ATOMIC (TRANSACTION)
	tx := config.DB.Begin()

	// Update Status Pesanan
	if err := tx.Model(&pesanan).Update("status_pesanan", "Menunggu Pembayaran").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal update pesanan"})
		return
	}

	// Buat Record Pembayaran Induk
	pembayaran := models.Pembayaran{
		PesananID:          pesanan.IdPesanan,
		OrderIdMidtrans:    orderID,
		PaymentType:        metodeInput,
		StatusTransaksi:    "pending",
		MetodePembayaranID: &metodeDb.IdMetodePembayaran, // 🔥 Relasi ke tabel metode pembayaran!
	}
	if err := tx.Create(&pembayaran).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal membuat pembayaran"})
		return
	}

	// Buat Record Detail Pembayaran (Buat Nota Kasir Nanti)
	detailPembayaran := models.DetailPembayaran{
		PembayaranID:    pembayaran.IdPembayaran,
		KanalPembayaran: metodeInput,
	}
	
	if metodeInput == "qris" || metodeInput == "gopay" {
		detailPembayaran.QrCode = qrUrl // Masukin link gambarnya
	} else {
		detailPembayaran.NomorVA = vaNumber // Masukin nomor rekeningnya
		detailPembayaran.NamaBank = strings.ToUpper(metodeInput)
	}

	if err := tx.Create(&detailPembayaran).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menyimpan detail pembayaran"})
		return
	}

	tx.Commit() // Berhasil semua, patenkan datanya!

	// 8. Kirim Response ke Flutter
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Transaksi berhasil digenerate",
		"data": gin.H{
			"order_id":  orderID,
			"metode":    metodeInput,
			"qr_url":    qrUrl,
			"va_number": vaNumber,
		},
	})
}

func CekStatusPembayaran(c *gin.Context) {
	orderId := c.Param("order_id") // Kita cek berdasarkan order_id dari Midtrans (contoh: MID-257-1781411165)

	var pembayaran models.Pembayaran
	// Cari data pembayaran di database
	if err := config.DB.Where("order_id_midtrans = ?", orderId).First(&pembayaran).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Pembayaran tidak ditemukan"})
		return
	}

	// Kasih tau Flutter kalau statusnya udah lunas (settlement/capture)
	isLunas := pembayaran.StatusTransaksi == "settlement" || pembayaran.StatusTransaksi == "capture"

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"is_lunas": isLunas,
		"data": gin.H{
			"status_transaksi": pembayaran.StatusTransaksi,
		},
	})
}

func hitungPajak(subtotal int) int {
	return int(float64(subtotal)*0.11 + 0.5)
}