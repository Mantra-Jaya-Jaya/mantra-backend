package routes

import (
	"backend-mantra/controllers"
	"backend-mantra/controllers/auth"
	"backend-mantra/controllers/customer"
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
		v1.GET("/scan/:kode_barcode", customer.GetDetailBarangByScan)

		// Customer Routes
		customerGroup := v1.Group("/customer")
		{
			customerGroup.GET("/promo", customer.GetPromo)
			customerGroup.GET("/kategori", customer.GetKategori)
			customerGroup.GET("/barang", customer.GetDaftarBarang)
			customerGroup.POST("/keranjang", customer.TambahKeKeranjang)
			customerGroup.PATCH("/keranjang/:id_keranjang", customer.UpdateKeranjang)
			customerGroup.DELETE("/keranjang/:id_keranjang", customer.HapusItemKeranjang)
			customerGroup.GET("/notifikasi", customer.GetNotifikasiCustomer)
			customerGroup.GET("/pesanan", customer.GetDaftarPesananCustomer)
			customerGroup.POST("/pesanan/checkout", customer.CheckoutPesanan)
			customerGroup.PATCH("/pesanan/:id_pesanan/batal", customer.BatalkanPesanan)
			customerGroup.GET("/pesanan/:id_pesanan", customer.GetDetailPesananCustomer)
			customerGroup.GET("/pesanan/:id_pesanan/lacak", customer.LacakPesanan)
			customerGroup.GET("/profil", customer.GetProfilCustomer)
			customerGroup.PUT("/akun", customer.UpdateAkunCustomer)
			customerGroup.POST("/alamat", customer.TambahAlamat)
			customerGroup.PUT("/alamat/:id_alamat", customer.UpdateAlamat)
			customerGroup.DELETE("/alamat/:id_alamat", customer.HapusAlamat)
		}

		// Kasir Routes
		kasir := v1.Group("/kasir")
		{
			kasir.GET("/dashboard", controllers.GetDashboardKasir)
			kasir.GET("/laporan", controllers.GetLaporanRingkasan)
			kasir.GET("/laporan/produk/:id_produk", controllers.GetDetailLaporanProduk)
			kasir.GET("/laporan/produk/:id_produk/:id_pesanan", controllers.GetDetailPesananDariLaporan)
			kasir.GET("/pesanan", controllers.GetDaftarPesananKasir)
			kasir.GET("/pesanan/:id_order", controllers.GetDetailPesananKasir)
			kasir.POST("/transaksi/produk", controllers.CariProdukTransaksi)
			kasir.PATCH("/transaksi/item/update", controllers.UpdateQuantityItem)
			kasir.GET("/transaksi/checkout", controllers.GetRingkasanCheckout)
			kasir.POST("/transaksi/bayar/tunai", controllers.BayarTunai)
			kasir.POST("/transaksi/bayar/non-tunai", controllers.BayarNonTunai)
			kasir.GET("/profil", controllers.GetProfilKasir)
			kasir.GET("/notifikasi", controllers.GetNotifikasiKasir)
		}

		// Admin Routes
		admin := v1.Group("/admin")
		{
			admin.GET("/dashboard", controllers.GetDashboardAdmin)
			admin.GET("/barang", controllers.GetDaftarBarangAdmin)
			admin.POST("/barang", controllers.TambahBarang)
			admin.GET("/barang/:id_barang", controllers.GetDetailBarangAdmin)
			admin.PUT("/barang/:id_barang", controllers.UpdateBarang)
			admin.DELETE("/barang/:id_barang", controllers.HapusBarang)
			admin.POST("/barang/:id_barang/diskon", controllers.TambahDiskonBarang)
			admin.GET("/karyawan", controllers.GetDaftarKaryawan)
			admin.GET("/karyawan/kasir", controllers.GetDaftarKasir)
			admin.POST("/karyawan/kasir", controllers.TambahKasir)
			admin.GET("/karyawan/kasir/:id_kasir", controllers.GetDetailKasir)
			admin.PUT("/karyawan/kasir/:id_kasir", controllers.UpdateKasir)
			admin.DELETE("/karyawan/kasir/:id_kasir", controllers.HapusKasir)
			admin.GET("/karyawan/kurir", controllers.GetDaftarKurir)
			admin.POST("/karyawan/kurir", controllers.TambahKurir)
			admin.GET("/karyawan/kurir/:id_kurir", controllers.GetDetailKurir)
			admin.PUT("/karyawan/kurir/:id_kurir", controllers.UpdateKurir)
			admin.DELETE("/karyawan/kurir/:id_kurir", controllers.HapusKurir)
			admin.GET("/notifikasi", controllers.GetNotifikasiAdmin)
			admin.GET("/profil", controllers.GetProfilAdmin)
			admin.PUT("/profil", controllers.UpdateProfilAdmin)
		}
	}
}
