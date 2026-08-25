package domain

import "time"

type NutritionMeasurement struct {
	ID             int       `gorm:"primaryKey;autoIncrement;column:nutrition_id" json:"id"`
	UserID         int       `gorm:"not null;index;column:user_id" json:"user_id"`
	MeasuredAt     time.Time `gorm:"not null;column:measured_at" json:"measured_at"`
	CaloriesInKcal *int      `gorm:"column:calories_in_kcal" json:"calories_in_kcal"`
	CarbsG         *float64  `gorm:"column:carbs_g" json:"carbs_g"`
	ProteinG       *float64  `gorm:"column:protein_g" json:"protein_g"`
	FatG           *float64  `gorm:"column:fat_g" json:"fat_g"`
	FiberG         *float64  `gorm:"column:fiber_g" json:"fiber_g"`
	WaterIntakeL   *float64  `gorm:"column:water_intake_l" json:"water_intake_l"`
	CaloriesOutKcal *int     `gorm:"column:calories_out_kcal" json:"calories_out_kcal"`
	StepCount      *int      `gorm:"column:step_count" json:"step_count"`
	MealTimingNote *string   `gorm:"column:meal_timing_note" json:"meal_timing_note"`
}

// TableName specifies the table name for NutritionMeasurement model
func (NutritionMeasurement) TableName() string {
	return "measurements_nutrition"
}

type NutritionMeasurementRepository interface {
	Create(nm *NutritionMeasurement) error
	GetByID(id int) (*NutritionMeasurement, error)
	ListByUser(userID int) ([]NutritionMeasurement, error)
	ListByUserWithDateRange(userID int, startDate, endDate time.Time) ([]NutritionMeasurement, error)
	GetDailySummary(userID int, date time.Time) (*NutritionMeasurement, error)
	Update(nm *NutritionMeasurement) error
	Delete(id int) error
}
