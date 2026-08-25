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

// WellbeingHandler handles wellbeing measurement HTTP requests
type WellbeingHandler struct {
	wellbeingUC *usecase.WellbeingUsecase
	logger      *logger.Logger
}

// NewWellbeingHandler creates a new wellbeing handler
func NewWellbeingHandler(wellbeingUC *usecase.WellbeingUsecase, log *logger.Logger) *WellbeingHandler {
	return &WellbeingHandler{
		wellbeingUC: wellbeingUC,
		logger:      log,
	}
}

// Create handles wellbeing measurement creation
func (h *WellbeingHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		middleware.ErrorResponder(w, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}

	var req dto.WellbeingMeasurementRequest
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

	measurement, err := h.wellbeingUC.Create(userID, req.EnergyLevel, req.MoodLevel, req.HungerLevel, req.SleepQuality, req.StressLevel, req.SleepHours, req.DigestionNote, measuredAt)
	if err != nil {
		h.logger.Error("Failed to create wellbeing measurement", err, map[string]interface{}{"user_id": userID})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to create wellbeing measurement", "INTERNAL_ERROR")
		return
	}

	response := h.toResponse(measurement)
	middleware.JSONResponder(w, http.StatusCreated, response)
}

// GetByID handles getting a wellbeing measurement by ID
func (h *WellbeingHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid ID", "INVALID_ID")
		return
	}

	measurement, err := h.wellbeingUC.GetByID(id)
	if err != nil {
		h.logger.Error("Wellbeing measurement not found", err, map[string]interface{}{"id": id})
		middleware.ErrorResponder(w, http.StatusNotFound, "Wellbeing measurement not found", "NOT_FOUND")
		return
	}

	response := h.toResponse(measurement)
	middleware.JSONResponder(w, http.StatusOK, response)
}

// List handles listing wellbeing measurements for the authenticated user
func (h *WellbeingHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		middleware.ErrorResponder(w, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}

	// Check for date range filters
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var measurements []domain.WellbeingMeasurement
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
		measurements, err = h.wellbeingUC.ListByUserWithDateRange(userID, startDate, endDate)
	} else {
		measurements, err = h.wellbeingUC.ListByUser(userID)
	}

	if err != nil {
		h.logger.Error("Failed to list wellbeing measurements", err, map[string]interface{}{"user_id": userID})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to list wellbeing measurements", "INTERNAL_ERROR")
		return
	}

	response := make([]dto.WellbeingMeasurementResponse, len(measurements))
	for i, m := range measurements {
		response[i] = h.toResponse(&m)
	}

	middleware.JSONResponder(w, http.StatusOK, response)
}

// Update handles updating a wellbeing measurement
func (h *WellbeingHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid ID", "INVALID_ID")
		return
	}

	var req dto.WellbeingMeasurementRequest
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

	measurement, err := h.wellbeingUC.Update(id, req.EnergyLevel, req.MoodLevel, req.HungerLevel, req.SleepQuality, req.StressLevel, req.SleepHours, req.DigestionNote, measuredAt)
	if err != nil {
		h.logger.Error("Failed to update wellbeing measurement", err, map[string]interface{}{"id": id})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to update wellbeing measurement", "INTERNAL_ERROR")
		return
	}

	response := h.toResponse(measurement)
	middleware.JSONResponder(w, http.StatusOK, response)
}

// Delete handles deleting a wellbeing measurement
func (h *WellbeingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid ID", "INVALID_ID")
		return
	}

	if err := h.wellbeingUC.Delete(id); err != nil {
		h.logger.Error("Failed to delete wellbeing measurement", err, map[string]interface{}{"id": id})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to delete wellbeing measurement", "INTERNAL_ERROR")
		return
	}

	middleware.JSONResponder(w, http.StatusNoContent, nil)
}

func (h *WellbeingHandler) toResponse(m *domain.WellbeingMeasurement) dto.WellbeingMeasurementResponse {
	return dto.WellbeingMeasurementResponse{
		ID:            m.ID,
		UserID:        m.UserID,
		MeasuredAt:    m.MeasuredAt,
		EnergyLevel:   m.EnergyLevel,
		MoodLevel:     m.MoodLevel,
		HungerLevel:   m.HungerLevel,
		SleepHours:    m.SleepHours,
		SleepQuality:  m.SleepQuality,
		DigestionNote: m.DigestionNote,
		StressLevel:   m.StressLevel,
	}
}