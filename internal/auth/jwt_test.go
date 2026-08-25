package auth

import (
	"testing"
	"time"
)

func TestNewJWTManager(t *testing.T) {
	secret := "test-secret-key"
	expiration := 24 * time.Hour
	
	manager := NewJWTManager(secret, expiration)
	
	if manager == nil {
		t.Fatal("NewJWTManager should not return nil")
	}
	
	if string(manager.secretKey) != secret {
		t.Error("Secret key should match the provided secret")
	}
	
	if manager.expiration != expiration {
		t.Error("Expiration should match the provided expiration")
	}
}

func TestGenerateToken(t *testing.T) {
	manager := NewJWTManager("test-secret-key", 24*time.Hour)
	
	userID := 123
	email := "test@example.com"
	
	token, err := manager.GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	
	if token == "" {
		t.Error("Token should not be empty")
	}
	
	if len(token) < 50 {
		t.Error("Token should be reasonably long")
	}
}

func TestValidateToken(t *testing.T) {
	manager := NewJWTManager("test-secret-key", 24*time.Hour)
	
	userID := 123
	email := "test@example.com"
	
	token, err := manager.GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	
	claims, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}
	
	if claims.UserID != userID {
		t.Errorf("UserID should be %d, got %d", userID, claims.UserID)
	}
	
	if claims.Email != email {
		t.Errorf("Email should be %s, got %s", email, claims.Email)
	}
}

func TestValidateInvalidToken(t *testing.T) {
	manager := NewJWTManager("test-secret-key", 24*time.Hour)
	
	invalidToken := "invalid.token.string"
	
	_, err := manager.ValidateToken(invalidToken)
	if err == nil {
		t.Error("ValidateToken should return error for invalid token")
	}
	
	if err != ErrInvalidToken {
		t.Errorf("Error should be ErrInvalidToken, got %v", err)
	}
}

func TestRefreshToken(t *testing.T) {
	manager := NewJWTManager("test-secret-key", 24*time.Hour)
	
	userID := 123
	email := "test@example.com"
	
	originalToken, err := manager.GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	
	newToken, err := manager.RefreshToken(originalToken)
	if err != nil {
		t.Fatalf("RefreshToken failed: %v", err)
	}
	
	if newToken == "" {
		t.Error("New token should not be empty")
	}
	
	if newToken == originalToken {
		t.Error("New token should be different from original token")
	}
	
	// Validate the new token
	claims, err := manager.ValidateToken(newToken)
	if err != nil {
		t.Fatalf("ValidateToken failed for new token: %v", err)
	}
	
	if claims.UserID != userID {
		t.Errorf("UserID should be %d, got %d", userID, claims.UserID)
	}
}