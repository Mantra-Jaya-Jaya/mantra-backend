package customer

import (
	"net/http"

	"backend-mantra/config"
	"backend-mantra/models"

	"github.com/gin-gonic/gin"
)

// TambahKeKeranjang handles adding items to cart
func TambahKeKeranjang(c *gin.Context) {
	type TambahKeranjangInput struct {
		IdSpesifikasiBarang uint `json:"id_spesifikasi_barang" binding:"required"`
		Quantity            int  `json:"quantity" binding:"required"`
	}

	var input TambahKeranjangInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Input tidak valid, pastikan id_spesifikasi_barang dan quantity diisi",
		})
		return
	}

	// TODO: Ambil ID Customer dari context JWT jika middleware sudah aktif
	// Untuk sementara kita hardcode ID Customer = 1
	customerID := uint(1)

	var keranjang models.Keranjang
	// Cek apakah barang dengan varian tersebut sudah ada di keranjang
	err := config.DB.Where("id_customer = ? AND id_spesifikasi_barang = ?", customerID, input.IdSpesifikasiBarang).First(&keranjang).Error

	if err == nil {
		// Kalau sudah ada, tinggal tambahkan quantity-nya
		keranjang.Quantity += input.Quantity
		if err := config.DB.Save(&keranjang).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Gagal memperbarui item di keranjang",
			})
			return
		}
	} else {
		// Kalau belum ada, buat baris baru
		newItem := models.Keranjang{
			CustomerID:          customerID,
			SpesifikasiBarangID: input.IdSpesifikasiBarang,
			Quantity:            input.Quantity,
		}
		if err := config.DB.Create(&newItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Gagal menambahkan item ke keranjang",
			})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Barang berhasil ditambahkan ke keranjang",
	})
}

// UpdateKeranjang handles updating cart item quantity
func UpdateKeranjang(c *gin.Context) {
	idKeranjang := c.Param("id_keranjang")

	type UpdateKeranjangInput struct {
		Quantity int `json:"quantity" binding:"required"`
	}

	var input UpdateKeranjangInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Input tidak valid, pastikan quantity diisi",
		})
		return
	}

	// TODO: Ambil ID Customer dari context JWT jika middleware sudah aktif
	customerID := uint(1)

	var keranjang models.Keranjang
	// Cari item keranjang berdasarkan ID dan pastikan milik customer yang bersangkutan
	if err := config.DB.Where("id_keranjang = ? AND id_customer = ?", idKeranjang, customerID).First(&keranjang).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Item keranjang tidak ditemukan",
		})
		return
	}

	// Update quantity
	keranjang.Quantity = input.Quantity
	if err := config.DB.Save(&keranjang).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memperbarui item di keranjang",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Keranjang berhasil diperbarui",
	})
}

// HapusItemKeranjang handles deleting item from cart
func HapusItemKeranjang(c *gin.Context) {
	idKeranjang := c.Param("id_keranjang")

	// TODO: Ambil ID Customer dari context JWT jika middleware sudah aktif
	customerID := uint(1)

	// Hapus item dari database jika ID dan Customer ID cocok
	result := config.DB.Where("id_keranjang = ? AND id_customer = ?", idKeranjang, customerID).Delete(&models.Keranjang{})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menghapus item dari keranjang",
		})
		return
	}

	// Jika tidak ada baris yang terhapus, berarti item tidak ditemukan atau bukan milik customer
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Item keranjang tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Item keranjang berhasil dihapus",
	})
}
