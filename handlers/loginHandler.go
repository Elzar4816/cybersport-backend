package handlers

import (
	"cybersport-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"os"
	_ "os"
	"time"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func LoginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&creds); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		var user models.User
		if err := db.Where("username = ?", creds.Username).First(&user).Error; err != nil {
			c.JSON(401, gin.H{"error": "User not found"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
			c.JSON(401, gin.H{"error": "Incorrect password"})
			return
		}

		// Создаём JWT
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user": user.Username,
			"exp":  time.Now().Add(24 * time.Hour).Unix(),
		})

		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			c.JSON(500, gin.H{"error": "Token creation failed"})
			return
		}

		c.JSON(200, gin.H{"token": tokenString})
	}
}
