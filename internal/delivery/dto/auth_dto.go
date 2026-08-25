package dto

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// RegisterRequest represents user registration request
type RegisterRequest struct {
	FullName    string `json:"full_name" validate:"required,min=2,max=100"`
	Email       string `json:"email" validate:"required,email,max=120"`
	Password    string `json:"password" validate:"required,min=8"`
	Gender      string `json:"gender" validate:"omitempty,oneof=M F"`
	BirthDate   string `json:"birth_date" validate:"omitempty"`
	HeightCm    *float64 `json:"height_cm" validate:"omitempty,gt=0"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         UserResponse `json:"user"`
}

// UserResponse represents user response
type UserResponse struct {
	ID        int        `json:"id"`
	FullName  string     `json:"full_name"`
	Email     string     `json:"email"`
	Gender    string     `json:"gender"`
	BirthDate *time.Time `json:"birth_date"`
	HeightCm  *float64   `json:"height_cm"`
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
}

// UpdateUserRequest represents user update request
type UpdateUserRequest struct {
	FullName  string     `json:"full_name" validate:"omitempty,min=2,max=100"`
	Gender    string     `json:"gender" validate:"omitempty,oneof=M F"`
	BirthDate string     `json:"birth_date" validate:"omitempty"`
	HeightCm  *float64   `json:"height_cm" validate:"omitempty,gt=0"`
}

// Validate validates the struct using validator
func (r *RegisterRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *LoginRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *RefreshTokenRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *UpdateUserRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}