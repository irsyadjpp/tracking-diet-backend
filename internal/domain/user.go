package domain

import "time"

type User struct {
    ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
    Name      string    `gorm:"size:100;not null" json:"name"`
    Email     string    `gorm:"size:100;unique;not null" json:"email"`
    Password  string    `gorm:"size:255;not null" json:"-"`
    Age       int       `json:"age"`
    Gender    string    `gorm:"size:10" json:"gender"`
    HeightCm  float64   `json:"height_cm"`
    WeightKg  float64   `json:"weight_kg"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
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
