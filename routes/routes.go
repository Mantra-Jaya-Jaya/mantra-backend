package routes

import (
	"backend-mantra/controllers"
	// "net/http" hapus // kalau mau pakai net/http
	"github.com/gin-gonic/gin"
)

func SetupRoutes() {
	// Daftarin endpoint di sini
	http.HandleFunc("/customer/kategori", controllers.GetKategori)
	http.HandleFunc("/admin/kategori", controllers.CreateKategori)
	http.HandleFunc("/customer/diskon", controllers.GetPromoCustomer)
	
	// Nanti endpoint lain nyusul di bawahnya
	// http.HandleFunc("/customer/barang", controllers.GetBarang)
}
// SetupRoutes registers all API routes
func SetupRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		// Auth Routes
		v1.POST("/login", controllers.Login)
		v1.POST("/register", controllers.RegisterCustomer)
		v1.POST("/auth/refresh", controllers.RefreshToken)
		v1.POST("/logout", controllers.Logout)
		v1.PUT("/change-password", controllers.ChangePassword)

		// Shared Routes
		v1.GET("/scan/:kode_barcode", controllers.GetDetailBarangByScan)

		// Customer Routes
		customer := v1.Group("/customer")
		{
			customer.GET("/promo", controllers.GetPromo)
			customer.GET("/kategori", controllers.GetKategori)
			customer.GET("/barang", controllers.GetDaftarBarang)
			customer.POST("/keranjang", controllers.TambahKeKeranjang)
			customer.PATCH("/keranjang/:id_keranjang", controllers.UpdateKeranjang)
			customer.DELETE("/keranjang/:id_keranjang", controllers.HapusItemKeranjang)
			customer.GET("/notifikasi", controllers.GetNotifikasiCustomer)
			customer.GET("/pesanan", controllers.GetDaftarPesananCustomer)
			customer.POST("/pesanan/checkout", controllers.CheckoutPesanan)
			customer.PATCH("/pesanan/:id_pesanan/batal", controllers.BatalkanPesanan)
			customer.GET("/pesanan/:id_pesanan", controllers.GetDetailPesananCustomer)
			customer.GET("/pesanan/:id_pesanan/lacak", controllers.LacakPesanan)
			customer.GET("/profil", controllers.GetProfilCustomer)
			customer.PUT("/akun", controllers.UpdateAkunCustomer)
			customer.POST("/alamat", controllers.TambahAlamat)
			customer.PUT("/alamat/:id_alamat", controllers.UpdateAlamat)
			customer.DELETE("/alamat/:id_alamat", controllers.HapusAlamat)
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
