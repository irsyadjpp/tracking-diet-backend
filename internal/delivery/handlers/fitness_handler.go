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

// FitnessHandler handles fitness measurement HTTP requests
type FitnessHandler struct {
	fitnessUC *usecase.FitnessUsecase
	logger    *logger.Logger
}

// NewFitnessHandler creates a new fitness handler
func NewFitnessHandler(fitnessUC *usecase.FitnessUsecase, log *logger.Logger) *FitnessHandler {
	return &FitnessHandler{
		fitnessUC: fitnessUC,
		logger:    log,
	}
}

// Create handles fitness measurement creation
func (h *FitnessHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		middleware.ErrorResponder(w, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}

	var req dto.FitnessMeasurementRequest
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

	measurement, err := h.fitnessUC.Create(userID, req.VO2MaxMlKgMin, req.Strength1RmKg, req.PushupCount, req.SquatCount, req.EnduranceNote, req.FlexibilityNote, measuredAt)
	if err != nil {
		h.logger.Error("Failed to create fitness measurement", err, map[string]interface{}{"user_id": userID})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to create fitness measurement", "INTERNAL_ERROR")
		return
	}

	response := h.toResponse(measurement)
	middleware.JSONResponder(w, http.StatusCreated, response)
}

// GetByID handles getting a fitness measurement by ID
func (h *FitnessHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid ID", "INVALID_ID")
		return
	}

	measurement, err := h.fitnessUC.GetByID(id)
	if err != nil {
		h.logger.Error("Fitness measurement not found", err, map[string]interface{}{"id": id})
		middleware.ErrorResponder(w, http.StatusNotFound, "Fitness measurement not found", "NOT_FOUND")
		return
	}

	response := h.toResponse(measurement)
	middleware.JSONResponder(w, http.StatusOK, response)
}

// List handles listing fitness measurements for the authenticated user
func (h *FitnessHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		middleware.ErrorResponder(w, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}

	// Check for date range filters
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var measurements []domain.FitnessMeasurement
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
		measurements, err = h.fitnessUC.ListByUserWithDateRange(userID, startDate, endDate)
	} else {
		measurements, err = h.fitnessUC.ListByUser(userID)
	}

	if err != nil {
		h.logger.Error("Failed to list fitness measurements", err, map[string]interface{}{"user_id": userID})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to list fitness measurements", "INTERNAL_ERROR")
		return
	}

	response := make([]dto.FitnessMeasurementResponse, len(measurements))
	for i, m := range measurements {
		response[i] = h.toResponse(&m)
	}

	middleware.JSONResponder(w, http.StatusOK, response)
}

// Update handles updating a fitness measurement
func (h *FitnessHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid ID", "INVALID_ID")
		return
	}

	var req dto.FitnessMeasurementRequest
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

	measurement, err := h.fitnessUC.Update(id, req.VO2MaxMlKgMin, req.Strength1RmKg, req.PushupCount, req.SquatCount, req.EnduranceNote, req.FlexibilityNote, measuredAt)
	if err != nil {
		h.logger.Error("Failed to update fitness measurement", err, map[string]interface{}{"id": id})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to update fitness measurement", "INTERNAL_ERROR")
		return
	}

	response := h.toResponse(measurement)
	middleware.JSONResponder(w, http.StatusOK, response)
}

// Delete handles deleting a fitness measurement
func (h *FitnessHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid ID", "INVALID_ID")
		return
	}

	if err := h.fitnessUC.Delete(id); err != nil {
		h.logger.Error("Failed to delete fitness measurement", err, map[string]interface{}{"id": id})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to delete fitness measurement", "INTERNAL_ERROR")
		return
	}

	middleware.JSONResponder(w, http.StatusNoContent, nil)
}

func (h *FitnessHandler) toResponse(m *domain.FitnessMeasurement) dto.FitnessMeasurementResponse {
	return dto.FitnessMeasurementResponse{
		ID:              m.ID,
		UserID:          m.UserID,
		MeasuredAt:      m.MeasuredAt,
		VO2MaxMlKgMin:   m.VO2MaxMlKgMin,
		Strength1RmKg:   m.Strength1RmKg,
		PushupCount:     m.PushupCount,
		SquatCount:      m.SquatCount,
		EnduranceNote:   m.EnduranceNote,
		FlexibilityNote: m.FlexibilityNote,
	}
}