package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/google/uuid"
	"raven.go.invoice-builder/internal/contextutil"
	"raven.go.invoice-builder/internal/models"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Print(trace)
	app.apiHandler.ServerErrorResponse(w, nil, err)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	app.apiHandler.ErrorResponse(w, nil, status, http.StatusText(status))
}

func (app *application) notFound(w http.ResponseWriter) {
	app.apiHandler.NotFoundResponse(w, nil)
}

func (app *application) methodNotAllowed(w http.ResponseWriter) {
	app.apiHandler.ErrorResponse(w, nil, http.StatusMethodNotAllowed, "method not allowed")
}

func (app *application) rateLimitExceeded(w http.ResponseWriter, r *http.Request) {
	app.apiHandler.ErrorResponse(w, r, http.StatusTooManyRequests, "rate limit exceeded")
}

func (app *application) contextGetUser(r *http.Request) *models.User {
	return contextutil.ContextGetUser(r)
}

func (app *application) isAuthenticated(r *http.Request) bool {
	return app.contextGetUser(r) != nil
}

func (app *application) auditLog(r *http.Request, action string, entityType string, entityID string, metadata map[string]any) {
	user := app.contextGetUser(r)

	var userID *uuid.UUID
	if user != nil {
		id := user.ID
		userID = &id
	}

	metadataJSON := "{}"
	if metadata != nil {
		if bytes, err := json.Marshal(metadata); err == nil {
			metadataJSON = string(bytes)
		}
	}

	ip := getClientIP(r)

	log := &models.AuditLog{
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		IPAddress:  ip,
		UserAgent:  r.UserAgent(),
		Metadata:   metadataJSON,
	}

	go func() {
		_ = app.models.AuditLog.Record(context.Background(), log)
	}()
}

// getClientIP extracts the real visitor IP from Cloudflare or reverse proxy headers,
// falling back to r.RemoteAddr with port stripped.
func getClientIP(r *http.Request) string {
	// 1. Cloudflare header
	if cfIP := r.Header.Get("CF-Connecting-IP"); cfIP != "" {
		return strings.TrimSpace(cfIP)
	}

	// 2. Standard X-Real-IP header from Nginx
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return strings.TrimSpace(realIP)
	}

	// 3. X-Forwarded-For header (comma-separated list: client, proxy1, proxy2)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			ip := strings.TrimSpace(parts[0])
			if ip != "" {
				return ip
			}
		}
	}

	// 4. Direct socket address (strip port if present)
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}
	return r.RemoteAddr
}
