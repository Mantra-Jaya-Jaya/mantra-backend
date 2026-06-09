package routes

import (
	"backend-mantra/controllers/auth"
	"backend-mantra/controllers/katalog"
	"backend-mantra/controllers/keranjang"
	"backend-mantra/controllers/notifikasi"
	"backend-mantra/controllers/transaksi"
	"backend-mantra/controllers/user"
	"backend-mantra/controllers/pengantaran"
	"backend-mantra/controllers/pemesanan"
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
			customerGroup.PATCH("/keranjang/:public_id", keranjang.UpdateKeranjang)
			customerGroup.DELETE("/keranjang/:public_id", keranjang.HapusItemKeranjang)
			customerGroup.GET("/notifikasi", notifikasi.GetNotifikasi)
			customerGroup.GET("/pesanan", transaksi.GetDaftarPesanan)
			customerGroup.POST("/pesanan/checkout", transaksi.CheckoutPesanan)
			customerGroup.PATCH("/pesanan/:public_id/batal", transaksi.BatalkanPesanan)
			customerGroup.GET("/pesanan/:public_id", transaksi.GetDetailPesanan)
			customerGroup.GET("/pesanan/:public_id/lacak", transaksi.LacakPesanan)
			customerGroup.GET("/profil", user.GetProfilCustomer)
			customerGroup.PUT("/akun", user.EditAkunCustomer)
			customerGroup.GET("/alamat", user.GetAlamat)
			customerGroup.POST("/alamat", user.TambahAlamat)
			customerGroup.PUT("/alamat/:public_id", user.UpdateAlamat)
			customerGroup.DELETE("/alamat/:public_id", user.HapusAlamat)
		}

		// Kasir Routes
		kasirGroup := v1.Group("/kasir")
		kasirGroup.Use(middleware.AuthMiddleware())
		{
			kasirGroup.GET("/dashboard", transaksi.GetDashboardKasir)
			kasirGroup.GET("/laporan", transaksi.GetLaporanRingkasan)
			kasirGroup.GET("/laporan/produk/:public_id", transaksi.GetDetailLaporanProduk)
			kasirGroup.GET("/laporan/produk/:public_id/:pesanan_id", transaksi.GetDetailPesananDariLaporan)
			kasirGroup.GET("/pesanan", transaksi.GetDaftarPesanan)
			kasirGroup.GET("/pesanan/:public_id", transaksi.GetDetailPesanan)
			kasirGroup.GET("/kategori", katalog.GetKategori)
			kasirGroup.POST("/transaksi/produk", katalog.CariProdukTransaksi)
			kasirGroup.PATCH("/transaksi/item/update", transaksi.UpdateQuantityItem)
			kasirGroup.GET("/transaksi/checkout", transaksi.GetRingkasanCheckout)
			kasirGroup.POST("/transaksi/bayar/tunai", transaksi.BayarTunai)
			kasirGroup.POST("/transaksi/bayar/non-tunai", transaksi.BayarNonTunai)
			kasirGroup.GET("/profil", user.GetProfilKasir)
			kasirGroup.GET("/notifikasi", notifikasi.GetNotifikasi)
		}

		//Kurir Routes
		kurirGroup := v1.Group("/kurir")
		kurirGroup.Use(middleware.AuthMiddleware())
		{
			kurirGroup.GET("/tugas", pengantaran.GetDaftarPengantaran)
			kurirGroup.GET("/laporan", pengantaran.GetLaporanHariIni)
			kurirGroup.GET("/profile", user.GetProfilKurir)
			kurirGroup.GET("/pesanan/new", pemesanan.GetPesananTerbaru)
			kurirGroup.GET("/pesanan", pemesanan.GetAllPesananOnline)
			kurirGroup.GET("/pengantaran/:id_pengantaran/detail", pengantaran.GetDetailPengantaran)
		}

		// Admin Routes
		adminGroup := v1.Group("/admin")
		adminGroup.Use(middleware.AuthMiddleware())
		{
			adminGroup.GET("/dashboard", user.GetDashboardAdmin)
			adminGroup.GET("/dashboard/chart", user.GetChartDashboardAdmin)
			adminGroup.GET("/kategori", katalog.GetKategori)
			adminGroup.POST("/kategori", katalog.TambahKategori)
			adminGroup.PUT("/kategori/:public_id", katalog.UpdateKategori)
			adminGroup.DELETE("/kategori/:public_id", katalog.HapusKategori)
			adminGroup.POST("/kategori/upload", katalog.UploadIconKategori)
			adminGroup.GET("/barang", katalog.GetDaftarBarang)
			adminGroup.POST("/barang", katalog.TambahBarang)
			adminGroup.GET("/barang/detail/:public_id", katalog.GetDetailBarang)
			adminGroup.POST("/barang/upload", katalog.UploadGambarBarang)
			adminGroup.PUT("/barang/:public_id", katalog.UpdateBarang)
			adminGroup.DELETE("/barang/:public_id", katalog.HapusBarang)
			adminGroup.GET("/satuan", katalog.GetSatuan)
			adminGroup.POST("/diskon/upload", katalog.UploadBannerDiskon)
			adminGroup.POST("/diskon", katalog.TambahDiskon)
			adminGroup.DELETE("/diskon/:public_id", katalog.HapusDiskon)
			adminGroup.GET("/diskon/semua", katalog.GetAllDiskon)
			adminGroup.GET("/diskon", katalog.GetPromo)
			adminGroup.GET("/karyawan", user.GetDaftarKaryawan)
			adminGroup.POST("/karyawan", user.TambahKaryawan)
			adminGroup.POST("/karyawan/upload", user.UploadFotoKaryawan)
			adminGroup.GET("/karyawan/:public_id", user.GetDetailKaryawan)
			adminGroup.PUT("/karyawan/:public_id", user.UpdateKaryawan)
			adminGroup.DELETE("/karyawan/:public_id", user.HapusKaryawan)
			adminGroup.GET("/notifikasi", notifikasi.GetNotifikasiAdmin)
			adminGroup.GET("/profil", user.GetProfilAdmin)
			adminGroup.PUT("/profil", user.UpdateProfilAdmin)
		}
	}
}
