package models

import "errors"

var (
	ErrNoRecord = errors.New("models: no matchinng record found")

	ErrInvalidCredentials = errors.New("models: invalid credentials")

	ErrDuplicateEmail = errors.New("models: duplicate email")

	ErrInvalidPackageID = errors.New("models: invalid package ID")
	ErrInvalidDateFormat = errors.New("models: invalid date format, expected YYYY-MM-DD")
	ErrInvoiceNil = errors.New("models: invoice is nil")
	ErrClientEmailRequired = errors.New("models: client email is required for dispatch")
	ErrAccountNotActivated = errors.New("models: user account is not activated")
	ErrActivationLinkCooldown = errors.New("models: please wait 5 minutes before requesting a new activation link")
	ErrJWTUninitialized = errors.New("models: jwtManager uninitialized")
	ErrWebhookNotFound = errors.New("models: webhook not found")
	ErrTokenRevoked = errors.New("models: refresh token has been revoked")
	ErrTokenExpired = errors.New("models: refresh token has expired")
	ErrInvoiceNotPaid = errors.New("models: invoice must be paid before generating a receipt")
	ErrReceiptAlreadyExists = errors.New("models: receipt already exists for this invoice")
	ErrInvoiceDraftOnly = errors.New("models: only draft invoices can be edited")
	ErrDraftCannotBePaid = errors.New("models: draft invoices cannot be marked as paid until dispatched")
	ErrTransactionAlreadyProcessed = errors.New("models: transaction has already been processed")
)
