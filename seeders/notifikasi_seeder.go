package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"backend-mantra/utils"
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
				{Judul: "Pesanan Diterima", Pesan: "Pesanan Anda sedang diproses oleh kasir.", StatusNotifikasiID: utils.GetStatusNotifikasiID("unread")},
				{Judul: "Pesanan Dikirim", Pesan: "Kurir sedang menuju ke alamat Anda.", StatusNotifikasiID: utils.GetStatusNotifikasiID("unread")},
				{Judul: "Promo Spesial", Pesan: "Dapatkan diskon 50% untuk produk baru!", StatusNotifikasiID: utils.GetStatusNotifikasiID("unread")},
			}
		case "Kasir":
			notifs = []models.Notifikasi{
				{Judul: "Stok Menipis", Pesan: "Beberapa barang hampir habis, segera lakukan restock.", StatusNotifikasiID: utils.GetStatusNotifikasiID("unread")},
				{Judul: "Transaksi Baru", Pesan: "Ada pesanan online baru yang perlu dikonfirmasi.", StatusNotifikasiID: utils.GetStatusNotifikasiID("unread")},
			}
		case "Admin":
			notifs = []models.Notifikasi{
				{Judul: "Laporan Harian", Pesan: "Laporan penjualan hari ini sudah tersedia.", StatusNotifikasiID: utils.GetStatusNotifikasiID("unread")},
				{Judul: "Peringatan Sistem", Pesan: "Perlu pengecekan stok opname bulan ini.", StatusNotifikasiID: utils.GetStatusNotifikasiID("unread")},
				{Judul: "Karyawan Baru", Pesan: "Ada pendaftaran kasir baru yang menunggu persetujuan.", StatusNotifikasiID: utils.GetStatusNotifikasiID("unread")},
			}
		case "Kurir":
			notifs = []models.Notifikasi{
				{Judul: "Tugas Pickup", Pesan: "Ada pesanan baru yang harus diambil di toko.", StatusNotifikasiID: utils.GetStatusNotifikasiID("unread")},
				{Judul: "Rute Terupdate", Pesan: "Perhatikan rute pengantaran karena ada penutupan jalan.", StatusNotifikasiID: utils.GetStatusNotifikasiID("unread")},
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
