package services

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
	"raven.go.invoice-builder/internal/models"
)

type FinanceService interface {
	PreseedCategories(ctx context.Context, userID uuid.UUID) error
	GetCategories(ctx context.Context, userID uuid.UUID) ([]*models.FinancialCategory, error)
	CreateCategory(ctx context.Context, userID uuid.UUID, name, catType, color, icon string) (*models.FinancialCategory, error)
	CreateTransaction(ctx context.Context, txn *models.FinancialTransaction) error
	UpdateTransaction(ctx context.Context, txn *models.FinancialTransaction) error
	DeleteTransaction(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	GetTransactionByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.FinancialTransaction, error)
	GetTransactionByInvoiceID(ctx context.Context, invoiceID uuid.UUID) (*models.FinancialTransaction, error)
	GetAllTransactions(ctx context.Context, userID uuid.UUID, txnType string, categoryID string, startDate string, endDate string, search string) ([]*models.FinancialTransaction, error)
	GetFinancialStats(ctx context.Context, userID uuid.UUID, startDate string, endDate string) (totalIncome float64, totalExpenses float64, netProfit float64, err error)
	ExportCSVLedger(ctx context.Context, userID uuid.UUID, txnType, categoryID, startDate, endDate, search string, out io.Writer) error
}

type financeService struct {
	models models.Models
}

func NewFinanceService(models models.Models) FinanceService {
	return &financeService{models: models}
}

func (s *financeService) PreseedCategories(ctx context.Context, userID uuid.UUID) error {
	return s.models.Finance.PreseedDefaultCategories(ctx, userID)
}

func (s *financeService) GetCategories(ctx context.Context, userID uuid.UUID) ([]*models.FinancialCategory, error) {
	return s.models.Finance.GetCategoriesByUserID(ctx, userID)
}

func (s *financeService) CreateCategory(ctx context.Context, userID uuid.UUID, name, catType, color, icon string) (*models.FinancialCategory, error) {
	if color == "" {
		color = "primary"
	}
	if icon == "" {
		icon = "payments"
	}
	cat := &models.FinancialCategory{
		ID:        uuid.New(),
		UserID:    &userID,
		Name:      strings.TrimSpace(name),
		Type:      strings.TrimSpace(catType),
		Color:     color,
		Icon:      icon,
		IsDefault: false,
	}
	err := s.models.Finance.InsertCategory(ctx, cat)
	if err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *financeService) CreateTransaction(ctx context.Context, txn *models.FinancialTransaction) error {
	if txn.ID == uuid.Nil {
		txn.ID = uuid.New()
	}
	if txn.TransactionDate.IsZero() {
		txn.TransactionDate = time.Now()
	}
	if strings.TrimSpace(txn.Currency) == "" {
		profile, _ := s.models.BusinessProfiles.GetByUserID(ctx, txn.UserID)
		if profile != nil && profile.DefaultCurrency != "" {
			txn.Currency = profile.DefaultCurrency
		} else {
			txn.Currency = "USD"
		}
	}
	return s.models.Finance.InsertTransaction(ctx, txn)
}

func (s *financeService) UpdateTransaction(ctx context.Context, txn *models.FinancialTransaction) error {
	return s.models.Finance.UpdateTransaction(ctx, txn)
}

func (s *financeService) DeleteTransaction(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return s.models.Finance.DeleteTransaction(ctx, id, userID)
}

func (s *financeService) GetTransactionByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.FinancialTransaction, error) {
	return s.models.Finance.GetTransactionByID(ctx, id, userID)
}

func (s *financeService) GetTransactionByInvoiceID(ctx context.Context, invoiceID uuid.UUID) (*models.FinancialTransaction, error) {
	return s.models.Finance.GetTransactionByInvoiceID(ctx, invoiceID)
}

func (s *financeService) GetAllTransactions(ctx context.Context, userID uuid.UUID, txnType string, categoryID string, startDate string, endDate string, search string) ([]*models.FinancialTransaction, error) {
	return s.models.Finance.GetAllTransactionsByUserID(ctx, userID, txnType, categoryID, startDate, endDate, search)
}

func (s *financeService) GetFinancialStats(ctx context.Context, userID uuid.UUID, startDate string, endDate string) (totalIncome float64, totalExpenses float64, netProfit float64, err error) {
	return s.models.Finance.GetFinancialStats(ctx, userID, startDate, endDate)
}

func (s *financeService) ExportCSVLedger(ctx context.Context, userID uuid.UUID, txnType, categoryID, startDate, endDate, search string, out io.Writer) error {
	transactions, err := s.GetAllTransactions(ctx, userID, txnType, categoryID, startDate, endDate, search)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(out)
	_ = writer.Write([]string{"Date", "Type", "Title", "Category", "Payee/Payer", "Currency", "Amount", "Status", "Notes"})

	for _, txn := range transactions {
		catName := "Uncategorized"
		if txn.Category != nil {
			catName = txn.Category.Name
		}
		_ = writer.Write([]string{
			txn.TransactionDate.Format("2006-01-02"),
			strings.ToUpper(txn.Type),
			txn.Title,
			catName,
			txn.PayeeOrPayer,
			txn.Currency,
			fmt.Sprintf("%.2f", txn.Amount),
			txn.Status,
			txn.Description,
		})
	}
	writer.Flush()
	return writer.Error()
}
