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
		{NamaMetode: "QRIS", KodeMetode: "qris", Penyedia: "midtrans", Icon: "qr_code_scanner", Urutan: 2, IsActive: true},
		{NamaMetode: "COD", KodeMetode: "cod", Penyedia: "internal", Icon: "package_outlined", Urutan: 3, IsActive: true},

		// Kodenya sama-sama "va", tapi NamaMetode dan Icon-nya berbeda-beda!
		{NamaMetode: "BNI Virtual Account", KodeMetode: "va", Penyedia: "midtrans", Icon: "https://upload.wikimedia.org/wikipedia/commons/thumb/f/f0/Bank_Negara_Indonesia_logo_%282004%29.svg/3840px-Bank_Negara_Indonesia_logo_%282004%29.svg.png", Urutan: 4, IsActive: true},
		{NamaMetode: "BCA Virtual Account", KodeMetode: "va", Penyedia: "midtrans", Icon: "https://www.bca.co.id/-/media/Feature/Card/List-Card/Tentang-BCA/Brand-Assets/Logo-BCA/Logo-BCA_Biru.png", Urutan: 5, IsActive: true},
		{NamaMetode: "BRI Virtual Account", KodeMetode: "va", Penyedia: "midtrans", Icon: "https://upload.wikimedia.org/wikipedia/commons/thumb/9/97/Logo_BRI.png/1920px-Logo_BRI.png", Urutan: 7, IsActive: true},
	}

	for _, m := range daftarMetode {
		// Assign memaksa semua field di-update jika record sudah ada
		err := config.DB.Where("nama_metode = ?", m.NamaMetode).Assign(&m).FirstOrCreate(&m).Error
		if err != nil {
			fmt.Println("❌ Error insert metode", m.NamaMetode, ":", err)
			continue
		}
	}

	fmt.Println("✅ Berhasil melakukan seeder metode pembayaran!")
}
