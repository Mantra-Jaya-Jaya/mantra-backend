package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"
	"net/url"

	"golang.org/x/crypto/bcrypt"
)

const BCRYPT_COST = 12

func hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), BCRYPT_COST)
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
		// 1 Admin
		{"admin_mantra", "admin@mantra.test", "Terra Admin", "Admin", "Admin123!"},
		// 15 Customers
		{"customer_01", "customer01@mantra.test", "Alya Putri", "Customer", "Customer123!"},
		{"customer_02", "customer02@mantra.test", "Budi Santoso", "Customer", "Customer123!"},
		{"customer_03", "customer03@mantra.test", "Citra Dewi", "Customer", "Customer123!"},
		{"customer_04", "customer04@mantra.test", "Dinda Lestari", "Customer", "Customer123!"},
		{"customer_05", "customer05@mantra.test", "Eko Prasetyo", "Customer", "Customer123!"},
		{"customer_06", "customer06@mantra.test", "Farah Amalia", "Customer", "Customer123!"},
		{"customer_07", "customer07@mantra.test", "Gilang Ramadhan", "Customer", "Customer123!"},
		{"customer_08", "customer08@mantra.test", "Hana Safira", "Customer", "Customer123!"},
		{"customer_09", "customer09@mantra.test", "Irfan Hakim", "Customer", "Customer123!"},
		{"customer_10", "customer10@mantra.test", "Jihan Nabila", "Customer", "Customer123!"},
		{"customer_11", "customer11@mantra.test", "Kevin Wijaya", "Customer", "Customer123!"},
		{"customer_12", "customer12@mantra.test", "Laras Ayu", "Customer", "Customer123!"},
		{"customer_13", "customer13@mantra.test", "Miko Ardiansyah", "Customer", "Customer123!"},
		{"customer_14", "customer14@mantra.test", "Nindi Putri", "Customer", "Customer123!"},
		{"customer_15", "customer15@mantra.test", "Oka Sanjaya", "Customer", "Customer123!"},
		// 5 Kasirs
		{"kasir_01", "kasir01@mantra.test", "Eka Saputra", "Kasir", "Kasir123!"},
		{"kasir_02", "kasir02@mantra.test", "Fajar Nugraha", "Kasir", "Kasir123!"},
		{"kasir_03", "kasir03@mantra.test", "Gita Permata", "Kasir", "Kasir123!"},
		{"kasir_04", "kasir04@mantra.test", "Hesti Wulandari", "Kasir", "Kasir123!"},
		{"kasir_05", "kasir05@mantra.test", "Irwan Saputra", "Kasir", "Kasir123!"},
		// 3 Kurirs
		{"kurir_01", "kurir01@mantra.test", "Hadi Kurnia", "Kurir", "Kurir123!"},
		{"kurir_02", "kurir02@mantra.test", "Intan Maharani", "Kurir", "Kurir123!"},
		{"kurir_03", "kurir03@mantra.test", "Joko Susilo", "Kurir", "Kurir123!"},
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

	fmt.Println("Yeyy, berhasil seed 24 user accounts!")
}
