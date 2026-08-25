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

// LabTestHandler handles lab test HTTP requests
type LabTestHandler struct {
	labTestUC *usecase.LabTestUsecase
	logger    *logger.Logger
}

// NewLabTestHandler creates a new lab test handler
func NewLabTestHandler(labTestUC *usecase.LabTestUsecase, log *logger.Logger) *LabTestHandler {
	return &LabTestHandler{
		labTestUC: labTestUC,
		logger:    log,
	}
}

// Create handles lab test creation
func (h *LabTestHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		middleware.ErrorResponder(w, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}

	var req dto.LabTestRequest
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

	labTest, err := h.labTestUC.Create(userID, req.TestName, req.ResultValue, req.Unit, req.ReferenceRange, measuredAt)
	if err != nil {
		h.logger.Error("Failed to create lab test", err, map[string]interface{}{"user_id": userID})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to create lab test", "INTERNAL_ERROR")
		return
	}

	response := h.toResponse(labTest)
	middleware.JSONResponder(w, http.StatusCreated, response)
}

// GetByID handles getting a lab test by ID
func (h *LabTestHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid ID", "INVALID_ID")
		return
	}

	labTest, err := h.labTestUC.GetByID(id)
	if err != nil {
		h.logger.Error("Lab test not found", err, map[string]interface{}{"id": id})
		middleware.ErrorResponder(w, http.StatusNotFound, "Lab test not found", "NOT_FOUND")
		return
	}

	response := h.toResponse(labTest)
	middleware.JSONResponder(w, http.StatusOK, response)
}

// List handles listing lab tests for the authenticated user
func (h *LabTestHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		middleware.ErrorResponder(w, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}

	// Check for test name filter
	testName := r.URL.Query().Get("test_name")

	// Check for date range filters
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var labTests []domain.LabTest
	var err error

	if testName != "" {
		labTests, err = h.labTestUC.ListByUserWithTestName(userID, testName)
	} else if startDateStr != "" && endDateStr != "" {
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
		labTests, err = h.labTestUC.ListByUserWithDateRange(userID, startDate, endDate)
	} else {
		labTests, err = h.labTestUC.ListByUser(userID)
	}

	if err != nil {
		h.logger.Error("Failed to list lab tests", err, map[string]interface{}{"user_id": userID})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to list lab tests", "INTERNAL_ERROR")
		return
	}

	response := make([]dto.LabTestResponse, len(labTests))
	for i, lt := range labTests {
		response[i] = h.toResponse(&lt)
	}

	middleware.JSONResponder(w, http.StatusOK, response)
}

// Update handles updating a lab test
func (h *LabTestHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid ID", "INVALID_ID")
		return
	}

	var req dto.LabTestRequest
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

	labTest, err := h.labTestUC.Update(id, &req.TestName, req.ResultValue, req.Unit, req.ReferenceRange, measuredAt)
	if err != nil {
		h.logger.Error("Failed to update lab test", err, map[string]interface{}{"id": id})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to update lab test", "INTERNAL_ERROR")
		return
	}

	response := h.toResponse(labTest)
	middleware.JSONResponder(w, http.StatusOK, response)
}

// Delete handles deleting a lab test
func (h *LabTestHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid ID", "INVALID_ID")
		return
	}

	if err := h.labTestUC.Delete(id); err != nil {
		h.logger.Error("Failed to delete lab test", err, map[string]interface{}{"id": id})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to delete lab test", "INTERNAL_ERROR")
		return
	}

	middleware.JSONResponder(w, http.StatusNoContent, nil)
}

func (h *LabTestHandler) toResponse(lt *domain.LabTest) dto.LabTestResponse {
	return dto.LabTestResponse{
		ID:             lt.ID,
		UserID:         lt.UserID,
		TestName:       lt.TestName,
		ResultValue:    lt.ResultValue,
		Unit:           lt.Unit,
		ReferenceRange: lt.ReferenceRange,
		MeasuredAt:     lt.MeasuredAt,
	}
}