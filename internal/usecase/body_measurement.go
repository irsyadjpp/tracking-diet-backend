package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
)

var (
	ErrBodyMeasurementNotFound = errors.New("body measurement not found")
	ErrInvalidMeasurement      = errors.New("invalid measurement value")
)

// BodyMeasurementUsecase handles body measurement business logic
type BodyMeasurementUsecase struct {
	bodyRepo domain.BodyMeasurementRepository
	userRepo domain.UserRepository
}

// NewBodyMeasurementUsecase creates a new body measurement use case
func NewBodyMeasurementUsecase(bodyRepo domain.BodyMeasurementRepository, userRepo domain.UserRepository) *BodyMeasurementUsecase {
	return &BodyMeasurementUsecase{
		bodyRepo: bodyRepo,
		userRepo: userRepo,
	}
}

// Create creates a new body measurement with automatic calculations
func (uc *BodyMeasurementUsecase) Create(userID int, weightKg, waistCm, hipCm, chestCm, thighCm, armCm, bodyFatPct, muscleMassKg, visceralFat, skinfoldMm *float64, measuredAt time.Time) (*domain.BodyMeasurement, error) {
	// Validate user exists
	_, err := uc.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	measurement := &domain.BodyMeasurement{
		UserID:        userID,
		MeasuredAt:    measuredAt,
		WeightKg:      weightKg,
		WaistCm:       waistCm,
		HipCm:         hipCm,
		ChestCm:       chestCm,
		ThighCm:       thighCm,
		ArmCm:         armCm,
		BodyFatPct:    bodyFatPct,
		MuscleMassKg:  muscleMassKg,
		VisceralFat:   visceralFat,
		SkinfoldMm:    skinfoldMm,
	}

	// Calculate BMI if weight is provided
	if weightKg != nil {
		user, err := uc.userRepo.GetByID(userID)
		if err == nil && user.HeightCm != nil {
			heightM := *user.HeightCm / 100.0
			bmi := *weightKg / (heightM * heightM)
			measurement.BMI = &bmi
		}
	}

	// Calculate waist-hip ratio if both values are provided
	if waistCm != nil && hipCm != nil && *hipCm > 0 {
		whr := *waistCm / *hipCm
		measurement.WHR = &whr
	}

	if err := uc.bodyRepo.Create(measurement); err != nil {
		return nil, fmt.Errorf("failed to create body measurement: %w", err)
	}

	return measurement, nil
}

// GetByID retrieves a body measurement by ID
func (uc *BodyMeasurementUsecase) GetByID(id int) (*domain.BodyMeasurement, error) {
	measurement, err := uc.bodyRepo.GetByID(id)
	if err != nil {
		return nil, ErrBodyMeasurementNotFound
	}
	return measurement, nil
}

// ListByUser retrieves all body measurements for a user
func (uc *BodyMeasurementUsecase) ListByUser(userID int) ([]domain.BodyMeasurement, error) {
	measurements, err := uc.bodyRepo.ListByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list body measurements: %w", err)
	}
	return measurements, nil
}

// ListByUserWithDateRange retrieves body measurements for a user within a date range
func (uc *BodyMeasurementUsecase) ListByUserWithDateRange(userID int, startDate, endDate time.Time) ([]domain.BodyMeasurement, error) {
	measurements, err := uc.bodyRepo.ListByUserWithDateRange(userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to list body measurements: %w", err)
	}
	return measurements, nil
}

// Update updates a body measurement
func (uc *BodyMeasurementUsecase) Update(id int, weightKg, waistCm, hipCm, chestCm, thighCm, armCm, bodyFatPct, muscleMassKg, visceralFat, skinfoldMm *float64, measuredAt time.Time) (*domain.BodyMeasurement, error) {
	measurement, err := uc.bodyRepo.GetByID(id)
	if err != nil {
		return nil, ErrBodyMeasurementNotFound
	}

	// Update fields if provided
	if weightKg != nil {
		measurement.WeightKg = weightKg
	}
	if waistCm != nil {
		measurement.WaistCm = waistCm
	}
	if hipCm != nil {
		measurement.HipCm = hipCm
	}
	if chestCm != nil {
		measurement.ChestCm = chestCm
	}
	if thighCm != nil {
		measurement.ThighCm = thighCm
	}
	if armCm != nil {
		measurement.ArmCm = armCm
	}
	if bodyFatPct != nil {
		measurement.BodyFatPct = bodyFatPct
	}
	if muscleMassKg != nil {
		measurement.MuscleMassKg = muscleMassKg
	}
	if visceralFat != nil {
		measurement.VisceralFat = visceralFat
	}
	if skinfoldMm != nil {
		measurement.SkinfoldMm = skinfoldMm
	}
	if !measuredAt.IsZero() {
		measurement.MeasuredAt = measuredAt
	}

	// Recalculate BMI if weight changed
	if weightKg != nil {
		user, err := uc.userRepo.GetByID(measurement.UserID)
		if err == nil && user.HeightCm != nil {
			heightM := *user.HeightCm / 100.0
			bmi := *weightKg / (heightM * heightM)
			measurement.BMI = &bmi
		}
	}

	// Recalculate waist-hip ratio if values changed
	if waistCm != nil || hipCm != nil {
		waist := measurement.WaistCm
		hip := measurement.HipCm
		if waistCm != nil {
			waist = waistCm
		}
		if hipCm != nil {
			hip = hipCm
		}
		if waist != nil && hip != nil && *hip > 0 {
			whr := *waist / *hip
			measurement.WHR = &whr
		}
	}

	if err := uc.bodyRepo.Update(measurement); err != nil {
		return nil, fmt.Errorf("failed to update body measurement: %w", err)
	}

	return measurement, nil
}

// Delete deletes a body measurement
func (uc *BodyMeasurementUsecase) Delete(id int) error {
	if err := uc.bodyRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete body measurement: %w", err)
	}
	return nil
}

// GetBMIStatus returns BMI status category
func (uc *BodyMeasurementUsecase) GetBMIStatus(bmi float64) string {
	switch {
	case bmi < 18.5:
		return "Underweight"
	case bmi < 25:
		return "Normal"
	case bmi < 30:
		return "Overweight"
	default:
		return "Obese"
	}
}

// GetWHRRisk returns waist-hip ratio risk assessment
func (uc *BodyMeasurementUsecase) GetWHRRisk(whr float64, gender string) string {
	// Different thresholds for men and women
	if gender == "M" {
		switch {
		case whr < 0.9:
			return "Low risk"
		case whr < 1.0:
			return "Moderate risk"
		default:
			return "High risk"
		}
	} else {
		switch {
		case whr < 0.8:
			return "Low risk"
		case whr < 0.85:
			return "Moderate risk"
		default:
			return "High risk"
		}
	}
}