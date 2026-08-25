package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/irsyadjpp/tracking-diet-backend/internal/delivery/dto"
	"github.com/irsyadjpp/tracking-diet-backend/internal/delivery/middleware"
	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
	"github.com/irsyadjpp/tracking-diet-backend/internal/usecase"
	"github.com/irsyadjpp/tracking-diet-backend/pkg/logger"
)

// UserHandler handles user HTTP requests
type UserHandler struct {
	userUC *usecase.UserUsecase
	logger *logger.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUC *usecase.UserUsecase, log *logger.Logger) *UserHandler {
	return &UserHandler{
		userUC: userUC,
		logger: log,
	}
}

// GetProfile handles getting current user profile
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		middleware.ErrorResponder(w, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}

	user, err := h.userUC.GetByID(userID)
	if err != nil {
		h.logger.Error("User not found", err, map[string]interface{}{"user_id": userID})
		middleware.ErrorResponder(w, http.StatusNotFound, "User not found", "NOT_FOUND")
		return
	}

	response := h.toResponse(user)
	middleware.JSONResponder(w, http.StatusOK, response)
}

// UpdateProfile handles updating user profile
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		middleware.ErrorResponder(w, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST")
		return
	}

	if err := req.Validate(); err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Validation failed: "+err.Error(), "VALIDATION_ERROR")
		return
	}

	user, err := h.userUC.UpdateProfile(userID, req.FullName, req.Gender, &req.BirthDate, req.HeightCm)
	if err != nil {
		h.logger.Error("Failed to update user profile", err, map[string]interface{}{"user_id": userID})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to update user profile", "INTERNAL_ERROR")
		return
	}

	response := h.toResponse(user)
	middleware.JSONResponder(w, http.StatusOK, response)
}

// GetUserByID handles getting a user by ID (admin only)
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid ID", "INVALID_ID")
		return
	}

	user, err := h.userUC.GetByID(id)
	if err != nil {
		h.logger.Error("User not found", err, map[string]interface{}{"id": id})
		middleware.ErrorResponder(w, http.StatusNotFound, "User not found", "NOT_FOUND")
		return
	}

	response := h.toResponse(user)
	middleware.JSONResponder(w, http.StatusOK, response)
}

// ListUsers handles listing all users (admin only)
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userUC.List()
	if err != nil {
		h.logger.Error("Failed to list users", err, nil)
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to list users", "INTERNAL_ERROR")
		return
	}

	response := make([]dto.UserResponse, len(users))
	for i, user := range users {
		response[i] = h.toResponse(&user)
	}

	middleware.JSONResponder(w, http.StatusOK, response)
}

// DeleteUser handles deleting a user (admin only)
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid ID", "INVALID_ID")
		return
	}

	if err := h.userUC.Delete(id); err != nil {
		h.logger.Error("Failed to delete user", err, map[string]interface{}{"id": id})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to delete user", "INTERNAL_ERROR")
		return
	}

	middleware.JSONResponder(w, http.StatusNoContent, nil)
}

func (h *UserHandler) toResponse(user *domain.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		Gender:    user.Gender,
		BirthDate: user.BirthDate,
		HeightCm:  user.HeightCm,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
	}
}