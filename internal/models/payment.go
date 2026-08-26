package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusPaid    PaymentStatus = "completed"
	PaymentStatusFailed  PaymentStatus = "failed"
)

type Payment struct {
	ID             uuid.UUID     `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`
	UserID         uuid.UUID     `gorm:"type:uuid;not null;index" json:"user_id"`
	InvoiceID      *uuid.UUID    `gorm:"type:uuid;index" json:"invoice_id"`
	Amount         float64       `gorm:"type:decimal(12,2);not null" json:"amount"`
	Currency       string        `gorm:"size:10;default:'USD'" json:"currency"`
	PaymentMethod  string        `gorm:"size:50" json:"payment_method"`
	TransactionRef string        `gorm:"size:255;unique" json:"transaction_ref"`
	Status         PaymentStatus `gorm:"size:20;not null;default:'pending'" json:"status"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

type PaymentInterface interface {
	Insert(ctx context.Context, payment *Payment) error
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Payment, error)
	GetByReference(ctx context.Context, reference string) (*Payment, error)
	GetByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]*Payment, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*Payment, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status PaymentStatus) error
	GetSystemRevenueStats(ctx context.Context) (totalRevenue float64, totalTxns int64, err error)
	GetAllSystemPayments(ctx context.Context, status string, page, limit int) ([]*Payment, int64, error)
}

type PaymentModel struct {
	DB *gorm.DB
}

func (m *PaymentModel) Insert(ctx context.Context, payment *Payment) error {
	if payment.ID == uuid.Nil {
		payment.ID = uuid.New()
	}
	return m.DB.WithContext(ctx).Create(payment).Error
}

func (m *PaymentModel) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Payment, error) {
	var payment Payment
	err := m.DB.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&payment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &payment, nil
}

func (m *PaymentModel) GetByReference(ctx context.Context, reference string) (*Payment, error) {
	var payment Payment
	err := m.DB.WithContext(ctx).Where("transaction_ref = ?", reference).First(&payment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &payment, nil
}

func (m *PaymentModel) GetByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]*Payment, error) {
	var payments []*Payment
	err := m.DB.WithContext(ctx).Where("invoice_id = ?", invoiceID).Order("created_at DESC").Find(&payments).Error
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (m *PaymentModel) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*Payment, error) {
	var payments []*Payment
	err := m.DB.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&payments).Error
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (m *PaymentModel) UpdateStatus(ctx context.Context, id uuid.UUID, status PaymentStatus) error {
	return m.DB.WithContext(ctx).Model(&Payment{}).Where("id = ?", id).Update("status", status).Error
}

func (m *PaymentModel) GetSystemRevenueStats(ctx context.Context) (totalRevenue float64, totalTxns int64, err error) {
	err = m.DB.WithContext(ctx).Model(&Payment{}).Where("status = ?", PaymentStatusPaid).Select("COALESCE(SUM(amount), 0)").Scan(&totalRevenue).Error
	if err != nil {
		return 0, 0, err
	}
	err = m.DB.WithContext(ctx).Model(&Payment{}).Count(&totalTxns).Error
	if err != nil {
		return 0, 0, err
	}
	return totalRevenue, totalTxns, nil
}

func (m *PaymentModel) GetAllSystemPayments(ctx context.Context, status string, page, limit int) ([]*Payment, int64, error) {
	var payments []*Payment
	var total int64

	query := m.DB.WithContext(ctx).Model(&Payment{})
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
		limit = 20
	}
	offset := (page - 1) * limit

	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&payments).Error
	if err != nil {
		return nil, 0, err
	}
	return payments, total, nil
}
