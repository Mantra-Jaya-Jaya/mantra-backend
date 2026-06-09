# Deployment

## Build Binary

```bash
make build
```

Hasil: `bin/mantra-backend` — binary statis, siap di-copy ke server.

## Environment Variables

Semua konfigurasi via environment variables (file `.env` di root proyek):

### Database
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=mantra_db
```

### JWT
```env
JWT_SECRET=generate_random_secret_here
JWT_REFRESH_SECRET=generate_different_secret_here
```

### Midtrans (Payment Gateway)
```env
MIDTRANS_SERVER_KEY=your_server_key
MIDTRANS_CLIENT_KEY=your_client_key
MIDTRANS_MERCHANT_ID=your_merchant_id
```

### Biteship (Shipping)
```env
BITESHIP_API_KEY=your_biteship_api_key
BITESHIP_MODE=sandbox
BITESHIP_STORE_NAME=Toko Anda
BITESHIP_STORE_ADDRESS=Alamat toko
BITESHIP_STORE_CITY=Kota
BITESHIP_STORE_POSTAL_CODE=kode_pos
BITESHIP_STORE_REGION=Provinsi
BITESHIP_STORE_COORDINATE_LAT=-6.9175
BITESHIP_STORE_COORDINATE_LONG=107.6191
```

### MinIO (File Storage)
```env
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=your_access_key
MINIO_SECRET_KEY=your_secret_key
MINIO_BUCKET=mantra-storage
```

### CORS
```env
ALLOWED_ORIGIN=https://admin.mantra.web.id
```

## Production Checklist

- [ ] Ganti `JWT_SECRET` dengan secret kuat (min 32 karakter)
- [ ] Ganti `JWT_REFRESH_SECRET` dengan secret berbeda
- [ ] Setup MinIO dengan bucket `mantra-storage` dan akses publik untuk folder upload
- [ ] Setup Midtrans dengan server key production, merchant ID
- [ ] Setup Biteship dengan API key production, koordinat toko, alamat toko
- [ ] Konfigurasi CORS `ALLOWED_ORIGIN` sesuai domain admin web
- [ ] Nonaktifkan `gin.Default()` recovery atau ganti dengan custom recovery
- [ ] Setup reverse proxy (Nginx/Caddy) untuk TLS/SSL
- [ ] Jalankan dengan process manager (systemd/supervisor/pm2)
- [ ] Backup database rutin
- [ ] Monitoring: health check endpoint (belum diimplement)
