package tokens_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"raven.go.invoice-builder/internal/tokens"
)

func TestJWTManager_GenerateAndValidateToken(t *testing.T) {
	secretKey := "super-secret-key-123"
	duration := 1 * time.Hour
	manager := tokens.NewJWTManager(secretKey, duration)

	userID := uuid.New()
	email := "test@example.com"
	role := "admin"
	tokenVersion := 1

	tokenString, err := manager.GenerateToken(userID, email, role, tokenVersion)
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	if tokenString == "" {
		t.Fatal("expected non-empty token string")
	}

	claims, err := manager.ValidateToken(tokenString)
	if err != nil {
		t.Fatalf("unexpected error validating token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("expected UserID %s, got %s", userID, claims.UserID)
	}

	if claims.Email != email {
		t.Errorf("expected Email %s, got %s", email, claims.Email)
	}

	if claims.Role != role {
		t.Errorf("expected Role %s, got %s", role, claims.Role)
	}

	if claims.TokenVersion != tokenVersion {
		t.Errorf("expected TokenVersion %d, got %d", tokenVersion, claims.TokenVersion)
	}

	if claims.Issuer != tokens.TokenIssuer {
		t.Errorf("expected Issuer %s, got %s", tokens.TokenIssuer, claims.Issuer)
	}

	if claims.Subject != userID.String() {
		t.Errorf("expected Subject %s, got %s", userID.String(), claims.Subject)
	}

	if claims.ID == "" {
		t.Error("expected non-empty JWT ID (jti)")
	}
}

func TestJWTManager_DefaultTokenVersion(t *testing.T) {
	manager := tokens.NewJWTManager("secret", 1*time.Hour)
	userID := uuid.New()

	// GenerateToken without tokenVersion should default to 0.
	tokenString, err := manager.GenerateToken(userID, "user@test.com", "user")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	claims, err := manager.ValidateToken(tokenString)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if claims.TokenVersion != 0 {
		t.Errorf("expected default TokenVersion 0, got %d", claims.TokenVersion)
	}
}

func TestJWTManager_ExpiredToken(t *testing.T) {
	secretKey := "super-secret-key-123"
	duration := 1 * time.Millisecond // Expire almost immediately
	manager := tokens.NewJWTManager(secretKey, duration)

	userID := uuid.New()
	tokenString, err := manager.GenerateToken(userID, "expired@example.com", "user")
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	time.Sleep(5 * time.Millisecond)

	_, err = manager.ValidateToken(tokenString)
	if err == nil {
		t.Fatal("expected error validating expired token, got nil")
	}

	if !errors.Is(err, tokens.ErrExpiredToken) {
		t.Errorf("expected ErrExpiredToken, got: %v", err)
	}
}

func TestJWTManager_InvalidToken(t *testing.T) {
	secretKey := "super-secret-key-123"
	manager := tokens.NewJWTManager(secretKey, 1*time.Hour)

	// Malformed token
	_, err := manager.ValidateToken("invalid.token.value")
	if !errors.Is(err, tokens.ErrInvalidToken) {
		t.Errorf("expected ErrInvalidToken for malformed token, got: %v", err)
	}

	// Token signed with wrong secret key
	wrongManager := tokens.NewJWTManager("wrong-secret-key", 1*time.Hour)
	tokenString, err := wrongManager.GenerateToken(uuid.New(), "user@example.com", "user")
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	_, err = manager.ValidateToken(tokenString)
	if !errors.Is(err, tokens.ErrInvalidToken) {
		t.Errorf("expected ErrInvalidToken for wrong signature, got: %v", err)
	}
}
