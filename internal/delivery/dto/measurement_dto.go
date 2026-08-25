package dto

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// BodyMeasurementRequest represents body measurement request
type BodyMeasurementRequest struct {
	WeightKg      *float64 `json:"weight_kg" validate:"omitempty,gt=0"`
	WaistCm       *float64 `json:"waist_cm" validate:"omitempty,gt=0"`
	HipCm         *float64 `json:"hip_cm" validate:"omitempty,gt=0"`
	ChestCm       *float64 `json:"chest_cm" validate:"omitempty,gt=0"`
	ThighCm       *float64 `json:"thigh_cm" validate:"omitempty,gt=0"`
	ArmCm         *float64 `json:"arm_cm" validate:"omitempty,gt=0"`
	BodyFatPct    *float64 `json:"body_fat_pct" validate:"omitempty,gte=0,lte=100"`
	MuscleMassKg  *float64 `json:"muscle_mass_kg" validate:"omitempty,gt=0"`
	VisceralFat   *float64 `json:"visceral_fat" validate:"omitempty,gte=0"`
	SkinfoldMm    *float64 `json:"skinfold_mm" validate:"omitempty,gt=0"`
	MeasuredAt    string   `json:"measured_at" validate:"required"`
}

// BodyMeasurementResponse represents body measurement response
type BodyMeasurementResponse struct {
	ID            int        `json:"id"`
	UserID        int        `json:"user_id"`
	MeasuredAt    time.Time  `json:"measured_at"`
	WeightKg      *float64   `json:"weight_kg"`
	BMI           *float64   `json:"bmi"`
	WaistCm       *float64   `json:"waist_cm"`
	HipCm         *float64   `json:"hip_cm"`
	ChestCm       *float64   `json:"chest_cm"`
	ThighCm       *float64   `json:"thigh_cm"`
	ArmCm         *float64   `json:"arm_cm"`
	BodyFatPct    *float64   `json:"body_fat_pct"`
	MuscleMassKg  *float64   `json:"muscle_mass_kg"`
	VisceralFat   *float64   `json:"visceral_fat"`
	WHR           *float64   `json:"whr"`
	SkinfoldMm    *float64   `json:"skinfold_mm"`
}

// NutritionMeasurementRequest represents nutrition measurement request
type NutritionMeasurementRequest struct {
	CaloriesInKcal *int     `json:"calories_in_kcal" validate:"omitempty,gte=0"`
	CarbsG         *float64 `json:"carbs_g" validate:"omitempty,gte=0"`
	ProteinG       *float64 `json:"protein_g" validate:"omitempty,gte=0"`
	FatG           *float64 `json:"fat_g" validate:"omitempty,gte=0"`
	FiberG         *float64 `json:"fiber_g" validate:"omitempty,gte=0"`
	WaterIntakeL   *float64 `json:"water_intake_l" validate:"omitempty,gte=0"`
	CaloriesOutKcal *int    `json:"calories_out_kcal" validate:"omitempty,gte=0"`
	StepCount      *int     `json:"step_count" validate:"omitempty,gte=0"`
	MealTimingNote *string  `json:"meal_timing_note"`
	MeasuredAt     string   `json:"measured_at" validate:"required"`
}

// NutritionMeasurementResponse represents nutrition measurement response
type NutritionMeasurementResponse struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	MeasuredAt     time.Time `json:"measured_at"`
	CaloriesInKcal *int      `json:"calories_in_kcal"`
	CarbsG         *float64  `json:"carbs_g"`
	ProteinG       *float64  `json:"protein_g"`
	FatG           *float64  `json:"fat_g"`
	FiberG         *float64  `json:"fiber_g"`
	WaterIntakeL   *float64  `json:"water_intake_l"`
	CaloriesOutKcal *int      `json:"calories_out_kcal"`
	StepCount      *int      `json:"step_count"`
	MealTimingNote *string   `json:"meal_timing_note"`
}

// MetabolicMeasurementRequest represents metabolic measurement request
type MetabolicMeasurementRequest struct {
	SystolicBP          *int     `json:"systolic_bp" validate:"omitempty,gte=0,lte=300"`
	DiastolicBP         *int     `json:"diastolic_bp" validate:"omitempty,gte=0,lte=200"`
	FastingGlucose      *float64 `json:"fasting_glucose" validate:"omitempty,gte=0"`
	Hba1cPct            *float64 `json:"hba1c_pct" validate:"omitempty,gte=0,lte=20"`
	CholesterolTotal    *float64 `json:"cholesterol_total" validate:"omitempty,gte=0"`
	LDL                 *float64 `json:"ldl" validate:"omitempty,gte=0"`
	HDL                 *float64 `json:"hdl" validate:"omitempty,gte=0"`
	Triglycerides       *float64 `json:"triglycerides" validate:"omitempty,gte=0"`
	UricAcid            *float64 `json:"uric_acid" validate:"omitempty,gte=0"`
	LiverFunctionNote   *string  `json:"liver_function_note"`
	KidneyFunctionNote  *string  `json:"kidney_function_note"`
	MeasuredAt          string   `json:"measured_at" validate:"required"`
}

