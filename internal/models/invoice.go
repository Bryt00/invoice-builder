package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InvoiceStatus string

const (
	InvoiceStatusDraft   InvoiceStatus = "draft"
	InvoiceStatusSent    InvoiceStatus = "sent"
	InvoiceStatusPaid    InvoiceStatus = "paid"
	InvoiceStatusOverdue InvoiceStatus = "overdue"
)

type Invoice struct {
	ID            uuid.UUID     `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`
	UserID        uuid.UUID     `gorm:"type:uuid;not null;index;uniqueIndex:idx_user_invoice_number" json:"user_id"`
	User          *User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ClientID      *uuid.UUID    `gorm:"type:uuid;index" json:"client_id"`
	Client        *Client       `gorm:"foreignKey:ClientID;constraint:OnDelete:SET NULL;" json:"client,omitempty"`
	InvoiceNumber string        `gorm:"size:100;not null;uniqueIndex:idx_user_invoice_number" json:"invoice_number"`
	PublicToken   string        `gorm:"size:255;not null;unique" json:"public_token"`
	Status        InvoiceStatus `gorm:"size:20;not null;default:'draft'" json:"status"`
	IsPaid        bool          `gorm:"default:false;not null" json:"is_paid"`
	IssueDate     time.Time     `gorm:"not null" json:"issue_date"`
	DueDate       time.Time     `gorm:"not null" json:"due_date"`
	Subtotal      float64       `gorm:"type:decimal(12,2);not null;default:0.00" json:"subtotal"`
	Tax           float64       `gorm:"type:decimal(12,2);not null;default:0.00" json:"tax"`
	Discount      float64       `gorm:"type:decimal(12,2);not null;default:0.00" json:"discount"`
	Total         float64       `gorm:"type:decimal(12,2);not null;default:0.00" json:"total"`
	Currency      string        `gorm:"size:10;default:'USD'" json:"currency"`
	Notes         string        `gorm:"type:text" json:"notes"`
	LineItems     []LineItem    `gorm:"foreignKey:InvoiceID" json:"line_items,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type InvoiceInterface interface {
	Insert(ctx context.Context, invoice *Invoice) error
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Invoice, error)
	AdminGetByID(ctx context.Context, id uuid.UUID) (*Invoice, error)
	GetByPublicToken(ctx context.Context, token string) (*Invoice, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]*Invoice, int64, error)
	GetAllSystemInvoices(ctx context.Context, search, status string, page, limit int) ([]*Invoice, int64, error)
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
	GenerateInvoiceNumber(ctx context.Context, userID uuid.UUID) (string, error)
	Update(ctx context.Context, invoice *Invoice) error
	MarkAsPaid(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	AdminUpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	AdminDelete(ctx context.Context, id uuid.UUID) error
}

type InvoiceModel struct {
	DB *gorm.DB
}

func (m *InvoiceModel) GenerateInvoiceNumber(ctx context.Context, userID uuid.UUID) (string, error) {
	var count int64
	err := m.DB.WithContext(ctx).Model(&Invoice{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return "", err
	}
	year := time.Now().Format("2006")
	return fmt.Sprintf("INV-%s-%04d", year, count+1), nil
}

func (m *InvoiceModel) Insert(ctx context.Context, invoice *Invoice) error {
	if invoice.ID == uuid.Nil {
		invoice.ID = uuid.New()
	}
	return m.DB.WithContext(ctx).Create(invoice).Error
}

func (m *InvoiceModel) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Invoice, error) {
	var invoice Invoice
	err := m.DB.WithContext(ctx).
		Preload("LineItems").
		Preload("Client").
		Where("id = ? AND user_id = ?", id, userID).
		First(&invoice).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &invoice, nil
}

func (m *InvoiceModel) AdminGetByID(ctx context.Context, id uuid.UUID) (*Invoice, error) {
	var invoice Invoice
	err := m.DB.WithContext(ctx).
		Preload("LineItems").
		Preload("Client").
		Preload("User").
		Where("id = ?", id).
		First(&invoice).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &invoice, nil
}

func (m *InvoiceModel) GetByPublicToken(ctx context.Context, token string) (*Invoice, error) {
	var invoice Invoice
	err := m.DB.WithContext(ctx).
		Preload("LineItems").
		Preload("Client").
		Where("public_token = ?", token).
		First(&invoice).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &invoice, nil
}

func (m *InvoiceModel) GetAllByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]*Invoice, int64, error) {
	var invoices []*Invoice
	var total int64

	query := m.DB.WithContext(ctx).Model(&Invoice{}).Where("user_id = ?", userID)

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

	err := query.Preload("Client").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&invoices).Error

	return invoices, total, err
}

func (m *InvoiceModel) Update(ctx context.Context, invoice *Invoice) error {
	return m.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("invoice_id = ?", invoice.ID).Delete(&LineItem{}).Error; err != nil {
			return err
		}
		if err := tx.Save(invoice).Error; err != nil {
			return err
		}
		return nil
	})
}

func (m *InvoiceModel) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := m.DB.WithContext(ctx).Model(&Invoice{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (m *InvoiceModel) MarkAsPaid(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	var inv Invoice
	err := m.DB.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&inv).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNoRecord
		}
		return err
	}

	if inv.Status == InvoiceStatusDraft {
		return ErrDraftCannotBePaid
	}

	result := m.DB.WithContext(ctx).
		Model(&Invoice{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]interface{}{
			"is_paid": true,
			"status":  InvoiceStatusPaid,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNoRecord
	}
	return nil
}

func (m *InvoiceModel) AdminUpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	isPaid := (status == "paid")
	result := m.DB.WithContext(ctx).
		Model(&Invoice{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_paid": isPaid,
			"status":  status,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNoRecord
	}
	return nil
}

func (m *InvoiceModel) AdminDelete(ctx context.Context, id uuid.UUID) error {
	result := m.DB.WithContext(ctx).Where("id = ?", id).Delete(&Invoice{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNoRecord
	}
	return nil
}

func (m *InvoiceModel) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	result := m.DB.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&Invoice{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNoRecord
	}
	return nil
}

func (m *InvoiceModel) GetAllSystemInvoices(ctx context.Context, search, status string, page, limit int) ([]*Invoice, int64, error) {
	var invoices []*Invoice
	var total int64

	query := m.DB.WithContext(ctx).Model(&Invoice{}).Preload("User").Preload("Client")
	if search != "" {
		s := "%" + search + "%"
		query = query.Where("invoice_number ILIKE ?", s)
	}
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

	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&invoices).Error
	if err != nil {
		return nil, 0, err
	}
	return invoices, total, nil
}

