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

// BodyMeasurementHandler handles body measurement HTTP requests
type BodyMeasurementHandler struct {
	bodyUC *usecase.BodyMeasurementUsecase
	logger *logger.Logger
}

// NewBodyMeasurementHandler creates a new body measurement handler
func NewBodyMeasurementHandler(bodyUC *usecase.BodyMeasurementUsecase, log *logger.Logger) *BodyMeasurementHandler {
	return &BodyMeasurementHandler{
		bodyUC: bodyUC,
		logger: log,
	}
}

// Create handles body measurement creation
func (h *BodyMeasurementHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		middleware.ErrorResponder(w, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}

	var req dto.BodyMeasurementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST")
		return
	}

	if err := req.Validate(); err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Validation failed: "+err.Error(), "VALIDATION_ERROR")
		return
	}

	measuredAt, err := time.Parse("2006-01-02T15:04:05Z", req.MeasuredAt)
	if err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid measured_at format", "INVALID_DATE_FORMAT")
		return
	}

	measurement, err := h.bodyUC.Create(userID, req.WeightKg, req.WaistCm, req.HipCm, req.ChestCm, req.ThighCm, req.ArmCm, req.BodyFatPct, req.MuscleMassKg, req.VisceralFat, req.SkinfoldMm, measuredAt)
	if err != nil {
		h.logger.Error("Failed to create body measurement", err, map[string]interface{}{"user_id": userID})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to create body measurement", "INTERNAL_ERROR")
		return
	}

	response := h.toResponse(measurement)
	middleware.JSONResponder(w, http.StatusCreated, response)
}

// GetByID handles getting a body measurement by ID
func (h *BodyMeasurementHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid ID", "INVALID_ID")
		return
	}

	measurement, err := h.bodyUC.GetByID(id)
	if err != nil {
		h.logger.Error("Body measurement not found", err, map[string]interface{}{"id": id})
		middleware.ErrorResponder(w, http.StatusNotFound, "Body measurement not found", "NOT_FOUND")
		return
	}

	response := h.toResponse(measurement)
	middleware.JSONResponder(w, http.StatusOK, response)
}

// List handles listing body measurements for the authenticated user
func (h *BodyMeasurementHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		middleware.ErrorResponder(w, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}

	// Check for date range filters
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var measurements []domain.BodyMeasurement
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
		measurements, err = h.bodyUC.ListByUserWithDateRange(userID, startDate, endDate)
	} else {
		measurements, err = h.bodyUC.ListByUser(userID)
	}

	if err != nil {
		h.logger.Error("Failed to list body measurements", err, map[string]interface{}{"user_id": userID})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to list body measurements", "INTERNAL_ERROR")
		return
	}

	response := make([]dto.BodyMeasurementResponse, len(measurements))
	for i, m := range measurements {
		response[i] = h.toResponse(&m)
	}

	middleware.JSONResponder(w, http.StatusOK, response)
}

// Update handles updating a body measurement
func (h *BodyMeasurementHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid ID", "INVALID_ID")
		return
	}

	var req dto.BodyMeasurementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST")
		return
	}

	if err := req.Validate(); err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Validation failed: "+err.Error(), "VALIDATION_ERROR")
		return
	}

	measuredAt, err := time.Parse("2006-01-02T15:04:05Z", req.MeasuredAt)
	if err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid measured_at format", "INVALID_DATE_FORMAT")
		return
	}

	measurement, err := h.bodyUC.Update(id, req.WeightKg, req.WaistCm, req.HipCm, req.ChestCm, req.ThighCm, req.ArmCm, req.BodyFatPct, req.MuscleMassKg, req.VisceralFat, req.SkinfoldMm, measuredAt)
	if err != nil {
		h.logger.Error("Failed to update body measurement", err, map[string]interface{}{"id": id})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to update body measurement", "INTERNAL_ERROR")
		return
	}

	response := h.toResponse(measurement)
	middleware.JSONResponder(w, http.StatusOK, response)
}

// Delete handles deleting a body measurement
func (h *BodyMeasurementHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid ID", "INVALID_ID")
		return
	}

	if err := h.bodyUC.Delete(id); err != nil {
		h.logger.Error("Failed to delete body measurement", err, map[string]interface{}{"id": id})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to delete body measurement", "INTERNAL_ERROR")
		return
	}

	middleware.JSONResponder(w, http.StatusNoContent, nil)
}

func (h *BodyMeasurementHandler) toResponse(m *domain.BodyMeasurement) dto.BodyMeasurementResponse {
	return dto.BodyMeasurementResponse{
		ID:            m.ID,
		UserID:        m.UserID,
		MeasuredAt:    m.MeasuredAt,
		WeightKg:      m.WeightKg,
		BMI:           m.BMI,
		WaistCm:       m.WaistCm,
		HipCm:         m.HipCm,
		ChestCm:       m.ChestCm,
		ThighCm:       m.ThighCm,
		ArmCm:         m.ArmCm,
		BodyFatPct:    m.BodyFatPct,
		MuscleMassKg:  m.MuscleMassKg,
		VisceralFat:   m.VisceralFat,
		WHR:           m.WHR,
		SkinfoldMm:    m.SkinfoldMm,
	}
}