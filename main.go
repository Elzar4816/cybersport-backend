package main

import (
	"cybersport-backend/db"
	"cybersport-backend/handlers"
	"cybersport-backend/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

func main() {
	// Подключаемся к БД
	gormDB := db.ConnectDB()
	if gormDB == nil {
		log.Fatal("Database connection failed")
	}

	// Инициализация Gin
	r := gin.Default()
	r.Use(gin.Recovery())

	// Подключаем статику и отдачу frontend'а
	setupStatic(r)

	// Подключаем маршруты
	setupRoutes(r, gormDB)

	// Запуск сервера
	log.Println("Server started on http://localhost:8000")
	if err := r.Run(":8000"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
func setupStatic(r *gin.Engine) {
	// Раздача изображений
	r.Static("/uploads", "./uploads")

}

func setupRoutes(r *gin.Engine, gormDB *gorm.DB) {
	// Группа публичных API
	api := r.Group("/api")
	{
		// Аутентификация
		api.POST("/login", handlers.LoginHandler(gormDB))

		// Новости (публичные)
		api.GET("/news", handlers.GetAllNews(gormDB))
		api.GET("/news/:id", handlers.GetNewsByID(gormDB))

		api.GET("/games", handlers.GetGames(gormDB))

		// Турниры
		tournament := api.Group("/tournaments")
		{
			tournament.POST("/", handlers.CreateTournament(gormDB))
			tournament.GET("/", handlers.GetTournaments(gormDB))
		}
	}

	// Защищённая пресс-группа (только авторизованные пользователи)
	press := r.Group("/press")
	press.Use(middleware.AuthMiddleware_forLogin())
	{
		press.GET("/profile", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Пресс-панель доступна"})
		})
		press.POST("/news", handlers.CreateNewsHandler(gormDB))
		press.PUT("/news/:id", handlers.UpdateNewsHandler(gormDB))
		press.DELETE("/news/:id", handlers.DeleteNewsHandler(gormDB))
	}
}
