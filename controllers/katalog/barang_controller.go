package katalog

import (
	"net/http"
	"strconv"
	"time"

	"backend-mantra/config"
	"backend-mantra/models"
	"backend-mantra/utils"

	"github.com/gin-gonic/gin"
)

// GetDaftarBarang mengambil daftar barang dari katalog dengan pagination.
// Dipakai oleh: customer (GET /customer/katalog/barang), kasir (GET /kasir/katalog/barang), admin (GET /admin/katalog/barang)
// Auth: Wajib login, semua role boleh akses (dikontrol di route)
func GetDaftarBarang(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "50")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	offset := (page - 1) * limit

	// 🔥 Struct yang sudah di-upgrade kebal Badai NULL
	type ProductResult struct {
		IdBarang       uint       `gorm:"column:id_barang"`
		PublicId       string     `gorm:"column:public_id"`
		NamaBarang     string     `gorm:"column:nama_barang"`
		GambarBarang   *string    `gorm:"column:gambar_barang"` // Pakai pointer biar aman dari NULL
		HargaTerendah  int        `gorm:"column:harga_terendah"`
		HargaTertinggi int        `gorm:"column:harga_tertinggi"`
		BesarDiskon    int        `gorm:"column:besar_diskon"`
		TglMulai       *time.Time `gorm:"column:tgl_mulai"`
		TglSelesai     *time.Time `gorm:"column:tgl_selesai"`
		NamaKategori   string     `gorm:"column:nama_kategori"`
		TotalStok      int        `gorm:"column:total_stok"`
	}

	var results []ProductResult
	var total int64

	kategoriIdStr := c.Query("id_kategori")

	queryCount := config.DB.Model(&models.Barang{})
	querySelect := config.DB.Table("barang")

	if kategoriIdStr != "" {
		kategoriId, err := strconv.Atoi(kategoriIdStr)
		if err == nil && kategoriId > 0 {
			queryCount = queryCount.Where("id_kategori = ?", kategoriId)
			querySelect = querySelect.Where("barang.id_kategori = ?", kategoriId)
		}
	}

	// Hitung total barang
	queryCount.Count(&total)

	// Ambil data barang dengan join spesifikasi, diskon, dan kategori
	querySelect.
		Select("barang.id_barang, barang.public_id, barang.nama_barang, barang.gambar_barang, MIN(spesifikasi_barang.harga_barang) as harga_terendah, MAX(spesifikasi_barang.harga_barang) as harga_tertinggi, SUM(spesifikasi_barang.jumlah) as total_stok, diskon.besar_diskon, diskon.tgl_mulai, diskon.tgl_selesai, kategori.nama_kategori").
		Joins("LEFT JOIN spesifikasi_barang ON spesifikasi_barang.id_barang = barang.id_barang").
		Joins("LEFT JOIN diskon ON diskon.id_diskon = barang.id_diskon").
		Joins("LEFT JOIN kategori ON kategori.id_kategori = barang.id_kategori").
		Group("barang.id_barang, barang.public_id, barang.nama_barang, barang.gambar_barang, diskon.besar_diskon, diskon.tgl_mulai, diskon.tgl_selesai, kategori.nama_kategori").
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

		// 🔥 Ekstrak URL gambar dengan aman dari Pointer
		gambarValid := ""
		if r.GambarBarang != nil {
			gambarValid = *r.GambarBarang
		}

		responseData = append(responseData, gin.H{
			"id_barang":       r.IdBarang,
			"public_id":       r.PublicId,
			"nama_barang":     r.NamaBarang,
			"harga_terendah":  r.HargaTerendah,
			"harga_tertinggi": r.HargaTertinggi,
			"harga_diskon":    hargaDiskon,
			"punya_diskon":    punyaDiskon,
			"gambar_barang":   gambarValid,
			"kategori":        r.NamaKategori,
			"stok":            r.TotalStok,
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
  // 🚀 1. Ambil param public_id (String/UUID)
  publicId := c.Param("public_id")

  var barang models.Barang
  // 🚀 2. Cari berdasarkan kolom public_id
  if err := config.DB.
    Preload("Kategori").
    Preload("Satuan").
    Preload("Diskon").
    First(&barang, "public_id = ?", publicId).Error; err != nil {
    c.JSON(http.StatusNotFound, gin.H{
      "status":  "error",
      "message": "Barang tidak ditemukan",
    })
    return
  }

  // 🚀 3. Ambil semua varian
  var varians []models.SpesifikasiBarang
  config.DB.
    Preload("DetailSpesifikasi.Spesifikasi").
    Where("id_barang = ?", barang.IdBarang).
    Find(&varians)

  now := time.Now()
  var responseVarian []gin.H
  for _, v := range varians {
    hargaDiskon := v.HargaBarang
    if barang.DiskonId != nil && barang.Diskon.IdDiskon != 0 {
      if barang.Diskon.TglMulai.Before(now) && barang.Diskon.TglSelesai.After(now) {
        hargaDiskon = v.HargaBarang - (v.HargaBarang * barang.Diskon.BesarDiskon / 100)
      }
    }

    // 🚀 4. Ambil Barcode untuk varian ini
    var barcodes []models.Barcode
    // Pastikan nama field string "id_spesifikasi_barang" sesuai dengan nama kolom di database Postgres lu
    config.DB.Where("id_spesifikasi_barang = ?", v.IdSpesifikasiBarang).Find(&barcodes)
    
    var responseBarcodes []gin.H
    for _, b := range barcodes {
      responseBarcodes = append(responseBarcodes, gin.H{
        "id_barcode":   b.IdBarcode,
        "kode_barcode": b.KodeBarcode, // 🔥 WAJIB ADA: Biar frontend bisa nampilin teks barcodenya
        "kuantitas":    b.Kuantitas,   // 🔥 FIX: Hapus SatuanId karena emang gak ada di tabel barcode
      })
    }
    if responseBarcodes == nil {
      responseBarcodes = []gin.H{}
    }

    responseVarian = append(responseVarian, gin.H{
      "id_spesifikasi_barang": v.IdSpesifikasiBarang,
      "nama_spesifikasi":      v.DetailSpesifikasi.Spesifikasi.NamaSpesifikasi,
      "nama_detail":           v.DetailSpesifikasi.NamaDetailSpesifikasi,
      "harga_barang":          v.HargaBarang,
      "harga_diskon":          hargaDiskon,
      "stok":                  v.Jumlah,
      "barcodes":              responseBarcodes, // 🔥 SUNTIKKAN BARCODE KE DALAM VARIAN
    })
  }
  
  if responseVarian == nil {
    responseVarian = []gin.H{}
  }

  var diskonData interface{} = nil
  if barang.DiskonId != nil && barang.Diskon.IdDiskon != 0 {
    diskonData = gin.H{
      "nama_diskon":  barang.Diskon.NamaDiskon,
      "besar_diskon": barang.Diskon.BesarDiskon,
      "tgl_selesai":  barang.Diskon.TglSelesai,
    }
  }

  // 🚀 5. Kirim response balik
  c.JSON(http.StatusOK, gin.H{
    "status":  "success",
    "message": "Detail barang berhasil diambil",
    "data": gin.H{
      "id_barang":     barang.IdBarang,
      "public_id":     barang.PublicId,
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
    Media           string `json:"media"`
    InformasiBarang struct {
      Nama      string `json:"nama" binding:"required"`
      HargaBeli string `json:"hargaBeli"`
      Kategori  string `json:"kategori" binding:"required"`
      Satuan    string `json:"satuan" binding:"required"`
      Deskripsi string `json:"deskripsi"`
    } `json:"informasi_barang"`
    Spesifikasi []struct {
      Atribut   string `json:"atribut"`
      Nilai     string `json:"nilai"`
      Stok      string `json:"stok"`
      HargaJual string `json:"hargaJual"`
      Barcodes  []struct {
        Code string `json:"code"`
        Qty  string `json:"qty"`
      } `json:"barcodes"`
    } `json:"spesifikasi"`
  }

  if err := c.ShouldBindJSON(&input); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format JSON tidak sesuai: " + err.Error()})
    return
  }

  // MULAI TRANSAKSI DATABASE
  tx := config.DB.Begin()
  if tx.Error != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal memulai transaksi DB"})
    return
  }

  // PROSES KATEGORI & SATUAN
  var kategori models.Kategori
  if err := tx.FirstOrCreate(&kategori, models.Kategori{NamaKategori: input.InformasiBarang.Kategori}).Error; err != nil {
    tx.Rollback()
    c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal memproses kategori"})
    return
  }

  var satuan models.Satuan
  if err := tx.FirstOrCreate(&satuan, models.Satuan{NamaSatuan: input.InformasiBarang.Satuan}).Error; err != nil {
    tx.Rollback()
    c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal memproses satuan"})
    return
  }

  // SIMPAN INDUK BARANG
  barang := models.Barang{
    NamaBarang:   input.InformasiBarang.Nama,
    GambarBarang: input.Media,
    Deskripsi:    input.InformasiBarang.Deskripsi,
    KategoriId:   kategori.IdKategori,
    SatuanId:     satuan.IdSatuan,
  }

  if err := tx.Create(&barang).Error; err != nil {
    tx.Rollback()
    c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menyimpan induk barang"})
    return
  }

  hargaBeliInt, _ := strconv.Atoi(input.InformasiBarang.HargaBeli)

  // LOOPING PEMETAAN VARIAN & BARCODE SESUAI MODELS LU
  for _, spek := range input.Spesifikasi {
    
    // A. Atribut Induk
    var spesifikasi models.Spesifikasi
    if err := tx.FirstOrCreate(&spesifikasi, models.Spesifikasi{NamaSpesifikasi: spek.Atribut}).Error; err != nil {
      tx.Rollback()
      c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal memproses tabel spesifikasi"})
      return
    }

    // B. Detail Nilai Atribut
    var detail models.DetailSpesifikasi
    if err := tx.FirstOrCreate(&detail, models.DetailSpesifikasi{
      SpesifikasiID:         spesifikasi.IdSpesifikasi, // 🔥 FIX: Menggunakan SpesifikasiID sesuai model lu
      NamaDetailSpesifikasi: spek.Nilai,
    }).Error; err != nil {
      tx.Rollback()
      c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal memproses tabel detail_spesifikasi"})
      return
    }

    stokInt, _ := strconv.Atoi(spek.Stok)
    hargaJualInt, _ := strconv.Atoi(spek.HargaJual)

    // C. Gabungkan ke SpesifikasiBarang
    spekBarang := models.SpesifikasiBarang{
      BarangID:            barang.IdBarang,            // 🔥 FIX: Menggunakan BarangID
      DetailSpesifikasiID: detail.IdDetailSpesifikasi, // 🔥 FIX: Menggunakan DetailSpesifikasiID
      Jumlah:              stokInt,
      HargaBarang:         hargaJualInt,
    }
    if err := tx.Create(&spekBarang).Error; err != nil {
      tx.Rollback()
      c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menyimpan spesifikasi_barang"})
      return
    }

    // D. Catat Saldo Awal ke StokOpname
    stokOpname := models.StokOpname{
      SpesifikasiBarangID: spekBarang.IdSpesifikasiBarang, // 🔥 FIX: Menggunakan SpesifikasiBarangID
      HargaBeli:           hargaBeliInt,
      Status:              true, 
      JumlahStok:          stokInt,
      Keterangan:          "Stok awal dari tambah barang baru",
      Tanggal:             time.Now(),
    }
    if err := tx.Create(&stokOpname).Error; err != nil {
      tx.Rollback()
      c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mencatat stok_opname"})
      return
    }

    // E. Looping Barcode
    for _, bc := range spek.Barcodes {
      qtyInt, _ := strconv.Atoi(bc.Qty)
      
      barcodeData := models.Barcode{
        SpesifikasiBarangID: spekBarang.IdSpesifikasiBarang, // 🔥 FIX: Menggunakan SpesifikasiBarangID
        KodeBarcode:         bc.Code,                         // 🔥 Cek catatan di bawah jika baris ini masih merah
        Kuantitas:           uint(qtyInt),                    // 🔥 FIX: Konversi aman ke uint
      }
      if err := tx.Create(&barcodeData).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menyimpan barcode"})
        return
      }
    }
  }

  // COMMIT TRANSAKSI
  if err := tx.Commit().Error; err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menyelesaikan transaksi database"})
    return
  }

  config.DB.First(&barang, barang.IdBarang)

  c.JSON(http.StatusCreated, gin.H{
    "status":  "success",
    "message": "Barang baru berhasil ditambahkan!",
    "data": gin.H{
      "id_barang":   barang.IdBarang,
      "public_id":   barang.PublicId,
      "nama_barang": barang.NamaBarang,
    },
  })
}

