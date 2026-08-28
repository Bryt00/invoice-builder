package main

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"raven.go.invoice-builder/internal/models"
	"raven.go.invoice-builder/internal/validator"
)

func (h *ApiHandler) GetInvoices(w http.ResponseWriter, r *http.Request) {
	user := h.ContextGetUser(r)
	if user == nil {
		h.AuthenticationRequiredResponse(w, r)
		return
	}

	page, limit := h.ParsePagination(r)

	invoices, total, err := h.Services.Invoice.GetAllByUserID(r.Context(), user.ID, page, limit)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status":   "success",
		"invoices": invoices,
		"meta": Envelope{
			"current_page":  page,
			"page_size":     limit,
			"total_records": total,
			"total_pages":   totalPages,
		},
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) GetInvoice(w http.ResponseWriter, r *http.Request) {
	user := h.ContextGetUser(r)
	if user == nil {
		h.AuthenticationRequiredResponse(w, r)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.BadRequestResponse(w, r, ErrMissingID)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		h.BadRequestResponse(w, r, ErrInvalidID)
		return
	}

	invoice, err := h.Services.Invoice.GetByID(r.Context(), id, user.ID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			h.NotFoundResponse(w, r)
		} else {
			h.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status":  "success",
		"invoice": invoice,
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) DownloadInvoice(w http.ResponseWriter, r *http.Request) {
	user := h.ContextGetUser(r)
	if user == nil {
		h.AuthenticationRequiredResponse(w, r)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.BadRequestResponse(w, r, ErrMissingID)
		return
	}

	id, parseErr := uuid.Parse(idStr)
	if parseErr != nil {
		h.BadRequestResponse(w, r, ErrInvalidID)
		return
	}

	invoice, err := h.Services.Invoice.GetByID(r.Context(), id, user.ID)
	if err != nil || invoice == nil {
		h.NotFoundResponse(w, r)
		return
	}

	size := r.URL.Query().Get("size")
	if size == "" {
		size = "a4"
	}

	// Deduct 1 credit for PDF download
	err = h.Services.Credit.DeductCredits(r.Context(), user.ID, 1, fmt.Sprintf("PDF Export for Invoice %s", invoice.InvoiceNumber))
	if err != nil {
		h.BadRequestResponse(w, r, fmt.Errorf("insufficient credits for PDF export: %w", err))
		return
	}

	pdfBytes, err := h.Services.Invoice.GeneratePDF(r.Context(), invoice, size)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.pdf\"", invoice.InvoiceNumber))
	w.Header().Set("Content-Length", strconv.Itoa(len(pdfBytes)))
	_, _ = w.Write(pdfBytes)
}

func (h *ApiHandler) DownloadPublicInvoice(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		h.BadRequestResponse(w, r, ErrMissingToken)
		return
	}

	invoice, err := h.Services.Invoice.GetByPublicToken(r.Context(), tokenStr)
	if err != nil || invoice == nil {
		h.NotFoundResponse(w, r)
		return
	}

	size := r.URL.Query().Get("size")
	if size == "" {
		size = "a4"
	}

	pdfBytes, err := h.Services.Invoice.GeneratePDF(r.Context(), invoice, size)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.pdf\"", invoice.InvoiceNumber))
	w.Header().Set("Content-Length", strconv.Itoa(len(pdfBytes)))
	_, _ = w.Write(pdfBytes)
}

