package usecase

import (
	"errors"
	"fmt"

	"github.com/irsyadjpp/tracking-diet-backend/internal/auth"
	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// UserUsecase handles user business logic
type UserUsecase struct {
	userRepo domain.UserRepository
}

// NewUserUsecase creates a new user use case
func NewUserUsecase(userRepo domain.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

// Register registers a new user
func (uc *UserUsecase) Register(fullName, email, password, gender string, birthDate *string, heightCm *float64) (*domain.User, error) {
	// Check if email already exists
	existingUser, err := uc.userRepo.GetByEmail(email)
	if err == nil && existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Validate password strength
	if err := auth.ValidatePasswordStrength(password); err != nil {
		return nil, err
	}

	// Create user
	user := &domain.User{
		FullName:     fullName,
		Email:        email,
		PasswordHash: password, // Will be hashed in repository
		Gender:       gender,
		Role:         domain.RoleUser,
	}

	if birthDate != nil {
		// Parse birth date if provided
		// Assuming birthDate is in "2006-01-02" format
		// Add proper date parsing logic here
	}

	if heightCm != nil {
		user.HeightCm = heightCm
	}

	if err := uc.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Login authenticates a user
func (uc *UserUsecase) Login(email, password string) (*domain.User, error) {
	user, err := uc.userRepo.GetByEmail(email)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Verify password
	if !auth.CheckPassword(password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (uc *UserUsecase) GetByID(id int) (*domain.User, error) {
	user, err := uc.userRepo.GetByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// UpdateProfile updates user profile information
func (uc *UserUsecase) UpdateProfile(id int, fullName, gender string, birthDate *string, heightCm *float64) (*domain.User, error) {
	user, err := uc.userRepo.GetByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if fullName != "" {
		user.FullName = fullName
	}

	if gender != "" {
		user.Gender = gender
	}

	if birthDate != nil {
		// Parse birth date if provided
		// Add proper date parsing logic here
	}

	if heightCm != nil {
		user.HeightCm = heightCm
	}

	if err := uc.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// Delete deletes a user (cascade delete will handle related records)
func (uc *UserUsecase) Delete(id int) error {
	if err := uc.userRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// List retrieves all users (admin only)
func (uc *UserUsecase) List() ([]domain.User, error) {
	users, err := uc.userRepo.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}