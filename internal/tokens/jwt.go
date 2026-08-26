package tokens

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	// TokenIssuer identifies this service as the JWT issuer.
	TokenIssuer = "raven.go.invoice-builder"
)

var (
	ErrInvalidToken = errors.New("invalid authentication token")
	ErrExpiredToken = errors.New("authentication token has expired")
)

type Claims struct {
	UserID       uuid.UUID `json:"user_id"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	TokenVersion int       `json:"token_version"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secretKey     []byte
	tokenDuration time.Duration
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	if tokenDuration <= 0 {
		tokenDuration = 30 * 24 * time.Hour // Default 30 days
	}
	return &JWTManager{
		secretKey:     []byte(secretKey),
		tokenDuration: tokenDuration,
	}
}

// GenerateToken creates a signed JWT for a given user.
// tokenVersion should match the user's current TokenVersion from the database.
func (m *JWTManager) GenerateToken(userID uuid.UUID, email, role string, tokenVersion ...int) (string, error) {
	now := time.Now().UTC()
	tv := 0
	if len(tokenVersion) > 0 {
		tv = tokenVersion[0]
	}
	claims := Claims{
		UserID:       userID,
		Email:        email,
		Role:         role,
		TokenVersion: tv,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    TokenIssuer,
			Subject:   userID.String(),
			ID:        uuid.New().String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.tokenDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
}

// ValidateToken verifies the JWT signature and expiration, returning Claims.
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
