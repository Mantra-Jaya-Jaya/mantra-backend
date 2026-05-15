package customer

import (
	"net/http"
	"strconv"
	"time"

	"backend-mantra/config"
	"backend-mantra/models"

	"github.com/gin-gonic/gin"
)

// GetPromo handles fetching active promos
func GetPromo(c *gin.Context) {
	var diskons []models.Diskon
	now := time.Now()

	// Ambil diskon yang aktif (tgl_mulai <= now <= tgl_selesai)
	if err := config.DB.Where("tgl_mulai <= ? AND tgl_selesai >= ?", now, now).Find(&diskons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data promo",
		})
		return
	}

	var responseData []gin.H
	for _, d := range diskons {
		responseData = append(responseData, gin.H{
			"id_diskon":   d.IdDiskon,
			"nama_diskon": d.NamaDiskon,
			"banner_url":  d.BannerDiskon,
			"tgl_selesai": d.TglSelesai,
		})
	}

	// Jika data kosong, pastikan return array kosong bukan null
	if responseData == nil {
		responseData = []gin.H{}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berhasil mengambil data promo",
		"data":    responseData,
	})
}

// GetKategori handles fetching product categories
func GetKategori(c *gin.Context) {
	kategori := []models.Kategori{}

	if err := config.DB.Find(&kategori).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil daftar kategori",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berhasil mengambil daftar kategori",
		"data":    kategori,
	})
}

// GetDaftarBarang handles fetching catalog products
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

	// Ambil data barang dengan join
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

// GetDetailBarangByScan handles fetching product by barcode
func GetDetailBarangByScan(c *gin.Context) {
	kodeBarcode := c.Param("kode_barcode")
	var barcode models.Barcode

	// Cari barcode berdasarkan ID (karena kolom kode_barcode tidak ada di schema)
	if err := config.DB.Preload("Barang.Kategori").Preload("Barang.Diskon").Preload("Barang.Satuan").Where("id_barcode = ?", kodeBarcode).First(&barcode).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data barang tidak ditemukan",
		})
		return
	}

	var varians []models.SpesifikasiBarang
	config.DB.Preload("DetailSpesifikasi.Spesifikasi").Where("id_barang = ?", barcode.BarangId).Find(&varians)

	var responseVarian []gin.H
	now := time.Now()
	for _, v := range varians {
		hargaDiskon := v.HargaBarang
		if barcode.Barang.DiskonId != 0 && barcode.Barang.Diskon.IdDiskon != 0 {
			if barcode.Barang.Diskon.TglMulai.Before(now) && barcode.Barang.Diskon.TglSelesai.After(now) {
				hargaDiskon = v.HargaBarang - (v.HargaBarang * barcode.Barang.Diskon.BesarDiskon / 100)
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
	if barcode.Barang.DiskonId != 0 && barcode.Barang.Diskon.IdDiskon != 0 {
		diskonData = gin.H{
			"nama_diskon":  barcode.Barang.Diskon.NamaDiskon,
			"besar_diskon": barcode.Barang.Diskon.BesarDiskon,
			"tgl_selesai":  barcode.Barang.Diskon.TglSelesai,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data barang ditemukan",
		"data": gin.H{
			"id_barang":     barcode.BarangId,
			"nama_barang":   barcode.Barang.NamaBarang,
			"kode_barcode":  kodeBarcode,
			"gambar_barang": barcode.Barang.GambarBarang,
			"kategori":      barcode.Barang.Kategori.NamaKategori,
			"satuan":        barcode.Barang.Satuan.NamaSatuan,
			"diskon":        diskonData,
			"varian":        responseVarian,
		},
	})
}
