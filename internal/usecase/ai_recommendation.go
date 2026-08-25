package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/irsyadjpp/tracking-diet-backend/internal/ai"
	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
)

var (
	ErrAIRecommendationNotFound = errors.New("AI recommendation not found")
	ErrAIServiceUnavailable     = errors.New("AI service unavailable")
)

// AIRecommendationUsecase handles AI recommendation business logic
type AIRecommendationUsecase struct {
	aiRepo   domain.AIRecommendationRepository
	userRepo domain.UserRepository
	aiService *ai.AIService
}

// NewAIRecommendationUsecase creates a new AI recommendation use case
func NewAIRecommendationUsecase(aiRepo domain.AIRecommendationRepository, userRepo domain.UserRepository, aiService *ai.AIService) *AIRecommendationUsecase {
	return &AIRecommendationUsecase{
		aiRepo:   aiRepo,
		userRepo: userRepo,
		aiService: aiService,
	}
}

// Create creates a new AI recommendation
func (uc *AIRecommendationUsecase) Create(userID int, category string, inputSummary interface{}, outputText string, confidenceScore float64) (*domain.AIRecommendation, error) {
	// Validate user exists
	_, err := uc.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if category == "" {
		return nil, fmt.Errorf("category is required")
	}

	if outputText == "" {
		return nil, fmt.Errorf("output text is required")
	}

	if confidenceScore < 0 || confidenceScore > 1 {
		return nil, fmt.Errorf("confidence score must be between 0 and 1")
	}

	recommendation := &domain.AIRecommendation{
		UserID:          userID,
		Category:        &category,
		InputSummary:    domain.JSONB{Data: inputSummary},
		OutputText:      &outputText,
		ConfidenceScore: &confidenceScore,
	}

	if err := uc.aiRepo.Create(recommendation); err != nil {
		return nil, fmt.Errorf("failed to create AI recommendation: %w", err)
	}

	return recommendation, nil
}

// GetByID retrieves an AI recommendation by ID
func (uc *AIRecommendationUsecase) GetByID(id int) (*domain.AIRecommendation, error) {
	recommendation, err := uc.aiRepo.GetByID(id)
	if err != nil {
		return nil, ErrAIRecommendationNotFound
	}
	return recommendation, nil
}

// ListByUser retrieves all AI recommendations for a user
func (uc *AIRecommendationUsecase) ListByUser(userID int) ([]domain.AIRecommendation, error) {
	recommendations, err := uc.aiRepo.ListByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list AI recommendations: %w", err)
	}
	return recommendations, nil
}

// ListByUserWithCategory retrieves AI recommendations for a user with a specific category
func (uc *AIRecommendationUsecase) ListByUserWithCategory(userID int, category string) ([]domain.AIRecommendation, error) {
	recommendations, err := uc.aiRepo.ListByUserWithCategory(userID, category)
	if err != nil {
		return nil, fmt.Errorf("failed to list AI recommendations: %w", err)
	}
	return recommendations, nil
}

// GenerateRecommendation triggers AI recommendation generation using the AI service
func (uc *AIRecommendationUsecase) GenerateRecommendation(ctx context.Context, userID int, category string, userData interface{}) (*domain.AIRecommendation, error) {
	// Validate user exists
	user, err := uc.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Build user profile for AI
	userProfile := map[string]interface{}{
		"age":         calculateAge(user.BirthDate),
		"gender":      user.Gender,
		"height_cm":   user.HeightCm,
		"weight_kg":   getUserWeight(user), // Would need to get from body measurements
		"activity_level": "moderate", // Would need to get from fitness measurements
	}

	// Prepare request based on category
	var outputText string
	var confidenceScore float64

	switch category {
	case "meal_plan":
		req := &ai.MealPlanRequest{
			UserProfile:         userProfile,
			Goals:              []string{"weight loss", "muscle gain"},
			DietaryRestrictions: []string{},
			CalorieTarget:       2000,
			MealsPerDay:         3,
		}
		response, err := uc.aiService.GenerateMealPlan(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to generate meal plan: %w", err)
		}
		outputText = response.MealPlan
		confidenceScore = response.Confidence

	case "progress_report":
		req := &ai.ProgressReportRequest{
			UserProfile:  userProfile,
			Measurements: []map[string]interface{}{}, // Would need to fetch actual measurements
			TimeRange:    "last 30 days",
			FocusAreas:   []string{"weight", "nutrition", "fitness"},
		}
		response, err := uc.aiService.GenerateProgressReport(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to generate progress report: %w", err)
		}
		outputText = response.Report
		confidenceScore = response.Confidence

	default:
		req := &ai.RecommendationRequest{
			UserProfile:   userProfile,
			CurrentStatus: map[string]interface{}{},
			Concerns:      []string{},
			Category:      category,
		}
		response, err := uc.aiService.GenerateRecommendation(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to generate recommendation: %w", err)
		}
		outputText = response.Recommendation
		confidenceScore = response.Confidence
	}

	return uc.Create(userID, category, userData, outputText, confidenceScore)
}

