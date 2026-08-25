package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/irsyadjpp/tracking-diet-backend/internal/auth"
	"github.com/irsyadjpp/tracking-diet-backend/internal/delivery/dto"
	"github.com/irsyadjpp/tracking-diet-backend/internal/delivery/middleware"
	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
	"github.com/irsyadjpp/tracking-diet-backend/pkg/logger"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	userRepo   domain.UserRepository
	jwtManager *auth.JWTManager
	logger     *logger.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userRepo domain.UserRepository, jwtManager *auth.JWTManager, log *logger.Logger) *AuthHandler {
	return &AuthHandler{
		userRepo:   userRepo,
		jwtManager: jwtManager,
		logger:     log,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST")
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Validation failed: "+err.Error(), "VALIDATION_ERROR")
		return
	}

	// Check if user already exists
	existingUser, err := h.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		middleware.ErrorResponder(w, http.StatusConflict, "User with this email already exists", "USER_EXISTS")
		return
	}

	// Create user
	user := &domain.User{
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: req.Password, // Will be hashed in repository
		Gender:       req.Gender,
		Role:         domain.RoleUser,
	}

	if req.BirthDate != "" {
		birthDate, err := time.Parse("2006-01-02", req.BirthDate)
		if err == nil {
			user.BirthDate = &birthDate
		}
	}

	if req.HeightCm != nil {
		user.HeightCm = req.HeightCm
	}

	if err := h.userRepo.Create(user); err != nil {
		h.logger.Error("Failed to create user", err, map[string]interface{}{"email": req.Email})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to create user", "INTERNAL_ERROR")
		return
	}

	// Generate JWT token
	token, err := h.jwtManager.GenerateToken(user.ID, user.Email)
	if err != nil {
		h.logger.Error("Failed to generate token", err, map[string]interface{}{"user_id": user.ID})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to generate token", "INTERNAL_ERROR")
		return
	}

	// Prepare response
	response := dto.AuthResponse{
		Token:        token,
		RefreshToken: token, // For simplicity, using same token as refresh token
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		User: dto.UserResponse{
			ID:        user.ID,
			FullName:  user.FullName,
			Email:     user.Email,
			Gender:    user.Gender,
			BirthDate: user.BirthDate,
			HeightCm:  user.HeightCm,
			Role:      string(user.Role),
			CreatedAt: user.CreatedAt,
		},
	}

	h.logger.Info("User registered successfully", map[string]interface{}{"user_id": user.ID, "email": user.Email})
	middleware.JSONResponder(w, http.StatusCreated, response)
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST")
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Validation failed: "+err.Error(), "VALIDATION_ERROR")
		return
	}

	// Get user by email
	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		h.logger.Error("User not found", err, map[string]interface{}{"email": req.Email})
		middleware.ErrorResponder(w, http.StatusUnauthorized, "Invalid credentials", "INVALID_CREDENTIALS")
		return
	}

	// Verify password
	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		h.logger.Warn("Invalid password attempt", map[string]interface{}{"email": req.Email})
		middleware.ErrorResponder(w, http.StatusUnauthorized, "Invalid credentials", "INVALID_CREDENTIALS")
		return
	}

	// Generate JWT token
	token, err := h.jwtManager.GenerateToken(user.ID, user.Email)
	if err != nil {
		h.logger.Error("Failed to generate token", err, map[string]interface{}{"user_id": user.ID})
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to generate token", "INTERNAL_ERROR")
		return
	}

	// Prepare response
	response := dto.AuthResponse{
		Token:        token,
		RefreshToken: token,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		User: dto.UserResponse{
			ID:        user.ID,
			FullName:  user.FullName,
			Email:     user.Email,
			Gender:    user.Gender,
			BirthDate: user.BirthDate,
			HeightCm:  user.HeightCm,
			Role:      string(user.Role),
			CreatedAt: user.CreatedAt,
		},
	}

	h.logger.Info("User logged in successfully", map[string]interface{}{"user_id": user.ID, "email": user.Email})
	middleware.JSONResponder(w, http.StatusOK, response)
}

// Refresh handles token refresh
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST")
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		middleware.ErrorResponder(w, http.StatusBadRequest, "Validation failed: "+err.Error(), "VALIDATION_ERROR")
		return
	}

	// Refresh token
	newToken, err := h.jwtManager.RefreshToken(req.RefreshToken)
	if err != nil {
		h.logger.Error("Failed to refresh token", err, nil)
		middleware.ErrorResponder(w, http.StatusUnauthorized, "Invalid or expired token", "INVALID_TOKEN")
		return
	}

	// Get user info from token
	claims, err := h.jwtManager.ValidateToken(newToken)
	if err != nil {
		h.logger.Error("Failed to validate new token", err, nil)
		middleware.ErrorResponder(w, http.StatusInternalServerError, "Failed to generate token", "INTERNAL_ERROR")
		return
	}

	user, err := h.userRepo.GetByID(claims.UserID)
	if err != nil {
		h.logger.Error("User not found", err, map[string]interface{}{"user_id": claims.UserID})
		middleware.ErrorResponder(w, http.StatusNotFound, "User not found", "USER_NOT_FOUND")
		return
	}

	// Prepare response
	response := dto.AuthResponse{
		Token:        newToken,
		RefreshToken: newToken,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		User: dto.UserResponse{
			ID:        user.ID,
			FullName:  user.FullName,
			Email:     user.Email,
			Gender:    user.Gender,
			BirthDate: user.BirthDate,
			HeightCm:  user.HeightCm,
			Role:      string(user.Role),
			CreatedAt: user.CreatedAt,
		},
	}

	h.logger.Info("Token refreshed successfully", map[string]interface{}{"user_id": user.ID})
	middleware.JSONResponder(w, http.StatusOK, response)
}

// Me handles getting current user profile
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		middleware.ErrorResponder(w, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		h.logger.Error("User not found", err, map[string]interface{}{"user_id": userID})
		middleware.ErrorResponder(w, http.StatusNotFound, "User not found", "USER_NOT_FOUND")
		return
	}

	response := dto.UserResponse{
		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		Gender:    user.Gender,
		BirthDate: user.BirthDate,
		HeightCm:  user.HeightCm,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
	}

	middleware.JSONResponder(w, http.StatusOK, response)
}