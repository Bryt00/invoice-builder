package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"raven.go.invoice-builder/internal/jsonlog"
)

func TestSPA_RoutesIntegration(t *testing.T) {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelOff)
	stdLog := log.New(os.Stdout, "", 0)

	app := &application{
		errorLog:   stdLog,
		infoLog:    stdLog,
		jsonLogger: logger,
		apiHandler: &ApiHandler{
			JsonLogger: logger,
		},
	}
	app.config.admin.secretPath = "admin"

	routes := app.routes()

	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
		expectInBody   string
	}{
		{
			name:           "Root GET / serves Vue index.html",
			method:         http.MethodGet,
			url:            "/",
			expectedStatus: http.StatusOK,
			expectInBody:   "<div id=\"app\"",
		},
		{
			name:           "SPA client route GET /user/dashboard serves index.html",
			method:         http.MethodGet,
			url:            "/user/dashboard",
			expectedStatus: http.StatusOK,
			expectInBody:   "<div id=\"app\"",
		},
		{
			name:           "PWA webmanifest is served directly",
			method:         http.MethodGet,
			url:            "/manifest.webmanifest",
			expectedStatus: http.StatusOK,
			expectInBody:   "Teks-Invoice",
		},
		{
			name:           "Unknown API path returns JSON 404, not HTML",
			method:         http.MethodGet,
			url:            "/api/v1/unknown-endpoint",
			expectedStatus: http.StatusNotFound,
			expectInBody:   "the requested resource could not be found",
		},
		{
			name:           "Healthz endpoint works",
			method:         http.MethodGet,
			url:            "/healthz",
			expectedStatus: http.StatusOK,
			expectInBody:   "available",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			rec := httptest.NewRecorder()

			routes.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d. Body: %s", tt.expectedStatus, rec.Code, rec.Body.String())
			}

			if tt.expectInBody != "" && !strings.Contains(rec.Body.String(), tt.expectInBody) {
				t.Errorf("expected body to contain %q, but got %q", tt.expectInBody, rec.Body.String())
			}

			// Verify Content-Security-Policy header does NOT block scripts
			csp := rec.Header().Get("Content-Security-Policy")
			if strings.Contains(csp, "default-src 'none'") {
				t.Errorf("CSP should not be default-src 'none', got %q", csp)
			}
		})
	}
}
