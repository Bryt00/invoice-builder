package models

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Receipt struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`
	UserID        uuid.UUID  `gorm:"type:uuid;not null;index;uniqueIndex:idx_user_receipt_number" json:"user_id"`
	InvoiceID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"invoice_id"`
	Invoice       *Invoice   `gorm:"foreignKey:InvoiceID" json:"invoice,omitempty"`
	PaymentID     *uuid.UUID `gorm:"type:uuid;unique" json:"payment_id,omitempty"`
	ReceiptNumber string     `gorm:"size:100;not null;uniqueIndex:idx_user_receipt_number" json:"receipt_number"`
	Amount        float64    `gorm:"type:decimal(12,2);not null" json:"amount"`
	Currency      string     `gorm:"size:10;default:'USD'" json:"currency"`
	IssuedAt      time.Time  `json:"issued_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

type ReceiptInterface interface {
	Insert(ctx context.Context, receipt *Receipt) error
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Receipt, error)
	GetByInvoiceID(ctx context.Context, invoiceID uuid.UUID) (*Receipt, error)
	GetByPaymentID(ctx context.Context, paymentID uuid.UUID) (*Receipt, error)
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
	GenerateReceiptNumber(ctx context.Context, userID uuid.UUID) (string, error)
}

type ReceiptModel struct {
	DB *gorm.DB
}

func (m *ReceiptModel) GenerateReceiptNumber(ctx context.Context, userID uuid.UUID) (string, error) {
	var count int64
	err := m.DB.WithContext(ctx).Model(&Receipt{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return "", err
	}
	year := time.Now().Format("2006")
	return fmt.Sprintf("RCT-%s-%04d", year, count+1), nil
}

func (m *ReceiptModel) Insert(ctx context.Context, receipt *Receipt) error {
	if receipt.ID == uuid.Nil {
		receipt.ID = uuid.New()
	}
	return m.DB.WithContext(ctx).Create(receipt).Error
}

func (m *ReceiptModel) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Receipt, error) {
	var receipt Receipt
	err := m.DB.WithContext(ctx).Preload("Invoice").Preload("Invoice.Client").Where("id = ? AND user_id = ?", id, userID).First(&receipt).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &receipt, nil
}

func (m *ReceiptModel) GetByInvoiceID(ctx context.Context, invoiceID uuid.UUID) (*Receipt, error) {
	var receipt Receipt
	err := m.DB.WithContext(ctx).Preload("Invoice").Preload("Invoice.Client").Where("invoice_id = ?", invoiceID).First(&receipt).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &receipt, nil
}

func (m *ReceiptModel) GetByPaymentID(ctx context.Context, paymentID uuid.UUID) (*Receipt, error) {
	var receipt Receipt
	err := m.DB.WithContext(ctx).Where("payment_id = ?", paymentID).First(&receipt).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &receipt, nil
}

func (m *ReceiptModel) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := m.DB.WithContext(ctx).Model(&Receipt{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}
