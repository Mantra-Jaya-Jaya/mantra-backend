package user

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"backend-mantra/config"
	"backend-mantra/models"

	"github.com/gin-gonic/gin"
)

// GetProfilAdmin mengambil data profil admin yang sedang login.
// Dipakai oleh: admin (GET /admin/profil)
// Auth: Wajib login, role admin
func GetProfilAdmin(c *gin.Context) {
	uid, exists := getUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User belum login",
			"error":   gin.H{"code": "AUTH_001", "detail": "Token tidak valid"},
		})
		return
	}

	var admin models.User
	if err := config.DB.Preload("Role").Where("id_user = ?", uid).First(&admin).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data admin tidak ditemukan",
			"error":   gin.H{"code": "DATA_004", "detail": "Admin tidak ditemukan di database"},
		})
		return
	}

	fotoProfil := admin.FotoProfil
	if fotoProfil != "" && !strings.HasPrefix(fotoProfil, "http") {
		baseURL := os.Getenv("BASE_URL")
		if baseURL != "" {
			fotoProfil = strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(fotoProfil, "/")
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data profil berhasil diambil",
		"data": gin.H{
			"id_user":      admin.IdUser,
			"public_id":    admin.PublicId,
			"username":     admin.Username,
			"email":        admin.Email,
			"nama_lengkap": admin.NamaLengkap,
			"foto_profil":  fotoProfil,
			"nama_role":    admin.Role.NamaRole,
		},
	})
}

