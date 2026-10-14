package frontend

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSPAHandler_ServeHTTP(t *testing.T) {
	handler := NewHandler("dist")

	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
		expectInBody   string
		expectHeader   string
	}{
		{
			name:           "Root path serves index.html",
			method:         http.MethodGet,
			url:            "/",
			expectedStatus: http.StatusOK,
			expectInBody:   "<div id=\"app\"",
			expectHeader:   "no-cache",
		},
		{
			name:           "Client-side route serves index.html",
			method:         http.MethodGet,
			url:            "/user/dashboard",
			expectedStatus: http.StatusOK,
			expectInBody:   "<div id=\"app\"",
			expectHeader:   "no-cache",
		},
		{
			name:           "Webmanifest is served with proper content type",
			method:         http.MethodGet,
			url:            "/manifest.webmanifest",
			expectedStatus: http.StatusOK,
			expectInBody:   "Teks-Invoice",
			expectHeader:   "no-cache",
		},
		{
			name:           "Method POST is rejected on static handler",
			method:         http.MethodPost,
			url:            "/user/dashboard",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.expectInBody != "" && !strings.Contains(rec.Body.String(), tt.expectInBody) {
				t.Errorf("expected body to contain %q, but got %q", tt.expectInBody, rec.Body.String())
			}

			if tt.expectHeader != "" {
				cacheControl := rec.Header().Get("Cache-Control")
				if !strings.Contains(cacheControl, tt.expectHeader) {
					t.Errorf("expected Cache-Control to contain %q, got %q", tt.expectHeader, cacheControl)
				}
			}
		})
	}
}
