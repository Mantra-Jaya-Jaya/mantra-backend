package utils

import (
	"fmt"
	"os"
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

	// 4. Return URL publiknya (Menggunakan dynamic env STORAGE_PUBLIC_URL)
	storageURL := os.Getenv("STORAGE_PUBLIC_URL")
	if storageURL == "" {
		storageURL = "https://storage.mantra.web.id" // Fallback
	}
	fileUrl := fmt.Sprintf("%s/%s/%s", storageURL, bucketName, objectName)

	return fileUrl, nil
}