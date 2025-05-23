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
	Description string           `json:"description"` // краткое описание
	Game        string           `json:"game"`
	Region      string           `json:"region"` // например, "EMEA", "Americas", "Student"
	StartDate   time.Time        `json:"start_date"`
	EndDate     time.Time        `json:"end_date"`
	Status      TournamentStatus `json:"status"`     // upcoming, ongoing, completed
	PrizePool   *int             `json:"prize_pool"` // может быть nil, если не указан
	Stage       string           `json:"stage"`      // например, "Stage 1 2025"
	IsOpen      bool             `json:"is_open"`    // можно ли зарегистрироваться
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

func (Tournament) TableName() string {
	return "tournaments"
}
