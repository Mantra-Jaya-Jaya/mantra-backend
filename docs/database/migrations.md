# Database Migrations (Atlas CLI)

## Prasyarat

- PostgreSQL running dengan database `mantra_db` dan `mantra_dev`
- Atlas CLI terinstall
- Atlas Provider GORM terinstall
- File `.env` sudah dikonfigurasi

## Konsep

Atlas membandingkan **GORM Structs** (Single Source of Truth) dengan **database sandbox** (`mantra_dev`), lalu menghasilkan SQL untuk menyinkronkan **database utama** (`mantra_db`).

```text
┌──────────────┐    diff     ┌─────────────┐    apply    ┌───────────┐
│  GORM Models │───────────▶│ Sandbox DB  │────────────▶│  Main DB  │
│  (models/)   │            │ (mantra_dev)│             │(mantra_db)│
└──────────────┘            └─────────────┘             └───────────┘
```

## Workflow Harian

### 1. Buat / ubah struct di `models/`
Buat file baru atau edit struct yang sudah ada. Pastikan:
- Nama struct Capital (exported)
- Tag `gorm:"column:nama_kolom"` snake_case
- Method `TableName() string` mengembalikan nama tabel

### 2. Cek perubahan
```bash
make db-diff
```

### 3. Lihat SQL yang akan dijalankan (dry run)
```bash
make db-plan
```

### 4. Apply ke database utama
```bash
make db-apply
```

### 5. Update DBML
Setelah apply, update file `docs/mantra.dbml` agar sinkron dengan kode.

## Perintah Lengkap

| Perintah | Fungsi |
|----------|--------|
| `make db-diff` | Deteksi perubahan model → sandbox |
| `make db-plan` | Lihat SQL sebelum apply |
| `make db-apply` | Apply perubahan ke DB utama |
| `make db-inspect` | Lihat skema DB utama dalam HCL |
| `make db-ui` | Buka ERD visual di browser |
| `make db-clean` | Drop semua tabel di DB utama (⚠️ reset) |
| `make db-clean-dev` | Drop semua tabel di sandbox |
| `make db-clean-all` | Clean DB utama + sandbox |

## Cara Tambah Tabel Baru

1. Buat file `models/baru.go`:
   ```go
   package models

   func (Baru) TableName() string {
       return "baru"
   }

   type Baru struct {
       IdBaru   uint   `gorm:"primaryKey;column:id_baru"`
       NamaBaru string `gorm:"column:nama_baru"`
       // FK ke tabel lain
       UserID   uint   `gorm:"column:id_user"`
       User     User   `gorm:"foreignKey:UserID;references:IdUser"`
   }
   ```

2. Daftarkan di `config/database.go` → `AutoMigrate()`

3. Jalankan:
   ```bash
   make db-diff
   make db-apply
   ```

4. Update `docs/mantra.dbml` dengan tabel baru dan relasinya.

## Troubleshooting

| Masalah | Solusi |
|---------|--------|
| `atlas` command not found | Install Atlas CLI, pastikan ada di PATH |
| `dial tcp ...:5432: connect: connection refused` | Pastikan PostgreSQL running |
| `mantra_dev` database not found | Buat database `mantra_dev` |
| GORM provider error | Jalankan `go mod tidy`, pastikan `atlas-provider-gorm` terinstall |
| `relation sudah ada` | Bersihkan sandbox: `make db-clean-dev` |
