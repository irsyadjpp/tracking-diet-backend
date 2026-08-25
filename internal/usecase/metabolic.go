package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
)

var (
	ErrMetabolicMeasurementNotFound = errors.New("metabolic measurement not found")
)

// MetabolicUsecase handles metabolic measurement business logic
type MetabolicUsecase struct {
	metabolicRepo domain.MetabolicMeasurementRepository
	userRepo      domain.UserRepository
}

// NewMetabolicUsecase creates a new metabolic use case
func NewMetabolicUsecase(metabolicRepo domain.MetabolicMeasurementRepository, userRepo domain.UserRepository) *MetabolicUsecase {
	return &MetabolicUsecase{
		metabolicRepo: metabolicRepo,
		userRepo:      userRepo,
	}
}

// Create creates a new metabolic measurement
func (uc *MetabolicUsecase) Create(userID int, systolicBP, diastolicBP *int, fastingGlucose, hba1cPct, cholesterolTotal, ldl, hdl, triglycerides, uricAcid *float64, liverFunctionNote, kidneyFunctionNote *string, measuredAt time.Time) (*domain.MetabolicMeasurement, error) {
	// Validate user exists
	_, err := uc.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	measurement := &domain.MetabolicMeasurement{
		UserID:              userID,
		MeasuredAt:          measuredAt,
		SystolicBP:          systolicBP,
		DiastolicBP:         diastolicBP,
		FastingGlucose:      fastingGlucose,
		Hba1cPct:            hba1cPct,
		CholesterolTotal:    cholesterolTotal,
		LDL:                 ldl,
		HDL:                 hdl,
		Triglycerides:       triglycerides,
		UricAcid:            uricAcid,
		LiverFunctionNote:   liverFunctionNote,
		KidneyFunctionNote:  kidneyFunctionNote,
	}

	if err := uc.metabolicRepo.Create(measurement); err != nil {
		return nil, fmt.Errorf("failed to create metabolic measurement: %w", err)
	}

	return measurement, nil
}

// GetByID retrieves a metabolic measurement by ID
func (uc *MetabolicUsecase) GetByID(id int) (*domain.MetabolicMeasurement, error) {
	measurement, err := uc.metabolicRepo.GetByID(id)
	if err != nil {
		return nil, ErrMetabolicMeasurementNotFound
	}
	return measurement, nil
}

// ListByUser retrieves all metabolic measurements for a user
func (uc *MetabolicUsecase) ListByUser(userID int) ([]domain.MetabolicMeasurement, error) {
	measurements, err := uc.metabolicRepo.ListByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list metabolic measurements: %w", err)
	}
	return measurements, nil
}

// ListByUserWithDateRange retrieves metabolic measurements for a user within a date range
func (uc *MetabolicUsecase) ListByUserWithDateRange(userID int, startDate, endDate time.Time) ([]domain.MetabolicMeasurement, error) {
	measurements, err := uc.metabolicRepo.ListByUserWithDateRange(userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to list metabolic measurements: %w", err)
	}
	return measurements, nil
}

// CalculateHealthScore calculates a health score based on metabolic metrics
func (uc *MetabolicUsecase) CalculateHealthScore(measurement *domain.MetabolicMeasurement) (float64, []string) {
	score := 100.0
	issues := []string{}

	// Blood pressure assessment
	if measurement.SystolicBP != nil && measurement.DiastolicBP != nil {
		if *measurement.SystolicBP >= 140 || *measurement.DiastolicBP >= 90 {
			score -= 15
			issues = append(issues, "High blood pressure")
		} else if *measurement.SystolicBP >= 120 || *measurement.DiastolicBP >= 80 {
			score -= 5
			issues = append(issues, "Elevated blood pressure")
		}
	}

	// Blood glucose assessment
	if measurement.FastingGlucose != nil {
		if *measurement.FastingGlucose >= 126 {
			score -= 20
			issues = append(issues, "High fasting glucose (diabetic range)")
		} else if *measurement.FastingGlucose >= 100 {
			score -= 10
			issues = append(issues, "Elevated fasting glucose (prediabetic range)")
		}
	}

	// Cholesterol assessment
	if measurement.LDL != nil {
		if *measurement.LDL >= 160 {
			score -= 15
			issues = append(issues, "High LDL cholesterol")
		} else if *measurement.LDL >= 130 {
			score -= 5
			issues = append(issues, "Borderline high LDL cholesterol")
		}
	}

	if measurement.HDL != nil {
		if *measurement.HDL < 40 {
			score -= 10
			issues = append(issues, "Low HDL cholesterol")
		}
	}

	// Ensure score doesn't go below 0
	if score < 0 {
		score = 0
	}

	return score, issues
}

