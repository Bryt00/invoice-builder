package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditLog struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`
	UserID     *uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	User       *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Action     string     `gorm:"type:varchar(100);not null;index" json:"action"`
	EntityType string     `gorm:"type:varchar(100);index" json:"entity_type"`
	EntityID   string     `gorm:"type:varchar(255)" json:"entity_id"`
	IPAddress  string     `gorm:"type:varchar(50)" json:"ip_address"`
	UserAgent  string     `gorm:"type:text" json:"user_agent"`
	Metadata   string     `gorm:"type:jsonb" json:"metadata"`
	CreatedAt  time.Time  `gorm:"index" json:"created_at"`
}

type AuditLogInterface interface {
	Record(ctx context.Context, log *AuditLog) error
	GetLogs(ctx context.Context, action string, page, limit int) ([]*AuditLog, int64, error)
}

type AuditLogModel struct {
	DB *gorm.DB
}

func (m *AuditLogModel) Record(ctx context.Context, log *AuditLog) error {
	if log.ID == uuid.Nil {
		log.ID = uuid.New()
	}
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now()
	}
	return m.DB.WithContext(ctx).Create(log).Error
}

func (m *AuditLogModel) GetLogs(ctx context.Context, action string, page, limit int) ([]*AuditLog, int64, error) {
	var logs []*AuditLog
	var total int64

	query := m.DB.WithContext(ctx).Model(&AuditLog{}).Preload("User")
	if action != "" {
		query = query.Where("action = ?", action)
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
