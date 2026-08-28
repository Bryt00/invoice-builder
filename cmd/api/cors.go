package main

import (
	"net/http"
	"strings"
)

// CORSMiddleware handles Cross-Origin Resource Sharing for the Vue SPA frontend
// and mobile clients.
func (h *ApiHandler) CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Allow known origins
		allowedOrigins := []string{
			"https://teks-invoice.com",
			"https://www.teks-invoice.com",
			"http://localhost:5173", // Vite dev server
			"http://localhost:3000",
		}

		for _, allowed := range allowedOrigins {
			if strings.EqualFold(origin, allowed) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Handle preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
