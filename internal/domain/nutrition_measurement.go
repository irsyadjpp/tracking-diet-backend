package domain

import "time"

type NutritionMeasurement struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int       `gorm:"not null;index" json:"user_id"`
	Calories   float64   `json:"calories"`
	Carbs      float64   `json:"carbs"`
	Protein    float64   `json:"protein"`
	Fat        float64   `json:"fat"`
	Fiber      float64   `json:"fiber"`
	Sugar      float64   `json:"sugar"`
	Sodium     float64   `json:"sodium"`
	MeasuredAt time.Time `gorm:"not null" json:"measured_at"`
}

type NutritionMeasurementRepository interface {
	Create(nm *NutritionMeasurement) error
	GetByID(id int) (*NutritionMeasurement, error)
	ListByUser(userID int) ([]NutritionMeasurement, error)
	Update(nm *NutritionMeasurement) error
	Delete(id int) error
}
