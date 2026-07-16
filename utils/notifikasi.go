package utils

import (
	"fmt"
	"backend-mantra/config"
	"backend-mantra/models"
)

// BuatNotifikasi adalah helper untuk membuat notifikasi ke satu user
func BuatNotifikasi(userID uint, judul, pesan string) error {
	idStatusUnread := GetStatusNotifikasiIDSafe("unread")
	if idStatusUnread == 0 {
		idStatusUnread = 1 // Fallback
	}

	notif := models.Notifikasi{
		UserID:             userID,
		Judul:              judul,
		Pesan:              pesan,
		StatusNotifikasiID: idStatusUnread,
	}

	if err := config.DB.Create(&notif).Error; err != nil {
		fmt.Println("Error Create Notifikasi (Single):", err)
		return err
	}
	return nil
}

// BuatNotifikasiRole adalah helper untuk menyebarkan notifikasi ke semua user dalam satu role
func BuatNotifikasiRole(namaRole string, judul, pesan string) error {
	var role models.Role
	if err := config.DB.Where("nama_role = ?", namaRole).First(&role).Error; err != nil {
		return err
	}

	var users []models.User
	if err := config.DB.Where("id_role = ?", role.IdRole).Find(&users).Error; err != nil {
		return err
	}

	idStatusUnread := GetStatusNotifikasiIDSafe("unread")
	if idStatusUnread == 0 {
		idStatusUnread = 1
	}

	var notifs []models.Notifikasi
	for _, u := range users {
		notifs = append(notifs, models.Notifikasi{
			UserID:             u.IdUser,
			Judul:              judul,
			Pesan:              pesan,
			StatusNotifikasiID: idStatusUnread,
		})
	}

	if len(notifs) > 0 {
		if err := config.DB.Create(&notifs).Error; err != nil {
			fmt.Println("Error Create Notifikasi Role:", err)
			return err
		}
		fmt.Printf("Berhasil buat %d notifikasi untuk role %s\n", len(notifs), namaRole)
	} else {
		fmt.Println("Tidak ada user dengan role", namaRole)
	}
	return nil
}
