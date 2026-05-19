package stok

import (
	"net/http"
	"time"

	"backend-mantra/config"
	"backend-mantra/models"

	"github.com/gin-gonic/gin"
)

// GetRiwayatStok mengambil riwayat perubahan stok barang.
// Dipakai oleh: admin (GET /admin/stok/riwayat)
// Auth: Wajib login, role admin
func GetRiwayatStok(c *gin.Context) {
	var riwayat []models.StokOpname
	if err := config.DB.
		Preload("SpesifikasiBarang.Barang").
		Preload("SpesifikasiBarang.DetailSpesifikasi.Spesifikasi").
		Order("tanggal desc").
		Find(&riwayat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil riwayat stok",
		})
		return
	}

	var responseData []gin.H
	for _, r := range riwayat {
		varianName := ""
		if r.SpesifikasiBarang.DetailSpesifikasi.Spesifikasi.NamaSpesifikasi != "" {
			varianName = r.SpesifikasiBarang.DetailSpesifikasi.Spesifikasi.NamaSpesifikasi + " " + r.SpesifikasiBarang.DetailSpesifikasi.NamaDetailSpesifikasi
		} else {
			varianName = "Default"
		}

		responseData = append(responseData, gin.H{
			"id_stok_opname":        r.IdStokOpname,
			"id_spesifikasi_barang": r.SpesifikasiBarangID,
			"nama_barang":           r.SpesifikasiBarang.Barang.NamaBarang,
			"varian":                varianName,
			"harga_beli":            r.HargaBeli,
			"status":                r.Status, // true = masuk, false = keluar
			"jumlah_stok":           r.JumlahStok,
			"keterangan":            r.Keterangan,
			"tanggal":               r.Tanggal,
		})
	}

	if responseData == nil {
		responseData = []gin.H{}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Riwayat stok berhasil diambil",
		"data":    responseData,
	})
}

// OpnameStok melakukan penyesuaian (opname) stok barang.
// Dipakai oleh: admin (POST /admin/stok/opname)
// Auth: Wajib login, role admin
func OpnameStok(c *gin.Context) {
	type OpnameInput struct {
		IdSpesifikasiBarang uint   `json:"id_spesifikasi_barang" binding:"required"`
		HargaBeli           int    `json:"harga_beli" binding:"required"`
		Status              *bool  `json:"status" binding:"required"` // Pointer agar booleans false ter-bind benar
		JumlahStok          int    `json:"jumlah_stok" binding:"required"`
		Keterangan          string `json:"keterangan"`
	}

	var input OpnameInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Format inputan salah, pastikan semua data terisi dengan benar",
		})
		return
	}

	tx := config.DB.Begin()

	// 1. Cari spesifikasi barang
	var spek models.SpesifikasiBarang
	if err := tx.Where("id_spesifikasi_barang = ?", input.IdSpesifikasiBarang).First(&spek).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Spesifikasi barang tidak ditemukan",
		})
		return
	}

	// 2. Hitung stok baru
	statusVal := *input.Status
	stokBaru := spek.Jumlah
	if statusVal {
		// Stok masuk
		stokBaru += input.JumlahStok
	} else {
		// Stok keluar
		if spek.Jumlah < input.JumlahStok {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Stok tidak mencukupi untuk melakukan penyesuaian keluar",
			})
			return
		}
		stokBaru -= input.JumlahStok
	}

	// 3. Update jumlah di spesifikasi_barang
	spek.Jumlah = stokBaru
	if err := tx.Save(&spek).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memperbarui stok barang",
		})
		return
	}

	// 4. Catat riwayat opname di stok_opname
	opnameLog := models.StokOpname{
		SpesifikasiBarangID: input.IdSpesifikasiBarang,
		HargaBeli:           input.HargaBeli,
		Status:              statusVal,
		JumlahStok:          input.JumlahStok,
		Keterangan:          input.Keterangan,
		Tanggal:             time.Now(),
	}

	if err := tx.Create(&opnameLog).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menyimpan riwayat opname stok",
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Opname stok berhasil disimpan",
		"data": gin.H{
			"id_stok_opname":        opnameLog.IdStokOpname,
			"id_spesifikasi_barang": opnameLog.SpesifikasiBarangID,
			"stok_baru":             stokBaru,
		},
	})
}
