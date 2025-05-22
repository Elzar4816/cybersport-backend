package handlers

import (
	"cybersport-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func GetAllNews(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var news []models.News
		if err := db.Order("date DESC").Find(&news).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch news"})
			return
		}
		c.JSON(200, news)
	}
}
func CreateNewsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.PostForm("title")
		content := c.PostForm("content")
		videoUrl := c.PostForm("videoUrl")

		// Обработка картинки
		file, err := c.FormFile("image")
		var imagePath string
		if err == nil && file != nil {
			imagePath, err = saveImage(c, file)
			if err != nil {
				c.JSON(500, gin.H{"error": "Не удалось сохранить изображение"})
				return
			}
		}

		news := models.News{
			Title:    title,
			Content:  content,
			VideoURL: videoUrl,
			ImageURL: imagePath,
			Date:     time.Now(),
		}

		if err := db.Create(&news).Error; err != nil {
			c.JSON(500, gin.H{"error": "Ошибка при сохранении новости"})
			return
		}

		c.JSON(200, gin.H{"message": "Новость успешно создана"})
	}
}

func saveImage(c *gin.Context, file *multipart.FileHeader) (string, error) {
	uploadsDir := "./uploads"
	if err := os.MkdirAll(uploadsDir, os.ModePerm); err != nil {
		log.Println("Ошибка при создании папки uploads:", err)
		return "", err
	}

	filename := time.Now().Format("20060102_150405") + "_" + filepath.Base(file.Filename)
	path := filepath.Join(uploadsDir, filename)

	if err := c.SaveUploadedFile(file, path); err != nil {
		return "", err
	}

	return "/uploads/" + filename, nil
}
