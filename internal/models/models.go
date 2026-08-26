package models

import (
	"context"

	"gorm.io/gorm"
)

type Models struct {
	DB               *gorm.DB
	Users            UserInterface
	Roles            RoleInterface
	BusinessProfiles BusinessProfileInterface
	Clients          ClientInterface
	Invoice          InvoiceInterface
	LineItem         LineItemInterface
	Receipt          ReceiptInterface
	Payment          PaymentInterface
	CreditTxn        CreditTxnInterface
	CreditPackages   CreditPackageInterface
	AuditLog         AuditLogInterface
	WebhookLogs      WebhookLogInterface
	SystemSettings   SystemSettingInterface
	Finance          FinanceInterface
	RefreshTokens    RefreshTokenInterface
}

func NewModel(db *gorm.DB) Models {
	return Models{
		DB:               db,
		Users:            &UserModel{DB: db},
		Roles:            &RoleModel{DB: db},
		BusinessProfiles: &BusinessProfileModel{DB: db},
		Clients:          &ClientModel{DB: db},
		Invoice:          &InvoiceModel{DB: db},
		LineItem:         &LineItemModel{DB: db},
		Receipt:          &ReceiptModel{DB: db},
		Payment:          &PaymentModel{DB: db},
		CreditTxn:        &CreditTxnModel{DB: db},
		CreditPackages:   &CreditPackageModel{DB: db},
		AuditLog:         &AuditLogModel{DB: db},
		WebhookLogs:      &WebhookLogModel{DB: db},
		SystemSettings:   &SystemSettingModel{DB: db},
		Finance:          &FinanceModel{DB: db},
		RefreshTokens:    &RefreshTokenModel{DB: db},
	}
}

func (m Models) Transaction(ctx context.Context, fn func(txModels Models) error) error {
	if m.DB == nil {
		return fn(m)
	}
	return m.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txModels := NewModel(tx)
		return fn(txModels)
	})
}
