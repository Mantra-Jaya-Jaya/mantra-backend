package transaksi

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"backend-mantra/config"
	"backend-mantra/models"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
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
func GetRingkasanCheckout(c *gin.Context) {
	idPesananStr := c.Query("id_pesanan")
	if idPesananStr == "" {
		idPesananStr = c.Query("id_transaksi")
	}

	var pesanan models.Pesanan
	var err error
	if idPesananStr != "" {
		err = config.DB.First(&pesanan, "id_pesanan = ?", idPesananStr).Error
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
		IdPesanan uint `json:"id_pesanan"`
		Simulasi  bool `json:"simulasi"`
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

	// Jika ini adalah simulasi (untuk demo PBL/localhost tanpa webhook)
	if input.Simulasi {
		tx := config.DB.Begin()
		pesanan.StatusPesanan = "Selesai"
		if err := tx.Save(&pesanan).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal update status simulasi"})
			return
		}

		pembayaran := models.Pembayaran{
			PesananID:       pesanan.IdPesanan,
			PaymentType:     "non-cash (simulation)",
			StatusTransaksi: "settlement",
			FraudStatus:     "accept",
			TotalDibayar:    pesanan.TotalPembayaran,
		}
		now := time.Now()
		pembayaran.WaktuPembayaran = &now
		
		if err := tx.Create(&pembayaran).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal simpan pembayaran simulasi"})
			return
		}
		tx.Commit()

		invoiceNum := "INV-" + pesanan.TanggalPesanan.Format("20060102") + "-" + strconv.Itoa(int(pesanan.IdPesanan))
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Pembayaran simulasi berhasil",
			"data": gin.H{
				"nomor_invoice": invoiceNum,
			},
		})
		return
	}

	// Jika bukan simulasi, panggil Midtrans
	totalAkhir := pesanan.TotalPembayaran

	var s snap.Client
	s.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	orderID := "MID-" + strconv.Itoa(int(pesanan.IdPesanan)) + "-" + strconv.FormatInt(time.Now().Unix(), 10)

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(totalAkhir),
		},
	}

	snapResp, err := s.CreateTransaction(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal terhubung ke Midtrans"})
		return
	}

	// Update status pesanan jadi "Menunggu Pembayaran"
	config.DB.Model(&pesanan).Update("status_pesanan", "Menunggu Pembayaran")

	// Simpan data pembayaran awal (status: pending) agar webhook bisa mencocokkan order_id_midtrans
	pembayaran := models.Pembayaran{
		PesananID:       pesanan.IdPesanan,
		OrderIdMidtrans: orderID,
		PaymentType:     "non-cash",
		StatusTransaksi: "pending",
	}
	config.DB.Create(&pembayaran)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   gin.H{"midtrans_data": gin.H{"token": snapResp.Token}},
	})
}
