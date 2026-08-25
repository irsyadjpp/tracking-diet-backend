package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
)

var (
	ErrFitnessMeasurementNotFound = errors.New("fitness measurement not found")
)

// FitnessUsecase handles fitness measurement business logic
type FitnessUsecase struct {
	fitnessRepo domain.FitnessMeasurementRepository
	userRepo    domain.UserRepository
}

// NewFitnessUsecase creates a new fitness use case
func NewFitnessUsecase(fitnessRepo domain.FitnessMeasurementRepository, userRepo domain.UserRepository) *FitnessUsecase {
	return &FitnessUsecase{
		fitnessRepo: fitnessRepo,
		userRepo:    userRepo,
	}
}

// Create creates a new fitness measurement
func (uc *FitnessUsecase) Create(userID int, vo2MaxMlKgMin, strength1RmKg *float64, pushupCount, squatCount *int, enduranceNote, flexibilityNote *string, measuredAt time.Time) (*domain.FitnessMeasurement, error) {
	// Validate user exists
	_, err := uc.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	measurement := &domain.FitnessMeasurement{
		UserID:          userID,
		MeasuredAt:      measuredAt,
		VO2MaxMlKgMin:   vo2MaxMlKgMin,
		Strength1RmKg:   strength1RmKg,
		PushupCount:     pushupCount,
		SquatCount:      squatCount,
		EnduranceNote:   enduranceNote,
		FlexibilityNote: flexibilityNote,
	}

	if err := uc.fitnessRepo.Create(measurement); err != nil {
		return nil, fmt.Errorf("failed to create fitness measurement: %w", err)
	}

	return measurement, nil
}

// GetByID retrieves a fitness measurement by ID
func (uc *FitnessUsecase) GetByID(id int) (*domain.FitnessMeasurement, error) {
	measurement, err := uc.fitnessRepo.GetByID(id)
	if err != nil {
		return nil, ErrFitnessMeasurementNotFound
	}
	return measurement, nil
}

// ListByUser retrieves all fitness measurements for a user
func (uc *FitnessUsecase) ListByUser(userID int) ([]domain.FitnessMeasurement, error) {
	measurements, err := uc.fitnessRepo.ListByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list fitness measurements: %w", err)
	}
	return measurements, nil
}

// ListByUserWithDateRange retrieves fitness measurements for a user within a date range
func (uc *FitnessUsecase) ListByUserWithDateRange(userID int, startDate, endDate time.Time) ([]domain.FitnessMeasurement, error) {
	measurements, err := uc.fitnessRepo.ListByUserWithDateRange(userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to list fitness measurements: %w", err)
	}
	return measurements, nil
}

// ClassifyActivityLevel classifies activity level based on VO2Max
func (uc *FitnessUsecase) ClassifyActivityLevel(vo2Max *float64) string {
	if vo2Max == nil {
		return "Unknown"
	}

	vo2 := *vo2Max
	switch {
	case vo2Max < 20:
		return "Sedentary"
	case vo2Max < 35:
		return "Low"
	case vo2Max < 45:
		return "Moderate"
	case vo2Max < 55:
		return "High"
	default:
		return "Elite"
	}
}

// EstimateCalorieBurn estimates calorie burn based on workout data
func (uc *FitnessUsecase) EstimateCalorieBurn(measurement *domain.FitnessMeasurement, durationMinutes int) *float64 {
	if durationMinutes <= 0 {
		return nil
	}

	caloriesPerMinute := 8.0 // Base estimate
	if measurement.VO2MaxMlKgMin != nil {
		// Adjust based on VO2Max (higher VO2Max = more efficient calorie burn)
		caloriesPerMinute = *measurement.VO2MaxMlKgMin * 0.15
	}

	totalCalories := caloriesPerMinute * float64(durationMinutes)
	return &totalCalories
}

// TrackWorkoutConsistency tracks workout consistency over time
func (uc *FitnessUsecase) TrackWorkoutConsistency(userID int, days int) (int, []time.Time, error) {
	startDate := time.Now().AddDate(0, 0, -days)
	measurements, err := uc.fitnessRepo.ListByUserWithDateRange(userID, startDate, time.Now())
	if err != nil {
		return 0, nil, fmt.Errorf("failed to track workout consistency: %w", err)
	}

	workoutDates := make([]time.Time, 0, len(measurements))
	for _, m := range measurements {
		workoutDates = append(workoutDates, m.MeasuredAt)
	}

	return len(measurements), workoutDates, nil
}

// Update updates a fitness measurement
func (uc *FitnessUsecase) Update(id int, vo2MaxMlKgMin, strength1RmKg *float64, pushupCount, squatCount *int, enduranceNote, flexibilityNote *string, measuredAt time.Time) (*domain.FitnessMeasurement, error) {
	measurement, err := uc.fitnessRepo.GetByID(id)
	if err != nil {
		return nil, ErrFitnessMeasurementNotFound
	}

	// Update fields if provided
	if vo2MaxMlKgMin != nil {
		measurement.VO2MaxMlKgMin = vo2MaxMlKgMin
	}
	if strength1RmKg != nil {
		measurement.Strength1RmKg = strength1RmKg
	}
	if pushupCount != nil {
		measurement.PushupCount = pushupCount
	}
	if squatCount != nil {
		measurement.SquatCount = squatCount
	}
	if enduranceNote != nil {
		measurement.EnduranceNote = enduranceNote
	}
	if flexibilityNote != nil {
		measurement.FlexibilityNote = flexibilityNote
	}
	if !measuredAt.IsZero() {
		measurement.MeasuredAt = measuredAt
	}

	if err := uc.fitnessRepo.Update(measurement); err != nil {
		return nil, fmt.Errorf("failed to update fitness measurement: %w", err)
	}

	return measurement, nil
}

// Delete deletes a fitness measurement
func (uc *FitnessUsecase) Delete(id int) error {
	if err := uc.fitnessRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete fitness measurement: %w", err)
	}
	return nil
}