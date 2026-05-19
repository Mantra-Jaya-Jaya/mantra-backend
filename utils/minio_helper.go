package utils

import (
	"fmt"
	"path/filepath"
	"time"

	"backend-mantra/config" 

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

// UploadFileToMinio adalah fungsi sakti buat upload apapun ke MinIO
// fileKey: nama field di form-data (misal: "gambar" atau "foto_profil")
// folderTarget: nama folder di dalam bucket (misal: "produk", "profile", "kategori")
func UploadFileToMinio(c *gin.Context, fileKey string, folderTarget string) (string, error) {
	// 1. Tangkap file berdasarkan key-nya
	file, header, err := c.Request.FormFile(fileKey)
	if err != nil {
		return "", fmt.Errorf("file tidak ditemukan: %v", err)
	}
	defer file.Close()

	// 2. Generate path dinamis (folder/tahun/bulan/uuid.ext)
	now := time.Now()
	objectName := fmt.Sprintf("%s/%d/%02d/%s%s",
		folderTarget,
		now.Year(),
		now.Month(), // Tambahin %02d biar bulannya jadi 05, bukan 5
		uuid.New().String(),
		filepath.Ext(header.Filename),
	)

	// 3. Upload ke MinIO (Pastikan nama bucket sesuai di config/env lu)
	bucketName := "mantra-storage"
	_, err = config.MinioClient.PutObject(c, bucketName, objectName, file, header.Size, minio.PutObjectOptions{
		ContentType: header.Header.Get("Content-Type"),
	})

	if err != nil {
		return "", fmt.Errorf("gagal upload ke MinIO: %v", err)
	}

	// 4. Return URL publiknya (Sesuaikan dengan domain lu)
	// Kalau belum pake domain, ganti jadi IP lu misal: http://192.168.10.36:9000/mantra-storage/...
	fileUrl := fmt.Sprintf("https://storage.mantra.web.id/%s/%s", bucketName, objectName)

	return fileUrl, nil
}