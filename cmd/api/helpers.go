package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

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

	ip := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ip = forwarded
	}

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
