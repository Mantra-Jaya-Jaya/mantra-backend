package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"
)

func SeedStatusKaryawan() {
	fmt.Println("Seeding status karyawan...")

	statusList := []string{"Aktif", "Nonaktif", "Probation"}
	for _, nama := range statusList {
		status := models.StatusKaryawan{NamaStatus: nama}
		if err := config.DB.Where("nama_status = ?", nama).FirstOrCreate(&status).Error; err != nil {
			fmt.Printf("Error seeding status karyawan '%s': %v\n", nama, err)
		}
	}
	fmt.Println("Status karyawan selesai di-seed!")
}
