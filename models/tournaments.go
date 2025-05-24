package models

import "time"

type TournamentStatus string

const (
	StatusUpcoming  TournamentStatus = "upcoming"
	StatusOngoing   TournamentStatus = "ongoing"
	StatusCompleted TournamentStatus = "completed"
)

type Tournament struct {
	ID          uint64           `json:"id" gorm:"primaryKey"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	GameID      uint             `json:"game_id"` // внешний ключ
	Game        Game             `json:"game" gorm:"foreignKey:GameID"`
	Region      string           `json:"region"`
	StartDate   time.Time        `json:"start_date"`
	EndDate     time.Time        `json:"end_date"`
	Status      TournamentStatus `json:"status"`
	PrizePool   *int             `json:"prize_pool"`
	Stage       string           `json:"stage"`
	IsOpen      bool             `json:"is_open"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

func (Tournament) TableName() string {
	return "tournaments"
}

type Game struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	Name    string `json:"name" gorm:"uniqueIndex;not null"`
	LogoURL string `json:"logo_url" gorm:"not null"`
}

func (Game) TableName() string {
	return "games"
}
