package ai

import (
	"context"
	"fmt"

	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
)

// AIService handles AI operations using Google Genkit
type AIService struct {
	genkitInitialized bool
	apiKey           string
}

// NewAIService creates a new AI service
func NewAIService(apiKey string) *AIService {
	return &AIService{
		apiKey: apiKey,
	}
}

// Initialize initializes the Genkit framework with Google AI plugin
func (s *AIService) Initialize(ctx context.Context) error {
	if s.genkitInitialized {
		return nil
	}

	if s.apiKey == "" {
		return fmt.Errorf("Google Genkit API key is required")
	}

	// Initialize Genkit with Google AI plugin
	if err := genkit.Init(ctx, nil); err != nil {
		return fmt.Errorf("failed to initialize Genkit: %w", err)
	}

	// Configure Google AI plugin
	if err := googlegenai.Init(ctx, &googlegenai.Config{
		APIKey: s.apiKey,
	}); err != nil {
		return fmt.Errorf("failed to initialize Google AI plugin: %w", err)
	}

	s.genkitInitialized = true
	return nil
}

// IsInitialized returns whether the AI service has been initialized
func (s *AIService) IsInitialized() bool {
	return s.genkitInitialized
}