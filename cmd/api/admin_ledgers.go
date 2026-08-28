package main

import (
	"net/http"

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
