package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/irsyadjpp/tracking-diet-backend/internal/auth"
)

// AuthMiddleware validates JWT tokens and adds user context to requests
type AuthMiddleware struct {
	jwtManager *auth.JWTManager
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(jwtManager *auth.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

// Authenticate validates the JWT token from the Authorization header
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Extract token from "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		claims, err := m.jwtManager.ValidateToken(tokenString)
		if err != nil {
			if err == auth.ErrExpiredToken {
				http.Error(w, "Token has expired", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user context to request
		ctx := r.Context()
		ctx = contextWithUserID(ctx, claims.UserID)
		ctx = contextWithEmail(ctx, claims.Email)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuthenticate is like Authenticate but doesn't fail if no token is provided
func (m *AuthMiddleware) OptionalAuthenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// No token provided, continue without authentication
			next.ServeHTTP(w, r)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			next.ServeHTTP(w, r)
			return
		}

		tokenString := parts[1]
		claims, err := m.jwtManager.ValidateToken(tokenString)
		if err != nil {
			// Invalid token, continue without authentication
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		ctx = contextWithUserID(ctx, claims.UserID)
		ctx = contextWithEmail(ctx, claims.Email)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Context key types
type contextKey string

const (
	userIDKey contextKey = "userID"
	emailKey  contextKey = "email"
)

func contextWithUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func contextWithEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, emailKey, email)
}

// GetUserIDFromContext retrieves the user ID from the request context
func GetUserIDFromContext(r *http.Request) (int, bool) {
	userID, ok := r.Context().Value(userIDKey).(int)
	return userID, ok
}

// GetEmailFromContext retrieves the email from the request context
func GetEmailFromContext(r *http.Request) (string, bool) {
	email, ok := r.Context().Value(emailKey).(string)
	return email, ok
}