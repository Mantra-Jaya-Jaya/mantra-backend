package seeders

import (
	"backend-mantra/config"
	"backend-mantra/models"
	"fmt"
)

func SeedMetodePembayaran() {
	fmt.Println("⏳ Menyiapkan data metode pembayaran...")

	daftarMetode := []models.MetodePembayaran{
		{NamaMetode: "Cash", KodeMetode: "cash", Penyedia: "internal", Icon: "payments_outlined", Urutan: 1, IsActive: true},
		{NamaMetode: "QRIS", KodeMetode: "qris", Penyedia: "midtrans", Icon: "https://xendit.co/wp-content/uploads/2020/03/iconQris.png", Urutan: 2, IsActive: true},
		{NamaMetode: "COD", KodeMetode: "cod", Penyedia: "internal", Icon: "https://xendit.co/wp-content/uploads/2020/03/iconCod.png", Urutan: 3, IsActive: true},

		// Kodenya sama-sama "va", tapi NamaMetode dan Icon-nya berbeda-beda!
		{NamaMetode: "BNI Virtual Account", KodeMetode: "va", Penyedia: "midtrans", Icon: "https://upload.wikimedia.org/wikipedia/commons/thumb/f/f0/Bank_Negara_Indonesia_logo_%282004%29.svg/3840px-Bank_Negara_Indonesia_logo_%282004%29.svg.png", Urutan: 4, IsActive: true},
		{NamaMetode: "BCA Virtual Account", KodeMetode: "va", Penyedia: "midtrans", Icon: "https://www.bca.co.id/-/media/Feature/Card/List-Card/Tentang-BCA/Brand-Assets/Logo-BCA/Logo-BCA_Biru.png", Urutan: 5, IsActive: true},
		{NamaMetode: "BRI Virtual Account", KodeMetode: "va", Penyedia: "midtrans", Icon: "https://upload.wikimedia.org/wikipedia/commons/thumb/9/97/Logo_BRI.png/1920px-Logo_BRI.png", Urutan: 7, IsActive: true},
	}

	for _, m := range daftarMetode {
		// Cari berdasarkan NamaMetode agar data antar bank tidak saling menimpa
		err := config.DB.Where("nama_metode = ?", m.NamaMetode).FirstOrCreate(&m).Error
		if err != nil {
			fmt.Println("❌ Error insert metode", m.NamaMetode, ":", err)
			continue
		}

		// Update icon dan kode_metode terbaru ke database
		config.DB.Model(&m).Updates(map[string]interface{}{
			"kode_metode": m.KodeMetode,
			"icon":        m.Icon,
			"urutan":      m.Urutan,
			"is_active":   m.IsActive,
		})
	}

	fmt.Println("✅ Berhasil melakukan seeder metode pembayaran!")
}
