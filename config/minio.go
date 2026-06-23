package config

import (
	"context"
	"log"
	"os"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func InitMinio() {
	// 1. Load file .env
	if err := godotenv.Load(); err != nil {
		log.Println("Gak nemu file .env, pake environment variable sistem aja")
	}

	// 2. Ambil data dari env
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	bucketName := os.Getenv("MINIO_BUCKET")
	if bucketName == "" {
		bucketName = "mantra-storage"
	}

	// 3. Inisialisasi client
	var err error
	MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false, 
	})
	
	if err != nil {
		log.Fatalln("Gagal inisialisasi MinIO:", err)
	}
	log.Println("MinIO berhasil tersambung!")

	// 4. Pastikan bucket sudah terbuat secara otomatis
	ctx := context.Background()
	exists, err := MinioClient.BucketExists(ctx, bucketName)
	if err != nil {
		log.Println("Gagal memeriksa keberadaan bucket MinIO:", err)
		return
	}
	if !exists {
		log.Printf("Bucket '%s' belum ada, mencoba membuat otomatis...", bucketName)
		err = MinioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Printf("Gagal membuat bucket '%s' otomatis: %v", bucketName, err)
		} else {
			log.Printf("Bucket '%s' berhasil dibuat otomatis!", bucketName)
		}
	}
}