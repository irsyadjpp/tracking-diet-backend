package domain

import "time"

// NotificationType represents different types of notifications
type NotificationType string

const (
	NotificationTypeGoalReminder     NotificationType = "goal_reminder"
	NotificationTypeAchievement       NotificationType = "achievement"
	NotificationTypeMealReminder     NotificationType = "meal_reminder"
	NotificationTypeWeeklyReport    NotificationType = "weekly_report"
	NotificationTypeAlert            NotificationType = "alert"
	NotificationTypeReminder        NotificationType = "reminder"
)

// NotificationChannel represents notification delivery channels
type NotificationChannel string

const (
	NotificationChannelEmail NotificationChannel = "email"
	NotificationChannelPush  NotificationChannel = "push"
	NotificationChannelSMS   NotificationChannel = "sms"
	NotificationChannelInApp NotificationChannel = "in_app"
)

// NotificationPriority represents urgency level
type NotificationPriority string

const (
	NotificationPriorityLow    NotificationPriority = "low"
	NotificationPriorityMedium NotificationPriority = "medium"
	NotificationPriorityHigh   NotificationPriority = "high"
	NotificationPriorityUrgent NotificationPriority = "urgent"
)

// Notification represents a user notification
type Notification struct {
	ID          int                `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int                `gorm:"not null;index" json:"user_id"`
	Type        NotificationType    `gorm:"size:50;not null;index" json:"type"`
	Title       string             `gorm:"size:200;not null" json:"title"`
	Message     string             `gorm:"type:text;not null" json:"message"`
	Data        interface{}         `gorm:"type:jsonb" json:"data"`
	Channels    []NotificationChannel `gorm:"type:jsonb" json:"channels"`
	Priority    NotificationPriority `gorm:"size:20;default:medium" json:"priority"`
	ReadAt      *time.Time          `json:"read_at"`
	ExpiresAt   *time.Time          `json:"expires_at"`
	SentAt      time.Time          `gorm:"autoCreateTime" json:"sent_at"`
	ActionURL   *string            `json:"action_url"`
	ActionLabel *string            `json:"action_label"`
}

// TableName specifies the table name for Notification model
func (Notification) TableName() string {
	return "notifications"
}

// NotificationPreference represents user's notification preferences
type NotificationPreference struct {
	ID                      int                `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID                  int                `gorm:"not null;unique" json:"user_id"`
	EmailEnabled            bool               `gorm:"default:true" json:"email_enabled"`
	PushEnabled             bool               `gorm:"default:true" json:"push_enabled"`
	SMSEnabled              bool               `gorm:"default:false" json:"sms_enabled"`
	InAppEnabled            bool               `default:true" json:"in_app_enabled"`
	QuietHoursStart         *int               `gorm:"column:quiet_hours_start" json:"quiet_hours_start"`
	QuietHoursEnd           *int               `gorm:"column:quiet_hours_end" json:"quiet_hours_end"`
	GoalRemindersEnabled    bool               `gorm:"default:true" json:"goal_reminders_enabled"`
	MealRemindersEnabled    bool               `gorm:"default:true" json:"meal_reminders_enabled"`
	WeeklyReportsEnabled    bool               `gorm:"default:true" json:"weekly_reports_enabled"`
	CreatedAt               time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt               time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for NotificationPreference model
func (NotificationPreference) TableName() string {
	return "notification_preferences"
}

// NotificationRepository defines operations for managing notifications
type NotificationRepository interface {
	Create(notification *Notification) error
	GetByID(id int) (*Notification, error)
	GetByUser(userID int) ([]Notification, error)
	GetUnread(userID int) ([]Notification, error)
	MarkAsRead(id int) error
	Delete(id int) error
	GetPreferences(userID int) (*NotificationPreference, error)
	UpdatePreferences(pref *NotificationPreference) error
	CreateDefaultPreferences(userID int) error
}