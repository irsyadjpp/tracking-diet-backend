package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/irsyadjpp/tracking-diet-backend/internal/delivery/dto"
	"github.com/irsyadjpp/tracking-diet-backend/internal/delivery/middleware"
	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
	"github.com/irsyadjpp/tracking-diet-backend/internal/usecase"
	"github.com/irsyadjpp/tracking-diet-backend/pkg/logger"
)

// AIRecommendationHandler handles AI recommendation HTTP requests
type AIRecommendationHandler struct {
	aiUC   *usecase.AIRecommendationUsecase
	logger *logger.Logger
}

// NewAIRecommendationHandler creates a new AI recommendation handler
func NewAIRecommendationHandler(aiUC *usecase.AIRecommendationUsecase, log *logger.Logger) *AIRecommendationHandler {
	return &AIRecommendationHandler{
		aiUC:   aiUC,
		logger: log,
	}
}

// GenerateRequest represents AI recommendation generation request
type GenerateRequest struct {
	Category string                 `json:"category" validate:"required"`
	UserData map[string]interface{} `json:"user_data"`
}

// Generate handles AI recommendation generation
func (h *AIRecommendationHandler) Generate(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		middleware.ErrorResponder(w, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}

	var req GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST")
		return
	}

	if req.Category == "" {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Category is required", "VALIDATION_ERROR")
		return
	}

	// Generate recommendation using AI
	ctx := context.Background()
	recommendation, err := h.aiUC.GenerateRecommendation(ctx, userID, req.Category, req.UserData)
	if err != nil {
		h.logger.Error("Failed to generate AI recommendation", err, map[string]interface{}{"user_id": userID, "category": req.Category})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to generate AI recommendation", "INTERNAL_ERROR")
		return
	}

	response := h.toResponse(recommendation)
	middleware.JSONResponder(w, http.StatusCreated, response)
}

// GetByID handles getting an AI recommendation by ID
func (h *AIRecommendationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid ID", "INVALID_ID")
		return
	}

	recommendation, err := h.aiUC.GetByID(id)
	if err != nil {
		h.logger.Error("AI recommendation not found", err, map[string]interface{}{"id": id})
		middleware.ErrorResponder(w, http.StatusNotFound, "AI recommendation not found", "NOT_FOUND")
		return
	}

	response := h.toResponse(recommendation)
	middleware.JSONResponder(w, http.StatusOK, response)
}

// List handles listing AI recommendations for the authenticated user
func (h *AIRecommendationHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		middleware.ErrorResponder(w, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}

	// Check for category filter
	category := r.URL.Query().Get("category")

	var recommendations []domain.AIRecommendation
	var err error

	if category != "" {
		recommendations, err = h.aiUC.ListByUserWithCategory(userID, category)
	} else {
		recommendations, err = h.aiUC.ListByUser(userID)
	}

	if err != nil {
		h.logger.Error("Failed to list AI recommendations", err, map[string]interface{}{"user_id": userID})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to list AI recommendations", "INTERNAL_ERROR")
		return
	}

	response := make([]dto.AIRecommendationResponse, len(recommendations))
	for i, rec := range recommendations {
		response[i] = h.toResponse(&rec)
	}

	middleware.JSONResponder(w, http.StatusOK, response)
}

// Delete handles deleting an AI recommendation
func (h *AIRecommendationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid ID", "INVALID_ID")
		return
	}

	if err := h.aiUC.Delete(id); err != nil {
		h.logger.Error("Failed to delete AI recommendation", err, map[string]interface{}{"id": id})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to delete AI recommendation", "INTERNAL_ERROR")
		return
	}

	middleware.JSONResponder(w, http.StatusNoContent, nil)
}

func (h *AIRecommendationHandler) toResponse(rec *domain.AIRecommendation) dto.AIRecommendationResponse {
	var category string
	if rec.Category != nil {
		category = *rec.Category
	}

	var confidenceScore float64
	if rec.ConfidenceScore != nil {
		confidenceScore = *rec.ConfidenceScore
	}

	return dto.AIRecommendationResponse{
		ID:              rec.ID,
		UserID:          rec.UserID,
		GeneratedAt:     rec.GeneratedAt,
		Category:        category,
		InputSummary:    rec.InputSummary.Data,
		OutputText:      rec.OutputText,
		ConfidenceScore: confidenceScore,
	}
}