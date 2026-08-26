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
}

func NewServices(models models.Models, mailer mailer.Mailer, paystack *paystack.Client, jwtManager *tokens.JWTManager) Services {
	financeService := NewFinanceService(models)
	invoiceService := NewInvoiceService(models, financeService)
	authService := NewAuthService(models, mailer, jwtManager)
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
	}
}
