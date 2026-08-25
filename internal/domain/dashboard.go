package domain

import "time"

// DashboardStats represents overall dashboard statistics
type DashboardStats struct {
	UserID              int       `json:"user_id"`
	CurrentWeight       *float64  `json:"current_weight"`
	WeightChange        *float64  `json:"weight_change"` // positive = gain, negative = loss
	CurrentBMI          *float64  `json:"current_bmi"`
	HealthScore         *float64  `json:"health_score"`
	ActivityStreak      int       `json:"activity_streak"`
	MealsLoggedToday    int       `json:"meals_logged_today"`
	WorkoutsThisWeek    int       `json:"workouts_this_week"`
	SleepQualityAvg     *float64  `json:"sleep_quality_avg"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// ProgressTrend represents data trend over time
type ProgressTrend struct {
	UserID      int       `json:"user_id"`
	MetricType  string    `json:"metric_type"` // weight, bmi, calories, etc.
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	DataPoints  []TrendPoint `json:"data_points"`
}

// TrendPoint represents a single data point in a trend
type TrendPoint struct {
	Date   time.Time `json:"date"`
	Value  float64   `json:"value"`
	Label  string    `json:"label,omitempty"`
}

// WeeklySummary represents a weekly summary of user's progress
type WeeklySummary struct {
	UserID              int       `json:"user_id"`
	WeekStart           time.Time `json:"week_start"`
	WeekEnd             time.Time `json:"week_end"`
	AverageCalories     *float64  `json:"average_calories"`
	AverageSteps        *int      `json:"average_steps"`
	WeightChange        *float64  `json:"weight_change"`
	WorkoutCount        int       `json:"workout_count"`
	MealCompliance      float64   `json:"meal_compliance"` // percentage
	SleepQualityAvg     *float64  `json:"sleep_quality_avg"`
	MoodAvg             *float64  `json:"mood_avg"`
	Achievements        []string  `json:"achievements"`
	Recommendations     []string  `json:"recommendations"`
}

// GoalProgress represents progress towards user's goals
type GoalProgress struct {
	UserID          int       `json:"user_id"`
	GoalType        string    `json:"goal_type"` // weight_loss, muscle_gain, etc.
	TargetValue     float64   `json:"target_value"`
	CurrentValue    float64   `json:"current_value"`
	ProgressPercent float64   `json:"progress_percent"`
	StartDate       time.Time `json:"start_date"`
	TargetDate      time.Time `json:"target_date"`
	Status          string    `json:"status"` // on_track, behind, ahead, completed
	EstimatedCompletion *time.Time `json:"estimated_completion,omitempty"`
}

// ComparisonReport represents comparison with previous periods
type ComparisonReport struct {
	UserID            int       `json:"user_id"`
	ComparisonPeriod  string    `json:"comparison_period"` // week_over_week, month_over_month
	WeightChange      *float64  `json:"weight_change"`
	BMIChange         *float64  `json:"bmi_change"`
	CalorieChange     *float64  `json:"calorie_change"`
	ActivityChange    *float64  `json:"activity_change"`
	HealthScoreChange *float64  `json:"health_score_change"`
	Improvements      []string  `json:"improvements"`
	Concerns          []string  `json:"concerns"`
}

// DashboardRepository defines operations for dashboard data
type DashboardRepository interface {
	GetDashboardStats(userID int) (*DashboardStats, error)
	GetProgressTrend(userID int, metricType string, startDate, endDate time.Time) (*ProgressTrend, error)
	GetWeeklySummary(userID int, weekStart time.Time) (*WeeklySummary, error)
	GetGoalProgress(userID int) ([]GoalProgress, error)
	GetComparisonReport(userID int, period string) (*ComparisonReport, error)
}