// UpdateProfilAdmin memperbarui data profil admin.
// Dipakai oleh: admin (PUT /admin/profil)
// Auth: Wajib login, role admin
func UpdateProfilAdmin(c *gin.Context) {
	uid, exists := getUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User belum login",
			"error":   gin.H{"code": "AUTH_001", "detail": "Token tidak valid"},
		})
		return
	}

	var input struct {
		NamaLengkap string `json:"nama_lengkap"`
		Email       string `json:"email"`
		Username    string `json:"username"`
		FotoProfil  string `json:"foto_profil"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "error",
			"message": "Validasi gagal",
			"error":   gin.H{"code": "VAL_001", "detail": "Input tidak valid"},
		})
		return
	}

	var admin models.User
	if err := config.DB.Preload("Role").Where("id_user = ?", uid).First(&admin).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data admin tidak ditemukan",
			"error":   gin.H{"code": "DATA_004", "detail": "Admin tidak ditemukan di database"},
		})
		return
	}

	// Cek duplikasi email jika email diubah
	if input.Email != "" && input.Email != admin.Email {
		var existingUser models.User
		if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": "Email sudah terdaftar",
				"error":   gin.H{"code": "CONF_002", "detail": "Email telah digunakan"},
			})
			return
		}
		admin.Email = input.Email
	}

	// Cek duplikasi username jika username diubah
	if input.Username != "" && input.Username != admin.Username {
		var existingUser models.User
		if err := config.DB.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": "Username sudah terdaftar",
				"error":   gin.H{"code": "CONF_001", "detail": "Username telah digunakan"},
			})
			return
		}
		admin.Username = input.Username
	}

	if input.NamaLengkap != "" {
		admin.NamaLengkap = input.NamaLengkap
	}
	if input.FotoProfil != "" {
		admin.FotoProfil = input.FotoProfil
	}

	if err := config.DB.Save(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memperbarui profil admin",
			"error":   gin.H{"code": "SERVER_001", "detail": err.Error()},
		})
		return
	}

	fotoProfil := admin.FotoProfil
	if fotoProfil != "" && !strings.HasPrefix(fotoProfil, "http") {
		baseURL := os.Getenv("BASE_URL")
		if baseURL != "" {
			fotoProfil = strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(fotoProfil, "/")
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Profil admin berhasil diperbarui",
		"data": gin.H{
			"nama_lengkap": admin.NamaLengkap,
			"username":     admin.Username,
			"email":        admin.Email,
			"foto_profil":  fotoProfil,
		},
	})
}

// Helper: Format Nominal Rupiah
func formatNominalRupiah(amount int64) string {
	if amount >= 1000000 {
		return fmt.Sprintf("Rp %.1fM", float64(amount)/1000000.0)
	}
	if amount >= 1000 {
		return fmt.Sprintf("Rp %.1fK", float64(amount)/1000.0)
	}
	return fmt.Sprintf("Rp %d", amount)
}

// Helper: Hitung Trend Persen
func hitungTrendPersen(current, previous int64) string {
	if previous == 0 {
		if current > 0 {
			return "+100%"
		}
		return "0%"
	}
	diff := current - previous
	percent := (float64(diff) / float64(previous)) * 100.0
	if percent > 0 {
		return fmt.Sprintf("+%.1f%%", percent)
	}
	return fmt.Sprintf("%.1f%%", percent)
}

// GetDashboardAdmin mengambil data ringkasan dashboard admin.
// Dipakai oleh: admin (GET /admin/dashboard)
// Auth: Wajib login, role admin
func GetDashboardAdmin(c *gin.Context) {
	// Waktu bulan ini dan bulan lalu
	now := time.Now()
	startOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startOfLastMonth := startOfThisMonth.AddDate(0, -1, 0)
	
	// 1. Total Revenue
	var penjualanHariIni int64
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	config.DB.Model(&models.Pesanan{}).
		Where("status_pesanan = ? AND tanggal_pesanan >= ? AND tanggal_pesanan < ?", "Selesai", startOfDay, endOfDay).
		Select("COALESCE(SUM(total_pembayaran), 0)").Scan(&penjualanHariIni)

	var revenueBulanIni, revenueBulanLalu int64
	config.DB.Model(&models.Pesanan{}).
		Where("status_pesanan = ? AND tanggal_pesanan >= ?", "Selesai", startOfThisMonth).
		Select("COALESCE(SUM(total_pembayaran), 0)").Scan(&revenueBulanIni)
	config.DB.Model(&models.Pesanan{}).
		Where("status_pesanan = ? AND tanggal_pesanan >= ? AND tanggal_pesanan < ?", "Selesai", startOfLastMonth, startOfThisMonth).
		Select("COALESCE(SUM(total_pembayaran), 0)").Scan(&revenueBulanLalu)

	trendRevenue := hitungTrendPersen(revenueBulanIni, revenueBulanLalu)

	// 2. Total Orders
	var totalPesanan int64
	config.DB.Model(&models.Pesanan{}).Count(&totalPesanan)

	var pesananBulanIni, pesananBulanLalu int64
	config.DB.Model(&models.Pesanan{}).Where("tanggal_pesanan >= ?", startOfThisMonth).Count(&pesananBulanIni)
	config.DB.Model(&models.Pesanan{}).Where("tanggal_pesanan >= ? AND tanggal_pesanan < ?", startOfLastMonth, startOfThisMonth).Count(&pesananBulanLalu)
	trendPesanan := hitungTrendPersen(pesananBulanIni, pesananBulanLalu)

	// 3. Active Customers
	var totalCustomerAktif int64
	config.DB.Model(&models.Pesanan{}).Where("status_pesanan = ?", "Selesai").Select("COUNT(DISTINCT id_customer)").Scan(&totalCustomerAktif)

	var custBulanIni, custBulanLalu int64
	config.DB.Model(&models.Pesanan{}).Where("status_pesanan = ? AND tanggal_pesanan >= ?", "Selesai", startOfThisMonth).Select("COUNT(DISTINCT id_customer)").Scan(&custBulanIni)
	config.DB.Model(&models.Pesanan{}).Where("status_pesanan = ? AND tanggal_pesanan >= ? AND tanggal_pesanan < ?", "Selesai", startOfLastMonth, startOfThisMonth).Select("COUNT(DISTINCT id_customer)").Scan(&custBulanLalu)
	trendCustomer := hitungTrendPersen(custBulanIni, custBulanLalu)

	// 4. Low Stock Items
	var totalStokMenipis int64
	config.DB.Model(&models.SpesifikasiBarang{}).Where("jumlah <= 10").Count(&totalStokMenipis)

	var stokMenipisList []models.SpesifikasiBarang
	config.DB.
		Preload("Barang").
		Preload("DetailSpesifikasi.Spesifikasi").
		Where("jumlah <= 10").
		Find(&stokMenipisList)

	type LowStockData struct {
		Id     uint   `json:"id"`
		Nama   string `json:"nama"`
		Sisa   int    `json:"sisa"`
		Status string `json:"status"`
	}

	var stokMenipisResponse []LowStockData
	for _, s := range stokMenipisList {
		varianName := ""
		if s.DetailSpesifikasi.Spesifikasi.NamaSpesifikasi != "" {
			varianName = s.DetailSpesifikasi.Spesifikasi.NamaSpesifikasi + " " + s.DetailSpesifikasi.NamaDetailSpesifikasi
		} else {
			varianName = "Default"
		}
		
		statusStok := "warning"
		if s.Jumlah <= 5 {
			statusStok = "kritis"
		}

		stokMenipisResponse = append(stokMenipisResponse, LowStockData{
			Id:     s.IdSpesifikasiBarang,
			Nama:   s.Barang.NamaBarang + " " + varianName,
			Sisa:   s.Jumlah,
			Status: statusStok,
		})
	}
	if stokMenipisResponse == nil {
		stokMenipisResponse = []LowStockData{}
	}

	// 5. Transaksi Terbaru
	var pesananTerbaru []models.Pesanan
	config.DB.
		Preload("Kasir.User").
		Preload("Customer.User").
		Order("tanggal_pesanan DESC").
		Limit(50).
		Find(&pesananTerbaru)

	type TransaksiData struct {
		Id        string `json:"id"`
		Kasir     string `json:"kasir"`
		Pelanggan string `json:"pelanggan"`
		Tanggal   string `json:"tanggal"`
		Total     string `json:"total"`
		Status    string `json:"status"`
	}

	var transaksiResponse []TransaksiData
	for _, p := range pesananTerbaru {
		kasirName := "-"
		if p.Kasir.User.NamaLengkap != "" {
			kasirName = p.Kasir.User.NamaLengkap
		}
		custName := "-"
		if p.Customer.User.NamaLengkap != "" {
			custName = p.Customer.User.NamaLengkap
		}
		
		transaksiResponse = append(transaksiResponse, TransaksiData{
			Id:        p.PublicId.String(),
			Kasir:     kasirName,
			Pelanggan: custName,
			Tanggal:   p.TanggalPesanan.Format("02 Jan 2006, 15:04"),
			Total:     formatNominalRupiah(int64(p.TotalPembayaran)),
			Status:    p.StatusPesanan,
		})
	}
	if transaksiResponse == nil {
		transaksiResponse = []TransaksiData{}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data dashboard berhasil diambil",
		"data": gin.H{
			"penjualan_hari_ini":   penjualanHariIni,
			"total_pesanan":        totalPesanan,
			"total_customer_aktif": totalCustomerAktif,
			"total_stok_menipis":   totalStokMenipis,
			"trend_revenue":        trendRevenue,
			"trend_pesanan":        trendPesanan,
			"trend_customer":       trendCustomer,
			"stok_menipis":         stokMenipisResponse,
			"transaksi_terbaru":    transaksiResponse,
		},
	})
}

// GetChartDashboardAdmin mengambil data chart untuk dashboard admin.
// Dipakai oleh: admin (GET /admin/dashboard/chart)
// Auth: Wajib login, role admin
func GetChartDashboardAdmin(c *gin.Context) {
	periode := c.DefaultQuery("periode", "minggu")
	tanggalStr := c.Query("tanggal")
	
	now := time.Now()
	if tanggalStr != "" {
		if parsedTime, err := time.Parse("2006-01-02", tanggalStr); err == nil {
			now = parsedTime
		}
	}

	type BarData struct {
		Name  string `json:"name"`
		Total int64  `json:"total"`
	}
	var bars []BarData
	var label string
	var monthNames = []string{"Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}

	if periode == "minggu" {
		// Cari hari Senin
		day := int(now.Weekday())
		if day == 0 {
			day = 7
		}
		senin := now.AddDate(0, 0, -day+1)
		minggu := senin.AddDate(0, 0, 6)

		label = fmt.Sprintf("%d %s - %d %s %d", 
			senin.Day(), monthNames[senin.Month()-1][:3], 
			minggu.Day(), monthNames[minggu.Month()-1][:3], minggu.Year())

		namaHari := []string{"Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu", "Minggu"}
		for i := 0; i < 7; i++ {
			targetDate := senin.AddDate(0, 0, i)
			startOfDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, targetDate.Location())
			endOfDay := startOfDay.Add(24 * time.Hour)

			var total int64
			config.DB.Model(&models.Pesanan{}).
				Where("status_pesanan = ? AND tanggal_pesanan >= ? AND tanggal_pesanan < ?", "Selesai", startOfDay, endOfDay).
				Select("COALESCE(SUM(total_pembayaran), 0)").Scan(&total)

			bars = append(bars, BarData{Name: namaHari[i], Total: total})
		}
	} else if periode == "bulan" {
		label = fmt.Sprintf("%s %d", monthNames[now.Month()-1], now.Year())
		
		// 5 minggu
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		
		for i := 0; i < 5; i++ {
			startOfWeek := startOfMonth.AddDate(0, 0, i*7)
			endOfWeek := startOfWeek.AddDate(0, 0, 7)
			
			// Jika startOfWeek sudah beda bulan, skip
			if startOfWeek.Month() != now.Month() && i == 4 {
				continue
			}

			var total int64
			config.DB.Model(&models.Pesanan{}).
				Where("status_pesanan = ? AND tanggal_pesanan >= ? AND tanggal_pesanan < ?", "Selesai", startOfWeek, endOfWeek).
				Select("COALESCE(SUM(total_pembayaran), 0)").Scan(&total)

			bars = append(bars, BarData{Name: fmt.Sprintf("Minggu %d", i+1), Total: total})
		}
	} else if periode == "tahun" {
		label = fmt.Sprintf("%d", now.Year())
		
		shortMonthNames := []string{"Jan", "Feb", "Mar", "Apr", "Mei", "Jun", "Jul", "Ags", "Sep", "Okt", "Nov", "Des"}
		for i := 1; i <= 12; i++ {
			startOfMonth := time.Date(now.Year(), time.Month(i), 1, 0, 0, 0, 0, now.Location())
			endOfMonth := startOfMonth.AddDate(0, 1, 0)

			var total int64
			config.DB.Model(&models.Pesanan{}).
				Where("status_pesanan = ? AND tanggal_pesanan >= ? AND tanggal_pesanan < ?", "Selesai", startOfMonth, endOfMonth).
				Select("COALESCE(SUM(total_pembayaran), 0)").Scan(&total)

			bars = append(bars, BarData{Name: shortMonthNames[i-1], Total: total})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"bars": bars,
			"label": label,
		},
	})
}
