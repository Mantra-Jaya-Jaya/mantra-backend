package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"
	"net/url"

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

		// Tentukan warna latar belakang foto profil berdasarkan role
		bgColor := "fff" // default
		switch data.NamaRole {
		case "Admin":
			bgColor = "6366f1"
		case "Kasir":
			bgColor = "10b981"
		case "Kurir":
			bgColor = "f59e0b"
		case "Customer":
			bgColor = "3b82f6"
		}
		
		// Buat URL Foto Profil otomatis
		fotoProfil := fmt.Sprintf("https://ui-avatars.com/api/?name=%s&background=%s&color=fff", url.QueryEscape(data.NamaLengkap), bgColor)

		userBaru := models.User{
			Username:    data.Username,
			Email:       data.Email,
			Password:    hashedPassword, // Yang disimpen tetep yang acak (Hashed)
			NamaLengkap: data.NamaLengkap,
			FotoProfil:  fotoProfil,
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
