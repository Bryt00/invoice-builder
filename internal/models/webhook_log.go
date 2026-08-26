package models

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WebhookLog struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`
	Event       string     `gorm:"type:varchar(100);not null;index" json:"event"`
	Payload     string     `gorm:"type:text;not null" json:"payload"`
	Status      string     `gorm:"type:varchar(50);default:'processed';index" json:"status"` // 'processed', 'failed', 'ignored'
	ErrorMsg    string     `gorm:"type:text" json:"error_msg"`
	Attempts    int        `gorm:"default:1" json:"attempts"`
	CreatedAt   time.Time  `gorm:"index" json:"created_at"`
	ProcessedAt *time.Time `json:"processed_at"`
}

type WebhookLogInterface interface {
	Record(ctx context.Context, log *WebhookLog) error
	GetByID(ctx context.Context, id uuid.UUID) (*WebhookLog, error)
	GetAll(ctx context.Context, status string, page, limit int) ([]*WebhookLog, int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, errorMsg string) error
}

type WebhookLogModel struct {
	DB *gorm.DB
}

func (m *WebhookLogModel) Record(ctx context.Context, log *WebhookLog) error {
	if log.ID == uuid.Nil {
		log.ID = uuid.New()
	}
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now()
	}
	return m.DB.WithContext(ctx).Create(log).Error
}

func (m *WebhookLogModel) GetByID(ctx context.Context, id uuid.UUID) (*WebhookLog, error) {
	var log WebhookLog
	err := m.DB.WithContext(ctx).First(&log, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &log, nil
}

func (m *WebhookLogModel) GetAll(ctx context.Context, status string, page, limit int) ([]*WebhookLog, int64, error) {
	var logs []*WebhookLog
	var total int64

	query := m.DB.WithContext(ctx).Model(&WebhookLog{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 25
	}
	offset := (page - 1) * limit

	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

func (m *WebhookLogModel) UpdateStatus(ctx context.Context, id uuid.UUID, status string, errorMsg string) error {
	now := time.Now()
	return m.DB.WithContext(ctx).Model(&WebhookLog{}).Where("id = ?", id).Updates(map[string]any{
		"status":       status,
		"error_msg":    errorMsg,
		"attempts":     gorm.Expr("attempts + 1"),
		"processed_at": &now,
	}).Error
}
