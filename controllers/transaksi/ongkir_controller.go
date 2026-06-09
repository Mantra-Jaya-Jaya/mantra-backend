package transaksi

import (
	"net/http"

	"backend-mantra/config"
	"backend-mantra/models"
	"backend-mantra/services"

	"github.com/gin-gonic/gin"
)

func CekOngkir(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "User belum login"})
		return
	}

	var result struct{ IdCustomer uint }
	if err := config.DB.Raw("SELECT id_customer FROM customer WHERE id_user = ?", userID).Scan(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengidentifikasi customer"})
		return
	}

	var input struct {
		IdAlamat string `json:"id_alamat" binding:"required"`
		Items    []struct {
			IdSpesifikasiBarang uint `json:"id_spesifikasi_barang" binding:"required"`
			Quantity            int  `json:"quantity" binding:"required"`
		} `json:"items" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Input tidak valid"})
		return
	}

	var alamat models.Alamat
	if err := config.DB.Where("public_id = ?", input.IdAlamat).First(&alamat).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Alamat tidak ditemukan"})
		return
	}

	totalBerat := 0
	totalNilai := 0
	var daftarItem []services.OngkirItem

	for _, it := range input.Items {
		var spek models.SpesifikasiBarang
		if err := config.DB.Preload("Barang").First(&spek, it.IdSpesifikasiBarang).Error; err != nil {
			continue
		}
		berat := spek.BeratBarang * it.Quantity
		totalBerat += berat
		nilai := spek.HargaBarang * it.Quantity
		totalNilai += nilai

		daftarItem = append(daftarItem, services.OngkirItem{
			Name:   spek.Barang.NamaBarang,
			Weight: spek.BeratBarang,
			Length: spek.Barang.PanjangBarang,
			Width:  spek.Barang.LebarBarang,
			Height: spek.Barang.TinggiBarang,
			Value:  spek.HargaBarang,
		})
	}

	ongkirReq := services.OngkirRequest{
		DestPostal: alamat.KodePos,
		DestLat:    alamat.Latitude,
		DestLng:    alamat.Longitude,
		Items:      daftarItem,
	}

	adapter := services.NewBiteshipAdapter()
	results, err := adapter.CekOngkir(ongkirReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengecek ongkos kirim: " + err.Error(),
		})
		return
	}

	type LayananDTO struct {
		KodeLayanan string `json:"kode_layanan"`
		NamaLayanan string `json:"nama_layanan"`
		Deskripsi   string `json:"deskripsi"`
		Harga       int    `json:"harga"`
		EstimasiMin int    `json:"estimasi_min"`
		EstimasiMax int    `json:"estimasi_max"`
		Durasi      string `json:"durasi"`
	}

	type EkspedisiDTO struct {
		KodeEkspedisi string       `json:"kode_ekspedisi"`
		NamaEkspedisi string       `json:"nama_ekspedisi"`
		Layanan       []LayananDTO `json:"layanan"`
	}

	var data []EkspedisiDTO
	for _, r := range results {
		var layananList []LayananDTO
		for _, l := range r.Layanan {
			layananList = append(layananList, LayananDTO{
				KodeLayanan: l.KodeLayanan,
				NamaLayanan: l.NamaLayanan,
				Deskripsi:   l.Deskripsi,
				Harga:       l.Harga,
				EstimasiMin: l.EstimasiMin,
				EstimasiMax: l.EstimasiMax,
				Durasi:      l.Durasi,
			})
		}
		data = append(data, EkspedisiDTO{
			KodeEkspedisi: r.EkspedisiKode,
			NamaEkspedisi: r.NamaEkspedisi,
			Layanan:       layananList,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Daftar ongkos kirim berhasil diambil",
		"data":    data,
	})
}
