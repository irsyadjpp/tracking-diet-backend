package domain

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type AIRecommendation struct {
	ID              int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          int       `gorm:"not null;index;column:user_id" json:"user_id"`
	GeneratedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP;column:generated_at" json:"generated_at"`
	Category        *string   `gorm:"size:50;column:category" json:"category"`
	InputSummary    JSONB     `gorm:"type:jsonb;column:input_summary" json:"input_summary"`
	OutputText      *string   `gorm:"type:text;column:output_text" json:"output_text"`
	ConfidenceScore *float64  `gorm:"column:confidence_score" json:"confidence_score"`
}

// TableName specifies the table name for AIRecommendation model
func (AIRecommendation) TableName() string {
	return "ai_recommendations"
}

// JSONB type for handling JSONB data in PostgreSQL
type JSONB struct {
	Data interface{}
}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		j.Data = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, &j.Data)
}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if j.Data == nil {
		return nil, nil
	}
	return json.Marshal(j.Data)
}

type AIRecommendationRepository interface {
	Create(ar *AIRecommendation) error
	GetByID(id int) (*AIRecommendation, error)
	ListByUser(userID int) ([]AIRecommendation, error)
	ListByUserWithCategory(userID int, category string) ([]AIRecommendation, error)
	Update(ar *AIRecommendation) error
	Delete(id int) error
}
