package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FinancialCategory struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    *uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	Name      string     `gorm:"size:100;not null" json:"name"`
	Type      string     `gorm:"size:20;not null" json:"type"` // "income" or "expense"
	Color     string     `gorm:"size:50;default:'primary'" json:"color"`
	Icon      string     `gorm:"size:50;default:'payments'" json:"icon"`
	IsDefault bool       `gorm:"default:false" json:"is_default"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type FinancialTransaction struct {
	ID              uuid.UUID          `gorm:"type:uuid;primaryKey" json:"id"`
	UserID          uuid.UUID          `gorm:"type:uuid;not null;index" json:"user_id"`
	CategoryID      uuid.UUID          `gorm:"type:uuid;not null;index" json:"category_id"`
	Category        *FinancialCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	InvoiceID       *uuid.UUID         `gorm:"type:uuid;index" json:"invoice_id,omitempty"`
	Invoice         *Invoice           `gorm:"foreignKey:InvoiceID" json:"invoice,omitempty"`
	Type            string             `gorm:"size:20;not null" json:"type"` // "income" or "expense"
	Amount          float64            `gorm:"type:decimal(12,2);not null" json:"amount"`
	Currency        string             `gorm:"size:10;default:'USD'" json:"currency"`
	Title           string             `gorm:"size:255;not null" json:"title"`
	Description     string             `gorm:"type:text" json:"description"`
	TransactionDate time.Time          `gorm:"not null" json:"transaction_date"`
	PayeeOrPayer    string             `gorm:"size:255" json:"payee_or_payer"`
	Status          string             `gorm:"size:20;default:'completed'" json:"status"` // "completed", "pending"
	ReceiptURL      string             `gorm:"size:255" json:"receipt_url"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
}

type FinanceInterface interface {
	PreseedDefaultCategories(ctx context.Context, userID uuid.UUID) error
	GetCategoriesByUserID(ctx context.Context, userID uuid.UUID) ([]*FinancialCategory, error)
	InsertCategory(ctx context.Context, cat *FinancialCategory) error
	InsertTransaction(ctx context.Context, txn *FinancialTransaction) error
	UpdateTransaction(ctx context.Context, txn *FinancialTransaction) error
	DeleteTransaction(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	GetTransactionByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*FinancialTransaction, error)
	GetTransactionByInvoiceID(ctx context.Context, invoiceID uuid.UUID) (*FinancialTransaction, error)
	GetAllTransactionsByUserID(ctx context.Context, userID uuid.UUID, txnType string, categoryID string, startDate string, endDate string, search string) ([]*FinancialTransaction, error)
	GetFinancialStats(ctx context.Context, userID uuid.UUID, startDate string, endDate string) (totalIncome float64, totalExpenses float64, netProfit float64, err error)
}

type FinanceModel struct {
	DB *gorm.DB
}

var defaultCategories = []struct {
	Name  string
	Type  string
	Color string
	Icon  string
}{
	{"Invoice Payments", "income", "emerald", "receipt_long"},
	{"Consulting Services", "income", "teal", "work"},
	{"Other Income", "income", "cyan", "attach_money"},
	{"Software & SaaS", "expense", "indigo", "laptop_mac"},
	{"Salaries & Contractor Fees", "expense", "purple", "group"},
	{"Marketing & Advertising", "expense", "pink", "campaign"},
	{"Office & Supplies", "expense", "amber", "inventory_2"},
	{"Utilities & Internet", "expense", "blue", "wifi"},
	{"Travel & Entertainment", "expense", "orange", "flight_takeoff"},
	{"Taxes & Licensing", "expense", "rose", "policy"},
	{"Miscellaneous Expense", "expense", "slate", "more_horiz"},
}

