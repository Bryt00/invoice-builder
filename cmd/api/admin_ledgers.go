package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

func (h *ApiHandler) AdminGetInvoices(w http.ResponseWriter, r *http.Request) {
	page, limit := h.ParsePagination(r)
	search := r.URL.Query().Get("search")
	status := r.URL.Query().Get("status")

	invoices, totalCount, err := h.Services.Admin.GetAllSystemInvoices(r.Context(), search, status, page, limit)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status":   "success",
		"invoices": invoices,
		"meta": Envelope{
			"total_count": totalCount,
			"page":        page,
			"limit":       limit,
		},
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) AdminGetPayments(w http.ResponseWriter, r *http.Request) {
	page, limit := h.ParsePagination(r)
	status := r.URL.Query().Get("status")

	payments, totalCount, err := h.Services.Admin.GetAllSystemPayments(r.Context(), status, page, limit)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status":   "success",
		"payments": payments,
		"meta": Envelope{
			"total_count": totalCount,
			"page":        page,
			"limit":       limit,
		},
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) AdminGetCredits(w http.ResponseWriter, r *http.Request) {
	page, limit := h.ParsePagination(r)
	txnType := r.URL.Query().Get("type")

	txns, totalCount, err := h.Services.Admin.GetAllSystemCreditTxns(r.Context(), txnType, page, limit)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status":  "success",
		"credits": txns,
		"meta": Envelope{
			"total_count": totalCount,
			"page":        page,
			"limit":       limit,
		},
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) AdminGetAuditLogs(w http.ResponseWriter, r *http.Request) {
	page, limit := h.ParsePagination(r)
	search := r.URL.Query().Get("search")
	
	logs, totalCount, err := h.Services.Admin.GetAuditLogs(r.Context(), search, page, limit)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status": "success",
		"logs":   logs,
		"meta": Envelope{
			"total_count": totalCount,
			"page":        page,
			"limit":       limit,
		},
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) AdminGetWebhooks(w http.ResponseWriter, r *http.Request) {
	page, limit := h.ParsePagination(r)
	status := r.URL.Query().Get("status")

	logs, totalCount, err := h.Services.Admin.GetAllWebhookLogs(r.Context(), status, page, limit)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status":   "success",
		"webhooks": logs,
		"meta": Envelope{
			"total_count": totalCount,
			"page":        page,
			"limit":       limit,
		},
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) AdminReplayWebhook(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		h.NotFoundResponse(w, r)
		return
	}

	err = h.Services.Admin.ReplayWebhook(r.Context(), id)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{"status": "success", "message": "webhook queued for replay"}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) AdminDeleteInvoice(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		h.NotFoundResponse(w, r)
		return
	}

	err = h.Models.Invoice.AdminDelete(r.Context(), id)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{"status": "success", "message": "invoice deleted successfully"}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) AdminUpdateInvoiceStatus(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		h.NotFoundResponse(w, r)
		return
	}

	var input struct {
		Status string `json:"status"`
	}
	err = h.ReadJSON(w, r, &input)
	if err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	err = h.Models.Invoice.AdminUpdateStatus(r.Context(), id, input.Status)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{"status": "success", "message": "invoice status updated"}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) AdminExportInvoicesCSV(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	status := r.URL.Query().Get("status")

	invoices, _, err := h.Services.Admin.GetAllSystemInvoices(r.Context(), search, status, 1, 100000)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=\"invoices_export.csv\"")

	header := "ID,InvoiceNumber,UserID,ClientName,Amount,Currency,Status,IssueDate\n"
	_, _ = w.Write([]byte(header))

	for _, inv := range invoices {
		clientName := "Unknown"
		if inv.Client != nil {
			clientName = inv.Client.Name
		}
		clientName = strings.ReplaceAll(clientName, "\"", "\"\"")
		dateStr := inv.IssueDate.Format("2006-01-02")
		
		row := fmt.Sprintf("\"%s\",\"%s\",\"%s\",\"%s\",\"%.2f\",\"%s\",\"%s\",\"%s\"\n",
			inv.ID.String(), inv.InvoiceNumber, inv.UserID.String(), clientName, float64(inv.Total)/100.0, inv.Currency, inv.Status, dateStr)
		_, _ = w.Write([]byte(row))
	}
}

func (h *ApiHandler) AdminExportPaymentsCSV(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	payments, _, err := h.Services.Admin.GetAllSystemPayments(r.Context(), status, 1, 100000)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=\"payments_export.csv\"")

	header := "ID,InvoiceID,Amount,Currency,Provider,ProviderTxnID,Status,CreatedAt\n"
	_, _ = w.Write([]byte(header))

	for _, p := range payments {
		dateStr := p.CreatedAt.Format("2006-01-02 15:04:05")
		
		row := fmt.Sprintf("\"%s\",\"%s\",\"%.2f\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"\n",
			p.ID.String(), p.InvoiceID.String(), float64(p.Amount)/100.0, p.Currency, p.PaymentMethod, p.TransactionRef, p.Status, dateStr)
		_, _ = w.Write([]byte(row))
	}
}
