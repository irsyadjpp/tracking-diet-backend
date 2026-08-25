package ai

import (
	"context"
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// MealPlanRequest represents a request for meal plan generation
type MealPlanRequest struct {
	UserProfile     map[string]interface{} `json:"user_profile"`
	Goals          []string               `json:"goals"`
	DietaryRestrictions []string          `json:"dietary_restrictions"`
	CalorieTarget  int                    `json:"calorie_target"`
	MealsPerDay    int                    `json:"meals_per_day"`
}

// MealPlanResponse represents a generated meal plan
type MealPlanResponse struct {
	MealPlan      string  `json:"meal_plan"`
	Confidence    float64 `json:"confidence"`
	Reasoning     string  `json:"reasoning"`
}

// ProgressReportRequest represents a request for progress report generation
type ProgressReportRequest struct {
	UserProfile      map[string]interface{} `json:"user_profile"`
	Measurements     []map[string]interface{} `json:"measurements"`
	TimeRange        string                 `json:"time_range"`
	FocusAreas       []string               `json:"focus_areas"`
}

// ProgressReportResponse represents a generated progress report
type ProgressReportResponse struct {
	Report       string  `json:"report"`
	Confidence   float64 `json:"confidence"`
	Recommendations []string `json:"recommendations"`
}

// RecommendationRequest represents a request for personalized recommendations
type RecommendationRequest struct {
	UserProfile    map[string]interface{} `json:"user_profile"`
	CurrentStatus  map[string]interface{} `json:"current_status"`
	Concerns       []string               `json:"concerns"`
	Category       string                 `json:"category"`
}

// RecommendationResponse represents a generated recommendation
type RecommendationResponse struct {
	Recommendation string  `json:"recommendation"`
	Confidence     float64 `json:"confidence"`
	ActionItems    []string `json:"action_items"`
}

// GenerateMealPlan generates a personalized meal plan using AI
func (s *AIService) GenerateMealPlan(ctx context.Context, req *MealPlanRequest) (*MealPlanResponse, error) {
	if !s.genkitInitialized {
		return nil, fmt.Errorf("AI service not initialized")
	}

	// Build prompt for meal plan generation
	prompt := s.buildMealPlanPrompt(req)

	// Call AI model (using Gemini Flash)
	response, err := genkit.Generate(ctx, ai.Of("googleai/gemini-flash-1.5"), 
		ai.WithConfig(ai.GenerationConfig{
			Temperature: 0.7,
			MaxOutputTokens: 2048,
		}),
		ai.WithText(prompt),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate meal plan: %w", err)
	}

	return &MealPlanResponse{
		MealPlan:   response.Text(),
		Confidence: 0.85, // Placeholder confidence score
		Reasoning:  "Generated based on user profile and goals",
	}, nil
}

// GenerateProgressReport generates a progress report using AI
func (s *AIService) GenerateProgressReport(ctx context.Context, req *ProgressReportRequest) (*ProgressReportResponse, error) {
	if !s.genkitInitialized {
		return nil, fmt.Errorf("AI service not initialized")
	}

	// Build prompt for progress report generation
	prompt := s.buildProgressReportPrompt(req)

	// Call AI model
	response, err := genkit.Generate(ctx, ai.Of("googleai/gemini-flash-1.5"),
		ai.WithConfig(ai.GenerationConfig{
			Temperature: 0.6,
			MaxOutputTokens: 3072,
		}),
		ai.WithText(prompt),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate progress report: %w", err)
	}

	return &ProgressReportResponse{
		Report:         response.Text(),
		Confidence:     0.80,
		Recommendations: []string{"Continue current routine", "Increase protein intake"}, // Placeholder
	}, nil
}

// GenerateRecommendation generates personalized recommendations using AI
func (s *AIService) GenerateRecommendation(ctx context.Context, req *RecommendationRequest) (*RecommendationResponse, error) {
	if !s.genkitInitialized {
		return nil, fmt.Errorf("AI service not initialized")
	}

	// Build prompt for recommendation generation
	prompt := s.buildRecommendationPrompt(req)

	// Call AI model
	response, err := genkit.Generate(ctx, ai.Of("googleai/gemini-flash-1.5"),
		ai.WithConfig(ai.GenerationConfig{
			Temperature: 0.8,
			MaxOutputTokens: 1536,
		}),
		ai.WithText(prompt),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate recommendation: %w", err)
	}

	return &RecommendationResponse{
		Recommendation: response.Text(),
		Confidence:    0.75,
		ActionItems:    []string{"Track daily intake", "Schedule check-up"}, // Placeholder
	}, nil
}

// buildMealPlanPrompt builds the prompt for meal plan generation
func (s *AIService) buildMealPlanPrompt(req *MealPlanRequest) string {
	return fmt.Sprintf(`You are a professional nutritionist and dietitian. Create a personalized meal plan based on the following information:

User Profile:
- Age: %v
- Gender: %v
- Height: %v cm
- Weight: %v kg
- Activity Level: %v

Goals: %v
Dietary Restrictions: %v
Target Calories: %v kcal/day
Meals per day: %d

Please provide:
1. A detailed meal plan with breakfast, lunch, dinner, and snacks
2. Nutritional breakdown for each meal
3. Shopping list suggestions
4. Meal preparation tips

Format the response in a clear, structured way that is easy to follow.`,
		req.UserProfile["age"],
		req.UserProfile["gender"],
		req.UserProfile["height_cm"],
		req.UserProfile["weight_kg"],
		req.UserProfile["activity_level"],
		req.Goals,
		req.DietaryRestrictions,
		req.CalorieTarget,
		req.MealsPerDay,
	)
}

// buildProgressReportPrompt builds the prompt for progress report generation
func (s *AIService) buildProgressReportPrompt(req *ProgressReportRequest) string {
	return fmt.Sprintf(`You are a health and fitness coach. Analyze the following data and create a comprehensive progress report:

User Profile:
- Name: %v
- Age: %v
- Gender: %v
- Goals: %v

Recent Measurements (over %s):
%v

Focus Areas: %v

Please provide:
1. Overall progress assessment
2. Key achievements and improvements
3. Areas that need attention
4. Specific recommendations for each focus area
5. Motivational feedback

Format the response in a clear, encouraging manner.`,
		req.UserProfile["name"],
		req.UserProfile["age"],
		req.UserProfile["gender"],
		req.UserProfile["goals"],
		req.TimeRange,
		formatMeasurements(req.Measurements),
		req.FocusAreas,
	)
}

// buildRecommendationPrompt builds the prompt for recommendation generation
func (s *AIService) buildRecommendationPrompt(req *RecommendationRequest) string {
	return fmt.Sprintf(`You are a health advisor. Provide personalized recommendations based on the following information:

User Profile:
- Age: %v
- Gender: %v
- Current Status: %v

Current Concerns: %v
Category: %s

Please provide:
1. Main recommendation addressing the concerns
2. Actionable steps to implement the recommendation
3. Expected timeline for seeing results
4. Any precautions or considerations

Format the response in a clear, practical manner.`,
		req.UserProfile["age"],
		req.UserProfile["gender"],
		formatCurrentStatus(req.CurrentStatus),
		req.Concerns,
		req.Category,
	)
}

// Helper functions for formatting
func formatMeasurements(measurements []map[string]interface{}) string {
	if len(measurements) == 0 {
		return "No measurements available"
	}
	
	result := ""
	for i, m := range measurements {
		result += fmt.Sprintf("\nMeasurement %d:\n", i+1)
		for key, value := range m {
			result += fmt.Sprintf("- %s: %v\n", key, value)
		}
	}
	return result
}

func formatCurrentStatus(status map[string]interface{}) string {
	result := ""
	for key, value := range status {
		result += fmt.Sprintf("%s: %v, ", key, value)
	}
	return result[:len(result)-2] // Remove trailing comma and space
}