package db

import (
	"cybersport-backend/models"
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var db *gorm.DB

func ConnectDB() *gorm.DB {
	// Загружаем переменные окружения
	if err := godotenv.Load("configuration.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Формируем строку подключения
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable client_encoding=UTF8",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	log.Println("Database connected successfully.")

	db.AutoMigrate(
		&models.News{},
		&models.User{},
		&models.Tournament{},
	)
	SeedPressUser(db)
	return db
}

func SeedPressUser(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Where("username = ?", "press").Count(&count)
	if count == 0 {
		passwordHash, _ := bcrypt.GenerateFromPassword([]byte("press"), bcrypt.DefaultCost)
		db.Create(&models.User{
			Username: "press",
			Password: string(passwordHash),
		})
		log.Println("Пресс-пользователь создан: press / press")
	}
}
