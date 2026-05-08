package routes

import (
	"net/http"
	"backend-mantra/controllers"
)

func SetupRoutes() {
	// Daftarin endpoint di sini
	http.HandleFunc("/customer/kategori", controllers.GetKategori)
	http.HandleFunc("/admin/kategori", controllers.CreateKategori)
	
	// Nanti endpoint lain nyusul di bawahnya
	// http.HandleFunc("/customer/barang", controllers.GetBarang)
}