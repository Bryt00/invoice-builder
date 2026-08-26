package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Client struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	Email     string    `gorm:"size:255;not null" json:"email"`
	Phone     string    `gorm:"size:50" json:"phone"`
	Company   string    `gorm:"size:255" json:"company"`
	Address   string    `gorm:"type:text" json:"address"`
	TaxID     string    `gorm:"size:100" json:"tax_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ClientInterface interface {
	Insert(ctx context.Context, client *Client) error
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Client, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]*Client, int64, error)
	Update(ctx context.Context, client *Client) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type ClientModel struct {
	DB *gorm.DB
}

func (m *ClientModel) Insert(ctx context.Context, client *Client) error {
	if client.ID == uuid.Nil {
		client.ID = uuid.New()
	}
	return m.DB.WithContext(ctx).Create(client).Error
}

func (m *ClientModel) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Client, error) {
	var client Client
	err := m.DB.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&client).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &client, nil
}

func (m *ClientModel) GetAllByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]*Client, int64, error) {
	var clients []*Client
	var total int64

	query := m.DB.WithContext(ctx).Model(&Client{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	err := query.Order("name ASC").Offset(offset).Limit(limit).Find(&clients).Error
	if err != nil {
		return nil, 0, err
	}
	return clients, total, nil
}

func (m *ClientModel) Update(ctx context.Context, client *Client) error {
	return m.DB.WithContext(ctx).Save(client).Error
}

func (m *ClientModel) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	result := m.DB.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&Client{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNoRecord
	}
	return nil
}