// DetectAbnormalValues detects abnormal metabolic values
func (uc *MetabolicUsecase) DetectAbnormalValues(measurement *domain.MetabolicMeasurement) map[string]string {
	abnormal := map[string]string{}

	if measurement.SystolicBP != nil && *measurement.SystolicBP >= 140 {
		abnormal["systolic_bp"] = "High"
	}
	if measurement.DiastolicBP != nil && *measurement.DiastolicBP >= 90 {
		abnormal["diastolic_bp"] = "High"
	}
	if measurement.FastingGlucose != nil && *measurement.FastingGlucose >= 126 {
		abnormal["fasting_glucose"] = "High"
	}
	if measurement.Hba1cPct != nil && *measurement.Hba1cPct >= 6.5 {
		abnormal["hba1c"] = "High"
	}
	if measurement.LDL != nil && *measurement.LDL >= 160 {
		abnormal["ldl"] = "High"
	}
	if measurement.HDL != nil && *measurement.HDL < 40 {
		abnormal["hdl"] = "Low"
	}
	if measurement.Triglycerides != nil && *measurement.Triglycerides >= 200 {
		abnormal["triglycerides"] = "High"
	}
	if measurement.UricAcid != nil && *measurement.UricAcid >= 7.0 {
		abnormal["uric_acid"] = "High"
	}

	return abnormal
}

// Update updates a metabolic measurement
func (uc *MetabolicUsecase) Update(id int, systolicBP, diastolicBP *int, fastingGlucose, hba1cPct, cholesterolTotal, ldl, hdl, triglycerides, uricAcid *float64, liverFunctionNote, kidneyFunctionNote *string, measuredAt time.Time) (*domain.MetabolicMeasurement, error) {
	measurement, err := uc.metabolicRepo.GetByID(id)
	if err != nil {
		return nil, ErrMetabolicMeasurementNotFound
	}

	// Update fields if provided
	if systolicBP != nil {
		measurement.SystolicBP = systolicBP
	}
	if diastolicBP != nil {
		measurement.DiastolicBP = diastolicBP
	}
	if fastingGlucose != nil {
		measurement.FastingGlucose = fastingGlucose
	}
	if hba1cPct != nil {
		measurement.Hba1cPct = hba1cPct
	}
	if cholesterolTotal != nil {
		measurement.CholesterolTotal = cholesterolTotal
	}
	if ldl != nil {
		measurement.LDL = ldl
	}
	if hdl != nil {
		measurement.HDL = hdl
	}
	if triglycerides != nil {
		measurement.Triglycerides = triglycerides
	}
	if uricAcid != nil {
		measurement.UricAcid = uricAcid
	}
	if liverFunctionNote != nil {
		measurement.LiverFunctionNote = liverFunctionNote
	}
	if kidneyFunctionNote != nil {
		measurement.KidneyFunctionNote = kidneyFunctionNote
	}
	if !measuredAt.IsZero() {
		measurement.MeasuredAt = measuredAt
	}

	if err := uc.metabolicRepo.Update(measurement); err != nil {
		return nil, fmt.Errorf("failed to update metabolic measurement: %w", err)
	}

	return measurement, nil
}

// Delete deletes a metabolic measurement
func (uc *MetabolicUsecase) Delete(id int) error {
	if err := uc.metabolicRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete metabolic measurement: %w", err)
	}
	return nil
}