// MetabolicMeasurementResponse represents metabolic measurement response
type MetabolicMeasurementResponse struct {
	ID                  int       `json:"id"`
	UserID              int       `json:"user_id"`
	MeasuredAt          time.Time `json:"measured_at"`
	SystolicBP          *int      `json:"systolic_bp"`
	DiastolicBP         *int      `json:"diastolic_bp"`
	FastingGlucose      *float64  `json:"fasting_glucose"`
	Hba1cPct            *float64  `json:"hba1c_pct"`
	CholesterolTotal    *float64  `json:"cholesterol_total"`
	LDL                 *float64  `json:"ldl"`
	HDL                 *float64  `json:"hdl"`
	Triglycerides       *float64  `json:"triglycerides"`
	UricAcid            *float64  `json:"uric_acid"`
	LiverFunctionNote   *string   `json:"liver_function_note"`
	KidneyFunctionNote  *string   `json:"kidney_function_note"`
}

// FitnessMeasurementRequest represents fitness measurement request
type FitnessMeasurementRequest struct {
	VO2MaxMlKgMin   *float64 `json:"vo2max_ml_kg_min" validate:"omitempty,gte=0"`
	Strength1RmKg   *float64 `json:"strength_1rm_kg" validate:"omitempty,gte=0"`
	PushupCount     *int     `json:"pushup_count" validate:"omitempty,gte=0"`
	SquatCount      *int     `json:"squat_count" validate:"omitempty,gte=0"`
	EnduranceNote   *string  `json:"endurance_note"`
	FlexibilityNote *string  `json:"flexibility_note"`
	MeasuredAt      string   `json:"measured_at" validate:"required"`
}

// FitnessMeasurementResponse represents fitness measurement response
type FitnessMeasurementResponse struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	MeasuredAt      time.Time `json:"measured_at"`
	VO2MaxMlKgMin   *float64  `json:"vo2max_ml_kg_min"`
	Strength1RmKg   *float64  `json:"strength_1rm_kg"`
	PushupCount     *int      `json:"pushup_count"`
	SquatCount      *int      `json:"squat_count"`
	EnduranceNote   *string   `json:"endurance_note"`
	FlexibilityNote *string   `json:"flexibility_note"`
}

// WellbeingMeasurementRequest represents wellbeing measurement request
type WellbeingMeasurementRequest struct {
	EnergyLevel   *int     `json:"energy_level" validate:"omitempty,gte=1,lte=10"`
	MoodLevel     *int     `json:"mood_level" validate:"omitempty,gte=1,lte=10"`
	HungerLevel   *int     `json:"hunger_level" validate:"omitempty,gte=1,lte=10"`
	SleepHours    *float64 `json:"sleep_hours" validate:"omitempty,gte=0,lte=24"`
	SleepQuality  *int     `json:"sleep_quality" validate:"omitempty,gte=1,lte=10"`
	DigestionNote *string  `json:"digestion_note"`
	StressLevel   *int     `json:"stress_level" validate:"omitempty,gte=1,lte=10"`
	MeasuredAt    string   `json:"measured_at" validate:"required"`
}

// WellbeingMeasurementResponse represents wellbeing measurement response
type WellbeingMeasurementResponse struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	MeasuredAt    time.Time `json:"measured_at"`
	EnergyLevel   *int      `json:"energy_level"`
	MoodLevel     *int      `json:"mood_level"`
	HungerLevel   *int      `json:"hunger_level"`
	SleepHours    *float64  `json:"sleep_hours"`
	SleepQuality  *int      `json:"sleep_quality"`
	DigestionNote *string   `json:"digestion_note"`
	StressLevel   *int      `json:"stress_level"`
}

// LabTestRequest represents lab test request
type LabTestRequest struct {
	TestName       string  `json:"test_name" validate:"required,min=1,max=100"`
	ResultValue    *string `json:"result_value"`
	Unit           *string `json:"unit"`
	ReferenceRange *string `json:"reference_range"`
	MeasuredAt     string  `json:"measured_at" validate:"required"`
}

// LabTestResponse represents lab test response
type LabTestResponse struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	TestName       string    `json:"test_name"`
	ResultValue    *string   `json:"result_value"`
	Unit           *string   `json:"unit"`
	ReferenceRange *string   `json:"reference_range"`
	MeasuredAt     time.Time `json:"measured_at"`
}

// DateRangeRequest represents date range filter request
type DateRangeRequest struct {
	StartDate string `json:"start_date" validate:"required"`
	EndDate   string `json:"end_date" validate:"required"`
}

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page     int `json:"page" validate:"omitempty,gte=1"`
	PageSize int `json:"page_size" validate:"omitempty,gte=1,lte=100"`
}

// Validate methods
func (r *BodyMeasurementRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *NutritionMeasurementRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *MetabolicMeasurementRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *FitnessMeasurementRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *WellbeingMeasurementRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *LabTestRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *DateRangeRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *PaginationRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// AIRecommendationResponse represents AI recommendation response
type AIRecommendationResponse struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	GeneratedAt     time.Time `json:"generated_at"`
	Category        string    `json:"category"`
	InputSummary    interface{} `json:"input_summary"`
	OutputText      *string   `json:"output_text"`
	ConfidenceScore float64   `json:"confidence_score"`
}