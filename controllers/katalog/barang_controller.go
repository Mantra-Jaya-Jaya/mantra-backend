package katalog

import (
	"net/http"
	"strconv"
	"time"

	"backend-mantra/config"
	"backend-mantra/models"

	"github.com/gin-gonic/gin"
)

// GetDaftarBarang mengambil daftar barang dari katalog dengan pagination.
// Dipakai oleh: customer (GET /customer/katalog/barang), kasir (GET /kasir/katalog/barang), admin (GET /admin/katalog/barang)
// Auth: Wajib login, semua role boleh akses (dikontrol di route)
func GetDaftarBarang(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	offset := (page - 1) * limit

	type ProductResult struct {
		IdBarang       uint
		NamaBarang     string
		GambarBarang   string
		HargaTerendah  int
		HargaTertinggi int
		BesarDiskon    int
		TglMulai       *time.Time
		TglSelesai     *time.Time
	}

	var results []ProductResult
	var total int64

	// Hitung total barang
	config.DB.Model(&models.Barang{}).Count(&total)

	// Ambil data barang dengan join spesifikasi dan diskon
	config.DB.Table("barang").
		Select("barang.id_barang, barang.nama_barang, barang.gambar_barang, MIN(spesifikasi_barang.harga_barang) as harga_terendah, MAX(spesifikasi_barang.harga_barang) as harga_tertinggi, diskon.besar_diskon, diskon.tgl_mulai, diskon.tgl_selesai").
		Joins("LEFT JOIN spesifikasi_barang ON spesifikasi_barang.id_barang = barang.id_barang").
		Joins("LEFT JOIN diskon ON diskon.id_diskon = barang.id_diskon").
		Group("barang.id_barang, barang.nama_barang, barang.gambar_barang, diskon.besar_diskon, diskon.tgl_mulai, diskon.tgl_selesai").
		Limit(limit).Offset(offset).
		Scan(&results)

	var responseData []gin.H
	now := time.Now()
	for _, r := range results {
		punyaDiskon := false
		hargaDiskon := r.HargaTerendah

		if r.BesarDiskon > 0 && r.TglMulai != nil && r.TglSelesai != nil {
			if r.TglMulai.Before(now) && r.TglSelesai.After(now) {
				punyaDiskon = true
				hargaDiskon = r.HargaTerendah - (r.HargaTerendah * r.BesarDiskon / 100)
			}
		}

		responseData = append(responseData, gin.H{
			"id_barang":       r.IdBarang,
			"nama_barang":     r.NamaBarang,
			"harga_terendah":  r.HargaTerendah,
			"harga_tertinggi": r.HargaTertinggi,
			"harga_diskon":    hargaDiskon,
			"punya_diskon":    punyaDiskon,
			"gambar_barang":   r.GambarBarang,
		})
	}

	if responseData == nil {
		responseData = []gin.H{}
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berhasil mengambil daftar barang",
		"data":    responseData,
		"meta": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetDetailBarang mengambil detail satu barang berdasarkan ID.
// Dipakai oleh: admin (GET /admin/katalog/barang/:id_barang)
// Auth: Wajib login, role admin
func GetDetailBarang(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Detail barang berhasil diambil",
		"data": gin.H{
			"id_barang":   1,
			"nama_barang": "Laptop Gaming X",
			"kategori":    "Elektronik",
			"satuan":      "Unit",
			"gambar":      "https://api.mantra.com/storage/barang/laptop-x.jpg",
			"diskon": gin.H{
				"nama_diskon":  "Promo Awal Tahun",
				"besar_diskon": 10,
				"tgl_selesai":  "2026-12-31T23:59:59Z",
			},
			"varian": []gin.H{
				{
					"id_spesifikasi_barang": 5,
					"nama_spesifikasi":      "RAM",
					"nama_detail":           "16GB",
					"harga_barang":          15000000,
					"stok":                  5,
				},
			},
		},
	})
}

// TambahBarang menambahkan barang baru ke katalog.
// Dipakai oleh: admin (POST /admin/katalog/barang)
// Auth: Wajib login, role admin
func TambahBarang(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Barang berhasil ditambahkan",
		"data": gin.H{
			"id_barang":   51,
			"nama_barang": "Laptop Gaming X",
		},
	})
}

// UpdateBarang memperbarui detail barang yang ada.
// Dipakai oleh: admin (PUT /admin/katalog/barang/:id_barang)
// Auth: Wajib login, role admin
func UpdateBarang(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Barang berhasil diperbarui",
	})
}

