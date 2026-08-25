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

// MetabolicHandler handles metabolic measurement HTTP requests
type MetabolicHandler struct {
	metabolicUC *usecase.MetabolicUsecase
	logger      *logger.Logger
}

// NewMetabolicHandler creates a new metabolic handler
func NewMetabolicHandler(metabolicUC *usecase.MetabolicUsecase, log *logger.Logger) *MetabolicHandler {
	return &MetabolicHandler{
		metabolicUC: metabolicUC,
		logger:      log,
	}
}

// Create handles metabolic measurement creation
func (h *MetabolicHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		middleware.ErrorResponder(w, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}

	var req dto.MetabolicMeasurementRequest
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

	measurement, err := h.metabolicUC.Create(userID, req.SystolicBP, req.DiastolicBP, req.FastingGlucose, req.Hba1cPct, req.CholesterolTotal, req.LDL, req.HDL, req.Triglycerides, req.UricAcid, req.LiverFunctionNote, req.KidneyFunctionNote, measuredAt)
	if err != nil {
		h.logger.Error("Failed to create metabolic measurement", err, map[string]interface{}{"user_id": userID})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to create metabolic measurement", "INTERNAL_ERROR")
		return
	}

	response := h.toResponse(measurement)
	middleware.JSONResponder(w, http.StatusCreated, response)
}

// GetByID handles getting a metabolic measurement by ID
func (h *MetabolicHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid ID", "INVALID_ID")
		return
	}

	measurement, err := h.metabolicUC.GetByID(id)
	if err != nil {
		h.logger.Error("Metabolic measurement not found", err, map[string]interface{}{"id": id})
		middleware.ErrorResponder(w, http.StatusNotFound, "Metabolic measurement not found", "NOT_FOUND")
		return
	}

	response := h.toResponse(measurement)
	middleware.JSONResponder(w, http.StatusOK, response)
}

// List handles listing metabolic measurements for the authenticated user
func (h *MetabolicHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		middleware.ErrorResponder(w, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}

	// Check for date range filters
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var measurements []domain.MetabolicMeasurement
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
		measurements, err = h.metabolicUC.ListByUserWithDateRange(userID, startDate, endDate)
	} else {
		measurements, err = h.metabolicUC.ListByUser(userID)
	}

	if err != nil {
		h.logger.Error("Failed to list metabolic measurements", err, map[string]interface{}{"user_id": userID})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to list metabolic measurements", "INTERNAL_ERROR")
		return
	}

	response := make([]dto.MetabolicMeasurementResponse, len(measurements))
	for i, m := range measurements {
		response[i] = h.toResponse(&m)
	}

	middleware.JSONResponder(w, http.StatusOK, response)
}

// Update handles updating a metabolic measurement
func (h *MetabolicHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid ID", "INVALID_ID")
		return
	}

	var req dto.MetabolicMeasurementRequest
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

	measurement, err := h.metabolicUC.Update(id, req.SystolicBP, req.DiastolicBP, req.FastingGlucose, req.Hba1cPct, req.CholesterolTotal, req.LDL, req.HDL, req.Triglycerides, req.UricAcid, req.LiverFunctionNote, req.KidneyFunctionNote, measuredAt)
	if err != nil {
		h.logger.Error("Failed to update metabolic measurement", err, map[string]interface{}{"id": id})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to update metabolic measurement", "INTERNAL_ERROR")
		return
	}

	response := h.toResponse(measurement)
	middleware.JSONResponder(w, http.StatusOK, response)
}

// Delete handles deleting a metabolic measurement
func (h *MetabolicHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid ID", "INVALID_ID")
		return
	}

	if err := h.metabolicUC.Delete(id); err != nil {
		h.logger.Error("Failed to delete metabolic measurement", err, map[string]interface{}{"id": id})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to delete metabolic measurement", "INTERNAL_ERROR")
		return
	}

	middleware.JSONResponder(w, http.StatusNoContent, nil)
}

func (h *MetabolicHandler) toResponse(m *domain.MetabolicMeasurement) dto.MetabolicMeasurementResponse {
	return dto.MetabolicMeasurementResponse{
		ID:                  m.ID,
		UserID:              m.UserID,
		MeasuredAt:          m.MeasuredAt,
		SystolicBP:          m.SystolicBP,
		DiastolicBP:         m.DiastolicBP,
		FastingGlucose:      m.FastingGlucose,
		Hba1cPct:            m.Hba1cPct,
		CholesterolTotal:    m.CholesterolTotal,
		LDL:                 m.LDL,
		HDL:                 m.HDL,
		Triglycerides:       m.Triglycerides,
		UricAcid:            m.UricAcid,
		LiverFunctionNote:   m.LiverFunctionNote,
		KidneyFunctionNote:  m.KidneyFunctionNote,
	}
}