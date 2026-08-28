package main

import (
	"raven.go.invoice-builder/internal/jsonlog"
	"raven.go.invoice-builder/internal/models"
	"raven.go.invoice-builder/internal/services"
	"raven.go.invoice-builder/internal/tokens"
)

type Envelope map[string]any

type ApiHandler struct {
	Services   services.Services
	Models     models.Models
	JsonLogger *jsonlog.Logger
	JWTManager *tokens.JWTManager
	FrontendURL string
}

func NewApiHandler(srv services.Services, m models.Models, logger *jsonlog.Logger, jwtManager *tokens.JWTManager, frontendURL string) *ApiHandler {
	return &ApiHandler{
		Services:   srv,
		Models:     m,
		JsonLogger: logger,
		JWTManager: jwtManager,
		FrontendURL: frontendURL,
	}
}