// Helper function to calculate age from birth date
func calculateAge(birthDate *time.Time) int {
	if birthDate == nil {
		return 30 // Default age if not provided
	}
	
	now := time.Now()
	age := now.Year() - birthDate.Year()
	
	// Adjust if birthday hasn't occurred yet this year
	if now.Month() < birthDate.Month() || (now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
		age--
	}
	
	return age
}

// Helper function to get user weight (placeholder - would need to fetch from body measurements)
func getUserWeight(user *domain.User) float64 {
	// In a real implementation, this would fetch the latest body measurement
	// For now, return a default value
	return 70.0
}

// CategorizeRecommendation categorizes a recommendation based on content
func (uc *AIRecommendationUsecase) CategorizeRecommendation(content string) string {
	// Simple keyword-based categorization
	// In production, this could use more sophisticated NLP
	keywords := map[string]string{
		"meal":        "meal_plan",
		"food":        "meal_plan",
		"diet":        "meal_plan",
		"nutrition":   "meal_plan",
		"progress":    "progress_report",
		"report":      "progress_report",
		"summary":     "progress_report",
		"workout":     "fitness",
		"exercise":    "fitness",
		"training":    "fitness",
		"sleep":       "wellbeing",
		"stress":      "wellbeing",
		"mental":      "wellbeing",
		"health":      "general",
		"advice":      "general",
		"recommend":   "general",
	}

	for keyword, category := range keywords {
		if contains(content, keyword) {
			return category
		}
	}

	return "general"
}

// CalculateConfidenceScore calculates a confidence score for a recommendation
func (uc *AIRecommendationUsecase) CalculateConfidenceScore(quality, relevance, completeness float64) float64 {
	// Weighted average of quality factors
	weights := map[string]float64{
		"quality":      0.4,
		"relevance":    0.4,
		"completeness": 0.2,
	}

	score := (quality * weights["quality"]) + 
	         (relevance * weights["relevance"]) + 
	         (completeness * weights["completeness"])

	// Ensure score is between 0 and 1
	if score < 0 {
		score = 0
	} else if score > 1 {
		score = 1
	}

	return score
}

// Update updates an AI recommendation
func (uc *AIRecommendationUsecase) Update(id int, category *string, inputSummary interface{}, outputText *string, confidenceScore *float64) (*domain.AIRecommendation, error) {
	recommendation, err := uc.aiRepo.GetByID(id)
	if err != nil {
		return nil, ErrAIRecommendationNotFound
	}

	// Update fields if provided
	if category != nil {
		recommendation.Category = category
	}
	if inputSummary != nil {
		recommendation.InputSummary = domain.JSONB{Data: inputSummary}
	}
	if outputText != nil {
		recommendation.OutputText = outputText
	}
	if confidenceScore != nil {
		if *confidenceScore < 0 || *confidenceScore > 1 {
			return nil, fmt.Errorf("confidence score must be between 0 and 1")
		}
		recommendation.ConfidenceScore = confidenceScore
	}

	if err := uc.aiRepo.Update(recommendation); err != nil {
		return nil, fmt.Errorf("failed to update AI recommendation: %w", err)
	}

	return recommendation, nil
}

// Delete deletes an AI recommendation
func (uc *AIRecommendationUsecase) Delete(id int) error {
	if err := uc.aiRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete AI recommendation: %w", err)
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