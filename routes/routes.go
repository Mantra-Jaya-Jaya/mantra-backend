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

		// Public Webhook (tanpa auth)
		v1.POST("/payment/notification", transaksi.MidtransNotificationHandler)
		v1.POST("/webhook/biteship", transaksi.BiteshipWebhookHandler)

		// Customer Routes
		customerGroup := v1.Group("/customer")
		customerGroup.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("Customer"))
		{
			customerGroup.GET("/promo", katalog.GetPromo)
			customerGroup.GET("/promo/:public_id/barang", katalog.GetBarangByDiskon)
			customerGroup.GET("/kategori", katalog.GetKategori)
			customerGroup.GET("/barang", katalog.GetDaftarBarang)
			customerGroup.GET("/barang/detail/:public_id", katalog.GetDetailBarang)
			customerGroup.GET("/keranjang", keranjang.GetKeranjang)
			customerGroup.POST("/keranjang", keranjang.TambahKeKeranjang)
			customerGroup.PATCH("/keranjang/:public_id", keranjang.UpdateKeranjang)
			customerGroup.DELETE("/keranjang/:public_id", keranjang.HapusItemKeranjang)
			customerGroup.GET("/notifikasi", notifikasi.GetNotifikasi)
			customerGroup.PATCH("/notifikasi/:id/baca", notifikasi.BacaNotifikasi)
			customerGroup.GET("/pesanan", transaksi.GetDaftarPesanan)
			customerGroup.POST("/pesanan/checkout", transaksi.CheckoutPesanan)
			customerGroup.PATCH("/pesanan/:public_id/batal", transaksi.BatalkanPesanan)
			customerGroup.GET("/pesanan/:public_id", transaksi.GetDetailPesanan)
			customerGroup.GET("/pesanan/:public_id/lacak", transaksi.LacakPesanan)
			customerGroup.GET("/pesanan/:public_id/status-biteship", transaksi.GetBiteshipOrderStatus)
			customerGroup.POST("/ongkir/cek", transaksi.CekOngkir)
			customerGroup.POST("/ongkir/cek-radius", transaksi.CekRadius)
			customerGroup.GET("/metode-pembayaran", transaksi.GetMetodePembayaranAktif)
			customerGroup.GET("/profil", user.GetProfilCustomer)
			customerGroup.PUT("/akun", user.EditAkunCustomer)
			customerGroup.GET("/alamat", user.GetAlamat)
			customerGroup.POST("/alamat", user.TambahAlamat)
			customerGroup.PUT("/alamat/:public_id", user.UpdateAlamat)
			customerGroup.DELETE("/alamat/:public_id", user.HapusAlamat)
		}

		// Kasir Routes
		kasirGroup := v1.Group("/kasir")
		kasirGroup.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("Kasir"))
		{
			kasirGroup.GET("/dashboard", transaksi.GetDashboardKasir)
			kasirGroup.GET("/aktivitas-hari-ini", transaksi.GetSemuaAktivitasHariIni)
			kasirGroup.GET("/laporan", transaksi.GetLaporanRingkasan)
			kasirGroup.GET("/laporan/produk/:public_id", transaksi.GetDetailLaporanProduk)
			kasirGroup.GET("/laporan/produk/:public_id/:pesanan_id", transaksi.GetDetailPesananDariLaporan)
			kasirGroup.GET("/pesanan", transaksi.GetDaftarPesanan)
			kasirGroup.GET("/pesanan/:public_id", transaksi.GetDetailPesanan)
			kasirGroup.GET("/kategori", katalog.GetKategori)
			kasirGroup.GET("/transaksi/produk", katalog.CariProdukTransaksi)
			kasirGroup.PATCH("/transaksi/item/update", transaksi.UpdateQuantityItem)
			kasirGroup.GET("/transaksi/checkout", transaksi.GetRingkasanCheckout)
			kasirGroup.GET("/metode-pembayaran", transaksi.GetMetodePembayaranAktif)
			kasirGroup.POST("/transaksi/bayar/tunai", transaksi.BayarTunai)
			kasirGroup.POST("/transaksi/bayar/non-tunai", transaksi.BayarNonTunai)
			kasirGroup.GET("/transaksi/cek-status/:order_id", transaksi.CekStatusPembayaran)
			kasirGroup.GET("/profil", user.GetProfilKasir)
			kasirGroup.PUT("/profil", user.UpdateProfilKasir)
			kasirGroup.PATCH("/pesanan/:public_id/kirim", transaksi.KirimPesanan)
			kasirGroup.GET("/notifikasi", notifikasi.GetNotifikasi)
		}

		//Kurir Routes
		kurirGroup := v1.Group("/kurir")
		kurirGroup.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("Kurir"))
		{
			kurirGroup.GET("/tugas", pengantaran.GetDaftarPengantaran)
			kurirGroup.GET("/laporan", pengantaran.GetLaporanHariIni)
			kurirGroup.GET("/profile", user.GetProfilKurir)
			kurirGroup.GET("/pesanan/new", pemesanan.GetPesananTerbaru)
			kurirGroup.GET("/pesanan", pemesanan.GetAllPesananOnline)
			kurirGroup.GET("/pesanan/:public_id", pemesanan.GetDetailPesanan)
			kurirGroup.POST("/pesanan/:public_id/terima", pemesanan.TerimaPesanan)
			kurirGroup.GET("/pengantaran/:public_id/detail", pengantaran.GetDetailPengantaran)
			kurirGroup.PUT("/pengantaran/:public_id/lokasi", pengantaran.UpdateLokasiKurir)
			kurirGroup.POST("/pengantaran/:public_id/ambil", pengantaran.AmbilPesanan)
			kurirGroup.POST("/pengantaran/:public_id/status", pengantaran.UpdateStatusPengantaran)
			kurirGroup.PUT("/pengantaran/:public_id/selesai", pengantaran.UploadBuktiPengiriman)
			kurirGroup.PUT("/pengantaran/:public_id/pembayaran", pengantaran.KonfirmasiPembayaran)
			kurirGroup.GET("/notifikasi", notifikasi.GetNotifikasi)
		}

		// Admin Routes
		adminGroup := v1.Group("/admin")
		adminGroup.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("Admin"))
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
			adminGroup.PATCH("/notifikasi/:id/baca", notifikasi.BacaNotifikasi)
			adminGroup.DELETE("/notifikasi/:id", notifikasi.HapusNotifikasi)
			adminGroup.GET("/profil", user.GetProfilAdmin)
			adminGroup.PUT("/profil", user.UpdateProfilAdmin)
			adminGroup.GET("/ekspedisi", katalog.GetDaftarEkspedisi)
			adminGroup.POST("/ekspedisi", katalog.TambahEkspedisi)
			adminGroup.PUT("/ekspedisi/:public_id", katalog.UpdateEkspedisi)
			adminGroup.DELETE("/ekspedisi/:public_id", katalog.HapusEkspedisi)
			adminGroup.POST("/ekspedisi/sync", katalog.SyncBiteshipCouriers)
			adminGroup.POST("/ekspedisi/layanan", katalog.TambahLayanan)
			adminGroup.PUT("/ekspedisi/layanan/:public_id", katalog.UpdateLayanan)
			adminGroup.DELETE("/ekspedisi/layanan/:public_id", katalog.HapusLayanan)
			adminGroup.GET("/metode-pembayaran", transaksi.GetMetodePembayaran)
			adminGroup.POST("/metode-pembayaran", transaksi.TambahMetodePembayaran)
			adminGroup.PUT("/metode-pembayaran/:public_id", transaksi.UpdateMetodePembayaran)
			adminGroup.DELETE("/metode-pembayaran/:public_id", transaksi.HapusMetodePembayaran)
			adminGroup.GET("/pengaturan", user.GetPengaturan)
			adminGroup.PUT("/pengaturan", user.UpdatePengaturan)
			adminGroup.GET("/pengantaran", transaksi.GetDaftarPengantaranAdmin)
			adminGroup.POST("/pesanan/:public_id/cancel-shipment", transaksi.CancelBiteshipOrder)
			adminGroup.GET("/pesanan/:public_id/status-biteship", transaksi.GetBiteshipOrderStatus)
		}
	}
}
