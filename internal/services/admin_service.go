package services

import (
	"context"
	"errors"
	"strconv"

	"github.com/google/uuid"
	"raven.go.invoice-builder/internal/models"
)

type AdminService interface {
	GetDashboardStats(ctx context.Context) (totalUsers, activeUsers, unverifiedUsers int64, totalRevenue float64, totalPayments int64, totalPurchasedCredits, totalUsedCredits int64, recentUsers []*models.User, recentLogs []*models.AuditLog, err error)
	GetAllUsers(ctx context.Context, search string, page, limit int) ([]*models.User, int64, error)
	GetRoles(ctx context.Context) ([]*models.Role, error)
	ToggleUserActivation(ctx context.Context, userID uuid.UUID, isActivated bool) error
	UpdateUserRole(ctx context.Context, userID uuid.UUID, roleName string) error
	CreateUser(ctx context.Context, name, email, password, roleName string, isActivated bool) (*models.User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, name, email, roleName string, isActivated bool) (*models.User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	GetAuditLogs(ctx context.Context, action string, page, limit int) ([]*models.AuditLog, int64, error)
	GetSystemSettings(ctx context.Context) (map[string]any, error)
	UpdateSystemSettings(ctx context.Context, bonus, maint, support, terms, privacy, refund, security string) error
	GetCreditPackages(ctx context.Context, includeInactive bool) ([]*models.CreditPackage, error)
	CreateCreditPackage(ctx context.Context, pkg *models.CreditPackage) error
	UpdateCreditPackage(ctx context.Context, pkg *models.CreditPackage) error
	ToggleCreditPackage(ctx context.Context, pkgID uuid.UUID, isActive bool) error
	DeleteCreditPackage(ctx context.Context, pkgID uuid.UUID) error
	GetAllSystemPayments(ctx context.Context, status string, page, limit int) ([]*models.Payment, int64, error)
	GetAllWebhookLogs(ctx context.Context, status string, page, limit int) ([]*models.WebhookLog, int64, error)
	ReplayWebhook(ctx context.Context, webhookID uuid.UUID) error
}

type adminService struct {
	models models.Models
}

func (s *adminService) GetRoles(ctx context.Context) ([]*models.Role, error) {
	return s.models.Roles.GetAll(ctx)
}

func NewAdminService(models models.Models) AdminService {
	return &adminService{models: models}
}

func (s *adminService) GetDashboardStats(ctx context.Context) (totalUsers, activeUsers, unverifiedUsers int64, totalRevenue float64, totalPayments int64, totalPurchasedCredits, totalUsedCredits int64, recentUsers []*models.User, recentLogs []*models.AuditLog, err error) {
	totalUsers, activeUsers, unverifiedUsers, err = s.models.Users.GetSystemUserStats(ctx)
	if err != nil {
		return
	}
	totalRevenue, totalPayments, err = s.models.Payment.GetSystemRevenueStats(ctx)
	if err != nil {
		return
	}
	purchased, used, err := s.models.CreditTxn.GetSystemCreditStats(ctx)
	if err != nil {
		return
	}
	totalPurchasedCredits = int64(purchased)
	totalUsedCredits = int64(used)
	recentUsers, _, _ = s.models.Users.GetAllUsers(ctx, "", 1, 5)
	recentLogs, _, _ = s.models.AuditLog.GetLogs(ctx, "", 1, 6)
	return
}

func (s *adminService) GetAllUsers(ctx context.Context, search string, page, limit int) ([]*models.User, int64, error) {
	return s.models.Users.GetAllUsers(ctx, search, page, limit)
}

func (s *adminService) ToggleUserActivation(ctx context.Context, userID uuid.UUID, isActivated bool) error {
	return s.models.Users.ToggleActivation(ctx, userID, isActivated)
}

func (s *adminService) UpdateUserRole(ctx context.Context, userID uuid.UUID, roleName string) error {
	role, err := s.models.Roles.GetByName(ctx, roleName)
	if err != nil {
		return err
	}
	return s.models.Users.UpdateRole(ctx, userID, role.ID)
}

func (s *adminService) CreateUser(ctx context.Context, name, email, password, roleName string, isActivated bool) (*models.User, error) {
	var newUser *models.User
	err := s.models.Transaction(ctx, func(txModels models.Models) error {
		role, err := txModels.Roles.GetByName(ctx, roleName)
		if err != nil {
			role = &models.Role{ID: uuid.New(), Name: roleName}
			_ = txModels.Roles.Create(ctx, role)
		}

		newUser = &models.User{
			BaseModel:         models.BaseModel{ID: uuid.New()},
			Name:              name,
			Email:             email,
			PasswordHash:      password,
			RoleID:            role.ID,
			IsActivated:       isActivated,
			IsProfileComplete: true,
		}

		err = txModels.Users.Insert(ctx, newUser)
		if err != nil {
			return err
		}

		if isActivated {
			bonusStr, _ := txModels.SystemSettings.Get(ctx, "default_signup_credits", "3")
			bonus, _ := strconv.Atoi(bonusStr)
			if bonus > 0 {
				_ = txModels.CreditTxn.AdminGrantCredits(ctx, newUser.ID, bonus, "Signup Bonus")
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return newUser, nil
}

func (s *adminService) UpdateUser(ctx context.Context, userID uuid.UUID, name, email, roleName string, isActivated bool) (*models.User, error) {
	user, err := s.models.Users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.Name = name
	user.Email = email
	user.IsActivated = isActivated

	if roleName != "" {
		role, err := s.models.Roles.GetByName(ctx, roleName)
		if err == nil {
			user.RoleID = role.ID
		}
	}

	err = s.models.Users.Update(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *adminService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return s.models.Users.Delete(ctx, userID)
}

func (s *adminService) GetAuditLogs(ctx context.Context, action string, page, limit int) ([]*models.AuditLog, int64, error) {
	return s.models.AuditLog.GetLogs(ctx, action, page, limit)
}

func (s *adminService) GetSystemSettings(ctx context.Context) (map[string]any, error) {
	settings, err := s.models.SystemSettings.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	signupBonus, _ := s.models.SystemSettings.Get(ctx, "default_signup_credits", "3")
	maintMode, _ := s.models.SystemSettings.Get(ctx, "maintenance_mode", "false")
	supportEmail, _ := s.models.SystemSettings.Get(ctx, "support_email", "support@teks-invoice.com")
	legalTerms, _ := s.models.SystemSettings.Get(ctx, "legal_terms", "")
	legalPrivacy, _ := s.models.SystemSettings.Get(ctx, "legal_privacy", "")
	legalRefund, _ := s.models.SystemSettings.Get(ctx, "legal_refund", "")
	legalSecurity, _ := s.models.SystemSettings.Get(ctx, "legal_security", "")

	return map[string]any{
		"Settings":            settings,
		"DefaultSignupBonus":  signupBonus,
		"MaintenanceMode":     maintMode == "true",
		"SupportContactEmail": supportEmail,
		"LegalTerms":          legalTerms,
		"LegalPrivacy":        legalPrivacy,
		"LegalRefund":         legalRefund,
		"LegalSecurity":       legalSecurity,
	}, nil
}

func (s *adminService) UpdateSystemSettings(ctx context.Context, bonus, maint, support, terms, privacy, refund, security string) error {
	_ = s.models.SystemSettings.Set(ctx, "default_signup_credits", bonus, "Number of free credits granted to new users upon signup")
	_ = s.models.SystemSettings.Set(ctx, "maintenance_mode", maint, "Toggle platform maintenance mode")
	_ = s.models.SystemSettings.Set(ctx, "support_email", support, "Platform support contact email")
	_ = s.models.SystemSettings.Set(ctx, "legal_terms", terms, "Terms of Service content")
	_ = s.models.SystemSettings.Set(ctx, "legal_privacy", privacy, "Privacy Policy content")
	_ = s.models.SystemSettings.Set(ctx, "legal_refund", refund, "Refund & Credit Policy content")
	_ = s.models.SystemSettings.Set(ctx, "legal_security", security, "Security & Data Integrity content")
	return nil
}

func (s *adminService) GetCreditPackages(ctx context.Context, includeInactive bool) ([]*models.CreditPackage, error) {
	return s.models.CreditPackages.GetAll(ctx, includeInactive)
}

func (s *adminService) CreateCreditPackage(ctx context.Context, pkg *models.CreditPackage) error {
	if pkg.ID == uuid.Nil {
		pkg.ID = uuid.New()
	}
	return s.models.CreditPackages.Create(ctx, pkg)
}

func (s *adminService) UpdateCreditPackage(ctx context.Context, pkg *models.CreditPackage) error {
	return s.models.CreditPackages.Update(ctx, pkg)
}

func (s *adminService) ToggleCreditPackage(ctx context.Context, pkgID uuid.UUID, isActive bool) error {
	return s.models.CreditPackages.ToggleActive(ctx, pkgID, isActive)
}

func (s *adminService) DeleteCreditPackage(ctx context.Context, pkgID uuid.UUID) error {
	return s.models.CreditPackages.Delete(ctx, pkgID)
}

func (s *adminService) GetAllSystemPayments(ctx context.Context, status string, page, limit int) ([]*models.Payment, int64, error) {
	return s.models.Payment.GetAllSystemPayments(ctx, status, page, limit)
}

func (s *adminService) GetAllWebhookLogs(ctx context.Context, status string, page, limit int) ([]*models.WebhookLog, int64, error) {
	return s.models.WebhookLogs.GetAll(ctx, status, page, limit)
}

func (s *adminService) ReplayWebhook(ctx context.Context, webhookID uuid.UUID) error {
	webhook, err := s.models.WebhookLogs.GetByID(ctx, webhookID)
	if err != nil {
		return err
	}
	if webhook == nil {
		return errors.New("webhook not found")
	}
	return s.models.WebhookLogs.UpdateStatus(ctx, webhook.ID, "processed", "")
}
