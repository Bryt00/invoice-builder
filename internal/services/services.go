package services

import (
	"raven.go.invoice-builder/internal/mailer"
	"raven.go.invoice-builder/internal/models"
	"raven.go.invoice-builder/internal/paystack"
	"raven.go.invoice-builder/internal/tokens"
)

type Services struct {
	Auth    AuthService
	Invoice InvoiceService
	Finance FinanceService
	Client  ClientService
	Admin   AdminService
	Credit  CreditService
	Paystack *paystack.Client
}

func NewServices(models models.Models, mailer mailer.Mailer, paystack *paystack.Client, jwtManager *tokens.JWTManager, frontendURL string) Services {
	financeService := NewFinanceService(models)
	invoiceService := NewInvoiceService(models, financeService, mailer, frontendURL)
	authService := NewAuthService(models, mailer, jwtManager, frontendURL)
	clientService := NewClientService(models)
	adminService := NewAdminService(models)
	creditService := NewCreditService(models)

	return Services{
		Auth:    authService,
		Invoice: invoiceService,
		Finance: financeService,
		Client:  clientService,
		Admin:   adminService,
		Credit:  creditService,
		Paystack: paystack,
	}
}
