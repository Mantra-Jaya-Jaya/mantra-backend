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
		NamaKategori   string
		TotalStok      int // 🚀 1. TAMBAH KEY BUAT NAMPUNG TOTAL STOK GABUNGAN
	}

	var results []ProductResult
	var total int64

	// Hitung total barang
	config.DB.Model(&models.Barang{}).Count(&total)

	// Ambil data barang dengan join spesifikasi, diskon, dan kategori
	config.DB.Table("barang").
		// 🚀 2. SUNTIK SUM(spesifikasi_barang.jumlah) DI SELECT NYA
		Select("barang.id_barang, barang.nama_barang, barang.gambar_barang, MIN(spesifikasi_barang.harga_barang) as harga_terendah, MAX(spesifikasi_barang.harga_barang) as harga_tertinggi, SUM(spesifikasi_barang.jumlah) as total_stok, diskon.besar_diskon, diskon.tgl_mulai, diskon.tgl_selesai, kategori.nama_kategori").
		Joins("LEFT JOIN spesifikasi_barang ON spesifikasi_barang.id_barang = barang.id_barang").
		Joins("LEFT JOIN diskon ON diskon.id_diskon = barang.id_diskon").
		Joins("LEFT JOIN kategori ON kategori.id_kategori = barang.id_kategori").
		Group("barang.id_barang, barang.nama_barang, barang.gambar_barang, diskon.besar_diskon, diskon.tgl_mulai, diskon.tgl_selesai, kategori.nama_kategori").
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
			"kategori":        r.NamaKategori,
			"stok":            r.TotalStok, // 🚀 3. OPER DATA STOK ASLI KE JSON FRONTEND
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
	idBarangStr := c.Param("id_barang")
	idBarang, err := strconv.Atoi(idBarangStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ID barang tidak valid",
		})
		return
	}

	var barang models.Barang
	if err := config.DB.
		Preload("Kategori").
		Preload("Satuan").
		Preload("Diskon").
		First(&barang, "id_barang = ?", idBarang).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Barang tidak ditemukan",
		})
		return
	}

	// Ambil semua varian (spesifikasi_barang) barang ini
	var varians []models.SpesifikasiBarang
	config.DB.
		Preload("DetailSpesifikasi.Spesifikasi").
		Where("id_barang = ?", idBarang).
		Find(&varians)

	now := time.Now()
	var responseVarian []gin.H
	for _, v := range varians {
		hargaDiskon := v.HargaBarang
		if barang.DiskonId != 0 && barang.Diskon.IdDiskon != 0 {
			if barang.Diskon.TglMulai.Before(now) && barang.Diskon.TglSelesai.After(now) {
				hargaDiskon = v.HargaBarang - (v.HargaBarang * barang.Diskon.BesarDiskon / 100)
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
	if responseVarian == nil {
		responseVarian = []gin.H{}
	}

	var diskonData interface{} = nil
	if barang.DiskonId != 0 && barang.Diskon.IdDiskon != 0 {
		diskonData = gin.H{
			"nama_diskon":  barang.Diskon.NamaDiskon,
			"besar_diskon": barang.Diskon.BesarDiskon,
			"tgl_selesai":  barang.Diskon.TglSelesai,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Detail barang berhasil diambil",
		"data": gin.H{
			"id_barang":     barang.IdBarang,
			"nama_barang":   barang.NamaBarang,
			"deskripsi":     barang.Deskripsi,
			"gambar_barang": barang.GambarBarang,
			"kategori":      barang.Kategori.NamaKategori,
			"satuan":        barang.Satuan.NamaSatuan,
			"diskon":        diskonData,
			"varian":        responseVarian,
		},
	})
}

// TambahBarang menambahkan barang baru ke katalog.
// Dipakai oleh: admin (POST /admin/katalog/barang)
// Auth: Wajib login, role admin
func TambahBarang(c *gin.Context) {
	var input struct {
		NamaBarang   string `json:"nama_barang" binding:"required"`
		GambarBarang string `json:"gambar_barang"`
		Deskripsi    string `json:"deskripsi"`
		KategoriId   uint   `json:"id_kategori" binding:"required"`
		SatuanId     uint   `json:"id_satuan" binding:"required"`
		DiskonId     uint   `json:"id_diskon"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Format inputan salah: " + err.Error(),
		})
		return
	}

	barang := models.Barang{
		NamaBarang:   input.NamaBarang,
		GambarBarang: input.GambarBarang,
		Deskripsi:    input.Deskripsi,
		KategoriId:   input.KategoriId,
		SatuanId:     input.SatuanId,
		DiskonId:     input.DiskonId,
	}

	if err := config.DB.Create(&barang).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menyimpan data barang",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Barang berhasil ditambahkan",
		"data": gin.H{
			"id_barang":   barang.IdBarang,
			"nama_barang": barang.NamaBarang,
		},
	})
}

// UpdateBarang memperbarui detail barang yang ada.
// Dipakai oleh: admin (PUT /admin/katalog/barang/:id_barang)
// Auth: Wajib login, role admin
func UpdateBarang(c *gin.Context) {
	idBarangStr := c.Param("id_barang")
	idBarang, err := strconv.Atoi(idBarangStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ID barang tidak valid",
		})
		return
	}

	var barang models.Barang
	if err := config.DB.First(&barang, "id_barang = ?", idBarang).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Barang tidak ditemukan",
		})
		return
	}

	var input struct {
		NamaBarang   string `json:"nama_barang"`
		GambarBarang string `json:"gambar_barang"`
		Deskripsi    string `json:"deskripsi"`
		KategoriId   uint   `json:"id_kategori"`
		SatuanId     uint   `json:"id_satuan"`
		DiskonId     uint   `json:"id_diskon"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Format inputan salah: " + err.Error(),
		})
		return
	}

	updates := map[string]interface{}{}
	if input.NamaBarang != "" {
		updates["nama_barang"] = input.NamaBarang
	}
	if input.GambarBarang != "" {
		updates["gambar_barang"] = input.GambarBarang
	}
	if input.Deskripsi != "" {
		updates["deskripsi"] = input.Deskripsi
	}
	if input.KategoriId != 0 {
		updates["id_kategori"] = input.KategoriId
	}
	if input.SatuanId != 0 {
		updates["id_satuan"] = input.SatuanId
	}
	if input.DiskonId != 0 {
		updates["id_diskon"] = input.DiskonId
	}

	if err := config.DB.Model(&barang).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memperbarui data barang",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Barang berhasil diperbarui",
	})
}

