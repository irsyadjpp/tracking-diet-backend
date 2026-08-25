package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
)

var (
	ErrWellbeingMeasurementNotFound = errors.New("wellbeing measurement not found")
)

// WellbeingUsecase handles wellbeing measurement business logic
type WellbeingUsecase struct {
	wellbeingRepo domain.WellbeingMeasurementRepository
	userRepo      domain.UserRepository
}

// NewWellbeingUsecase creates a new wellbeing use case
func NewWellbeingUsecase(wellbeingRepo domain.WellbeingMeasurementRepository, userRepo domain.UserRepository) *WellbeingUsecase {
	return &WellbeingUsecase{
		wellbeingRepo: wellbeingRepo,
		userRepo:      userRepo,
	}
}

// Create creates a new wellbeing measurement
func (uc *WellbeingUsecase) Create(userID int, energyLevel, moodLevel, hungerLevel, sleepQuality, stressLevel *int, sleepHours *float64, digestionNote *string, measuredAt time.Time) (*domain.WellbeingMeasurement, error) {
	// Validate user exists
	_, err := uc.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Validate ranges
	if energyLevel != nil && (*energyLevel < 1 || *energyLevel > 10) {
		return nil, fmt.Errorf("energy level must be between 1 and 10")
	}
	if moodLevel != nil && (*moodLevel < 1 || *moodLevel > 10) {
		return nil, fmt.Errorf("mood level must be between 1 and 10")
	}
	if hungerLevel != nil && (*hungerLevel < 1 || *hungerLevel > 10) {
		return nil, fmt.Errorf("hunger level must be between 1 and 10")
	}
	if sleepQuality != nil && (*sleepQuality < 1 || *sleepQuality > 10) {
		return nil, fmt.Errorf("sleep quality must be between 1 and 10")
	}
	if stressLevel != nil && (*stressLevel < 1 || *stressLevel > 10) {
		return nil, fmt.Errorf("stress level must be between 1 and 10")
	}

	measurement := &domain.WellbeingMeasurement{
		UserID:        userID,
		MeasuredAt:    measuredAt,
		EnergyLevel:   energyLevel,
		MoodLevel:     moodLevel,
		HungerLevel:   hungerLevel,
		SleepHours:    sleepHours,
		SleepQuality:  sleepQuality,
		DigestionNote: digestionNote,
		StressLevel:   stressLevel,
	}

	if err := uc.wellbeingRepo.Create(measurement); err != nil {
		return nil, fmt.Errorf("failed to create wellbeing measurement: %w", err)
	}

	return measurement, nil
}

// GetByID retrieves a wellbeing measurement by ID
func (uc *WellbeingUsecase) GetByID(id int) (*domain.WellbeingMeasurement, error) {
	measurement, err := uc.wellbeingRepo.GetByID(id)
	if err != nil {
		return nil, ErrWellbeingMeasurementNotFound
	}
	return measurement, nil
}

// ListByUser retrieves all wellbeing measurements for a user
func (uc *WellbeingUsecase) ListByUser(userID int) ([]domain.WellbeingMeasurement, error) {
	measurements, err := uc.wellbeingRepo.ListByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list wellbeing measurements: %w", err)
	}
	return measurements, nil
}

// ListByUserWithDateRange retrieves wellbeing measurements for a user within a date range
func (uc *WellbeingUsecase) ListByUserWithDateRange(userID int, startDate, endDate time.Time) ([]domain.WellbeingMeasurement, error) {
	measurements, err := uc.wellbeingRepo.ListByUserWithDateRange(userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to list wellbeing measurements: %w", err)
	}
	return measurements, nil
}

// AnalyzeMoodStressCorrelation analyzes correlation between mood and stress levels
func (uc *WellbeingUsecase) AnalyzeMoodStressCorrelation(userID int, days int) (float64, error) {
	startDate := time.Now().AddDate(0, 0, -days)
	measurements, err := uc.wellbeingRepo.ListByUserWithDateRange(userID, startDate, time.Now())
	if err != nil {
		return 0, fmt.Errorf("failed to analyze mood-stress correlation: %w", err)
	}

	if len(measurements) < 2 {
		return 0, nil // Not enough data
	}

	// Calculate correlation (simplified)
	// In a real implementation, use proper statistical correlation formula
	var moodSum, stressSum, moodStressSum int
	count := 0

	for _, m := range measurements {
		if m.MoodLevel != nil && m.StressLevel != nil {
			moodSum += *m.MoodLevel
			stressSum += *m.StressLevel
			moodStressSum += (*m.MoodLevel) * (*m.StressLevel)
			count++
		}
	}

	if count == 0 {
		return 0, nil
	}

	avgMood := float64(moodSum) / float64(count)
	avgStress := float64(stressSum) / float64(count)
	avgMoodStress := float64(moodStressSum) / float64(count)

	// Simplified correlation calculation
	correlation := (avgMoodStress - avgMood*avgStress) / (avgMood * avgStress)
	
	// Normalize to -1 to 1 range
	if correlation > 1 {
		correlation = 1
	} else if correlation < -1 {
		correlation = -1
	}

	return correlation, nil
}

