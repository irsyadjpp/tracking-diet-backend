package domain

import "time"

type WellbeingMeasurement struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int       `gorm:"not null;index" json:"user_id"`
	SleepHrs    float64   `json:"sleep_hrs"`
	MoodScore   int       `json:"mood_score"`   // e.g., scale 1–10
	StressLvl   int       `json:"stress_level"` // e.g., scale 1–10
	WaterIntake float64   `json:"water_intake"` // in liters
	MeasuredAt  time.Time `gorm:"not null" json:"measured_at"`
}

type WellbeingMeasurementRepository interface {
	Create(wm *WellbeingMeasurement) error
	GetByID(id int) (*WellbeingMeasurement, error)
	ListByUser(userID int) ([]WellbeingMeasurement, error)
	Update(wm *WellbeingMeasurement) error
	Delete(id int) error
}
