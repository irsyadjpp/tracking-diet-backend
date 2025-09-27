package domain

import "time"

type AIRecommendation struct {
	ID                 int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID             int       `gorm:"not null;index" json:"user_id"`
	RecommendationText string    `gorm:"type:text" json:"recommendation_text"`
	GeneratedAt        time.Time `gorm:"autoCreateTime" json:"generated_at"`
}

type AIRecommendationRepository interface {
	Create(ar *AIRecommendation) error
	GetByID(id int) (*AIRecommendation, error)
	ListByUser(userID int) ([]AIRecommendation, error)
	Delete(id int) error
}