// UpdateBarang memperbarui detail barang yang ada.
// Dipakai oleh: admin (PUT /admin/katalog/barang/:id_barang)
// Auth: Wajib login, role admin
func UpdateBarang(c *gin.Context) {
	idBarangStr := c.Param("id_barang")

	var barang models.Barang
	if err := config.DB.First(&barang, "public_id = ?", idBarangStr).Error; err != nil {
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

	var barang models.Barang
	if err := config.DB.First(&barang, "public_id = ?", idBarangStr).Error; err != nil {
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
		if barcode.SpesifikasiBarang.Barang.DiskonId != nil && barcode.SpesifikasiBarang.Barang.Diskon.IdDiskon != 0 {
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
	if barcode.SpesifikasiBarang.Barang.DiskonId != nil && barcode.SpesifikasiBarang.Barang.Diskon.IdDiskon != 0 {
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
		if spek.Barang.DiskonId != nil && spek.Barang.Diskon.IdDiskon != 0 {
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
			if b.DiskonId != nil && b.Diskon.IdDiskon != 0 {
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


// UploadGambarBarang handle proses upload dari Next.js ke MinIO
// Dipakai oleh: POST /api/v1/admin/upload
func UploadGambarBarang(c *gin.Context) {
  // 1. Panggil helper sakti yang kita buat kemarin
  // Kita set fileKey-nya "gambar" dan folderTarget-nya "produk"
  fileUrl, err := utils.UploadFileToMinio(c, "gambar", "produk")
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "status":  "error",
      "message": "Gagal mengunggah gambar: " + err.Error(),
    })
    return
  }

  // 2. Kembalikan URL publik MinIO ke Next.js
  c.JSON(http.StatusOK, gin.H{
    "status":  "success",
    "message": "Gambar berhasil diunggah ke server storage",
    "url":     fileUrl, // URL ini yang nanti dikirim Next.js ke TambahBarang
  })
}
