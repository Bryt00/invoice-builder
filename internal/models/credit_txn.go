package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreditTxnType string

const (
	CreditTxnPurchase CreditTxnType = "purchase"
	CreditTxnUsage    CreditTxnType = "usage"
	CreditTxnRefund   CreditTxnType = "refund"
)

type CreditTxn struct {
	ID          uint          `gorm:"primaryKey" json:"id"`
	UserID      uuid.UUID     `gorm:"type:uuid;not null;index" json:"user_id"`
	Amount      int           `gorm:"not null" json:"amount"`
	Type        CreditTxnType `gorm:"size:20;not null" json:"type"`
	Description string        `gorm:"size:255" json:"description"`
	CreatedAt   time.Time     `json:"created_at"`
}

type CreditTxnInterface interface {
	Insert(ctx context.Context, txn *CreditTxn) error
	GetBalanceByUserID(ctx context.Context, userID uuid.UUID) (int, error)
	GetCreditStatsByUserID(ctx context.Context, userID uuid.UUID) (balance int, totalPurchased int, totalUsed int, err error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*CreditTxn, error)
	GetSystemCreditStats(ctx context.Context) (totalPurchased, totalUsed int, err error)
	GetAllSystemTxns(ctx context.Context, txnType string, page, limit int) ([]*CreditTxn, int64, error)
	AdminGrantCredits(ctx context.Context, userID uuid.UUID, amount int, reason string) error
}

type CreditTxnModel struct {
	DB *gorm.DB
}

func (m *CreditTxnModel) Insert(ctx context.Context, txn *CreditTxn) error {
	return m.DB.WithContext(ctx).Create(txn).Error
}

func (m *CreditTxnModel) GetBalanceByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	var total int
	err := m.DB.WithContext(ctx).Model(&CreditTxn{}).Where("user_id = ?", userID).Select("COALESCE(SUM(amount), 0)").Scan(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (m *CreditTxnModel) GetCreditStatsByUserID(ctx context.Context, userID uuid.UUID) (balance int, totalPurchased int, totalUsed int, err error) {
	var txns []*CreditTxn
	err = m.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&txns).Error
	if err != nil {
		return 0, 0, 0, err
	}
	for _, t := range txns {
		balance += t.Amount
		if t.Amount > 0 {
			totalPurchased += t.Amount
		} else if t.Amount < 0 {
			totalUsed += -t.Amount
		}
	}
	return balance, totalPurchased, totalUsed, nil
}

func (m *CreditTxnModel) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*CreditTxn, error) {
	var txns []*CreditTxn
	err := m.DB.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&txns).Error
	if err != nil {
		return nil, err
	}
	return txns, nil
}

func (m *CreditTxnModel) GetSystemCreditStats(ctx context.Context) (totalPurchased, totalUsed int, err error) {
	var txns []*CreditTxn
	err = m.DB.WithContext(ctx).Find(&txns).Error
	if err != nil {
		return 0, 0, err
	}
	for _, t := range txns {
		if t.Amount > 0 {
			totalPurchased += t.Amount
		} else if t.Amount < 0 {
			totalUsed += -t.Amount
		}
	}
	return totalPurchased, totalUsed, nil
}

func (m *CreditTxnModel) GetAllSystemTxns(ctx context.Context, txnType string, page, limit int) ([]*CreditTxn, int64, error) {
	var txns []*CreditTxn
	var total int64

	query := m.DB.WithContext(ctx).Model(&CreditTxn{})
	if txnType != "" {
		query = query.Where("type = ?", txnType)
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

	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&txns).Error
	if err != nil {
		return nil, 0, err
	}
	return txns, total, nil
}

func (m *CreditTxnModel) AdminGrantCredits(ctx context.Context, userID uuid.UUID, amount int, reason string) error {
	desc := "Admin Grant: " + reason
	if amount < 0 {
		desc = "Admin Deduction: " + reason
	}
	txn := &CreditTxn{
		UserID:      userID,
		Amount:      amount,
		Type:        CreditTxnPurchase,
		Description: desc,
		CreatedAt:   time.Now(),
	}
	return m.Insert(ctx, txn)
}
