package main

import (
	"fmt"
	"net/http"
	"strconv"

	"raven.go.invoice-builder/internal/models"
	"raven.go.invoice-builder/internal/systemlog"
)

// AdminGetSystemLogs handles querying structured system logs from SQLite.
func (h *ApiHandler) AdminGetSystemLogs(w http.ResponseWriter, r *http.Request) {
	page, limit := h.ParsePagination(r)
	level := r.URL.Query().Get("level")
	source := r.URL.Query().Get("source")
	search := r.URL.Query().Get("search")
	statusType := r.URL.Query().Get("status_type")

	statusCode := 0
	if codeStr := r.URL.Query().Get("status_code"); codeStr != "" {
		if c, err := strconv.Atoi(codeStr); err == nil {
			statusCode = c
		}
	}

	filter := systemlog.Filter{
		Level:      level,
		Source:     source,
		Search:     search,
		StatusCode: statusCode,
		StatusType: statusType,
		Page:       page,
		Limit:      limit,
	}

	logger := systemlog.Get()
	if logger == nil {
		h.ServerErrorResponse(w, r, fmt.Errorf("system logger is not initialized"))
		return
	}

	logs, totalCount, err := logger.GetLogs(filter)
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

// AdminGetSystemLogStats retrieves metrics and 24h summary stats from the SQLite store.
func (h *ApiHandler) AdminGetSystemLogStats(w http.ResponseWriter, r *http.Request) {
	logger := systemlog.Get()
	if logger == nil {
		h.ServerErrorResponse(w, r, fmt.Errorf("system logger is not initialized"))
		return
	}

	stats, err := logger.GetStats()
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status": "success",
		"stats":  stats,
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

// AdminClearSystemLogs purges all system logs from SQLite.
func (h *ApiHandler) AdminClearSystemLogs(w http.ResponseWriter, r *http.Request) {
	user := h.ContextGetUser(r)
	if user == nil {
		h.AuthenticationRequiredResponse(w, r)
		return
	}

	logger := systemlog.Get()
	if logger == nil {
		h.ServerErrorResponse(w, r, fmt.Errorf("system logger is not initialized"))
		return
	}

	if err := logger.ClearLogs(); err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	// Audit the clear logs action in PostgreSQL
	_ = h.Models.AuditLog.Record(r.Context(), &models.AuditLog{
		UserID:     &user.ID,
		Action:     "clear_system_logs",
		EntityType: "system",
		EntityID:   "system_logs",
		IPAddress:  getClientIP(r),
		UserAgent:  r.UserAgent(),
		Metadata:   `{"cleared_by": "admin"}`,
	})

	err := h.WriteJSON(w, http.StatusOK, Envelope{
		"status":  "success",
		"message": "All system logs cleared successfully",
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}
