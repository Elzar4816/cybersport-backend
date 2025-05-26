package handlers

import (
	"net/http"
	"time"
	_ "time"

	"cybersport-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateTournamentInput struct {
	Title       string                  `json:"title" binding:"required"`
	Description string                  `json:"description"`
	Region      string                  `json:"region" binding:"required"`
	StartDate   string                  `json:"start_date" binding:"required"`
	EndDate     string                  `json:"end_date" binding:"required"`
	Status      models.TournamentStatus `json:"status" binding:"required,oneof=upcoming ongoing completed"`
	PrizePool   *int                    `json:"prize_pool"`
	Stage       string                  `json:"stage"`
	IsOpen      bool                    `json:"is_open"`
	GameID      uint                    `json:"game_id" binding:"required"` // ✅ обязательный внешний ключ
}

func CreateTournament(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreateTournamentInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// ✳️ Парсим дату вручную
		startDate, err := time.Parse("2006-01-02", input.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат start_date. Ожидается YYYY-MM-DD"})
			return
		}
		endDate, err := time.Parse("2006-01-02", input.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат end_date. Ожидается YYYY-MM-DD"})
			return
		}

		tournament := models.Tournament{
			Title:       input.Title,
			Description: input.Description,
			Region:      input.Region,
			StartDate:   startDate,
			EndDate:     endDate,
			Status:      input.Status,
			PrizePool:   input.PrizePool,
			Stage:       input.Stage,
			IsOpen:      input.IsOpen,
			GameID:      input.GameID,
		}

		if err := db.Create(&tournament).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании турнира"})
			return
		}

		c.JSON(http.StatusCreated, tournament)
	}
}

func GetTournaments(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tournaments []models.Tournament
		if err := db.Preload("Game").Order("start_date desc").Find(&tournaments).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch tournaments"})
			return
		}

		c.JSON(http.StatusOK, tournaments)
	}
}
func GetGames(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var games []models.Game
		if err := db.Find(&games).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch games"})
			return
		}
		c.JSON(http.StatusOK, games)
	}
}
func UpdateTournament(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var existing models.Tournament
		id := c.Param("id")

		if err := db.First(&existing, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Турнир не найден"})
			return
		}

		var input CreateTournamentInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		startDate, err := time.Parse("2006-01-02", input.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат start_date. Ожидается YYYY-MM-DD"})
			return
		}
		endDate, err := time.Parse("2006-01-02", input.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат end_date. Ожидается YYYY-MM-DD"})
			return
		}

		existing.Title = input.Title
		existing.Description = input.Description
		existing.Region = input.Region
		existing.StartDate = startDate
		existing.EndDate = endDate
		existing.Status = input.Status
		existing.PrizePool = input.PrizePool
		existing.Stage = input.Stage
		existing.IsOpen = input.IsOpen
		existing.GameID = input.GameID

		if err := db.Save(&existing).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении"})
			return
		}

		c.JSON(http.StatusOK, existing)
	}
}

func DeleteTournament(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := db.Delete(&models.Tournament{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Турнир удалён"})
	}
}
