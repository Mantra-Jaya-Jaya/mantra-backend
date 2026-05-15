package customer

import (
	"net/http"
	"strconv"
	"time"

	"backend-mantra/config"
	"backend-mantra/models"

	"github.com/gin-gonic/gin"
)

// GetDaftarPesananCustomer handles fetching customer order history
func GetDaftarPesananCustomer(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	offset := (page - 1) * limit

	// TODO: Ambil ID Customer dari context JWT jika middleware sudah aktif
	customerID := uint(1)

	var pesanan []models.Pesanan
	var total int64

	// Hitung total pesanan
	config.DB.Model(&models.Pesanan{}).Where("id_customer = ?", customerID).Count(&total)

	// Ambil data pesanan
	if err := config.DB.Where("id_customer = ?", customerID).Limit(limit).Offset(offset).Order("tanggal_pesanan desc").Find(&pesanan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil daftar pesanan",
		})
		return
	}

	// Jika tidak ada pesanan, langsung return kosong
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

	// Ambil semua ID pesanan untuk query detail
	var pesananIDs []uint
	for _, p := range pesanan {
		pesananIDs = append(pesananIDs, p.IdPesanan)
	}

	// Ambil detail pesanan beserta barangnya
	var details []models.DetailPesanan
	config.DB.Preload("SpesifikasiBarang.Barang").Where("id_pesanan IN ?", pesananIDs).Find(&details)

	// Kelompokkan detail berdasarkan ID Pesanan
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

	// Susun response data
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

// CheckoutPesanan handles creating an order
func CheckoutPesanan(c *gin.Context) {
	// TODO: Ambil ID Customer dari context JWT jika middleware sudah aktif
	customerID := uint(1)

	var cartItems []models.Keranjang
	// Ambil semua item di keranjang beserta data barang dan diskonnya
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

	// Hitung total pembayaran dengan memperhitungkan diskon
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

		// Cek apakah ada diskon aktif
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

	// Mulai Transaksi Database
	tx := config.DB.Begin()

	pesanan := models.Pesanan{
		CustomerId:      customerID,
		TotalPembayaran: totalBayar,
		TanggalPesanan:  now,
		TipePesanan:     "Online",
		StatusPesanan:   "Diproses", // Default status
	}

	if err := tx.Create(&pesanan).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal membuat pesanan",
		})
		return
	}

	// Insert ke DetailPesanan
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

	// Kosongkan keranjang
	if err := tx.Where("id_customer = ?", customerID).Delete(&models.Keranjang{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengosongkan keranjang",
		})
		return
	}

	// Commit Transaksi
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

// BatalkanPesanan handles canceling an order
func BatalkanPesanan(c *gin.Context) {
	idPesanan := c.Param("id_pesanan")

	// TODO: Ambil ID Customer dari context JWT jika middleware sudah aktif
	customerID := uint(1)

	var pesanan models.Pesanan
	// Cari pesanan berdasarkan ID dan pastikan milik customer yang bersangkutan
	if err := config.DB.Where("id_pesanan = ? AND id_customer = ?", idPesanan, customerID).First(&pesanan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Pesanan tidak ditemukan",
		})
		return
	}

	// Cek apakah pesanan bisa dibatalkan (misal: hanya jika statusnya masih "Diproses")
	if pesanan.StatusPesanan != "Diproses" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Pesanan tidak bisa dibatalkan karena status saat ini: " + pesanan.StatusPesanan,
		})
		return
	}

	// Update status menjadi Dibatalkan
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

// GetDetailPesananCustomer handles fetching single order details
func GetDetailPesananCustomer(c *gin.Context) {
	idPesanan := c.Param("id_pesanan")

	// TODO: Ambil ID Customer dari context JWT jika middleware sudah aktif
	customerID := uint(1)

	var pesanan models.Pesanan
	// Ambil data pesanan dan alamat
	if err := config.DB.Preload("Alamat").Where("id_pesanan = ? AND id_customer = ?", idPesanan, customerID).First(&pesanan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Pesanan tidak ditemukan",
		})
		return
	}

	// Ambil detail pesanan beserta barang dan variannya
	var details []models.DetailPesanan
	config.DB.Preload("SpesifikasiBarang.Barang").Preload("SpesifikasiBarang.DetailSpesifikasi").Where("id_pesanan = ?", pesanan.IdPesanan).Find(&details)

	// Ambil data pengantaran jika ada
	var pengantaran models.Pengantaran
	config.DB.Preload("Kurir.User").Preload("Ekspedisi").Where("id_pesanan = ?", pesanan.IdPesanan).First(&pengantaran)

	// Susun list items
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

	// Susun data tujuan pengantaran
	var tujuanPengantaran interface{} = nil
	if pesanan.Alamat != nil {
		tujuanPengantaran = gin.H{
			"nama_penerima":  pesanan.Alamat.NamaPenerima,
			"alamat_lengkap": pesanan.Alamat.AlamatLengkap,
		}
	}

	// Susun data kurir (beberapa field dimock karena tidak ada di schema)
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

// LacakPesanan handles fetching courier location
func LacakPesanan(c *gin.Context) {
	idPesanan := c.Param("id_pesanan")

	var pengantaran models.Pengantaran
	// Cari data pengantaran berdasarkan ID Pesanan
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
				"plat_nomor": "H 6582 TH", // Mock karena tidak ada di DB
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
