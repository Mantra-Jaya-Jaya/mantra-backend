package seeders

import (
	"fmt"
	"backend-mantra/config"
	"backend-mantra/models"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

func SeedUser() {
	// Buat list data user
	usersData := []struct {
		Username     string
		Email        string
		NamaLengkap  string
		NamaRole     string 
		PasswordAsli string 
	}{
		{"terra_admin", "admin@mantra.com", "Terra Surya", "Admin", "AdminMantra#1"},
		{"nabila_kasir", "kasir@mantra.com", "Nabila Az Zahra", "Kasir", "Kasir#123"},
		{"riztika_kurir", "kurir@mantra.com", "Riztika Merizta", "Kurir", "Ngebut#123"},
		{"hamim_customer", "customer@mantra.com", "Rajaba Hamim", "Customer", "Customer#123"},
	}

	for _, data := range usersData {
		var role models.Role
		if err := config.DB.Where("nama_role = ?", data.NamaRole).First(&role).Error; err != nil {
			fmt.Println("Hamdehh, Role", data.NamaRole, "gak ketemu!")
			continue 
		}

		// Hash password aslinya SATU PER SATU pas lagi di-loop
		hashedPassword := hashPassword(data.PasswordAsli)

		userBaru := models.User{
			Username:    data.Username,
			Email:       data.Email,
			Password:    hashedPassword, // Yang disimpen tetep yang acak (Hashed)
			NamaLengkap: data.NamaLengkap,
			RoleID:      role.IdRole, 
		}

		// Cek apakah email udah ada biar gak dobel
		if err := config.DB.Where("email = ?", data.Email).FirstOrCreate(&userBaru).Error; err != nil {
			fmt.Println("Error", data.Username, "Error:", err)
			return
		}
	}

	fmt.Println("Yeyy, berhasil seed user!")
}