package transaksi

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"

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
	if err := config.DB.Where("public_id = ? AND id_customer = ?", input.IdAlamat, result.IdCustomer).First(&alamat).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Alamat tidak ditemukan"})
		return
	}

	totalBerat := 0
	totalNilai := 0
	var daftarItem []services.OngkirItem

	for _, it := range input.Items {
		var spek models.SpesifikasiBarang
		if err := config.DB.Preload("Barang").First(&spek, it.IdSpesifikasiBarang).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("Item dengan ID %d tidak ditemukan", it.IdSpesifikasiBarang),
			})
			return
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
		OriginPostal: os.Getenv("BITESHIP_STORE_POSTAL_CODE"),
		DestPostal:   alamat.KodePos,
		DestLat:      alamat.Latitude,
		DestLng:      alamat.Longitude,
		Items:        daftarItem,
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

	// Build lookup map: kode_api → id_ekspedisi
	var semuaEkspedisi []models.Ekspedisi
	config.DB.Where("is_active = ?", true).Find(&semuaEkspedisi)
	ekspedisiMap := make(map[string]uint)
	for _, e := range semuaEkspedisi {
		ekspedisiMap[e.KodeApi] = e.IdEkspedisi
	}

	// Build lookup map: (kode_ekspedisi + nama_layanan) → id_ekspedisi_layanan
	type EkspedisiLayananRow struct {
		IdEkspedisiLayanan uint
		NamaLayanan        string
		EkspedisiKodeApi   string
	}
	var layananRows []EkspedisiLayananRow
	config.DB.Table("ekspedisi_layanan").
		Select("ekspedisi_layanan.id_ekspedisi_layanan, ekspedisi_layanan.nama_layanan, ekspedisi.kode_api as ekspedisi_kode_api").
		Joins("JOIN ekspedisi ON ekspedisi.id_ekspedisi = ekspedisi_layanan.id_ekspedisi").
		Where("ekspedisi_layanan.is_active = ?", true).
		Scan(&layananRows)

	layananMap := make(map[string]uint)
	for _, l := range layananRows {
		key := l.EkspedisiKodeApi + "|" + l.NamaLayanan
		layananMap[key] = l.IdEkspedisiLayanan
	}

	type LayananDTO struct {
		IdLayananEkspedisi uint   `json:"id_layanan_ekspedisi"`
		KodeLayanan        string `json:"kode_layanan"`
		NamaLayanan        string `json:"nama_layanan"`
		Deskripsi          string `json:"deskripsi"`
		Harga              int    `json:"harga"`
		EstimasiMin        int    `json:"estimasi_min"`
		EstimasiMax        int    `json:"estimasi_max"`
		Durasi             string `json:"durasi"`
	}

	type EkspedisiDTO struct {
		IdEkspedisi   uint         `json:"id_ekspedisi"`
		KodeEkspedisi string       `json:"kode_ekspedisi"`
		NamaEkspedisi string       `json:"nama_ekspedisi"`
		Layanan       []LayananDTO `json:"layanan"`
	}

	var data []EkspedisiDTO
	for _, r := range results {
		var layananList []LayananDTO
		for _, l := range r.Layanan {
			layananKey := fmt.Sprintf("%s|%s", r.EkspedisiKode, l.NamaLayanan)
			layananList = append(layananList, LayananDTO{
				IdLayananEkspedisi: layananMap[layananKey],
				KodeLayanan:        l.KodeLayanan,
				NamaLayanan:        l.NamaLayanan,
				Deskripsi:          l.Deskripsi,
				Harga:              l.Harga,
				EstimasiMin:        l.EstimasiMin,
				EstimasiMax:        l.EstimasiMax,
				Durasi:             l.Durasi,
			})
		}
		data = append(data, EkspedisiDTO{
			IdEkspedisi:   ekspedisiMap[r.EkspedisiKode],
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

// CekRadius memeriksa apakah alamat customer dalam radius kurir internal.
// Dipakai oleh: customer (POST /customer/ongkir/cek-radius)
// Auth: Wajib login, role customer
func CekRadius(c *gin.Context) {
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
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Input tidak valid"})
		return
	}

	var alamat models.Alamat
	if err := config.DB.Where("public_id = ? AND id_customer = ?", input.IdAlamat, result.IdCustomer).First(&alamat).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Alamat tidak ditemukan"})
		return
	}

	if alamat.Latitude == 0 && alamat.Longitude == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Alamat belum memiliki koordinat, tidak bisa cek radius",
			"data": gin.H{
				"within_radius": false,
				"distance_km":   0,
				"max_radius_km": 5,
				"note":          "Alamat belum memiliki koordinat",
			},
		})
		return
	}

	storeLatStr := os.Getenv("BITESHIP_STORE_COORDINATE_LAT")
	storeLngStr := os.Getenv("BITESHIP_STORE_COORDINATE_LONG")
	storeLat, _ := strconv.ParseFloat(storeLatStr, 64)
	storeLng, _ := strconv.ParseFloat(storeLngStr, 64)

	distance := haversine(storeLat, storeLng, alamat.Latitude, alamat.Longitude)

	var setting models.PengaturanToko
	maxRadius := 5.0
	if err := config.DB.Where("key = ?", "radius_kurir_internal").First(&setting).Error; err == nil {
		if val, err := strconv.ParseFloat(setting.Value, 64); err == nil {
			maxRadius = val
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Cek radius berhasil",
		"data": gin.H{
			"within_radius": distance <= maxRadius,
			"distance_km":   math.Round(distance*100) / 100,
			"max_radius_km": maxRadius,
		},
	})
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
