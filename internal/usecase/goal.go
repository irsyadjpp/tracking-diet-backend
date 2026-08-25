package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
)

var (
	ErrGoalNotFound = errors.New("goal not found")
)

// GoalUsecase handles goal business logic
type GoalUsecase struct {
	goalRepo domain.GoalRepository
}

// NewGoalUsecase creates a new goal use case
func NewGoalUsecase(goalRepo domain.GoalRepository) *GoalUsecase {
	return &GoalUsecase{
		goalRepo: goalRepo,
	}
}

// CreateGoal creates a new goal with validation
func (uc *GoalUsecase) CreateGoal(userID int, goalType domain.GoalType, title, description string, targetValue float64, unit string, startDate, targetDate time.Time) (*domain.Goal, error) {
	// Validate dates
	if targetDate.Before(startDate) {
		return nil, errors.New("target date must be after start date")
	}

	// Validate target value based on goal type
	if err := uc.validateTargetValue(goalType, targetValue); err != nil {
		return nil, err
	}

	goal := &domain.Goal{
		UserID:          userID,
		Type:            goalType,
		Title:           title,
		Description:     description,
		TargetValue:     targetValue,
		CurrentValue:    0,
		Unit:            unit,
		StartDate:       startDate,
		TargetDate:      targetDate,
		Status:          domain.GoalStatusActive,
		ProgressPercent: 0,
		IsPublic:        false,
	}

	if err := uc.goalRepo.Create(goal); err != nil {
		return nil, fmt.Errorf("failed to create goal: %w", err)
	}

	return goal, nil
}

// GetGoal retrieves a goal by ID
func (uc *GoalUsecase) GetGoal(id int) (*domain.Goal, error) {
	goal, err := uc.goalRepo.GetByID(id)
	if err != nil {
		return nil, ErrGoalNotFound
	}
	return goal, nil
}

// GetUserGoals retrieves all goals for a user
func (uc *GoalUsecase) GetUserGoals(userID int) ([]domain.Goal, error) {
	goals, err := uc.goalRepo.GetByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user goals: %w", err)
	}
	return goals, nil
}

// GetActiveGoals retrieves active goals for a user
func (uc *GoalUsecase) GetActiveGoals(userID int) ([]domain.Goal, error) {
	goals, err := uc.goalRepo.GetActiveGoals(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active goals: %w", err)
	}
	return goals, nil
}

// UpdateGoal updates an existing goal
func (uc *GoalUsecase) UpdateGoal(id int, title, description *string, targetValue *float64, targetDate *time.Time, status *domain.GoalStatus) (*domain.Goal, error) {
	goal, err := uc.goalRepo.GetByID(id)
	if err != nil {
		return nil, ErrGoalNotFound
	}

	if title != "" {
		goal.Title = title
	}
	if description != nil {
		goal.Description = *description
	}
	if targetValue != nil {
		goal.TargetValue = *targetValue
	}
	if targetDate != nil {
		goal.TargetDate = *targetDate
	}
	if status != nil {
		goal.Status = *status
		if *status == domain.GoalStatusCompleted {
			now := time.Now()
			goal.CompletedAt = &now
		}
	}

	// Recalculate progress
	goal.ProgressPercent = uc.calculateProgress(goal)

	if err := uc.goalRepo.Update(goal); err != nil {
		return nil, fmt.Errorf("failed to update goal: %w", err)
	}

	return goal, nil
}

// DeleteGoal deletes a goal
func (uc *GoalUsecase) DeleteGoal(id int) error {
	if err := uc.goalRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete goal: %w", err)
	}
	return nil
}

// UpdateProgress updates the current value and progress of a goal
func (uc *GoalUsecase) UpdateProgress(goalID int, value float64) error {
	goal, err := uc.goalRepo.GetByID(goalID)
	if err != nil {
		return ErrGoalNotFound
	}

	goal.CurrentValue = value
	goal.ProgressPercent = uc.calculateProgress(goal)

	// Check if goal is completed
	if goal.ProgressPercent >= 100 {
		goal.Status = domain.GoalStatusCompleted
		now := time.Now()
		goal.CompletedAt = &now
	}

	if err := uc.goalRepo.Update(goalal); err != nil {
		return fmt.Errorf("failed to update goal: %w", err)
	}

	// Record progress history
	if err := uc.goalRepo.UpdateProgress(goalID, value); err != nil {
		// Log error but don't fail the main operation
		fmt.Printf("Failed to record progress history: %v\n", err)
	}

	return nil
}

// GetProgressHistory retrieves progress history for a goal
func (uc *GoalUsecase) GetProgressHistory(goalID int) ([]domain.GoalProgress, error) {
	history, err := uc.goalRepo.GetProgressHistory(goalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get progress history: %w", err)
	}
	return history, nil
}

// CreateReminder creates a reminder for a goal
func (uc *GoalUsecase) CreateReminder(goalID int, reminderType, reminderTime string) (*domain.GoalReminder, error) {
	reminder := &domain.GoalReminder{
		GoalID:       goalID,
		ReminderType: reminderType,
		ReminderTime: reminderTime,
		IsActive:     true,
	}

	if err := uc.goalRepo.CreateReminder(reminder); err != nil {
		return nil, fmt.Errorf("failed to create reminder: %w", err)
	}

	return reminder, nil
}

// GetReminders retrieves reminders for a goal
func (uc *GoalUsecase) GetReminders(goalID int) ([]domain.GoalReminder, error) {
	reminders, err := uc.goalRepo.GetReminders(goalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reminders: %w", err)
	}
	return reminders, nil
}

// Helper functions
func (uc *GoalUsecase) validateTargetValue(goalType domain.GoalType, targetValue float64) error {
	switch goalType {
	case domain.GoalTypeWeightLoss:
		if targetValue <= 0 || targetValue > 200 {
			return errors.New("invalid weight target")
		}
	case domain.GoalTypeWeightGain:
		if targetValue <= 0 || targetValue > 100 {
			return errors.New("invalid weight gain target")
		}
	case domain.GoalTypeMuscleGain:
		if targetValue <= 0 || targetValue > 50 {
			return errors.New("invalid muscle gain target")
		}
	case domain.GoalTypeCalorieTarget:
		if targetValue < 1000 || targetValue > 10000 {
			return errors.New("invalid calorie target")
		}
	case domain.GoalTypeStepsPerDay:
		if targetValue < 1000 || targetValue > 50000 {
			return errors.New("invalid steps target")
		}
	case domain.GoalTypeSleepHours:
		if targetValue < 4 || targetValue > 12 {
			return errors.New("invalid sleep hours target")
		}
	case domain.GoalTypeWaterIntake:
		if targetValue < 1 || targetValue > 10 {
			return errors.New("invalid water intake target")
		}
	}
	return nil
}

func (uc *GoalUsecase) calculateProgress(goal *domain.Goal) float64 {
	if goal.TargetValue == 0 {
		return 0
	}

	progress := (goal.CurrentValue / goal.TargetValue) * 100
	if progress > 100 {
		progress = 100
	}
	return progress
}