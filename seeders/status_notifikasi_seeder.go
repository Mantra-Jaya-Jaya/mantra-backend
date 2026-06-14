package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"
)

func SeedStatusNotifikasi() {
	fmt.Println("Seeding status notifikasi...")

	statusList := []string{"unread", "read", "aktif"}
	for _, nama := range statusList {
		status := models.StatusNotifikasi{NamaStatus: nama}
		if err := config.DB.Where("nama_status = ?", nama).FirstOrCreate(&status).Error; err != nil {
			fmt.Printf("Error seeding status notifikasi '%s': %v\n", nama, err)
		}
	}
	fmt.Println("Status notifikasi selesai di-seed!")
}
