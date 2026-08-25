package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
)

var (
	ErrNutritionMeasurementNotFound = errors.New("nutrition measurement not found")
)

// NutritionUsecase handles nutrition measurement business logic
type NutritionUsecase struct {
	nutritionRepo domain.NutritionMeasurementRepository
	userRepo      domain.UserRepository
}

// NewNutritionUsecase creates a new nutrition use case
func NewNutritionUsecase(nutritionRepo domain.NutritionMeasurementRepository, userRepo domain.UserRepository) *NutritionUsecase {
	return &NutritionUsecase{
		nutritionRepo: nutritionRepo,
		userRepo:      userRepo,
	}
}

// Create creates a new nutrition measurement
func (uc *NutritionUsecase) Create(userID int, caloriesInKcal *int, carbsG, proteinG, fatG, fiberG, waterIntakeL *float64, caloriesOutKcal, stepCount *int, mealTimingNote *string, measuredAt time.Time) (*domain.NutritionMeasurement, error) {
	// Validate user exists
	_, err := uc.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	measurement := &domain.NutritionMeasurement{
		UserID:         userID,
		MeasuredAt:     measuredAt,
		CaloriesInKcal: caloriesInKcal,
		CarbsG:         carbsG,
		ProteinG:       proteinG,
		FatG:           fatG,
		FiberG:         fiberG,
		WaterIntakeL:   waterIntakeL,
		CaloriesOutKcal: caloriesOutKcal,
		StepCount:      stepCount,
		MealTimingNote: mealTimingNote,
	}

	if err := uc.nutritionRepo.Create(measurement); err != nil {
		return nil, fmt.Errorf("failed to create nutrition measurement: %w", err)
	}

	return measurement, nil
}

// GetByID retrieves a nutrition measurement by ID
func (uc *NutritionUsecase) GetByID(id int) (*domain.NutritionMeasurement, error) {
	measurement, err := uc.nutritionRepo.GetByID(id)
	if err != nil {
		return nil, ErrNutritionMeasurementNotFound
	}
	return measurement, nil
}

// ListByUser retrieves all nutrition measurements for a user
func (uc *NutritionUsecase) ListByUser(userID int) ([]domain.NutritionMeasurement, error) {
	measurements, err := uc.nutritionRepo.ListByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list nutrition measurements: %w", err)
	}
	return measurements, nil
}

// ListByUserWithDateRange retrieves nutrition measurements for a user within a date range
func (uc *NutritionUsecase) ListByUserWithDateRange(userID int, startDate, endDate time.Time) ([]domain.NutritionMeasurement, error) {
	measurements, err := uc.nutritionRepo.ListByUserWithDateRange(userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to list nutrition measurements: %w", err)
	}
	return measurements, nil
}

// GetDailySummary retrieves daily nutrition summary for a user
func (uc *NutritionUsecase) GetDailySummary(userID int, date time.Time) (*domain.NutritionMeasurement, error) {
	summary, err := uc.nutritionRepo.GetDailySummary(userID, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily summary: %w", err)
	}
	return summary, nil
}

// CalculateCalorieBalance calculates net calorie balance (in - out)
func (uc *NutritionUsecase) CalculateCalorieBalance(measurement *domain.NutritionMeasurement) *int {
	if measurement.CaloriesInKcal == nil || measurement.CaloriesOutKcal == nil {
		return nil
	}
	balance := *measurement.CaloriesInKcal - *measurement.CaloriesOutKcal
	return &balance
}

// CalculateMacroRatios calculates macronutrient ratios
func (uc *NutritionUsecase) CalculateMacroRatios(measurement *domain.NutritionMeasurement) map[string]float64 {
	if measurement.CarbsG == nil || measurement.ProteinG == nil || measurement.FatG == nil {
		return nil
	}

	carbs := *measurement.CarbsG * 4  // 4 cal/g
	protein := *measurement.ProteinG * 4  // 4 cal/g
	fat := *measurement.FatG * 9  // 9 cal/g
	total := carbs + protein + fat

	if total == 0 {
		return nil
	}

	return map[string]float64{
		"carbs_pct":   (carbs / total) * 100,
		"protein_pct": (protein / total) * 100,
		"fat_pct":     (fat / total) * 100,
	}
}

// Update updates a nutrition measurement
func (uc *NutritionUsecase) Update(id int, caloriesInKcal *int, carbsG, proteinG, fatG, fiberG, waterIntakeL *float64, caloriesOutKcal, stepCount *int, mealTimingNote *string, measuredAt time.Time) (*domain.NutritionMeasurement, error) {
	measurement, err := uc.nutritionRepo.GetByID(id)
	if err != nil {
		return nil, ErrNutritionMeasurementNotFound
	}

	// Update fields if provided
	if caloriesInKcal != nil {
		measurement.CaloriesInKcal = caloriesInKcal
	}
	if carbsG != nil {
		measurement.CarbsG = carbsG
	}
	if proteinG != nil {
		measurement.ProteinG = proteinG
	}
	if fatG != nil {
		measurement.FatG = fatG
	}
	if fiberG != nil {
		measurement.FiberG = fiberG
	}
	if waterIntakeL != nil {
		measurement.WaterIntakeL = waterIntakeL
	}
	if caloriesOutKcal != nil {
		measurement.CaloriesOutKcal = caloriesOutKcal
	}
	if stepCount != nil {
		measurement.StepCount = stepCount
	}
	if mealTimingNote != nil {
		measurement.MealTimingNote = mealTimingNote
	}
	if !measuredAt.IsZero() {
		measurement.MeasuredAt = measuredAt
	}

	if err := uc.nutritionRepo.Update(measurement); err != nil {
		return nil, fmt.Errorf("failed to update nutrition measurement: %w", err)
	}

	return measurement, nil
}

// Delete deletes a nutrition measurement
func (uc *NutritionUsecase) Delete(id int) error {
	if err := uc.nutritionRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete nutrition measurement: %w", err)
	}
	return nil
}