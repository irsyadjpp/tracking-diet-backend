package domain

import "time"

type MetabolicMeasurement struct {
	ID                  int       `gorm:"primaryKey;autoIncrement;column:metabolic_id" json:"id"`
	UserID              int       `gorm:"not null;index;column:user_id" json:"user_id"`
	MeasuredAt          time.Time `gorm:"not null;column:measured_at" json:"measured_at"`
	SystolicBP          *int      `gorm:"column:systolic_bp" json:"systolic_bp"`
	DiastolicBP         *int      `gorm:"column:diastolic_bp" json:"diastolic_bp"`
	FastingGlucose      *float64  `gorm:"column:fasting_glucose" json:"fasting_glucose"`
	Hba1cPct            *float64  `gorm:"column:hba1c_pct" json:"hba1c_pct"`
	CholesterolTotal   *float64  `gorm:"column:cholesterol_total" json:"cholesterol_total"`
	LDL                 *float64  `gorm:"column:ldl" json:"ldl"`
	HDL                 *float64  `gorm:"column:hdl" json:"hdl"`
	Triglycerides       *float64  `gorm:"column:triglycerides" json:"triglycerides"`
	UricAcid            *float64  `gorm:"column:uric_acid" json:"uric_acid"`
	LiverFunctionNote  *string   `gorm:"column:liver_function_note" json:"liver_function_note"`
	KidneyFunctionNote *string   `gorm:"column:kidney_function_note" json:"kidney_function_note"`
}

// TableName specifies the table name for MetabolicMeasurement model
func (MetabolicMeasurement) TableName() string {
	return "measurements_metabolic"
}

type MetabolicMeasurementRepository interface {
	Create(mm *MetabolicMeasurement) error
	GetByID(id int) (*MetabolicMeasurement, error)
	ListByUser(userID int) ([]MetabolicMeasurement, error)
	ListByUserWithDateRange(userID int, startDate, endDate time.Time) ([]MetabolicMeasurement, error)
	Update(mm *MetabolicMeasurement) error
	Delete(id int) error
}
