package routes

import (
	"backend-mantra/controllers/auth"
	"backend-mantra/controllers/katalog"
	"backend-mantra/controllers/keranjang"
	"backend-mantra/controllers/notifikasi"
	"backend-mantra/controllers/transaksi"
	"backend-mantra/controllers/user"
	"backend-mantra/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		// Public Auth Routes
		v1.POST("/login", auth.Login)
		v1.POST("/register", auth.RegisterCustomer)
		v1.POST("/auth/refresh", auth.RefreshToken)

		// Protected Auth Routes
		authGroup := v1.Group("/")
		authGroup.Use(middleware.AuthMiddleware())
		{
			authGroup.POST("/logout", auth.Logout)
			authGroup.PUT("/change-password", auth.ChangePassword)
		}

		// Shared Routes
		v1.GET("/scan/:kode_barcode", katalog.GetDetailBarangByScan)

		// Customer Routes
		customerGroup := v1.Group("/customer")
		customerGroup.Use(middleware.AuthMiddleware())
		{
			customerGroup.GET("/promo", katalog.GetPromo)
			customerGroup.GET("/kategori", katalog.GetKategori)
			customerGroup.GET("/barang", katalog.GetDaftarBarang)
			customerGroup.POST("/keranjang", keranjang.TambahKeKeranjang)
			customerGroup.PATCH("/keranjang/:id_keranjang", keranjang.UpdateKeranjang)
			customerGroup.DELETE("/keranjang/:id_keranjang", keranjang.HapusItemKeranjang)
			customerGroup.GET("/notifikasi", notifikasi.GetNotifikasi)
			customerGroup.GET("/pesanan", transaksi.GetDaftarPesanan)
			customerGroup.POST("/pesanan/checkout", transaksi.CheckoutPesanan)
			customerGroup.PATCH("/pesanan/:id_pesanan/batal", transaksi.BatalkanPesanan)
			customerGroup.GET("/pesanan/:id_pesanan", transaksi.GetDetailPesanan)
			customerGroup.GET("/pesanan/:id_pesanan/lacak", transaksi.LacakPesanan)
			customerGroup.GET("/profil", user.GetProfilCustomer)
			customerGroup.PUT("/akun", user.EditAkunCustomer)
			customerGroup.GET("/alamat", user.GetAlamat)
			customerGroup.POST("/alamat", user.TambahAlamat)
			customerGroup.PUT("/alamat/:id_alamat", user.UpdateAlamat)
			customerGroup.DELETE("/alamat/:id_alamat", user.HapusAlamat)
		}

		// Kasir Routes
		kasirGroup := v1.Group("/kasir")
		kasirGroup.Use(middleware.AuthMiddleware())
		{
			kasirGroup.GET("/dashboard", transaksi.GetDashboardKasir)
			kasirGroup.GET("/laporan", transaksi.GetLaporanRingkasan)
			kasirGroup.GET("/laporan/produk/:id_produk", transaksi.GetDetailLaporanProduk)
			kasirGroup.GET("/laporan/produk/:id_produk/:id_pesanan", transaksi.GetDetailPesananDariLaporan)
			kasirGroup.GET("/pesanan", transaksi.GetDaftarPesanan)
			kasirGroup.GET("/pesanan/:id_order", transaksi.GetDetailPesanan)
			kasirGroup.GET("/kategori", katalog.GetKategori)
			kasirGroup.POST("/transaksi/produk", katalog.CariProdukTransaksi)
			kasirGroup.PATCH("/transaksi/item/update", transaksi.UpdateQuantityItem)
			kasirGroup.GET("/transaksi/checkout", transaksi.GetRingkasanCheckout)
			kasirGroup.POST("/transaksi/bayar/tunai", transaksi.BayarTunai)
			kasirGroup.POST("/transaksi/bayar/non-tunai", transaksi.BayarNonTunai)
			kasirGroup.GET("/profil", user.GetProfilKasir)
			kasirGroup.GET("/notifikasi", notifikasi.GetNotifikasi)
		}

		// Admin Routes
		adminGroup := v1.Group("/admin")
		adminGroup.Use(middleware.AuthMiddleware())
		{
			adminGroup.GET("/dashboard", user.GetDashboardAdmin)
<<<<<<< HEAD
=======
			adminGroup.GET("/dashboard/chart", user.GetChartDashboardAdmin)
			adminGroup.GET("/kategori", katalog.GetKategori)
			adminGroup.POST("/kategori", katalog.TambahKategori)
			adminGroup.PUT("/kategori/:id_kategori", katalog.UpdateKategori)
			adminGroup.DELETE("/kategori/:id_kategori", katalog.HapusKategori)
>>>>>>> c7388ad36c44b06dee7d3b6348a99f5903761e59
			adminGroup.GET("/barang", katalog.GetDaftarBarang)
			adminGroup.POST("/barang", katalog.TambahBarang)
			adminGroup.GET("/barang/detail/:public_id", katalog.GetDetailBarang)
			adminGroup.POST("/barang/upload", katalog.UploadGambarBarang)
			adminGroup.PUT("/barang/:id_barang", katalog.UpdateBarang)
			adminGroup.DELETE("/barang/:id_barang", katalog.HapusBarang)
			adminGroup.POST("/barang/:id_barang/diskon", katalog.TambahDiskon)
			adminGroup.GET("/karyawan", user.GetDaftarKaryawan)
			adminGroup.GET("/karyawan/kasir", user.GetDaftarKasir)
			adminGroup.POST("/karyawan/kasir", user.TambahKasir)
			adminGroup.GET("/karyawan/kasir/:id_kasir", user.GetDetailKasir)
			adminGroup.PUT("/karyawan/kasir/:id_kasir", user.UpdateKasir)
			adminGroup.DELETE("/karyawan/kasir/:id_kasir", user.HapusKasir)
			adminGroup.GET("/karyawan/kurir", user.GetDaftarKurir)
			adminGroup.POST("/karyawan/kurir", user.TambahKurir)
			adminGroup.GET("/karyawan/kurir/:id_kurir", user.GetDetailKurir)
			adminGroup.PUT("/karyawan/kurir/:id_kurir", user.UpdateKurir)
			adminGroup.DELETE("/karyawan/kurir/:id_kurir", user.HapusKurir)
			adminGroup.GET("/notifikasi", notifikasi.GetNotifikasiAdmin)
			adminGroup.GET("/profil", user.GetProfilAdmin)
			adminGroup.PUT("/profil", user.UpdateProfilAdmin)
		}
	}
}
