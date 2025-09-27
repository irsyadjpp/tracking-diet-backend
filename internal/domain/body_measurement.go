package domain

import "time"

type BodyMeasurement struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int       `gorm:"not null;index" json:"user_id"`
	WeightKg   float64   `json:"weight_kg"`
	BMI        float64   `json:"bmi"`
	WaistCm    float64   `json:"waist_cm"`
	HipCm      float64   `json:"hip_cm"`
	ChestCm    float64   `json:"chest_cm"`
	ThighCm    float64   `json:"thigh_cm"`
	ArmCm      float64   `json:"arm_cm"`
	MeasuredAt time.Time `gorm:"not null" json:"measured_at"`
}

type BodyMeasurementRepository interface {
	Create(bm *BodyMeasurement) error
	GetByID(id int) (*BodyMeasurement, error)
	ListByUser(userID int) ([]BodyMeasurement, error)
	Update(bm *BodyMeasurement) error
	Delete(id int) error
}
