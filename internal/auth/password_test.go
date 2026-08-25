package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "testPassword123"
	
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	
	if hashedPassword == password {
		t.Error("Hashed password should be different from original")
	}
	
	if len(hashedPassword) == 0 {
		t.Error("Hashed password should not be empty")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "testPassword123"
	
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	
	// Test correct password
	if !CheckPassword(password, hashedPassword) {
		t.Error("CheckPassword should return true for correct password")
	}
	
	// Test incorrect password
	if CheckPassword("wrongPassword", hashedPassword) {
		t.Error("CheckPassword should return false for incorrect password")
	}
}

func TestValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		name    string
		password string
		wantErr bool
	}{
		{
			name:    "valid password",
			password: "validPassword123",
			wantErr: false,
		},
		{
			name:    "too short",
			password: "short",
			wantErr: true,
		},
		{
			name:    "exactly 8 characters",
			password: "12345678",
			wantErr: false,
		},
		{
			name:    "empty password",
			password: "",
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePasswordStrength(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePasswordStrength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}