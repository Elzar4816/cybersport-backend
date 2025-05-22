package models

import "time"

type News struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	Title    string    `json:"title"`
	Date     time.Time `json:"date"`
	Content  string    `json:"content"`
	ImageURL string    `json:"imageUrl"` // путь до загруженного файла
	VideoURL string    `json:"videoUrl"` // ссылка на YouTube
}

func (News) TableName() string {
	return "news"
}
