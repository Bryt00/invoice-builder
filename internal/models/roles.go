package models

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Permission struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`
	Name        string    `gorm:"unique;type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
}

type Role struct {
	ID          uuid.UUID    `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`
	Name        string       `gorm:"unique;type:varchar(255);not null" json:"name"`
	Description string       `gorm:"type:text" json:"description"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
}

type RoleInterface interface {
	GetByName(ctx context.Context, name string) (*Role, error)
	GetAll(ctx context.Context) ([]*Role, error)
	Create(ctx context.Context, role *Role) error
}

type RoleModel struct {
	DB *gorm.DB
}

func (m *RoleModel) GetByName(ctx context.Context, name string) (*Role, error) {
	var role Role
	// Preload permissions so we can check what the role is allowed to do
	if err := m.DB.WithContext(ctx).Preload("Permissions").Where("name = ?", name).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &role, nil
}

func (m *RoleModel) GetAll(ctx context.Context) ([]*Role, error) {
	var roles []*Role
	err := m.DB.WithContext(ctx).Order("name ASC").Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (m *RoleModel) Create(ctx context.Context, role *Role) error {
	if role.ID == uuid.Nil {
		role.ID = uuid.New()
	}
	return m.DB.WithContext(ctx).Create(role).Error
}
