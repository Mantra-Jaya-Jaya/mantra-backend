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
- [ ] Setup Midtrans dengan server key production
- [ ] Konfigurasi CORS `ALLOWED_ORIGIN` sesuai domain admin web
- [ ] Nonaktifkan `gin.Default()` recovery atau ganti dengan custom recovery
- [ ] Setup reverse proxy (Nginx/Caddy) untuk TLS/SSL
- [ ] Jalankan dengan process manager (systemd/supervisor/pm2)
- [ ] Backup database rutin
- [ ] Monitoring: health check endpoint (belum diimplement)
