package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"raven.go.invoice-builder/internal/mailer"
	"raven.go.invoice-builder/internal/models"
	"raven.go.invoice-builder/internal/pdf"
)

type InvoiceService interface {
	GetAllByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]*models.Invoice, int64, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.Invoice, error)
	GetByPublicToken(ctx context.Context, token string) (*models.Invoice, error)
	MarkAsPaid(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	CreateInvoice(ctx context.Context, invoice *models.Invoice) error
	UpdateInvoice(ctx context.Context, invoice *models.Invoice) error
	GeneratePDF(ctx context.Context, invoice *models.Invoice, paperSize string) ([]byte, error)
	DispatchInvoiceEmail(ctx context.Context, invoice *models.Invoice, profile *models.BusinessProfile, client *models.Client) error
}

type invoiceService struct {
	models         models.Models
	financeService FinanceService
	mailer         mailer.Mailer
	frontendURL    string
}

func NewInvoiceService(models models.Models, financeService FinanceService, mailer mailer.Mailer, frontendURL string) InvoiceService {
	return &invoiceService{
		models:         models,
		financeService: financeService,
		mailer:         mailer,
		frontendURL:    frontendURL,
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

func (s *invoiceService) CreateInvoice(ctx context.Context, invoice *models.Invoice) error {
	if invoice.ID == uuid.Nil {
		invoice.ID = uuid.New()
	}
	if invoice.PublicToken == "" {
		invoice.PublicToken = uuid.New().String()
	}
	
	// Pre-process items
	for i := range invoice.LineItems {
		if invoice.LineItems[i].ID == uuid.Nil {
			invoice.LineItems[i].ID = uuid.New()
		}
		invoice.LineItems[i].InvoiceID = invoice.ID
	}

	return s.models.Invoice.Insert(ctx, invoice)
}

func (s *invoiceService) UpdateInvoice(ctx context.Context, invoice *models.Invoice) error {
	return s.models.Invoice.Update(ctx, invoice)
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
		return nil, models.ErrInvoiceNil
	}
	profile, _ := s.models.BusinessProfiles.GetByUserID(ctx, invoice.UserID)
	return pdf.GenerateInvoicePDF(invoice, profile, paperSize)
}

func (s *invoiceService) DispatchInvoiceEmail(ctx context.Context, invoice *models.Invoice, profile *models.BusinessProfile, client *models.Client) error {
	if client == nil || client.Email == "" {
		return models.ErrClientEmailRequired
	}

	// Generate PDF
	pdfBytes, err := s.GeneratePDF(ctx, invoice, "a4")
	if err != nil {
		return fmt.Errorf("failed to generate invoice pdf: %w", err)
	}

	// Prepare email data
	companyName := "Teks-Invoice"
	if profile != nil && profile.CompanyName != "" {
		companyName = profile.CompanyName
	}

	data := map[string]any{
		"CompanyName":   companyName,
		"ClientName":    client.Name,
		"InvoiceNumber": invoice.InvoiceNumber,
		"TotalAmount":   fmt.Sprintf("%s%.2f", invoice.Currency, invoice.Total),
		"DueDate":       invoice.DueDate.Format("Jan 02, 2006"),
		"InvoiceURL":    fmt.Sprintf("%s/invoice/public/%s", s.frontendURL, invoice.PublicToken),
	}

	fileName := fmt.Sprintf("%s.pdf", invoice.InvoiceNumber)

	// Send email with attachment
	err = s.mailer.SendMailWithAttachment(client.Email, "invoice_dispatch.tmpl", data, fileName, pdfBytes)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
