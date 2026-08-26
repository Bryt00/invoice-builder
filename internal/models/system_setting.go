package models

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type SystemSetting struct {
	Key         string    `gorm:"type:varchar(100);primaryKey" json:"key"`
	Value       string    `gorm:"type:text;not null" json:"value"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SystemSettingInterface interface {
	Get(ctx context.Context, key string, fallback string) (string, error)
	Set(ctx context.Context, key string, value string, description string) error
	GetAll(ctx context.Context) ([]*SystemSetting, error)
}

type SystemSettingModel struct {
	DB *gorm.DB
}

func (m *SystemSettingModel) Get(ctx context.Context, key string, fallback string) (string, error) {
	var setting SystemSetting
	err := m.DB.WithContext(ctx).Where("key = ?", key).Limit(1).Find(&setting).Error
	if err != nil {
		return fallback, err
	}
	if setting.Key == "" {
		return fallback, nil
	}
	return setting.Value, nil
}

func (m *SystemSettingModel) Set(ctx context.Context, key string, value string, description string) error {
	setting := SystemSetting{
		Key:         key,
		Value:       value,
		Description: description,
		UpdatedAt:   time.Now(),
	}
	return m.DB.WithContext(ctx).Save(&setting).Error
}

func (m *SystemSettingModel) GetAll(ctx context.Context) ([]*SystemSetting, error) {
	var settings []*SystemSetting
	err := m.DB.WithContext(ctx).Order("key ASC").Find(&settings).Error
	return settings, err
}
