package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"raven.go.invoice-builder/internal/models"
	"raven.go.invoice-builder/internal/pdf"
)

func (h *ApiHandler) GenerateReceipt(w http.ResponseWriter, r *http.Request) {
	user := h.ContextGetUser(r)
	if user == nil {
		h.AuthenticationRequiredResponse(w, r)
		return
	}

	var input struct {
		InvoiceID string `json:"invoice_id"`
	}

	err := h.ReadJSON(w, r, &input)
	if err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	invoiceID, err := uuid.Parse(input.InvoiceID)
	if err != nil {
		h.BadRequestResponse(w, r, ErrInvalidID)
		return
	}

	invoice, err := h.Services.Invoice.GetByID(r.Context(), invoiceID, user.ID)
	if err != nil {
		h.NotFoundResponse(w, r)
		return
	}

	if !invoice.IsPaid || invoice.Status == models.InvoiceStatusDraft {
		h.BadRequestResponse(w, r, models.ErrInvoiceNotPaid)
		return
	}

	// Check if receipt already exists
	existing, err := h.Models.Receipt.GetByInvoiceID(r.Context(), invoiceID)
	if err == nil && existing != nil {
		// Deduct 1 credit for regenerating
		err = h.Services.Credit.DeductCredits(r.Context(), user.ID, 1, "Receipt Generation (Duplicate)")
		if err != nil {
			h.BadRequestResponse(w, r, fmt.Errorf("failed to deduct credits: %w", err))
			return
		}
		
		_ = h.WriteJSON(w, http.StatusOK, map[string]any{"receipt": existing}, nil)
		return
	}

	// Deduct 1 credit
	err = h.Services.Credit.DeductCredits(r.Context(), user.ID, 1, "Receipt Generation")
	if err != nil {
		h.BadRequestResponse(w, r, fmt.Errorf("failed to deduct credits: %w", err))
		return
	}

	// Generate Receipt Number
	receiptNumber, err := h.Models.Receipt.GenerateReceiptNumber(r.Context(), user.ID)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	receipt := &models.Receipt{
		UserID:        user.ID,
		InvoiceID:     invoice.ID,
		ReceiptNumber: receiptNumber,
		Amount:        invoice.Total,
		Currency:      invoice.Currency,
		IssuedAt:      invoice.UpdatedAt,
	}

	err = h.Models.Receipt.Insert(r.Context(), receipt)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	_ = h.WriteJSON(w, http.StatusOK, map[string]any{"receipt": receipt}, nil)
}

func (h *ApiHandler) GetReceipt(w http.ResponseWriter, r *http.Request) {
	user := h.ContextGetUser(r)
	if user == nil {
		h.AuthenticationRequiredResponse(w, r)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.BadRequestResponse(w, r, ErrInvalidID)
		return
	}

	receipt, err := h.Models.Receipt.GetByID(r.Context(), id, user.ID)
	if err != nil || receipt == nil {
		h.NotFoundResponse(w, r)
		return
	}

	_ = h.WriteJSON(w, http.StatusOK, Envelope{"receipt": receipt}, nil)
}

func (h *ApiHandler) DownloadReceipt(w http.ResponseWriter, r *http.Request) {
	user := h.ContextGetUser(r)
	if user == nil {
		h.AuthenticationRequiredResponse(w, r)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.BadRequestResponse(w, r, ErrInvalidID)
		return
	}

	receipt, err := h.Models.Receipt.GetByID(r.Context(), id, user.ID)
	if err != nil || receipt == nil {
		h.NotFoundResponse(w, r)
		return
	}

	profile, _ := h.Models.BusinessProfiles.GetByUserID(r.Context(), user.ID)
	size := r.URL.Query().Get("size")
	if size == "" {
		size = "a4"
	}

	pdfBytes, err := pdf.GenerateReceiptPDF(receipt, receipt.Invoice, profile, size)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.pdf\"", receipt.ReceiptNumber))
	w.Header().Set("Content-Length", strconv.Itoa(len(pdfBytes)))
	_, _ = w.Write(pdfBytes)
}

func (h *ApiHandler) DispatchReceipt(w http.ResponseWriter, r *http.Request) {
	user := h.ContextGetUser(r)
	if user == nil {
		h.AuthenticationRequiredResponse(w, r)
		return
	}

	var input struct {
		ReceiptID   string `json:"receipt_id"`
		TargetEmail string `json:"email"`
	}

	err := h.ReadJSON(w, r, &input)
	if err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	receiptID, err := uuid.Parse(input.ReceiptID)
	if err != nil {
		h.BadRequestResponse(w, r, ErrInvalidID)
		return
	}

	receipt, err := h.Models.Receipt.GetByID(r.Context(), receiptID, user.ID)
	if err != nil || receipt == nil {
		h.NotFoundResponse(w, r)
		return
	}

	profile, _ := h.Models.BusinessProfiles.GetByUserID(r.Context(), user.ID)

	targetEmail := input.TargetEmail
	if targetEmail == "" && receipt.Invoice != nil && receipt.Invoice.Client != nil {
		targetEmail = receipt.Invoice.Client.Email
	}

	if targetEmail == "" {
		h.BadRequestResponse(w, r, fmt.Errorf("recipient email is required for receipt dispatch"))
		return
	}

	err = h.Services.Invoice.DispatchReceiptEmail(r.Context(), receipt, receipt.Invoice, profile, targetEmail)
	if err != nil {
		h.BadRequestResponse(w, r, fmt.Errorf("failed to send receipt email dispatch: %w", err))
		return
	}

	_ = h.WriteJSON(w, http.StatusOK, map[string]any{"status": "success", "message": "Receipt dispatched successfully"}, nil)
}
