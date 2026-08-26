package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BusinessProfile struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	UserID              uuid.UUID `gorm:"type:uuid;not null;unique" json:"user_id"`
	CompanyName         string    `gorm:"size:255" json:"company_name"`
	Role                string    `gorm:"size:100" json:"role"`
	Address             string    `gorm:"type:text" json:"address"`
	LogoURL             string    `gorm:"size:512" json:"logo_url"`
	TaxID               string    `gorm:"size:100" json:"tax_id"`
	RegistrationNumber  string    `gorm:"size:100" json:"registration_number"`
	RegistrationDate    string    `gorm:"size:50" json:"registration_date"`
	BusinessType        string    `gorm:"size:100" json:"business_type"`
	RegisteredAddress   string    `gorm:"type:text" json:"registered_address"`
	DefaultCurrency     string    `gorm:"size:10;default:'USD'" json:"default_currency"`
	PaymentInstructions string    `gorm:"type:text" json:"payment_instructions"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type BusinessProfileInterface interface {
	Upsert(ctx context.Context, profile *BusinessProfile) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*BusinessProfile, error)
}

type BusinessProfileModel struct {
	DB *gorm.DB
}

func (m *BusinessProfileModel) Upsert(ctx context.Context, profile *BusinessProfile) error {
	return m.DB.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"company_name", "role", "address", "logo_url", "tax_id", "registration_number", "registration_date", "business_type", "registered_address", "default_currency", "payment_instructions", "updated_at"}),
	}).Create(profile).Error
}

func (m *BusinessProfileModel) GetByUserID(ctx context.Context, userID uuid.UUID) (*BusinessProfile, error) {
	var profile BusinessProfile
	err := m.DB.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &profile, nil
}
