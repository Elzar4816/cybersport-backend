package handlers

import (
	"net/http"
	"time"

	"cybersport-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateTournamentInput struct {
	Title       string                  `json:"title" binding:"required"`
	Description string                  `json:"description"`
	Region      string                  `json:"region" binding:"required"`
	StartDate   time.Time               `json:"start_date" binding:"required"`
	EndDate     time.Time               `json:"end_date" binding:"required"`
	Status      models.TournamentStatus `json:"status" binding:"required,oneof=upcoming ongoing completed"`
	PrizePool   *int                    `json:"prize_pool"`
	Stage       string                  `json:"stage"`
	IsOpen      bool                    `json:"is_open"`
}

func CreateTournament(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreateTournamentInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tournament := models.Tournament{
			Title:       input.Title,
			Description: input.Description,
			Region:      input.Region,
			StartDate:   input.StartDate,
			EndDate:     input.EndDate,
			Status:      input.Status,
			PrizePool:   input.PrizePool,
			Stage:       input.Stage,
			IsOpen:      input.IsOpen,
		}

		if err := db.Create(&tournament).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create tournament"})
			return
		}

		c.JSON(http.StatusCreated, tournament)
	}
}

func GetTournaments(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tournaments []models.Tournament
		if err := db.Order("start_date desc").Find(&tournaments).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch tournaments"})
			return
		}

		c.JSON(http.StatusOK, tournaments)
	}
}
