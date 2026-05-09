package seeders

import (
	"fmt"
	"backend-mantra/config"
	"backend-mantra/models"
)

func SeedRole() {
	daftarRole := []string{"Customer", "Admin", "Kasir", "Kurir"}

	for _, nama := range daftarRole {
		roleBaru := models.Role{
			NamaRole: nama,
		}

		// FirstOrCreate
		if err := config.DB.FirstOrCreate(&roleBaru, models.Role{NamaRole: nama}).Error; err != nil {
			fmt.Println("Error", nama, "Error:", err)
			return
		}
	}

	fmt.Println("Yeyy, berhasil seed role!")
}