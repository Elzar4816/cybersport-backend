package handlers

import (
	"cybersport-backend/models"
	"cybersport-backend/storage"
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

		// добавляем presigned URL
		for i := range news {
			if news[i].ImageURL != "" {
				url, err := storage.GetPresignedURL(news[i].ImageURL)
				if err == nil {
					news[i].ImageURL = url
				}
			}
		}

		c.JSON(200, news)
	}
}

func CreateNewsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.PostForm("title")
		content := c.PostForm("content")
		videoUrl := c.PostForm("videoUrl")

		file, err := c.FormFile("image")
		var imagePath string
		if err == nil && file != nil {
			imagePath, err = storage.SaveImageToMinio(file)

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

func UpdateNewsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var news models.News
		id := c.Param("id")

		if err := db.First(&news, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "Новость не найдена"})
			return
		}

		title := c.PostForm("title")
		content := c.PostForm("content")
		videoUrl := c.PostForm("videoUrl")
		removeImage := c.PostForm("removeImage") // из FormData приходит строка

		if title != "" {
			news.Title = title
		}
		if content != "" {
			news.Content = content
		}
		if videoUrl != "" {
			news.VideoURL = videoUrl
		}

		file, err := c.FormFile("image")
		if err == nil && file != nil {
			// Загружается новое изображение → удалить старое и сохранить новое
			if news.ImageURL != "" {
				_ = storage.DeleteImageFromMinio(news.ImageURL)
			}

			newImageName, err := storage.SaveImageToMinio(file)
			if err != nil {
				c.JSON(500, gin.H{"error": "Ошибка при загрузке нового изображения"})
				return
			}
			news.ImageURL = newImageName
		} else if removeImage == "true" && news.ImageURL != "" {
			// Флаг удаления изображения активен, нового файла нет → удалить
			_ = storage.DeleteImageFromMinio(news.ImageURL)
			news.ImageURL = ""
		}

		if err := db.Save(&news).Error; err != nil {
			c.JSON(500, gin.H{"error": "Ошибка при обновлении новости"})
			return
		}

		c.JSON(200, gin.H{"message": "Новость обновлена"})
	}
}

func DeleteNewsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var news models.News
		if err := db.First(&news, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "Новость не найдена"})
			return
		}

		if news.ImageURL != "" {
			err := storage.DeleteImageFromMinio(news.ImageURL)
			if err != nil {
				log.Printf("⚠️ Не удалось удалить изображение из MinIO: %v", err)
			}
		}

		if err := db.Delete(&news).Error; err != nil {
			c.JSON(500, gin.H{"error": "Не удалось удалить новость"})
			return
		}

		c.JSON(200, gin.H{"message": "Новость удалена"})
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
func GetNewsByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var news models.News

		if err := db.First(&news, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "Новость не найдена"})
			return
		}

		if news.ImageURL != "" {
			url, err := storage.GetPresignedURL(news.ImageURL)
			if err == nil {
				news.ImageURL = url
			}
		}

		c.JSON(200, news)
	}
}
