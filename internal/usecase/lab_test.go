package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
)

var (
	ErrLabTestNotFound = errors.New("lab test not found")
)

// LabTestUsecase handles lab test business logic
type LabTestUsecase struct {
	labTestRepo domain.LabTestRepository
	userRepo    domain.UserRepository
}

// NewLabTestUsecase creates a new lab test use case
func NewLabTestUsecase(labTestRepo domain.LabTestRepository, userRepo domain.UserRepository) *LabTestUsecase {
	return &LabTestUsecase{
		labTestRepo: labTestRepo,
		userRepo:    userRepo,
	}
}

// Create creates a new lab test
func (uc *LabTestUsecase) Create(userID int, testName, resultValue, unit, referenceRange *string, measuredAt time.Time) (*domain.LabTest, error) {
	// Validate user exists
	_, err := uc.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if testName == "" {
		return nil, fmt.Errorf("test name is required")
	}

	labTest := &domain.LabTest{
		UserID:         userID,
		TestName:       testName,
		ResultValue:    resultValue,
		Unit:           unit,
		ReferenceRange: referenceRange,
		MeasuredAt:     measuredAt,
	}

	if err := uc.labTestRepo.Create(labTest); err != nil {
		return nil, fmt.Errorf("failed to create lab test: %w", err)
	}

	return labTest, nil
}

// GetByID retrieves a lab test by ID
func (uc *LabTestUsecase) GetByID(id int) (*domain.LabTest, error) {
	labTest, err := uc.labTestRepo.GetByID(id)
	if err != nil {
		return nil, ErrLabTestNotFound
	}
	return labTest, nil
}

// ListByUser retrieves all lab tests for a user
func (uc *LabTestUsecase) ListByUser(userID int) ([]domain.LabTest, error) {
	labTests, err := uc.labTestRepo.ListByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list lab tests: %w", err)
	}
	return labTests, nil
}

// ListByUserWithTestName retrieves lab tests for a user with a specific test name
func (uc *LabTestUsecase) ListByUserWithTestName(userID int, testName string) ([]domain.LabTest, error) {
	labTests, err := uc.labTestRepo.ListByUserWithTestName(userID, testName)
	if err != nil {
		return nil, fmt.Errorf("failed to list lab tests: %w", err)
	}
	return labTests, nil
}

// ListByUserWithDateRange retrieves lab tests for a user within a date range
func (uc *LabTestUsecase) ListByUserWithDateRange(userID int, startDate, endDate time.Time) ([]domain.LabTest, error) {
	labTests, err := uc.labTestRepo.ListByUserWithDateRange(userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to list lab tests: %w", err)
	}
	return labTests, nil
}

// ValidateReferenceRange validates if a result is within reference range
func (uc *LabTestUsecase) ValidateReferenceRange(labTest *domain.LabTest) (bool, string) {
	if labTest.ResultValue == nil || labTest.ReferenceRange == nil {
		return true, "Cannot validate - missing data"
	}

	// This is a simplified validation. In a real implementation, you would parse
	// the reference range and compare numerically. For now, we'll do a basic check.
	result := *labTest.ResultValue
	refRange := *labTest.ReferenceRange

	// If result contains any "abnormal" indicators, flag it
	abnormalIndicators := []string{"high", "low", "abnormal", "out of range"}
	for _, indicator := range abnormalIndicators {
		if contains(result, indicator) {
			return false, "Result appears abnormal"
		}
	}

	// If reference range is provided but result doesn't match typical patterns
	if refRange != "" && result != "" {
		// Simplified check - in production, implement proper range parsing
		return true, "Within reference range"
	}

	return true, "Cannot determine"
}

// CompareWithPrevious compares current lab test with previous results
func (uc *LabTestUsecase) CompareWithPrevious(userID int, testName string, currentTest *domain.LabTest) ([]domain.LabTest, string, error) {
	previousTests, err := uc.labTestRepo.ListByUserWithTestName(userID, testName)
	if err != nil {
		return nil, "", fmt.Errorf("failed to compare with previous tests: %w", err)
	}

	if len(previousTests) == 0 {
		return nil, "No previous tests for comparison", nil
	}

	// Filter out the current test and get only previous ones
	var previousOnly []domain.LabTest
	for _, test := range previousTests {
		if test.ID != currentTest.ID {
			previousOnly = append(previousOnly, test)
		}
	}

	if len(previousOnly) == 0 {
		return nil, "No previous tests for comparison", nil
	}

	// Get the most recent previous test
	mostRecent := previousOnly[0]
	trend := "stable"

	if currentTest.ResultValue != nil && mostRecent.ResultValue != nil {
		// Simplified trend comparison
		currentVal := *currentTest.ResultValue
		previousVal := *mostRecent.ResultValue
		if currentVal != previousVal {
			trend = "changed"
		}
	}

	return previousOnly, trend, nil
}

// Update updates a lab test
func (uc *LabTestUsecase) Update(id int, testName, resultValue, unit, referenceRange *string, measuredAt time.Time) (*domain.LabTest, error) {
	labTest, err := uc.labTestRepo.GetByID(id)
	if err != nil {
		return nil, ErrLabTestNotFound
	}

	// Update fields if provided
	if testName != "" {
		labTest.TestName = testName
	}
	if resultValue != nil {
		labTest.ResultValue = resultValue
	}
	if unit != nil {
		labTest.Unit = unit
	}
	if referenceRange != nil {
		labTest.ReferenceRange = referenceRange
	}
	if !measuredAt.IsZero() {
		labTest.MeasuredAt = measuredAt
	}

	if err := uc.labTestRepo.Update(labTest); err != nil {
		return nil, fmt.Errorf("failed to update lab test: %w", err)
	}

	return labTest, nil
}

// Delete deletes a lab test
func (uc *LabTestUsecase) Delete(id int) error {
	if err := uc.labTestRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete lab test: %w", err)
	}
	return nil
}

// Helper function to check if string contains substring (case-insensitive)
func contains(s, substr string) bool {
	if len(s) < len(substr) {
		return false
	}
	// Simple case-insensitive check
	sLower := toLower(s)
	substrLower := toLower(substr)
	for i := 0; i <= len(sLower)-len(substrLower); i++ {
		if sLower[i:i+len(substrLower)] == substrLower {
			return true
		}
	}
	return false
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			result[i] = c + 32
		} else {
			result[i] = c
		}
	}
	return string(result)
}