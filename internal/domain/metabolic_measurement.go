package domain

import "time"

type MetabolicMeasurement struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int       `gorm:"not null;index" json:"user_id"`
	RMR        float64   `json:"rmr"`     // Resting Metabolic Rate
	BMR        float64   `json:"bmr"`     // Basal Metabolic Rate
	VO2Max     float64   `json:"vo2_max"` // Aerobic capacity
	MeasuredAt time.Time `gorm:"not null" json:"measured_at"`
}

type MetabolicMeasurementRepository interface {
	Create(mm *MetabolicMeasurement) error
	GetByID(id int) (*MetabolicMeasurement, error)
	ListByUser(userID int) ([]MetabolicMeasurement, error)
	Update(mm *MetabolicMeasurement) error
	Delete(id int) error
}