// CalculateSleepQualityScore calculates sleep quality score
func (uc *WellbeingUsecase) CalculateSleepQualityScore(measurement *domain.WellbeingMeasurement) (int, string) {
	score := 0
	status := "Unknown"

	if measurement.SleepHours != nil {
		hours := *measurement.SleepHours
		switch {
		case hours >= 7 && hours <= 9:
			score += 50
		case hours >= 6 && hours < 7:
			score += 30
		case hours >= 5 && hours < 6:
			score += 10
		default:
			score += 0
		}
	}

	if measurement.SleepQuality != nil {
		quality := *measurement.SleepQuality
		score += quality * 5
	}

	// Determine status
	switch {
	case score >= 80:
		status = "Excellent"
	case score >= 60:
		status = "Good"
	case score >= 40:
		status = "Fair"
	default:
		status = "Poor"
	}

	return score, status
}

// TrackWellbeingTrend tracks wellbeing trends over time
func (uc *WellbeingUsecase) TrackWellbeingTrend(userID int, days int) (map[string]float64, error) {
	startDate := time.Now().AddDate(0, 0, -days)
	measurements, err := uc.wellbeingRepo.ListByUserWithDateRange(userID, startDate, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to track wellbeing trend: %w", err)
	}

	trends := map[string]float64{
		"avg_energy":   0,
		"avg_mood":     0,
		"avg_stress":   0,
		"avg_sleep":    0,
		"avg_hunger":   0,
	}

	count := 0
	for _, m := range measurements {
		if m.EnergyLevel != nil {
			trends["avg_energy"] += float64(*m.EnergyLevel)
		}
		if m.MoodLevel != nil {
			trends["avg_mood"] += float64(*m.MoodLevel)
		}
		if m.StressLevel != nil {
			trends["avg_stress"] += float64(*m.StressLevel)
		}
		if m.SleepHours != nil {
			trends["avg_sleep"] += *m.SleepHours
		}
		if m.HungerLevel != nil {
			trends["avg_hunger"] += float64(*m.HungerLevel)
		}
		count++
	}

	if count > 0 {
		for key := range trends {
			trends[key] /= float64(count)
		}
	}

	return trends, nil
}

// Update updates a wellbeing measurement
func (uc *WellbeingUsecase) Update(id int, energyLevel, moodLevel, hungerLevel, sleepQuality, stressLevel *int, sleepHours *float64, digestionNote *string, measuredAt time.Time) (*domain.WellbeingMeasurement, error) {
	measurement, err := uc.wellbeingRepo.GetByID(id)
	if err != nil {
		return nil, ErrWellbeingMeasurementNotFound
	}

	// Validate ranges
	if energyLevel != nil && (*energyLevel < 1 || *energyLevel > 10) {
		return nil, fmt.Errorf("energy level must be between 1 and 10")
	}
	if moodLevel != nil && (*moodLevel < 1 || *moodLevel > 10) {
		return nil, fmt.Errorf("mood level must be between 1 and 10")
	}
	if hungerLevel != nil && (*hungerLevel < 1 || *hungerLevel > 10) {
		return nil, fmt.Errorf("hunger level must be between 1 and 10")
	}
	if sleepQuality != nil && (*sleepQuality < 1 || *sleepQuality > 10) {
		return nil, fmt.Errorf("sleep quality must be between 1 and 10")
	}
	if stressLevel != nil && (*stressLevel < 1 || *stressLevel > 10) {
		return nil, fmt.Errorf("stress level must be between 1 and 10")
	}

	// Update fields if provided
	if energyLevel != nil {
		measurement.EnergyLevel = energyLevel
	}
	if moodLevel != nil {
		measurement.MoodLevel = moodLevel
	}
	if hungerLevel != nil {
		measurement.HungerLevel = hungerLevel
	}
	if sleepHours != nil {
		measurement.SleepHours = sleepHours
	}
	if sleepQuality != nil {
		measurement.SleepQuality = sleepQuality
	}
	if digestionNote != nil {
		measurement.DigestionNote = digestionNote
	}
	if stressLevel != nil {
		measurement.StressLevel = stressLevel
	}
	if !measuredAt.IsZero() {
		measurement.MeasuredAt = measuredAt
	}

	if err := uc.wellbeingRepo.Update(measurement); err != nil {
		return nil, fmt.Errorf("failed to update wellbeing measurement: %w", err)
	}

	return measurement, nil
}

// Delete deletes a wellbeing measurement
func (uc *WellbeingUsecase) Delete(id int) error {
	if err := uc.wellbeingRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete wellbeing measurement: %w", err)
	}
	return nil
}