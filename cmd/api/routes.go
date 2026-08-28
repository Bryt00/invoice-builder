package main

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"raven.go.invoice-builder/ui"
)

func (app *application) routes() http.Handler {
	mux := httprouter.New()

	mux.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.apiHandler.NotFoundResponse(w, r)
	})

	mux.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.apiHandler.ErrorResponse(w, r, http.StatusMethodNotAllowed, "method not allowed")
	})

	// Static asset server (logos, uploaded files)
	fileServer := http.FileServer(http.FS(ui.Files))
	staticCache := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		w.Header().Set("Vary", "Accept-Encoding")
		
		// Serve user uploads dynamically from disk instead of embed.FS
		if strings.HasPrefix(r.URL.Path, "/asset/img/uploads/") {
			http.StripPrefix("/asset/img/uploads/", http.FileServer(http.Dir("./ui/asset/img/uploads"))).ServeHTTP(w, r)
			return
		}
		
		fileServer.ServeHTTP(w, r)
	})
	mux.Handler(http.MethodGet, "/asset/*filepath", staticCache)
	mux.HandlerFunc(http.MethodGet, "/healthz", app.healthCheck)

	// API Middleware Chains
	apiAuth := alice.New(app.rateLimitAuth)
	apiProtected := alice.New(app.apiHandler.RequireJWTAuthentication)

	// Auth API Endpoints
	mux.Handler(http.MethodPost, "/api/v1/auth/login", apiAuth.ThenFunc(app.apiHandler.LoginUser))
	mux.Handler(http.MethodPost, "/api/v1/auth/register", apiAuth.ThenFunc(app.apiHandler.RegisterUser))
	mux.Handler(http.MethodPost, "/api/v1/auth/refresh", apiAuth.ThenFunc(app.apiHandler.RefreshToken))
	mux.Handler(http.MethodGet, "/api/v1/auth/verify", apiAuth.ThenFunc(app.apiHandler.VerifyUser))
	mux.Handler(http.MethodPost, "/api/v1/auth/resend-verification", apiAuth.ThenFunc(app.apiHandler.ResendVerification))
	mux.Handler(http.MethodPost, "/api/v1/auth/forgot-password", apiAuth.ThenFunc(app.apiHandler.ForgotPassword))
	mux.Handler(http.MethodPost, "/api/v1/auth/reset-password", apiAuth.ThenFunc(app.apiHandler.ResetPassword))
	mux.Handler(http.MethodGet, "/api/v1/auth/me", apiProtected.ThenFunc(app.apiHandler.GetMe))

	// Profile API Endpoints
	mux.Handler(http.MethodGet, "/api/v1/profile", apiProtected.ThenFunc(app.apiHandler.GetProfile))
	mux.Handler(http.MethodPut, "/api/v1/profile", apiProtected.ThenFunc(app.apiHandler.UpdateProfile))
	mux.Handler(http.MethodPost, "/api/v1/profile/logo", apiProtected.ThenFunc(app.apiHandler.UploadProfileLogo))

	// Clients API Endpoints
	mux.Handler(http.MethodGet, "/api/v1/clients", apiProtected.ThenFunc(app.apiHandler.GetClients))
	mux.Handler(http.MethodPost, "/api/v1/clients", apiProtected.ThenFunc(app.apiHandler.CreateClient))
	mux.Handler(http.MethodPut, "/api/v1/clients", apiProtected.ThenFunc(app.apiHandler.UpdateClient))
	mux.Handler(http.MethodPost, "/api/v1/clients/delete", apiProtected.ThenFunc(app.apiHandler.DeleteClient))

	// Invoices API Endpoints
	mux.Handler(http.MethodGet, "/api/v1/invoices", apiProtected.ThenFunc(app.apiHandler.GetInvoices))
	mux.Handler(http.MethodPost, "/api/v1/invoices", apiProtected.ThenFunc(app.apiHandler.CreateInvoice))
	mux.Handler(http.MethodPut, "/api/v1/invoices", apiProtected.ThenFunc(app.apiHandler.UpdateInvoice))
	mux.Handler(http.MethodGet, "/api/v1/invoices/view", apiProtected.ThenFunc(app.apiHandler.GetInvoice))
	mux.Handler(http.MethodGet, "/api/v1/invoices/download", apiProtected.ThenFunc(app.apiHandler.DownloadInvoice))
	mux.Handler(http.MethodPost, "/api/v1/invoices/mark-paid", apiProtected.ThenFunc(app.apiHandler.MarkInvoicePaid))
	mux.Handler(http.MethodGet, "/api/v1/invoices/public", http.HandlerFunc(app.apiHandler.GetPublicInvoice))
	mux.Handler(http.MethodGet, "/api/v1/invoices/public/download", http.HandlerFunc(app.apiHandler.DownloadPublicInvoice))
	
	mux.Handler(http.MethodPost, "/api/v1/invoices/receipts", apiProtected.ThenFunc(app.apiHandler.GenerateReceipt))
	mux.Handler(http.MethodGet, "/api/v1/invoices/receipts/view", apiProtected.ThenFunc(app.apiHandler.GetReceipt))
	mux.Handler(http.MethodGet, "/api/v1/invoices/receipts/download", apiProtected.ThenFunc(app.apiHandler.DownloadReceipt))

	// Finance API Endpoints
	mux.Handler(http.MethodGet, "/api/v1/finance/transactions", apiProtected.ThenFunc(app.apiHandler.GetFinanceTransactions))
	mux.Handler(http.MethodGet, "/api/v1/finance/export", apiProtected.ThenFunc(app.apiHandler.ExportFinanceTransactionsCSV))
	mux.Handler(http.MethodPost, "/api/v1/finance/transactions", apiProtected.ThenFunc(app.apiHandler.CreateFinanceTransaction))
	mux.Handler(http.MethodPost, "/api/v1/finance/transactions/delete", apiProtected.ThenFunc(app.apiHandler.DeleteFinanceTransaction))
	mux.Handler(http.MethodGet, "/api/v1/finance/summary", apiProtected.ThenFunc(app.apiHandler.GetFinanceSummary))
	mux.Handler(http.MethodGet, "/api/v1/finance/categories", apiProtected.ThenFunc(app.apiHandler.GetFinanceCategories))
	mux.Handler(http.MethodPost, "/api/v1/finance/categories", apiProtected.ThenFunc(app.apiHandler.CreateFinanceCategory))

	// Credits API Endpoints
	mux.Handler(http.MethodGet, "/api/v1/credits/balance", apiProtected.ThenFunc(app.apiHandler.GetCreditsBalance))
	mux.Handler(http.MethodGet, "/api/v1/credits/history", apiProtected.ThenFunc(app.apiHandler.GetCreditsHistory))
	mux.Handler(http.MethodGet, "/api/v1/credits/packages", http.HandlerFunc(app.apiHandler.GetCreditsPackages))
	mux.Handler(http.MethodPost, "/api/v1/credits/topup/initialize", apiProtected.ThenFunc(app.apiHandler.InitializeCreditsTopup))
	mux.Handler(http.MethodGet, "/api/v1/credits/topup/verify", apiProtected.ThenFunc(app.apiHandler.VerifyCreditsTopup))

	// Admin API Endpoints (Obscured Path & Secret Header Safeguarded)
	apiAdmin := alice.New(app.rateLimitAdmin, app.requireAdminSecretHeader, app.apiHandler.RequireJWTAuthentication, app.apiHandler.RequireRole("Admin"))
	adminBasePath := "/api/v1/" + app.config.admin.secretPath
	
	// Admin Dashboard
	mux.Handler(http.MethodGet, adminBasePath+"/stats", apiAdmin.ThenFunc(app.apiHandler.AdminGetStats))
	
	// Admin Users
	mux.Handler(http.MethodGet, adminBasePath+"/users", apiAdmin.ThenFunc(app.apiHandler.AdminGetUsers))
	mux.Handler(http.MethodPost, adminBasePath+"/users", apiAdmin.ThenFunc(app.apiHandler.AdminCreateUser))
	mux.Handler(http.MethodGet, adminBasePath+"/users/:id", apiAdmin.ThenFunc(app.apiHandler.AdminGetUser))
	mux.Handler(http.MethodPut, adminBasePath+"/users", apiAdmin.ThenFunc(app.apiHandler.AdminUpdateUser))
	mux.Handler(http.MethodDelete, adminBasePath+"/users/:id", apiAdmin.ThenFunc(app.apiHandler.AdminDeleteUser))
	mux.Handler(http.MethodPost, adminBasePath+"/users/credits", apiAdmin.ThenFunc(app.apiHandler.AdminGrantCredits))

	// Admin Packages
	mux.Handler(http.MethodGet, adminBasePath+"/packages", apiAdmin.ThenFunc(app.apiHandler.AdminGetPackages))
	mux.Handler(http.MethodPost, adminBasePath+"/packages", apiAdmin.ThenFunc(app.apiHandler.AdminCreatePackage))
	mux.Handler(http.MethodPut, adminBasePath+"/packages", apiAdmin.ThenFunc(app.apiHandler.AdminUpdatePackage))
	mux.Handler(http.MethodDelete, adminBasePath+"/packages/:id", apiAdmin.ThenFunc(app.apiHandler.AdminDeletePackage))

	// Admin Ledgers
	mux.Handler(http.MethodGet, adminBasePath+"/invoices", apiAdmin.ThenFunc(app.apiHandler.AdminGetInvoices))
	mux.Handler(http.MethodGet, adminBasePath+"/payments", apiAdmin.ThenFunc(app.apiHandler.AdminGetPayments))
	mux.Handler(http.MethodGet, adminBasePath+"/credits", apiAdmin.ThenFunc(app.apiHandler.AdminGetCredits))
	mux.Handler(http.MethodGet, adminBasePath+"/audit-logs", apiAdmin.ThenFunc(app.apiHandler.AdminGetAuditLogs))
	
	// Admin Webhooks
	mux.Handler(http.MethodGet, adminBasePath+"/webhooks", apiAdmin.ThenFunc(app.apiHandler.AdminGetWebhooks))
	mux.Handler(http.MethodPost, adminBasePath+"/webhooks/:id/replay", apiAdmin.ThenFunc(app.apiHandler.AdminReplayWebhook))

	// Admin Settings
	mux.Handler(http.MethodGet, adminBasePath+"/settings", apiAdmin.ThenFunc(app.apiHandler.AdminGetSettings))
	mux.Handler(http.MethodPut, adminBasePath+"/settings", apiAdmin.ThenFunc(app.apiHandler.AdminUpdateSettings))

	// CORS & Global Security Middleware
	standard := alice.New(app.recoverPanic, app.logRequest, enableCORS, secureHeaders)

	return standard.Then(mux)
}
