package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"
)

func SeedFraudStatus() {
	fmt.Println("Seeding fraud status...")

	statusList := []string{"accept", "challenge", "deny"}
	for _, nama := range statusList {
		status := models.FraudStatus{NamaStatus: nama}
		if err := config.DB.Where("nama_status = ?", nama).FirstOrCreate(&status).Error; err != nil {
			fmt.Printf("Error seeding fraud status '%s': %v\n", nama, err)
		}
	}
	fmt.Println("Fraud status selesai di-seed!")
}
