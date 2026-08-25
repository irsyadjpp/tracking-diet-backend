package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/irsyadjpp/tracking-diet-backend/internal/delivery/dto"
	"github.com/irsyadjpp/tracking-diet-backend/internal/delivery/middleware"
	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
	"github.com/irsyadjpp/tracking-diet-backend/internal/usecase"
	"github.com/irsyadjpp/tracking-diet-backend/pkg/logger"
)

// NutritionHandler handles nutrition measurement HTTP requests
type NutritionHandler struct {
	nutritionUC *usecase.NutritionUsecase
	logger      *logger.Logger
}

// NewNutritionHandler creates a new nutrition handler
func NewNutritionHandler(nutritionUC *usecase.NutritionUsecase, log *logger.Logger) *NutritionHandler {
	return &NutritionHandler{
		nutritionUC: nutritionUC,
		logger:      log,
	}
}

// Create handles nutrition measurement creation
func (h *NutritionHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		middleware.ErrorResponder(w, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}

	var req dto.NutritionMeasurementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST")
		return
	}

	if err := req.Validate(); err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Validation failed: "+err.Error(), "VALIDATION_ERROR")
		return
	}

	measuredAt, err := time.Parse("2006-01-02", req.MeasuredAt)
	if err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid measured_at format", "INVALID_DATE_FORMAT")
		return
	}

	measurement, err := h.nutritionUC.Create(userID, req.CaloriesInKcal, req.CarbsG, req.ProteinG, req.FatG, req.FiberG, req.WaterIntakeL, req.CaloriesOutKcal, req.StepCount, req.MealTimingNote, measuredAt)
	if err != nil {
		h.logger.Error("Failed to create nutrition measurement", err, map[string]interface{}{"user_id": userID})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to create nutrition measurement", "INTERNAL_ERROR")
		return
	}

	response := h.toResponse(measurement)
	middleware.JSONResponder(w, http.StatusCreated, response)
}

// GetByID handles getting a nutrition measurement by ID
func (h *NutritionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid ID", "INVALID_ID")
		return
	}

	measurement, err := h.nutritionUC.GetByID(id)
	if err != nil {
		h.logger.Error("Nutrition measurement not found", err, map[string]interface{}{"id": id})
		middleware.ErrorResponder(w, http.StatusNotFound, "Nutrition measurement not found", "NOT_FOUND")
		return
	}

	response := h.toResponse(measurement)
	middleware.JSONResponder(w, http.StatusOK, response)
}

// List handles listing nutrition measurements for the authenticated user
func (h *NutritionHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		middleware.ErrorResponder(w, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}

	// Check for date range filters
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var measurements []domain.NutritionMeasurement
	var err error

	if startDateStr != "" && endDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid start_date format", "INVALID_DATE_FORMAT")
			return
		}
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid end_date format", "INVALID_DATE_FORMAT")
			return
		}
		measurements, err = h.nutritionUC.ListByUserWithDateRange(userID, startDate, endDate)
	} else {
		measurements, err = h.nutritionUC.ListByUser(userID)
	}

	if err != nil {
		h.logger.Error("Failed to list nutrition measurements", err, map[string]interface{}{"user_id": userID})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to list nutrition measurements", "INTERNAL_ERROR")
		return
	}

	response := make([]dto.NutritionMeasurementResponse, len(measurements))
	for i, m := range measurements {
		response[i] = h.toResponse(&m)
	}

	middleware.JSONResponder(w, http.StatusOK, response)
}

// GetDailySummary handles getting daily nutrition summary
func (h *NutritionHandler) GetDailySummary(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		middleware.ErrorResponder(w, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}

	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid date format", "INVALID_DATE_FORMAT")
		return
	}

	summary, err := h.nutritionUC.GetDailySummary(userID, date)
	if err != nil {
		h.logger.Error("Failed to get daily summary", err, map[string]interface{}{"user_id": userID, "date": dateStr})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to get daily summary", "INTERNAL_ERROR")
		return
	}

	response := h.toResponse(summary)
	middleware.JSONResponder(w, http.StatusOK, response)
}

// Update handles updating a nutrition measurement
func (h *NutritionHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid ID", "INVALID_ID")
		return
	}

	var req dto.NutritionMeasurementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST")
		return
	}

	if err := req.Validate(); err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Validation failed: "+err.Error(), "VALIDATION_ERROR")
		return
	}

	measuredAt, err := time.Parse("2006-01-02", req.MeasuredAt)
	if err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid measured_at format", "INVALID_DATE_FORMAT")
		return
	}

	measurement, err := h.nutritionUC.Update(id, req.CaloriesInKcal, req.CarbsG, req.ProteinG, req.FatG, req.FiberG, req.WaterIntakeL, req.CaloriesOutKcal, req.StepCount, req.MealTimingNote, measuredAt)
	if err != nil {
		h.logger.Error("Failed to update nutrition measurement", err, map[string]interface{}{"id": id})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to update nutrition measurement", "INTERNAL_ERROR")
		return
	}

	response := h.toResponse(measurement)
	middleware.JSONResponder(w, http.StatusOK, response)
}

// Delete handles deleting a nutrition measurement
func (h *NutritionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid ID", "INVALID_ID")
		return
	}

	if err := h.nutritionUC.Delete(id); err != nil {
		h.logger.Error("Failed to delete nutrition measurement", err, map[string]interface{}{"id": id})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to delete nutrition measurement", "INTERNAL_ERROR")
		return
	}

	middleware.JSONResponder(w, http.StatusNoContent, nil)
}

func (h *NutritionHandler) toResponse(m *domain.NutritionMeasurement) dto.NutritionMeasurementResponse {
	return dto.NutritionMeasurementResponse{
		ID:             m.ID,
		UserID:         m.UserID,
		MeasuredAt:     m.MeasuredAt,
		CaloriesInKcal: m.CaloriesInKcal,
		CarbsG:         m.CarbsG,
		ProteinG:       m.ProteinG,
		FatG:           m.FatG,
		FiberG:         m.FiberG,
		WaterIntakeL:   m.WaterIntakeL,
		CaloriesOutKcal: m.CaloriesOutKcal,
		StepCount:      m.StepCount,
		MealTimingNote: m.MealTimingNote,
	}
}