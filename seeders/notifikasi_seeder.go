package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"
)

func SeedNotifikasi() {
	fmt.Println("⏳ Menyiapkan data notifikasi...")

	var count int64
	config.DB.Model(&models.Notifikasi{}).Count(&count)
	if count > 0 {
		fmt.Println("Tabel notifikasi udah ada isinya, proses seeding dilewati.")
		return
	}

	var users []models.User
	if err := config.DB.Preload("Role").Find(&users).Error; err != nil || len(users) == 0 {
		fmt.Println("Gagal: Data User belum ada!")
		return
	}

	for _, user := range users {
		var notifs []models.Notifikasi

		switch user.Role.NamaRole {
		case "Customer":
			notifs = []models.Notifikasi{
				{Judul: "Pesanan Diterima", Pesan: "Pesanan Anda sedang diproses oleh kasir.", Status: "unread"},
				{Judul: "Pesanan Dikirim", Pesan: "Kurir sedang menuju ke alamat Anda.", Status: "unread"},
				{Judul: "Promo Spesial", Pesan: "Dapatkan diskon 50% untuk produk baru!", Status: "unread"},
			}
		case "Kasir":
			notifs = []models.Notifikasi{
				{Judul: "Stok Menipis", Pesan: "Beberapa barang hampir habis, segera lakukan restock.", Status: "unread"},
				{Judul: "Transaksi Baru", Pesan: "Ada pesanan online baru yang perlu dikonfirmasi.", Status: "unread"},
			}
		case "Admin":
			notifs = []models.Notifikasi{
				{Judul: "Laporan Harian", Pesan: "Laporan penjualan hari ini sudah tersedia.", Status: "unread"},
				{Judul: "Peringatan Sistem", Pesan: "Perlu pengecekan stok opname bulan ini.", Status: "unread"},
				{Judul: "Karyawan Baru", Pesan: "Ada pendaftaran kasir baru yang menunggu persetujuan.", Status: "unread"},
			}
		case "Kurir":
			notifs = []models.Notifikasi{
				{Judul: "Tugas Pickup", Pesan: "Ada pesanan baru yang harus diambil di toko.", Status: "unread"},
				{Judul: "Rute Terupdate", Pesan: "Perhatikan rute pengantaran karena ada penutupan jalan.", Status: "unread"},
			}
		}

		for _, notif := range notifs {
			notif.UserID = user.IdUser
			if err := config.DB.Create(&notif).Error; err != nil {
				fmt.Println("Error insert notifikasi:", err)
			}
		}
	}

	fmt.Println("Yeyy, berhasil seed notifikasi!")
}
