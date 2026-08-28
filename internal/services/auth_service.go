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
	ResendVerificationEmail(ctx context.Context, email string) error
	RequestPasswordReset(ctx context.Context, email string) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetBusinessProfile(ctx context.Context, userID uuid.UUID) (*models.BusinessProfile, error)
	UpdateBusinessProfile(ctx context.Context, profile *models.BusinessProfile) error
	UpdateUser(ctx context.Context, user *models.User) error
	GenerateToken(userID uuid.UUID, email, role string) (string, error)
	ValidateToken(tokenStr string) (*tokens.Claims, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (*models.User, *AuthTokenPair, error)
	SendAdminLoginAlert(ctx context.Context, user *models.User, ip, userAgent string) error
}

type authService struct {
	models     models.Models
	mailer     mailer.Mailer
	jwtManager *tokens.JWTManager
	frontendURL string
}

func NewAuthService(models models.Models, mailer mailer.Mailer, jwtManager *tokens.JWTManager, frontendURL string) AuthService {
	return &authService{
		models:     models,
		mailer:     mailer,
		jwtManager: jwtManager,
		frontendURL: frontendURL,
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
		return nil, nil, models.ErrAccountNotActivated
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
		return nil, nil, models.ErrAccountNotActivated
	}

	tokenPair, err := s.generateTokenPair(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokenPair, nil
}

func (s *authService) Register(ctx context.Context, name, email, password string) (*models.User, error) {
	activationToken := uuid.New().String()
	activationExpiry := time.Now().Add(1 * time.Hour)

	role, err := s.models.Roles.GetByName(ctx, "user")
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			role = &models.Role{
				Name:        "user",
				Description: "Standard User",
			}
			err = s.models.Roles.Create(ctx, role)
			if err != nil {
				return nil, fmt.Errorf("failed to create default role: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to retrieve default role: %w", err)
		}
	}

	user := &models.User{
		Name:             name,
		Email:            email,
		PasswordHash:     password,
		RoleID:           role.ID,
		IsActivated:      false,
		ActivationToken:  activationToken,
		ActivationExpiry: &activationExpiry,
	}

	err = s.models.Users.Insert(ctx, user)
	if err != nil {
		return nil, err
	}

	// Grant 5 free credits upon sign up
	_ = s.models.CreditTxn.AdminGrantCredits(ctx, user.ID, 5, "Sign-up Bonus")

	// Dispatch email in background
	userEmail := user.Email
	userName := user.Name
	go func() {
		err := s.mailer.SendMail(userEmail, "user_welcome.tmpl", map[string]any{
			"Name":            userName,
			"ActivationToken": activationToken,
			"FrontendURL":     s.frontendURL,
		})
		if err != nil {
			fmt.Printf("Failed to send welcome email to %s: %v\n", userEmail, err)
		}
	}()

	return user, nil
}

func (s *authService) ResendVerificationEmail(ctx context.Context, email string) error {
	user, err := s.models.Users.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user.IsActivated {
		return nil
	}

	if user.ActivationExpiry != nil {
		timeRemaining := time.Until(*user.ActivationExpiry)
		// Since tokens are valid for 1 hour, if timeRemaining is between 55-60 mins, it was just sent
		if timeRemaining > 55*time.Minute && timeRemaining <= 1*time.Hour {
			return models.ErrActivationLinkCooldown
		}
	}

	activationToken := uuid.New().String()
	activationExpiry := time.Now().Add(1 * time.Hour)
	user.ActivationToken = activationToken
	user.ActivationExpiry = &activationExpiry

	err = s.models.Users.Update(ctx, user)
	if err != nil {
		return err
	}

	go func() {
		err := s.mailer.SendMail(user.Email, "user_welcome.tmpl", map[string]any{
			"Name":            user.Name,
			"ActivationToken": activationToken,
			"FrontendURL":     s.frontendURL,
		})
		if err != nil {
			fmt.Printf("Failed to send verification email to %s: %v\n", user.Email, err)
		}
	}()
	return nil
}

func (s *authService) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := s.models.Users.GetByEmail(ctx, email)
	if err != nil {
		return err
	}

	resetToken := uuid.New().String()
	resetExpiry := time.Now().Add(1 * time.Hour)
	user.ResetToken = resetToken
	user.ResetExpiry = &resetExpiry

	err = s.models.Users.Update(ctx, user)
	if err != nil {
		return err
	}

	go func() {
		err := s.mailer.SendMail(user.Email, "password_reset.tmpl", map[string]any{
			"Name":       user.Name,
			"ResetToken": resetToken,
			"FrontendURL": s.frontendURL,
		})
		if err != nil {
			fmt.Printf("Failed to send password reset email to %s: %v\n", user.Email, err)
		}
	}()
	return nil
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

func (s *authService) UpdateBusinessProfile(ctx context.Context, profile *models.BusinessProfile) error {
	return s.models.BusinessProfiles.Upsert(ctx, profile)
}

func (s *authService) UpdateUser(ctx context.Context, user *models.User) error {
	return s.models.Users.Update(ctx, user)
}

func (s *authService) GenerateToken(userID uuid.UUID, email, role string) (string, error) {
	if s.jwtManager == nil {
		return "", models.ErrJWTUninitialized
	}
	return s.jwtManager.GenerateToken(userID, email, role)
}

func (s *authService) ValidateToken(tokenStr string) (*tokens.Claims, error) {
	if s.jwtManager == nil {
		return nil, models.ErrJWTUninitialized
	}
	return s.jwtManager.ValidateToken(tokenStr)
}

func (s *authService) SendAdminLoginAlert(ctx context.Context, user *models.User, ip, userAgent string) error {
	data := map[string]any{
		"Name":      user.Name,
		"Email":     user.Email,
		"IP":        ip,
		"UserAgent": userAgent,
		"Time":      time.Now().Format(time.RFC1123),
	}
	return s.mailer.SendMail(user.Email, "user_welcome.tmpl", data)
}