func (m *FinanceModel) PreseedDefaultCategories(ctx context.Context, userID uuid.UUID) error {
	var count int64
	err := m.DB.WithContext(ctx).Model(&FinancialCategory{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	for _, item := range defaultCategories {
		uid := userID
		cat := FinancialCategory{
			ID:        uuid.New(),
			UserID:    &uid,
			Name:      item.Name,
			Type:      item.Type,
			Color:     item.Color,
			Icon:      item.Icon,
			IsDefault: true,
		}
		_ = m.DB.WithContext(ctx).Create(&cat).Error
	}
	return nil
}

func (m *FinanceModel) GetCategoriesByUserID(ctx context.Context, userID uuid.UUID) ([]*FinancialCategory, error) {
	_ = m.PreseedDefaultCategories(ctx, userID)

	var categories []*FinancialCategory
	err := m.DB.WithContext(ctx).Where("user_id = ? OR user_id IS NULL", userID).Order("type ASC, name ASC").Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (m *FinanceModel) InsertCategory(ctx context.Context, cat *FinancialCategory) error {
	if cat.ID == uuid.Nil {
		cat.ID = uuid.New()
	}
	return m.DB.WithContext(ctx).Create(cat).Error
}

func (m *FinanceModel) InsertTransaction(ctx context.Context, txn *FinancialTransaction) error {
	if txn.ID == uuid.Nil {
		txn.ID = uuid.New()
	}
	if txn.TransactionDate.IsZero() {
		txn.TransactionDate = time.Now()
	}
	return m.DB.WithContext(ctx).Create(txn).Error
}

func (m *FinanceModel) UpdateTransaction(ctx context.Context, txn *FinancialTransaction) error {
	return m.DB.WithContext(ctx).Model(&FinancialTransaction{}).Where("id = ? AND user_id = ?", txn.ID, txn.UserID).Updates(map[string]any{
		"category_id":      txn.CategoryID,
		"type":             txn.Type,
		"amount":           txn.Amount,
		"currency":         txn.Currency,
		"title":            txn.Title,
		"description":      txn.Description,
		"transaction_date": txn.TransactionDate,
		"payee_or_payer":   txn.PayeeOrPayer,
		"status":           txn.Status,
		"receipt_url":      txn.ReceiptURL,
	}).Error
}

func (m *FinanceModel) DeleteTransaction(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return m.DB.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&FinancialTransaction{}).Error
}

func (m *FinanceModel) GetTransactionByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*FinancialTransaction, error) {
	var txn FinancialTransaction
	err := m.DB.WithContext(ctx).Preload("Category").Preload("Invoice").Where("id = ? AND user_id = ?", id, userID).First(&txn).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &txn, nil
}

func (m *FinanceModel) GetTransactionByInvoiceID(ctx context.Context, invoiceID uuid.UUID) (*FinancialTransaction, error) {
	var txn FinancialTransaction
	err := m.DB.WithContext(ctx).Where("invoice_id = ?", invoiceID).First(&txn).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &txn, nil
}

func (m *FinanceModel) GetAllTransactionsByUserID(ctx context.Context, userID uuid.UUID, txnType string, categoryID string, startDate string, endDate string, search string) ([]*FinancialTransaction, error) {
	var txns []*FinancialTransaction

	query := m.DB.WithContext(ctx).Preload("Category").Preload("Invoice").Where("user_id = ?", userID)

	if txnType != "" && txnType != "all" {
		query = query.Where("type = ?", txnType)
	}

	if categoryID != "" && categoryID != "all" {
		if cid, err := uuid.Parse(categoryID); err == nil {
			query = query.Where("category_id = ?", cid)
		}
	}

	if startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("transaction_date >= ?", t)
		}
	}

	if endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			query = query.Where("transaction_date <= ?", t.Add(23*time.Hour+59*time.Minute+59*time.Second))
		}
	}

	if search != "" {
		s := "%" + search + "%"
		query = query.Where("title ILIKE ? OR payee_or_payer ILIKE ? OR description ILIKE ?", s, s, s)
	}

	err := query.Order("transaction_date DESC, created_at DESC").Find(&txns).Error
	if err != nil {
		return nil, err
	}

	return txns, nil
}

func (m *FinanceModel) GetFinancialStats(ctx context.Context, userID uuid.UUID, startDate string, endDate string) (totalIncome float64, totalExpenses float64, netProfit float64, err error) {
	queryIncome := m.DB.WithContext(ctx).Model(&FinancialTransaction{}).Where("user_id = ? AND type = 'income'", userID)
	queryExpense := m.DB.WithContext(ctx).Model(&FinancialTransaction{}).Where("user_id = ? AND type = 'expense'", userID)

	if startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			queryIncome = queryIncome.Where("transaction_date >= ?", t)
			queryExpense = queryExpense.Where("transaction_date >= ?", t)
		}
	}

	if endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			endT := t.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			queryIncome = queryIncome.Where("transaction_date <= ?", endT)
			queryExpense = queryExpense.Where("transaction_date <= ?", endT)
		}
	}

	_ = queryIncome.Select("COALESCE(SUM(amount), 0)").Scan(&totalIncome).Error
	_ = queryExpense.Select("COALESCE(SUM(amount), 0)").Scan(&totalExpenses).Error

	netProfit = totalIncome - totalExpenses
	return totalIncome, totalExpenses, netProfit, nil
}
