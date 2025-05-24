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
		&models.Game{},
		&models.Tournament{},
	)
	SeedPressUser(db)
	SeedGames(db)

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
func SeedGames(db *gorm.DB) {
	games := []struct {
		Name    string
		LogoURL string
	}{
		{"PUBG Mobile", "/tournaments-logo/pubg-mobile.png"},
		{"Call of Duty: Mobile", "/tournaments-logo/call-of-duty-mobile.png"},
		{"Roblox", "/tournaments-logo/roblox.png"},
		{"World of Tanks", "/tournaments-logo/world-of-tanks.png"},
		{"Counter-Strike 1.6", "/tournaments-logo/cs-1-6.png"},
		{"Counter-Strike 2", "/tournaments-logo/cs-2.png"},
		{"Dota 2", "/tournaments-logo/dota-2.png"},
		{"Clash Royale", "/tournaments-logo/clash-royale.png"},
		{"Honor of Kings", "/tournaments-logo/honor-of-kings.png"},
	}

	for _, g := range games {
		var count int64
		db.Model(&models.Game{}).Where("name = ?", g.Name).Count(&count)
		if count == 0 {
			db.Create(&models.Game{Name: g.Name, LogoURL: g.LogoURL})
			log.Printf("Игра добавлена: %s\n", g.Name)
		}
	}
}
