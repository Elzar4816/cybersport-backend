package storage

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"mime/multipart"
	"net/url"
	_ "strings"
	"time"
)

var (
	MinioClient    *minio.Client
	MinioPublicURL string
)

func InitMinio(endpoint, accessKeyID, secretAccessKey string, useSSL bool) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к MinIO: %v", err)
	}
	MinioClient = client

	MinioPublicURL = endpoint // чтобы не подменять ничего

	bucketName := "news-images"
	location := "us-east-1"

	exists, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		log.Fatalf("❌ Ошибка при проверке bucket: %v", err)
	}
	if !exists {
		err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: location})
		if err != nil {
			log.Fatalf("❌ Не удалось создать bucket: %v", err)
		}
	}
	log.Println("✅ MinIO инициализирован")
}

func SaveImageToMinio(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	objectName := time.Now().Format("20060102_150405") + "_" + file.Filename
	contentType := file.Header.Get("Content-Type")

	_, err = MinioClient.PutObject(context.Background(), "news-images", objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	return objectName, nil
}

func GetPresignedURL(objectName string) (string, error) {
	presignedURL, err := MinioClient.PresignedGetObject(
		context.Background(),
		"news-images",
		objectName,
		time.Hour,
		url.Values{},
	)
	if err != nil {
		return "", err
	}

	parsedURL, err := url.Parse(presignedURL.String())
	if err != nil {
		return "", err
	}

	return parsedURL.String(), nil
}
func DeleteImageFromMinio(objectName string) error {
	err := MinioClient.RemoveObject(context.Background(), "news-images", objectName, minio.RemoveObjectOptions{})
	if err != nil {
		log.Printf("❌ Не удалось удалить изображение %s из MinIO: %v", objectName, err)
		return err
	}
	log.Printf("🗑️ Изображение %s удалено из MinIO", objectName)
	return nil
}
