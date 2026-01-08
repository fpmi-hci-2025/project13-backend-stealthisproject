package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	service := NewAuthService(nil, "test-secret")

	password := "testpassword123"
	hash, err := service.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hash == "" {
		t.Error("Hash should not be empty")
	}

	if hash == password {
		t.Error("Hash should not equal original password")
	}
}

func TestVerifyPassword(t *testing.T) {
	service := NewAuthService(nil, "test-secret")

	password := "testpassword123"
	hash, err := service.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	err = service.VerifyPassword(hash, password)
	if err != nil {
		t.Errorf("Password verification failed: %v", err)
	}

	err = service.VerifyPassword(hash, "wrongpassword")
	if err == nil {
		t.Error("Password verification should fail for wrong password")
	}
}

func TestGenerateToken(t *testing.T) {
	service := NewAuthService(nil, "test-secret")

	token, err := service.GenerateToken(1, "PASSENGER")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Token should not be empty")
	}
}

func TestValidateToken(t *testing.T) {
	service := NewAuthService(nil, "test-secret")

	userID := int64(123)
	role := "PASSENGER"

	token, err := service.GenerateToken(userID, role)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	claims, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected user ID %d, got %d", userID, claims.UserID)
	}

	if claims.Role != role {
		t.Errorf("Expected role %s, got %s", role, claims.Role)
	}
}

