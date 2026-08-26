package models

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RefreshToken stores hashed refresh tokens in the database for secure token rotation.
type RefreshToken struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	TokenHash string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"-"`
	ExpiresAt time.Time      `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time      `json:"created_at"`
	RevokedAt *time.Time     `json:"revoked_at,omitempty"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type RefreshTokenInterface interface {
	Create(ctx context.Context, userID uuid.UUID, duration time.Duration) (plainToken string, err error)
	Validate(ctx context.Context, plainToken string) (*RefreshToken, error)
	Revoke(ctx context.Context, tokenID uuid.UUID) error
	RevokeAllForUser(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

type RefreshTokenModel struct {
	DB *gorm.DB
}

// generateToken creates a cryptographically secure random token string.
func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// hashTokenSHA256 produces a hex-encoded SHA-256 hash of the plain token.
// Only the hash is stored in the database so a DB leak doesn't compromise tokens.
func hashTokenSHA256(plain string) string {
	h := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(h[:])
}

// Create generates a new refresh token for the user and stores its hash.
func (m *RefreshTokenModel) Create(ctx context.Context, userID uuid.UUID, duration time.Duration) (string, error) {
	plain, err := generateToken()
	if err != nil {
		return "", err
	}

	rt := &RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: hashTokenSHA256(plain),
		ExpiresAt: time.Now().Add(duration),
	}

	if err := m.DB.WithContext(ctx).Create(rt).Error; err != nil {
		return "", err
	}
	return plain, nil
}

// Validate looks up a refresh token by its hash and checks expiration/revocation.
func (m *RefreshTokenModel) Validate(ctx context.Context, plainToken string) (*RefreshToken, error) {
	hash := hashTokenSHA256(plainToken)
	var rt RefreshToken
	err := m.DB.WithContext(ctx).Where("token_hash = ?", hash).First(&rt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	if rt.RevokedAt != nil {
		return nil, errors.New("refresh token has been revoked")
	}

	if time.Now().After(rt.ExpiresAt) {
		return nil, errors.New("refresh token has expired")
	}

	return &rt, nil
}

// Revoke marks a single refresh token as revoked.
func (m *RefreshTokenModel) Revoke(ctx context.Context, tokenID uuid.UUID) error {
	now := time.Now()
	return m.DB.WithContext(ctx).Model(&RefreshToken{}).Where("id = ?", tokenID).
		Update("revoked_at", &now).Error
}

// RevokeAllForUser revokes all refresh tokens for a given user (e.g. on logout-all or password change).
func (m *RefreshTokenModel) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	return m.DB.WithContext(ctx).Model(&RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", &now).Error
}

// DeleteExpired permanently removes expired tokens (for cleanup cron jobs).
func (m *RefreshTokenModel) DeleteExpired(ctx context.Context) error {
	return m.DB.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&RefreshToken{}).Error
}
