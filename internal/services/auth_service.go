package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"raven.go.invoice-builder/internal/mailer"
	"raven.go.invoice-builder/internal/models"
	"raven.go.invoice-builder/internal/tokens"
)

const (
	// RefreshTokenDuration is how long a refresh token stays valid.
	RefreshTokenDuration = 30 * 24 * time.Hour // 30 days
)

// AuthTokenPair holds both the short-lived access token and the long-lived refresh token.
type AuthTokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthService interface {
	Authenticate(ctx context.Context, email, password string) (*models.User, *AuthTokenPair, error)
	Register(ctx context.Context, name, email, password string) (*models.User, error)
	Activate(ctx context.Context, token string) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetBusinessProfile(ctx context.Context, userID uuid.UUID) (*models.BusinessProfile, error)
	GenerateToken(userID uuid.UUID, email, role string) (string, error)
	ValidateToken(tokenStr string) (*tokens.Claims, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (*models.User, *AuthTokenPair, error)
}

type authService struct {
	models     models.Models
	mailer     mailer.Mailer
	jwtManager *tokens.JWTManager
}

func NewAuthService(models models.Models, mailer mailer.Mailer, jwtManager *tokens.JWTManager) AuthService {
	return &authService{
		models:     models,
		mailer:     mailer,
		jwtManager: jwtManager,
	}
}

func (s *authService) Authenticate(ctx context.Context, email, password string) (*models.User, *AuthTokenPair, error) {
	userID, err := s.models.Users.Authenticate(ctx, email, password)
	if err != nil {
		return nil, nil, err
	}

	user, err := s.models.Users.GetByID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	if !user.IsActivated {
		return nil, nil, errors.New("user account is not activated")
	}

	tokenPair, err := s.generateTokenPair(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokenPair, nil
}

// generateTokenPair creates both an access token (JWT) and a refresh token for the user.
func (s *authService) generateTokenPair(ctx context.Context, user *models.User) (*AuthTokenPair, error) {
	if s.jwtManager == nil {
		return nil, nil
	}

	roleName := "user"
	if user.IsAdmin() {
		roleName = "admin"
	}

	accessToken, err := s.jwtManager.GenerateToken(user.ID, user.Email, roleName, user.TokenVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.models.RefreshTokens.Create(ctx, user.ID, RefreshTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &AuthTokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshAccessToken validates a refresh token and issues a new access + refresh token pair.
// The old refresh token is revoked (rotation) to prevent reuse.
func (s *authService) RefreshAccessToken(ctx context.Context, refreshToken string) (*models.User, *AuthTokenPair, error) {
	rt, err := s.models.RefreshTokens.Validate(ctx, refreshToken)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Revoke the old refresh token (token rotation).
	_ = s.models.RefreshTokens.Revoke(ctx, rt.ID)

	user, err := s.models.Users.GetByID(ctx, rt.UserID)
	if err != nil {
		return nil, nil, err
	}

	if !user.IsActivated {
		return nil, nil, errors.New("user account is not activated")
	}

	tokenPair, err := s.generateTokenPair(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokenPair, nil
}

func (s *authService) Register(ctx context.Context, name, email, password string) (*models.User, error) {
	activationToken := uuid.New().String()
	activationExpiry := time.Now().Add(3 * 24 * time.Hour)

	user := &models.User{
		Name:             name,
		Email:            email,
		PasswordHash:     password,
		IsActivated:      false,
		ActivationToken:  activationToken,
		ActivationExpiry: &activationExpiry,
	}

	err := s.models.Users.Insert(ctx, user)
	if err != nil {
		return nil, err
	}

	// Dispatch email in background
	userEmail := user.Email
	userName := user.Name
	go func() {
		_ = s.mailer.SendMail(userEmail, "user_welcome.tmpl", map[string]any{
			"Name":            userName,
			"ActivationToken": activationToken,
		})
	}()

	return user, nil
}

func (s *authService) Activate(ctx context.Context, token string) (*models.User, error) {
	user, err := s.models.Users.GetByActivationToken(ctx, token)
	if err != nil {
		return nil, err
	}

	err = s.models.Users.ActivateUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.models.Users.GetByID(ctx, id)
}

func (s *authService) GetBusinessProfile(ctx context.Context, userID uuid.UUID) (*models.BusinessProfile, error) {
	return s.models.BusinessProfiles.GetByUserID(ctx, userID)
}

func (s *authService) GenerateToken(userID uuid.UUID, email, role string) (string, error) {
	if s.jwtManager == nil {
		return "", errors.New("jwtManager uninitialized")
	}
	return s.jwtManager.GenerateToken(userID, email, role)
}

func (s *authService) ValidateToken(tokenStr string) (*tokens.Claims, error) {
	if s.jwtManager == nil {
		return nil, errors.New("jwtManager uninitialized")
	}
	return s.jwtManager.ValidateToken(tokenStr)
}
