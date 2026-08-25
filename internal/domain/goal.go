package domain

import "time"

// GoalType represents different types of goals
type GoalType string

const (
	GoalTypeWeightLoss      GoalType = "weight_loss"
	GoalTypeWeightGain      GoalType = "weight_gain"
	GoalTypeMuscleGain     GoalType = "muscle_gain"
	GoalTypeCalorieTarget   GoalType = "calorie_target"
	GoalTypeStepsPerDay     GoalType = "steps_per_day"
	GoalTypeWorkoutFreq     GoalType = "workout_frequency"
	GoalTypeSleepHours      GoalType = "sleep_hours"
	GoalTypeWaterIntake     GoalType = "water_intake"
	GoalTypeCustom          GoalType = "custom"
)

// GoalStatus represents the status of a goal
type GoalStatus string

const (
	GoalStatusActive    GoalStatus = "active"
	GoalStatusCompleted GoalStatus = "completed"
	GoalStatusPaused    GoalStatus = "paused"
	GoalStatusFailed    GoalStatus = "failed"
)

// Goal represents a user's health and fitness goal
type Goal struct {
	ID              int         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          int         `gorm:"not null;index" json:"user_id"`
	Type            GoalType    `gorm:"size:50;not null" json:"type"`
	Title           string      `gorm:"size:200;not null" json:"title"`
	Description     string      `gorm:"type:text" json:"description"`
	TargetValue     float64     `gorm:"not null" json:"target_value"`
	CurrentValue    float64     `gorm:"default:0" json:"current_value"`
	Unit            string      `gorm:"size:50" json:"unit"`
	StartDate       time.Time   `gorm:"not null" json:"start_date"`
	TargetDate      time.Time   `gorm:"not null" json:"target_date"`
	Status          GoalStatus  `gorm:"size:20;default:active" json:"status"`
	ProgressPercent float64     `gorm:"default:0" json:"progress_percent"`
	IsPublic        bool        `gorm:"default:false" json:"is_public"`
	CreatedAt       time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
	CompletedAt     *time.Time  `json:"completed_at,omitempty"`
}

// TableName specifies the table name for Goal model
func (Goal) TableName() string {
	return "goals"
}

// GoalProgress represents detailed progress tracking for a goal
type GoalProgress struct {
	ID              int       `gorm:"primaryKey;autoIncrement" json:"id"`
	GoalID          int       `gorm:"not null;index" json:"goal_id"`
	RecordedAt      time.Time `gorm:"not null" json:"recorded_at"`
	Value           float64   `gorm:"not null" json:"value"`
	Notes           string    `gorm:"type:text" json:"notes"`
	MeasurementType string    `gorm:"size:50" json:"measurement_type"` // automatic, manual
}

// TableName specifies the table name for GoalProgress model
func (GoalProgress) TableName() string {
	return "goal_progress"
}

// GoalReminder represents reminder settings for goals
type GoalReminder struct {
	ID              int       `gorm:"primaryKey;autoIncrement" json:"id"`
	GoalID          int       `gorm:"not null;index" json:"goal_id"`
	ReminderType    string    `gorm:"size:50;not null" json:"reminder_type"` // daily, weekly, milestone
	ReminderTime    string    `gorm:"size:50" json:"reminder_time"`
	LastSentAt      *time.Time `json:"last_sent_at"`
	IsActive        bool      `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName specifies the table name for GoalReminder model
func (GoalReminder) TableName() string {
	return "goal_reminders"
}

// GoalRepository defines operations for managing goals
type GoalRepository interface {
	Create(goal *Goal) error
	GetByID(id int) (*Goal, error)
	GetByUser(userID int) ([]Goal, error)
	GetActiveGoals(userID int) ([]Goal, error)
	Update(goal *Goal) error
	Delete(id int) error
	UpdateProgress(goalID int, value float64) error
	GetProgressHistory(goalID int) ([]GoalProgress, error)
	CreateReminder(reminder *GoalReminder) error
	GetReminders(goalID int) ([]GoalReminder, error)
}