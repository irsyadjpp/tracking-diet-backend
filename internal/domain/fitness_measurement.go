package domain

import "time"

type FitnessMeasurement struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int       `gorm:"not null;index" json:"user_id"`
	Steps       int       `json:"steps"`
	DistanceKm  float64   `json:"distance_km"`
	CaloriesOut float64   `json:"calories_out"`
	WorkoutMin  int       `json:"workout_min"`
	MeasuredAt  time.Time `gorm:"not null" json:"measured_at"`
}

type FitnessMeasurementRepository interface {
	Create(fm *FitnessMeasurement) error
	GetByID(id int) (*FitnessMeasurement, error)
	ListByUser(userID int) ([]FitnessMeasurement, error)
	Update(fm *FitnessMeasurement) error
	Delete(id int) error
}
