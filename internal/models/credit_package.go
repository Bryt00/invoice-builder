package models

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreditPackage struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`
	Name           string         `gorm:"type:varchar(255);not null" json:"name"`
	Slug           string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Description    string         `gorm:"type:text" json:"description"`
	Price          int64          `gorm:"not null" json:"price"` // Price in smallest subunit (e.g. 5000 = 50.00)
	Currency       string         `gorm:"type:varchar(10);default:'GHS'" json:"currency"`
	CreditsGranted int            `gorm:"not null" json:"credits_granted"`
	BadgeTag       string         `gorm:"type:varchar(100)" json:"badge_tag"`
	DisplayOrder   int            `gorm:"default:0" json:"display_order"`
	IsActive       bool           `gorm:"default:true;index" json:"is_active"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type CreditPackageInterface interface {
	Create(ctx context.Context, pkg *CreditPackage) error
	GetByID(ctx context.Context, id uuid.UUID) (*CreditPackage, error)
	GetBySlug(ctx context.Context, slug string) (*CreditPackage, error)
	GetAll(ctx context.Context, includeInactive bool) ([]*CreditPackage, error)
	Update(ctx context.Context, pkg *CreditPackage) error
	ToggleActive(ctx context.Context, id uuid.UUID, isActive bool) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreditPackageModel struct {
	DB *gorm.DB
}

func (m *CreditPackageModel) Create(ctx context.Context, pkg *CreditPackage) error {
	if pkg.ID == uuid.Nil {
		pkg.ID = uuid.New()
	}
	return m.DB.WithContext(ctx).Create(pkg).Error
}

func (m *CreditPackageModel) GetByID(ctx context.Context, id uuid.UUID) (*CreditPackage, error) {
	var pkg CreditPackage
	err := m.DB.WithContext(ctx).First(&pkg, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &pkg, nil
}

func (m *CreditPackageModel) GetBySlug(ctx context.Context, slug string) (*CreditPackage, error) {
	var pkg CreditPackage
	err := m.DB.WithContext(ctx).First(&pkg, "slug = ?", slug).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &pkg, nil
}

func (m *CreditPackageModel) GetAll(ctx context.Context, includeInactive bool) ([]*CreditPackage, error) {
	var packages []*CreditPackage
	query := m.DB.WithContext(ctx).Order("display_order ASC, price ASC")
	if !includeInactive {
		query = query.Where("is_active = ?", true)
	}
	err := query.Find(&packages).Error
	return packages, err
}

func (m *CreditPackageModel) Update(ctx context.Context, pkg *CreditPackage) error {
	return m.DB.WithContext(ctx).Save(pkg).Error
}

func (m *CreditPackageModel) ToggleActive(ctx context.Context, id uuid.UUID, isActive bool) error {
	return m.DB.WithContext(ctx).Model(&CreditPackage{}).Where("id = ?", id).Update("is_active", isActive).Error
}

func (m *CreditPackageModel) Delete(ctx context.Context, id uuid.UUID) error {
	return m.DB.WithContext(ctx).Delete(&CreditPackage{}, "id = ?", id).Error
}
