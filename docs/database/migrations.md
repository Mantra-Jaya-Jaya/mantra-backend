# 🚀 Migrasi Database (Atlas CLI & GORM)

---
### 🧭 Navigasi Cepat
[🏠 Utama](../README.md) | [🏛️ Arsitektur](../architecture.md) | [🛠️ Deployment](../deployment.md) | [💳 Midtrans](../pembayaran.md) | [📦 Biteship](../biteship.md) | [📡 API Contract](../api/overview.md) | [🗄️ Database](erd.md) | [🔒 Keamanan](../security/README.md)
---

Dokumen ini menjelaskan alur kerja deklaratif untuk mengelola migrasi skema database PostgreSQL menggunakan **Atlas CLI** dan **GORM ORM**.

---

## 📋 1. Prasyarat System

Sebelum menjalankan perintah migrasi, pastikan prasyarat berikut sudah terpenuhi di sistem lokal Anda:
*   **PostgreSQL:** Sudah berjalan di localhost.
*   **Database:** Memiliki dua database terpisah di PostgreSQL:
    1.  `mantra_db` (Database utama untuk runtime aplikasi).
    2.  `mantra_dev` (Database sandbox/temporary untuk kalkulasi diff).
*   **Atlas CLI:** Sudah terinstall di komputer Anda.
*   **Go Dependencies:** `atlas-provider-gorm` terpasang di modul Go.
*   **Env File:** File `.env` sudah terisi dengan benar.

---

## 💡 2. Konsep Dasar Migrasi Deklaratif

Mantra menggunakan pendekatan migrasi deklaratif. Alih-alih menulis file migrasi SQL secara manual (seperti `UP` atau `DOWN`), Atlas akan membandingkan skema model struct Go Anda (**GORM Models** sebagai *Single Source of Truth*) dengan database sandbox (`mantra_dev`), kemudian menghasilkan instruksi SQL migrasi secara dinamis untuk memperbarui database utama (`mantra_db`).

```text
GORM Models (models/)
        │
        ▼ (perbandingan otomatis / diff)
Sandbox DB (mantra_dev)
        │
        ▼ (aplikasikan perubahan / apply)
Database Utama (mantra_db)
```

---

## 🔄 3. Alur Kerja Harian (Daily Workflow)

Ikuti 5 langkah berikut ketika Anda perlu menambah tabel atau mengubah kolom database:

### Langkah 1: Buat atau Ubah Struct di Folder `models/`
Buat file Go baru atau modifikasi struct yang sudah ada. Pastikan konvensi GORM terpenuhi:
```go
package models

type Baru struct {
    IdBaru   uint   `gorm:"primaryKey;column:id_baru"`
    NamaBaru string `gorm:"column:nama_baru"`
}

func (Baru) TableName() string {
    return "baru"
}
```

### Langkah 2: Daftarkan Model Baru
Buka file `config/database.go` dan tambahkan referensi pointer struct baru Anda ke dalam fungsi pembungkus GORM `AutoMigrate()`.

### Langkah 3: Deteksi Perubahan Skema
Jalankan perintah ini untuk memicu Atlas mendeteksi perbedaan kode model vs sandbox database:
```bash
make db-diff
```

### Langkah 4: Tinjau SQL (Dry Run)
Tinjau perintah SQL apa saja yang akan dieksekusi oleh Atlas untuk memperbarui database utama:
```bash
make db-plan
```

### Langkah 5: Terapkan Perubahan (Apply)
Terapkan perubahan skema database secara langsung ke database utama Anda:
```bash
make db-apply
```

---

## 💻 4. Daftar Perintah Migrasi (Makefile Commands)

Gunakan perintah `make` berikut di terminal direktori root backend:

| Perintah | Fungsi / Kegunaan |
| :--- | :--- |
| **`make db-diff`** | Deteksi perbedaan skema models vs database sandbox |
| **`make db-plan`** | Melihat rancangan query SQL migrasi yang akan berjalan (*dry-run*) |
| **`make db-apply`** | Menerapkan (apply) perubahan skema ke database utama |
| **`make db-inspect`** | Menghasilkan skema database saat ini dalam format representasi HCL |
| **`make db-ui`** | Membuka antarmuka ERD visual interaktif di browser lokal |
| **`make db-clean`** | Menghapus (drop) semua tabel di database utama **(⚠️ Data Hilang!)** |
| **`make db-clean-dev`** | Menghapus (drop) semua tabel di database sandbox/dev |
| **`make db-clean-all`** | Menghapus semua tabel baik di database utama maupun sandbox |

---

## 🐛 Troubleshooting

*   **Error: `atlas` command not found**
    *   *Solusi:* Install Atlas CLI melalui dokumentasi resmi [Atlasgo.io](https://atlasgo.io/) dan masukkan direktori instalasinya ke PATH OS Anda.
*   **Error: `dial tcp ...:5432: connect: connection refused`**
    *   *Solusi:* Pastikan service database PostgreSQL Anda sudah aktif dan port `5432` dapat diakses.
*   **Error: `database "mantra_dev" does not exist`**
    *   *Solusi:* Buat database kosong bernama `mantra_dev` di server PostgreSQL lokal Anda menggunakan pgAdmin, DBeaver, atau SQL Shell:
        ```sql
        CREATE DATABASE mantra_dev;
        ```
*   **Error: `relation "xxx" already exists`**
    *   *Solusi:* Terjadi ketidaksinkronan di sandbox. Bersihkan database sandbox dengan menjalankan perintah:
        ```bash
        make db-clean-dev
        ```