// HapusBarang menghapus barang dari katalog.
// Dipakai oleh: admin (DELETE /admin/katalog/barang/:id_barang)
// Auth: Wajib login, role admin
func HapusBarang(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Barang berhasil dihapus",
	})
}

// GetDetailBarangByScan mengambil detail barang berdasarkan kode barcode.
// Dipakai oleh: customer (GET /customer/katalog/scan/:kode_barcode), kasir (GET /kasir/katalog/scan/:kode_barcode)
// Auth: Wajib login, role customer atau kasir
func GetDetailBarangByScan(c *gin.Context) {
	kodeBarcode := c.Param("kode_barcode")
	var barcode models.Barcode

	// Cari barcode berdasarkan kode
	if err := config.DB.Preload("SpesifikasiBarang.Barang.Kategori").Preload("SpesifikasiBarang.Barang.Diskon").Preload("SpesifikasiBarang.Barang.Satuan").Where("id_barcode = ?", kodeBarcode).First(&barcode).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data barang tidak ditemukan",
		})
		return
	}

	var varians []models.SpesifikasiBarang
	config.DB.Preload("DetailSpesifikasi.Spesifikasi").Where("id_barang = ?", barcode.SpesifikasiBarang.BarangID).Find(&varians)

	var responseVarian []gin.H
	now := time.Now()
	for _, v := range varians {
		hargaDiskon := v.HargaBarang
		if barcode.SpesifikasiBarang.Barang.DiskonId != 0 && barcode.SpesifikasiBarang.Barang.Diskon.IdDiskon != 0 {
			if barcode.SpesifikasiBarang.Barang.Diskon.TglMulai.Before(now) && barcode.SpesifikasiBarang.Barang.Diskon.TglSelesai.After(now) {
				hargaDiskon = v.HargaBarang - (v.HargaBarang * barcode.SpesifikasiBarang.Barang.Diskon.BesarDiskon / 100)
			}
		}

		responseVarian = append(responseVarian, gin.H{
			"id_spesifikasi_barang": v.IdSpesifikasiBarang,
			"nama_spesifikasi":      v.DetailSpesifikasi.Spesifikasi.NamaSpesifikasi,
			"nama_detail":           v.DetailSpesifikasi.NamaDetailSpesifikasi,
			"harga_barang":          v.HargaBarang,
			"harga_diskon":          hargaDiskon,
			"stok":                  v.Jumlah,
		})
	}

	var diskonData interface{} = nil
	if barcode.SpesifikasiBarang.Barang.DiskonId != 0 && barcode.SpesifikasiBarang.Barang.Diskon.IdDiskon != 0 {
		diskonData = gin.H{
			"nama_diskon":  barcode.SpesifikasiBarang.Barang.Diskon.NamaDiskon,
			"besar_diskon": barcode.SpesifikasiBarang.Barang.Diskon.BesarDiskon,
			"tgl_selesai":  barcode.SpesifikasiBarang.Barang.Diskon.TglSelesai,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data barang ditemukan",
		"data": gin.H{
			"id_barang":     barcode.SpesifikasiBarang.BarangID,
			"nama_barang":   barcode.SpesifikasiBarang.Barang.NamaBarang,
			"kode_barcode":  kodeBarcode,
			"gambar_barang": barcode.SpesifikasiBarang.Barang.GambarBarang,
			"kategori":      barcode.SpesifikasiBarang.Barang.Kategori.NamaKategori,
			"satuan":        barcode.SpesifikasiBarang.Barang.Satuan.NamaSatuan,
			"diskon":        diskonData,
			"varian":        responseVarian,
		},
	})
}

// CariProdukTransaksi mencari produk berdasarkan barcode atau nama untuk keperluan transaksi POS.
// Dipakai oleh: kasir (GET /kasir/katalog/cari)
// Auth: Wajib login, role kasir
func CariProdukTransaksi(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Produk berhasil ditemukan",
		"data": gin.H{
			"id_produk":    1,
			"nama_produk":  "Laptop Gaming X",
			"harga_satuan": 13500000,
			"gambar":       "https://api.mantra.com/storage/barang/laptop-x.jpg",
			"varian": []gin.H{
				{
					"id_spesifikasi_barang": 5,
					"label":                 "16GB RAM",
					"stok":                  5,
				},
			},
		},
	})
}
