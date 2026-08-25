package domain

import "time"

type FitnessMeasurement struct {
	ID             int       `gorm:"primaryKey;autoIncrement;column:fitness_id" json:"id"`
	UserID         int       `gorm:"not null;index;column:user_id" json:"user_id"`
	MeasuredAt     time.Time `gorm:"not null;column:measured_at" json:"measured_at"`
	VO2MaxMlKgMin  *float64  `gorm:"column:vo2max_ml_kg_min" json:"vo2max_ml_kg_min"`
	Strength1RmKg  *float64  `gorm:"column:strength_1rm_kg" json:"strength_1rm_kg"`
	PushupCount    *int      `gorm:"column:pushup_count" json:"pushup_count"`
	SquatCount     *int      `gorm:"column:squat_count" json:"squat_count"`
	EnduranceNote  *string   `gorm:"column:endurance_note" json:"endurance_note"`
	FlexibilityNote *string  `gorm:"column:flexibility_note" json:"flexibility_note"`
}

// TableName specifies the table name for FitnessMeasurement model
func (FitnessMeasurement) TableName() string {
	return "measurements_fitness"
}

type FitnessMeasurementRepository interface {
	Create(fm *FitnessMeasurement) error
	GetByID(id int) (*FitnessMeasurement, error)
	ListByUser(userID int) ([]FitnessMeasurement, error)
	ListByUserWithDateRange(userID int, startDate, endDate time.Time) ([]FitnessMeasurement, error)
	Update(fm *FitnessMeasurement) error
	Delete(id int) error
}