// HapusBarang menghapus barang dari katalog.
// Dipakai oleh: admin (DELETE /admin/katalog/barang/:id_barang)
// Auth: Wajib login, role admin
func HapusBarang(c *gin.Context) {
	idBarangStr := c.Param("id_barang")
	idBarang, err := strconv.Atoi(idBarangStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ID barang tidak valid",
		})
		return
	}

	var barang models.Barang
	if err := config.DB.First(&barang, "id_barang = ?", idBarang).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Barang tidak ditemukan",
		})
		return
	}

	if err := config.DB.Delete(&barang).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menghapus barang",
		})
		return
	}

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
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Parameter pencarian 'q' wajib diisi",
		})
		return
	}

	// Cari di tabel barcode dulu (berdasarkan id_barcode = kode barcode)
	var barcode models.Barcode
	barcodeErr := config.DB.
		Preload("SpesifikasiBarang.Barang.Diskon").
		Preload("SpesifikasiBarang.DetailSpesifikasi.Spesifikasi").
		Where("id_barcode = ?", query).
		First(&barcode).Error

	if barcodeErr == nil {
		// Ditemukan via barcode
		spek := barcode.SpesifikasiBarang
		now := time.Now()
		harga := spek.HargaBarang
		if spek.Barang.DiskonId != 0 && spek.Barang.Diskon.IdDiskon != 0 {
			if spek.Barang.Diskon.TglMulai.Before(now) && spek.Barang.Diskon.TglSelesai.After(now) {
				harga = spek.HargaBarang - (spek.HargaBarang * spek.Barang.Diskon.BesarDiskon / 100)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Produk berhasil ditemukan",
			"data": gin.H{
				"id_barang":             spek.BarangID,
				"nama_barang":           spek.Barang.NamaBarang,
				"gambar_barang":         spek.Barang.GambarBarang,
				"id_spesifikasi_barang": spek.IdSpesifikasiBarang,
				"label":                 spek.DetailSpesifikasi.Spesifikasi.NamaSpesifikasi + " " + spek.DetailSpesifikasi.NamaDetailSpesifikasi,
				"harga_barang":          spek.HargaBarang,
				"harga_diskon":          harga,
				"stok":                  spek.Jumlah,
			},
		})
		return
	}

	// Cari berdasarkan nama barang (LIKE)
	var barangList []models.Barang
	config.DB.Where("nama_barang LIKE ?", "%"+query+"%").Limit(10).Find(&barangList)

	if len(barangList) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Produk tidak ditemukan",
		})
		return
	}

	now := time.Now()
	var responseData []gin.H
	for _, b := range barangList {
		var speks []models.SpesifikasiBarang
		config.DB.Preload("DetailSpesifikasi.Spesifikasi").Where("id_barang = ?", b.IdBarang).Find(&speks)
		config.DB.Preload("Diskon").First(&b, "id_barang = ?", b.IdBarang)

		var varianList []gin.H
		for _, v := range speks {
			harga := v.HargaBarang
			if b.DiskonId != 0 && b.Diskon.IdDiskon != 0 {
				if b.Diskon.TglMulai.Before(now) && b.Diskon.TglSelesai.After(now) {
					harga = v.HargaBarang - (v.HargaBarang * b.Diskon.BesarDiskon / 100)
				}
			}
			varianList = append(varianList, gin.H{
				"id_spesifikasi_barang": v.IdSpesifikasiBarang,
				"label":                 v.DetailSpesifikasi.Spesifikasi.NamaSpesifikasi + " " + v.DetailSpesifikasi.NamaDetailSpesifikasi,
				"harga_barang":          v.HargaBarang,
				"harga_diskon":          harga,
				"stok":                  v.Jumlah,
			})
		}

		responseData = append(responseData, gin.H{
			"id_barang":     b.IdBarang,
			"nama_barang":   b.NamaBarang,
			"gambar_barang": b.GambarBarang,
			"varian":        varianList,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Produk berhasil ditemukan",
		"data":    responseData,
	})
}
