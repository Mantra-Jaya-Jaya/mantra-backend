package stok

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetRiwayatStok mengambil riwayat perubahan stok barang.
// Dipakai oleh: admin (GET /admin/stok/riwayat)
// Auth: Wajib login, role admin
func GetRiwayatStok(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Riwayat stok berhasil diambil",
		"data":    []gin.H{},
	})
}

// OpnameStok melakukan penyesuaian (opname) stok barang.
// Dipakai oleh: admin (POST /admin/stok/opname)
// Auth: Wajib login, role admin
func OpnameStok(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Opname stok berhasil disimpan",
	})
}
