package domain

import "time"

// AchievementType represents different types of achievements
type AchievementType string

const (
	AchievementTypeStreak         AchievementType = "streak"
	AchievementTypeWeightLoss     AchievementType = "weight_loss"
	AchievementTypeWorkout        AchievementType = "workout"
	AchievementTypeNutrition      AchievementType = "nutrition"
	AchievementTypeConsistency    AchievementType = "consistency"
	AchievementTypeMilestone      AchievementType = "milestone"
)

// Achievement represents a badge or achievement
type Achievement struct {
	ID          int             `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string          `gorm:"size:100;not null" json:"name"`
	Description string          `gorm:"type:text" json:"description"`
	Type        AchievementType `gorm:"size:50;not null" json:"type"`
	Icon        string          `gorm:"size:200" json:"icon"`
	Points      int             `gorm:"default:0" json:"points"`
	Tier        string          `gorm:"size:20;default:bronze" json:"tier"` // bronze, silver, gold, platinum
	IsPublic    bool            `gorm:"default:true" json:"is_public"`
	CreatedAt   time.Time       `gorm:"autoCreateTime" json:"created_at"`
}

// TableName specifies the table name for Achievement model
func (Achievement) TableName() string {
	return "achievements"
}

// UserAchievement represents an achievement earned by a user
type UserAchievement struct {
	ID             int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         int       `gorm:"not null;index" json:"user_id"`
	AchievementID  int       `gorm:"not null;index" json:"achievement_id"`
	EarnedAt       time.Time `gorm:"autoCreateTime" json:"earned_at"`
	Progress       float64   `gorm:"default:0" json:"progress"` // for achievements with progress tracking
	IsCompleted    bool      `gorm:"default:true" json:"is_completed"`
}

// TableName specifies the table name for UserAchievement model
func (UserAchievement) TableName() string {
	return "user_achievements"
}

// Leaderboard represents a leaderboard entry
type Leaderboard struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int       `gorm:"not null;unique" json:"user_id"`
	Username    string    `gorm:"size:100" json:"username"`
	TotalPoints int       `gorm:"default:0" json:"total_points"`
	Rank        int       `gorm:"default:0" json:"rank"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for Leaderboard model
func (Leaderboard) TableName() string {
	return "leaderboard"
}

// Challenge represents a user challenge
type Challenge struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"size:200;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Type        string    `gorm:"size:50;not null" json:"type"` // individual, group, global
	StartDate   time.Time `gorm:"not null" json:"start_date"`
	EndDate     time.Time `gorm:"not null" json:"end_date"`
	TargetValue float64   `gorm:"not null" json:"target_value"`
	Points      int       `gorm:"default:0" json:"points"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName specifies the table name for Challenge model
func (Challenge) TableName() string {
	return "challenges"
}

// UserChallenge represents a user's participation in a challenge
type UserChallenge struct {
	ID            int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        int       `gorm:"not null;index" json:"user_id"`
	ChallengeID   int       `gorm:"not null;index" json:"challenge_id"`
	CurrentValue  float64   `gorm:"default:0" json:"current_value"`
	Progress      float64   `gorm:"default:0" json:"progress"`
	JoinedAt      time.Time `gorm:"autoCreateTime" json:"joined_at"`
	CompletedAt   *time.Time `json:"completed_at"`
}

// TableName specifies the table name for UserChallenge model
func (UserChallenge) TableName() string {
	return "user_challenges"
}

// Badge represents a badge that can be displayed on profile
type Badge struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	ImageURL    string    `gorm:"size:500" json:"image_url"`
	Category    string    `gorm:"size:50" json:"category"`
	Rarity      string    `gorm:"size:20;default:common" json:"rarity"` // common, rare, epic, legendary
	IsPublic    bool      `gorm:"default:true" json:"is_public"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName specifies the table name for Badge model
func (Badge) TableName() string {
	return "badges"
}

// UserBadge represents a badge earned by a user
type UserBadge struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int       `gorm:"not null;index" json:"user_id"`
	BadgeID   int       `gorm:"not null;index" json:"badge_id"`
	EarnedAt  time.Time `gorm:"autoCreateTime" json:"earned_at"`
	IsEquipped bool     `gorm:"default:false" json:"is_equipped"`
}

// TableName specifies the table name for UserBadge model
func (UserBadge) TableName() string {
	return "user_badges"
}

// AchievementRepository defines operations for managing achievements
type AchievementRepository interface {
	Create(achievement *Achievement) error
	GetByID(id int) (*Achievement, error)
	ListAll() ([]Achievement, error)
	ListByType(achievementType AchievementType) ([]Achievement, error)
	Update(achievement *Achievement) error
	Delete(id int) error
	GrantToUser(userID, achievementID int) error
	GetUserAchievements(userID int) ([]UserAchievement, error)
	GetLeaderboard(limit int) ([]Leaderboard, error)
	UpdateLeaderboard(userID int, points int) error
}