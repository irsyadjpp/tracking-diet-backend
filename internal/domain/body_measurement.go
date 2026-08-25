package domain

import "time"

type BodyMeasurement struct {
	ID            int       `gorm:"primaryKey;autoIncrement;column:body_id" json:"id"`
	UserID        int       `gorm:"not null;index;column:user_id" json:"user_id"`
	MeasuredAt    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;column:measured_at" json:"measured_at"`
	WeightKg      *float64  `gorm:"column:weight_kg" json:"weight_kg"`
	BMI           *float64  `gorm:"column:bmi" json:"bmi"`
	WaistCm       *float64  `gorm:"column:waist_cm" json:"waist_cm"`
	HipCm         *float64  `gorm:"column:hip_cm" json:"hip_cm"`
	ArmCm         *float64  `gorm:"column:arm_cm" json:"arm_cm"`
	ThighCm       *float64  `gorm:"column:thigh_cm" json:"thigh_cm"`
	BodyFatPct    *float64  `gorm:"column:body_fat_pct" json:"body_fat_pct"`
	MuscleMassKg  *float64  `gorm:"column:muscle_mass_kg" json:"muscle_mass_kg"`
	VisceralFat   *float64  `gorm:"column:visceral_fat" json:"visceral_fat"`
	WHR           *float64  `gorm:"column:whr" json:"whr"`
	SkinfoldMm    *float64  `gorm:"column:skinfold_mm" json:"skinfold_mm"`
}

// TableName specifies the table name for BodyMeasurement model
func (BodyMeasurement) TableName() string {
	return "measurements_body"
}

type BodyMeasurementRepository interface {
	Create(bm *BodyMeasurement) error
	GetByID(id int) (*BodyMeasurement, error)
	ListByUser(userID int) ([]BodyMeasurement, error)
	ListByUserWithDateRange(userID int, startDate, endDate time.Time) ([]BodyMeasurement, error)
	Update(bm *BodyMeasurement) error
	Delete(id int) error
}
