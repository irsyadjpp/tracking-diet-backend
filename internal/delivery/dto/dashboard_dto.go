package dto

import "time"

// DashboardStatsResponse represents dashboard statistics response
type DashboardStatsResponse struct {
	UserID              int       `json:"user_id"`
	CurrentWeight       *float64  `json:"current_weight"`
	WeightChange        *float64  `json:"weight_change"`
	CurrentBMI          *float64  `json:"current_bmi"`
	HealthScore         *float64  `json:"health_score"`
	ActivityStreak      int       `json:"activity_streak"`
	MealsLoggedToday    int       `json:"meals_logged_today"`
	WorkoutsThisWeek    int       `json:"workouts_this_week"`
	SleepQualityAvg     *float64  `json:"sleep_quality_avg"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// ProgressTrendResponse represents progress trend response
type ProgressTrendResponse struct {
	UserID     int         `json:"user_id"`
	MetricType string      `json:"metric_type"`
	StartDate  time.Time   `json:"start_date"`
	EndDate    time.Time   `json:"end_date"`
	DataPoints []TrendPoint `json:"data_points"`
}

// TrendPointResponse represents a trend point response
type TrendPointResponse struct {
	Date  time.Time `json:"date"`
	Value float64   `json:"value"`
	Label string    `json:"label,omitempty"`
}

// WeeklySummaryResponse represents weekly summary response
type WeeklySummaryResponse struct {
	UserID            int       `json:"user_id"`
	WeekStart         time.Time `json:"week_start"`
	WeekEnd           time.Time `json:"week_end"`
	AverageCalories   *float64  `json:"average_calories"`
	AverageSteps      *int      `json:"average_steps"`
	WeightChange      *float64  `json:"weight_change"`
	WorkoutCount      int       `json:"workout_count"`
	MealCompliance    float64   `json:"meal_compliance"`
	SleepQualityAvg   *float64  `json:"sleep_quality_avg"`
	MoodAvg           *float64  `json:"mood_avg"`
	Achievements      []string  `json:"achievements"`
	Recommendations   []string  `json:"recommendations"`
}

// GoalProgressResponse represents goal progress response
type GoalProgressResponse struct {
	UserID               int        `json:"user_id"`
	GoalType             string     `json:"goal_type"`
	TargetValue          float64    `json:"target_value"`
	CurrentValue         float64    `json:"current_value"`
	ProgressPercent      float64    `json:"progress_percent"`
	StartDate            time.Time  `json:"start_date"`
	TargetDate           time.Time  `json:"target_date"`
	Status               string     `json:"status"`
	EstimatedCompletion *time.Time `json:"estimated_completion,omitempty"`
}

// ComparisonReportResponse represents comparison report response
type ComparisonReportResponse struct {
	UserID            int       `json:"user_id"`
	ComparisonPeriod   string    `json:"comparison_period"`
	WeightChange       *float64  `json:"weight_change"`
	BMIChange          *float64  `json:"bmi_change"`
	CalorieChange      *float64  `json:"calorie_change"`
	ActivityChange     *float64  `json:"activity_change"`
	HealthScoreChange  *float64  `json:"health_score_change"`
	Improvements       []string  `json:"improvements"`
	Concerns           []string  `json:"concerns"`
}