package domain

import "time"

type WellbeingMeasurement struct {
	ID             int       `gorm:"primaryKey;autoIncrement;column:wellbeing_id" json:"id"`
	UserID         int       `gorm:"not null;index;column:user_id" json:"user_id"`
	MeasuredAt     time.Time `gorm:"not null;column:measured_at" json:"measured_at"`
	EnergyLevel    *int      `gorm:"column:energy_level" json:"energy_level"`
	MoodLevel      *int      `gorm:"column:mood_level" json:"mood_level"`
	HungerLevel    *int      `gorm:"column:hunger_level" json:"hunger_level"`
	SleepHours     *float64  `gorm:"column:sleep_hours" json:"sleep_hours"`
	SleepQuality   *int      `gorm:"column:sleep_quality" json:"sleep_quality"`
	DigestionNote  *string   `gorm:"column:digestion_note" json:"digestion_note"`
	StressLevel    *int      `gorm:"column:stress_level" json:"stress_level"`
}

// TableName specifies the table name for WellbeingMeasurement model
func (WellbeingMeasurement) TableName() string {
	return "measurements_wellbeing"
}

type WellbeingMeasurementRepository interface {
	Create(wm *WellbeingMeasurement) error
	GetByID(id int) (*WellbeingMeasurement, error)
	ListByUser(userID int) ([]WellbeingMeasurement, error)
	ListByUserWithDateRange(userID int, startDate, endDate time.Time) ([]WellbeingMeasurement, error)
	Update(wm *WellbeingMeasurement) error
	Delete(id int) error
}
