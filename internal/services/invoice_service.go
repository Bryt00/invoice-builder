package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"raven.go.invoice-builder/internal/models"
	"raven.go.invoice-builder/internal/pdf"
)

type InvoiceService interface {
	GetAllByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]*models.Invoice, int64, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.Invoice, error)
	GetByPublicToken(ctx context.Context, token string) (*models.Invoice, error)
	MarkAsPaid(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	GeneratePDF(ctx context.Context, invoice *models.Invoice, paperSize string) ([]byte, error)
}

type invoiceService struct {
	models         models.Models
	financeService FinanceService
}

func NewInvoiceService(models models.Models, financeService FinanceService) InvoiceService {
	return &invoiceService{
		models:         models,
		financeService: financeService,
	}
}

func (s *invoiceService) GetAllByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]*models.Invoice, int64, error) {
	return s.models.Invoice.GetAllByUserID(ctx, userID, page, limit)
}

func (s *invoiceService) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.Invoice, error) {
	return s.models.Invoice.GetByID(ctx, id, userID)
}

func (s *invoiceService) GetByPublicToken(ctx context.Context, token string) (*models.Invoice, error) {
	return s.models.Invoice.GetByPublicToken(ctx, token)
}

func (s *invoiceService) MarkAsPaid(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return s.models.Transaction(ctx, func(txModels models.Models) error {
		err := txModels.Invoice.MarkAsPaid(ctx, id, userID)
		if err != nil {
			return err
		}

		invoice, getErr := txModels.Invoice.GetByID(ctx, id, userID)
		if getErr == nil && invoice != nil {
			cats, _ := txModels.Finance.GetCategoriesByUserID(ctx, userID)
			var categoryID uuid.UUID
			for _, cat := range cats {
				if cat.Type == "income" {
					categoryID = cat.ID
					break
				}
			}

			payerName := "Direct Client"
			if invoice.Client != nil && invoice.Client.Name != "" {
				payerName = invoice.Client.Name
			}

			invID := invoice.ID
			txn := &models.FinancialTransaction{
				ID:              uuid.New(),
				UserID:          userID,
				CategoryID:      categoryID,
				InvoiceID:       &invID,
				Type:            "income",
				Amount:          invoice.Total,
				Currency:        invoice.Currency,
				Title:           fmt.Sprintf("Payment for Invoice #%s", invoice.InvoiceNumber),
				Description:     fmt.Sprintf("Auto-synced from Paid Invoice #%s", invoice.InvoiceNumber),
				TransactionDate: time.Now(),
				PayeeOrPayer:    payerName,
				Status:          "completed",
			}
			_ = txModels.Finance.InsertTransaction(ctx, txn)
		}
		return nil
	})
}

func (s *invoiceService) GeneratePDF(ctx context.Context, invoice *models.Invoice, paperSize string) ([]byte, error) {
	if invoice == nil {
		return nil, errors.New("invoice is nil")
	}
	profile, _ := s.models.BusinessProfiles.GetByUserID(ctx, invoice.UserID)
	return pdf.GenerateInvoicePDF(invoice, profile, paperSize)
}
