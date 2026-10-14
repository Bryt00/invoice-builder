package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"golang.org/x/time/rate"
	"raven.go.invoice-builder/internal/systemlog"
)

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *statusResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		allowedOrigins := []string{
			"https://teks-invoice.com",
			"https://www.teks-invoice.com",
			"http://localhost:5173", // Vite dev server
			"http://localhost:3000",
			"http://localhost:4000",
		}

		for _, allowed := range allowedOrigins {
			if strings.EqualFold(origin, allowed) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization, X-CSRF-Token, X-Admin-Secret-Key")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csp := "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' 'wasm-unsafe-eval'; " +
			"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; " +
			"font-src 'self' https://fonts.gstatic.com data:; " +
			"img-src 'self' data: blob: https:; " +
			"connect-src 'self' https://api.paystack.co https://api.paystack.com; " +
			"manifest-src 'self'; " +
			"frame-ancestors 'none';"
		w.Header().Set("Content-Security-Policy", csp)
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		latency := time.Since(start)
		ip := getClientIP(r)

		// Filter out static asset requests to keep logs high-signal and lean
		path := r.URL.Path
		if strings.HasPrefix(path, "/assets/") || strings.HasPrefix(path, "/asset/") ||
			path == "/favicon.ico" || path == "/sw.js" || strings.HasSuffix(path, ".webmanifest") {
			return
		}

		if sys := systemlog.Get(); sys != nil {
			sys.LogHTTP(r.Method, r.URL.RequestURI(), rw.statusCode, latency, ip, "")
		}
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				stack := string(debug.Stack())
				if sys := systemlog.Get(); sys != nil {
					sys.Error("PANIC", fmt.Sprintf("Recovered from panic on %s %s: %v", r.Method, r.URL.RequestURI(), err), nil, stack)
				}
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimitAuth(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.limiter.enabled {
			if !limiter.Allow() {
				app.rateLimitExceeded(w, r)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimitAdmin(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.limiter.enabled {
			if !limiter.Allow() {
				app.rateLimitExceeded(w, r)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAdminSecretHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := r.Header.Get("X-Admin-Secret-Key")
		if secret == "" || secret != app.config.admin.secretKey {
			app.apiHandler.ForbiddenResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) auditLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// We use a custom response writer to capture the status code
		// but for simplicity we can just log after next.ServeHTTP
		next.ServeHTTP(w, r)
		
		// Only log modifying requests
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch || r.Method == http.MethodDelete {
			// Skip logging for GET, OPTIONS etc.
			action := r.Method + " " + r.URL.Path
			// Don't log admin login alert as user activity? Actually it's fine.
			app.auditLog(r, action, "api_request", "", nil)
		}
	})
}

func (app *application) checkMaintenanceMode(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		maintMode, err := app.models.SystemSettings.Get(r.Context(), "maintenance_mode", "false")
		if err == nil && maintMode == "true" {
			user := app.apiHandler.ContextGetUser(r)
			if user == nil || !user.IsAdmin() {
				app.apiHandler.ErrorResponse(w, r, http.StatusServiceUnavailable, "The system is currently undergoing maintenance. Please try again later.")
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