func (h *ApiHandler) MarkInvoicePaid(w http.ResponseWriter, r *http.Request) {
	user := h.ContextGetUser(r)
	if user == nil {
		h.AuthenticationRequiredResponse(w, r)
		return
	}

	var input struct {
		ID string `json:"id"`
	}

	err := h.ReadJSON(w, r, &input)
	if err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	id, err := uuid.Parse(input.ID)
	if err != nil {
		h.BadRequestResponse(w, r, ErrInvalidID)
		return
	}

	err = h.Services.Invoice.MarkAsPaid(r.Context(), id, user.ID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			h.NotFoundResponse(w, r)
			return
		}
		if errors.Is(err, models.ErrDraftCannotBePaid) {
			h.BadRequestResponse(w, r, err)
			return
		}
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status":  "success",
		"message": "Invoice marked as paid successfully",
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) GetPublicInvoice(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		h.BadRequestResponse(w, r, ErrMissingToken)
		return
	}

	invoice, err := h.Services.Invoice.GetByPublicToken(r.Context(), token)
	if err != nil || invoice == nil {
		h.NotFoundResponse(w, r)
		return
	}

	profile, _ := h.Services.Auth.GetBusinessProfile(r.Context(), invoice.UserID)

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status":  "success",
		"invoice": invoice,
		"profile": profile,
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	user := h.ContextGetUser(r)
	if user == nil {
		h.AuthenticationRequiredResponse(w, r)
		return
	}

	var input struct {
		ID             string            `json:"id"`
		ClientID       string            `json:"client_id"`
		ClientEmail    string            `json:"client_email"`
		ClientAddress  string            `json:"client_address"`
		InvoiceNumber  string            `json:"invoice_number"`
		IssueDate      string            `json:"issue_date"`
		DueDate        string            `json:"due_date"`
		Currency       string            `json:"currency"`
		Notes          string            `json:"notes"`
		TaxRate        float64           `json:"tax_rate"`
		DiscountAmount float64           `json:"discount_amount"`
		Items          []models.LineItem `json:"items"`
		SaveAsDraft    bool              `json:"save_as_draft"`
		Action         string            `json:"action"`
	}

	err := h.ReadJSON(w, r, &input)
	if err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	var clientIDPtr *uuid.UUID
	if input.ClientID != "" {
		cID, err := uuid.Parse(input.ClientID)
		if err != nil {
			h.BadRequestResponse(w, r, ErrInvalidClientID)
			return
		}
		clientIDPtr = &cID
	} else if input.ClientEmail != "" {
		existingClient, err := h.Models.Clients.GetByEmail(r.Context(), input.ClientEmail, user.ID)
		if err == nil && existingClient != nil {
			clientIDPtr = &existingClient.ID
		}
	}

	v := validator.New()
	v.CheckField(len(input.Items) > 0, "items", "at least one item is required")
	v.CheckField(input.InvoiceNumber != "", "invoice_number", "invoice number is required")

	if !v.Valid() {
		h.FailedValidationResponse(w, r, v.FieldErrors)
		return
	}

	var issueDate, dueDate time.Time
	if input.IssueDate != "" {
		issueDate, _ = time.Parse("2006-01-02", input.IssueDate)
	} else {
		issueDate = time.Now()
	}
	if input.DueDate != "" {
		dueDate, _ = time.Parse("2006-01-02", input.DueDate)
	}

	var subtotal float64
	for i := range input.Items {
		input.Items[i].Amount = input.Items[i].Quantity * input.Items[i].UnitPrice
		subtotal += input.Items[i].Amount
	}
	taxAmount := (subtotal * input.TaxRate) / 100
	total := subtotal + taxAmount

	invoice := &models.Invoice{
		UserID:        user.ID,
		ClientID:      clientIDPtr,
		InvoiceNumber: input.InvoiceNumber,
		Status:        "draft",
		IssueDate:     issueDate,
		DueDate:       dueDate,
		Subtotal:      subtotal,
		Tax:           taxAmount,
		Discount:      input.DiscountAmount,
		Total:         total - input.DiscountAmount,
		Currency:      input.Currency,
		Notes:         strings.TrimSpace(input.Notes),
		LineItems:     input.Items,
	}

	if input.Action == "dispatch" {
		balance, err := h.Models.CreditTxn.GetBalanceByUserID(r.Context(), user.ID)
		if err != nil || balance < 1 {
			h.BadRequestResponse(w, r, errors.New("insufficient credits for email dispatch"))
			return
		}
	}

	err = h.Services.Invoice.CreateInvoice(r.Context(), invoice)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	// Dispatch if action is requested
	if input.Action == "dispatch" {
		// Fetch profile
		profile, err := h.Services.Auth.GetBusinessProfile(r.Context(), user.ID)
		if err != nil && !errors.Is(err, models.ErrNoRecord) {
			h.ServerErrorResponse(w, r, err)
			return
		}

		// Fetch client
		var client *models.Client
		if invoice.ClientID != nil {
			client, err = h.Services.Client.GetByID(r.Context(), *invoice.ClientID, user.ID)
			if err != nil && !errors.Is(err, models.ErrNoRecord) {
				h.ServerErrorResponse(w, r, err)
				return
			}
		} else if input.ClientEmail != "" {
			client = &models.Client{
				Name:    input.ClientEmail,
				Email:   input.ClientEmail,
				Address: input.ClientAddress,
			}
		}

		if client != nil && client.Email != "" {
			err = h.Services.Credit.DeductCredits(r.Context(), user.ID, 1, fmt.Sprintf("Email Dispatch for Invoice %s", invoice.InvoiceNumber))
			if err == nil {
				err = h.Services.Invoice.DispatchInvoiceEmail(r.Context(), invoice, profile, client)
				if err != nil {
					h.ServerErrorResponse(w, r, err)
					return
				}
				// Update status to sent
				invoice.Status = models.InvoiceStatusSent
				_ = h.Services.Invoice.UpdateInvoice(r.Context(), invoice)
			}
		}
	}

	err = h.WriteJSON(w, http.StatusCreated, Envelope{
		"status":  "success",
		"invoice": invoice,
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) UpdateInvoice(w http.ResponseWriter, r *http.Request) {
	user := h.ContextGetUser(r)
	if user == nil {
		h.AuthenticationRequiredResponse(w, r)
		return
	}

	var input struct {
		ID            string            `json:"id"`
		ClientID      string            `json:"client_id"`
		ClientEmail   string            `json:"client_email"`
		ClientAddress string            `json:"client_address"`
		Currency      string            `json:"currency"`
		InvoiceNumber string            `json:"invoice_number"`
		IssueDate     string            `json:"issue_date"`
		DueDate       string            `json:"due_date"`
		TaxRate       float64           `json:"tax_rate"`
		DiscountAmount float64          `json:"discount_amount"`
		Notes         string            `json:"notes"`
		SaveAsDraft   bool              `json:"save_as_draft"`
		Action        string            `json:"action"`
		Items         []struct {
			Description string  `json:"description"`
			Quantity    float64 `json:"quantity"`
			UnitPrice   float64 `json:"unit_price"`
		} `json:"items"`
	}

	err := h.ReadJSON(w, r, &input)
	if err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	id, err := uuid.Parse(input.ID)
	if err != nil {
		h.BadRequestResponse(w, r, ErrInvalidID)
		return
	}

	invoice, err := h.Services.Invoice.GetByID(r.Context(), id, user.ID)
	if err != nil {
		h.NotFoundResponse(w, r)
		return
	}

	// Only draft invoices can be fully edited
	if invoice.Status != "draft" {
		h.BadRequestResponse(w, r, models.ErrInvoiceDraftOnly)
		return
	}

	var clientIDPtr *uuid.UUID
	if input.ClientID != "" {
		cID, err := uuid.Parse(input.ClientID)
		if err != nil {
			h.BadRequestResponse(w, r, ErrInvalidClientID)
			return
		}
		clientIDPtr = &cID
	} else if input.ClientEmail != "" {
		existingClient, err := h.Models.Clients.GetByEmail(r.Context(), input.ClientEmail, user.ID)
		if err == nil && existingClient != nil {
			clientIDPtr = &existingClient.ID
		} else {
			name := "Direct Client"
			if parts := strings.Split(input.ClientEmail, "@"); len(parts) > 0 {
				name = parts[0]
			}
			newClient := &models.Client{
				UserID:  user.ID,
				Name:    name,
				Email:   input.ClientEmail,
				Address: input.ClientAddress,
			}
			err = h.Models.Clients.Insert(r.Context(), newClient)
			if err == nil {
				clientIDPtr = &newClient.ID
			}
		}
	}

	issueDate, err := time.Parse("2006-01-02", input.IssueDate)
	if err != nil {
		h.BadRequestResponse(w, r, models.ErrInvalidDateFormat)
		return
	}
	dueDate, err := time.Parse("2006-01-02", input.DueDate)
	if err != nil {
		h.BadRequestResponse(w, r, models.ErrInvalidDateFormat)
		return
	}

	invoice.ClientID = clientIDPtr
	invoice.Currency = input.Currency
	invoice.IssueDate = issueDate
	invoice.DueDate = dueDate
	invoice.Notes = input.Notes

	var subtotal float64
	var lineItems []models.LineItem
	for _, it := range input.Items {
		amt := it.Quantity * it.UnitPrice
		subtotal += amt
		lineItems = append(lineItems, models.LineItem{
			Description: it.Description,
			Quantity:    it.Quantity,
			UnitPrice:   it.UnitPrice,
			Amount:      amt,
		})
	}

	tax := 0.0
	if input.TaxRate > 0 {
		tax = subtotal * (input.TaxRate / 100.0)
	}
	discount := input.DiscountAmount
	total := subtotal + tax - discount
	if total < 0 {
		total = 0
	}

	invoice.Subtotal = subtotal
	invoice.Tax = tax
	invoice.Discount = discount
	invoice.Total = total
	invoice.LineItems = lineItems
	invoice.Status = "draft"
	if !input.SaveAsDraft && input.Action == "dispatch" {
		balance, err := h.Models.CreditTxn.GetBalanceByUserID(r.Context(), user.ID)
		if err != nil || balance < 1 {
			h.BadRequestResponse(w, r, errors.New("insufficient credits for email dispatch"))
			return
		}
		invoice.Status = "sent"
	}

	err = h.Services.Invoice.UpdateInvoice(r.Context(), invoice)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	if input.Action == "dispatch" {
		var client *models.Client
		if clientIDPtr != nil {
			client, _ = h.Services.Client.GetByID(r.Context(), *clientIDPtr, user.ID)
		} else if input.ClientEmail != "" {
			client = &models.Client{
				Name:    input.ClientEmail,
				Email:   input.ClientEmail,
				Address: input.ClientAddress,
			}
		}
		if client != nil && client.Email != "" {
			profile, err := h.Services.Auth.GetBusinessProfile(r.Context(), user.ID)
			if err == nil && profile != nil {
				err = h.Services.Credit.DeductCredits(r.Context(), user.ID, 1, fmt.Sprintf("Email Dispatch for Invoice %s", invoice.InvoiceNumber))
				if err == nil {
					_ = h.Services.Invoice.DispatchInvoiceEmail(r.Context(), invoice, profile, client)
				}
			}
		}
	}

	_ = h.WriteJSON(w, http.StatusOK, map[string]any{"invoice": invoice},nil)
}
