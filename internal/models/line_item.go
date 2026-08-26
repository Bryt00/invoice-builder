package models

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LineItem struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`
	InvoiceID   uuid.UUID `gorm:"type:uuid;not null;index" json:"invoice_id"`
	Description string    `gorm:"size:255;not null" json:"description"`
	Quantity    float64   `gorm:"type:decimal(10,2);not null;default:1.00" json:"quantity"`
	UnitPrice   float64   `gorm:"type:decimal(12,2);not null;default:0.00" json:"unit_price"`
	Amount      float64   `gorm:"type:decimal(12,2);not null;default:0.00" json:"amount"`
	OrderIndex  int       `gorm:"default:0" json:"order_index"`
}

type LineItemInterface interface {
	InsertBulk(ctx context.Context, items []LineItem) error
	GetByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]*LineItem, error)
	DeleteByInvoiceID(ctx context.Context, invoiceID uuid.UUID) error
}

type LineItemModel struct {
	DB *gorm.DB
}

func (m *LineItemModel) InsertBulk(ctx context.Context, items []LineItem) error {
	if len(items) == 0 {
		return nil
	}
	for i := range items {
		if items[i].ID == uuid.Nil {
			items[i].ID = uuid.New()
		}
	}
	return m.DB.WithContext(ctx).Create(&items).Error
}

func (m *LineItemModel) GetByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]*LineItem, error) {
	var items []*LineItem
	err := m.DB.WithContext(ctx).Where("invoice_id = ?", invoiceID).Order("order_index ASC").Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (m *LineItemModel) DeleteByInvoiceID(ctx context.Context, invoiceID uuid.UUID) error {
	return m.DB.WithContext(ctx).Where("invoice_id = ?", invoiceID).Delete(&LineItem{}).Error
}
