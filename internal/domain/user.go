package domain

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

type User struct {
	ID           int       `gorm:"primaryKey;autoIncrement;column:user_id" json:"id"`
	FullName     string    `gorm:"size:100;not null;column:full_name" json:"full_name"`
	Email        string    `gorm:"size:120;unique;not null;column:email" json:"email"`
	PasswordHash string    `gorm:"size:255;not null;column:password_hash" json:"-"`
	Gender       string    `gorm:"size:1;column:gender" json:"gender"`
	BirthDate    *time.Time `gorm:"column:birth_date" json:"birth_date"`
	HeightCm     *float64  `gorm:"column:height_cm" json:"height_cm"`
	Role         UserRole  `gorm:"size:20;default:user" json:"role"`
	CreatedAt    time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

// UserRepository defines operations for managing users.
type UserRepository interface {
    Create(u *User) error
    GetByID(id int) (*User, error)
    GetByEmail(email string) (*User, error)
    List() ([]User, error)
    Update(u *User) error
    Delete(id int) error
}
