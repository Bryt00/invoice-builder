package services

import (
	"context"

	"github.com/google/uuid"
	"raven.go.invoice-builder/internal/models"
)

type CreditService interface {
	GetBalance(ctx context.Context, userID uuid.UUID) (int, error)
	GetCreditStats(ctx context.Context, userID uuid.UUID) (balance, totalPurchased, totalUsed int, err error)
	GetCreditHistory(ctx context.Context, userID uuid.UUID) ([]*models.CreditTxn, error)
	AdminGrantCredits(ctx context.Context, userID uuid.UUID, amount int, reason string) error
	GetAvailablePackages(ctx context.Context) ([]*models.CreditPackage, error)
}

type creditService struct {
	models models.Models
}

func NewCreditService(models models.Models) CreditService {
	return &creditService{models: models}
}

func (s *creditService) GetBalance(ctx context.Context, userID uuid.UUID) (int, error) {
	return s.models.CreditTxn.GetBalanceByUserID(ctx, userID)
}

func (s *creditService) GetCreditStats(ctx context.Context, userID uuid.UUID) (balance, totalPurchased, totalUsed int, err error) {
	return s.models.CreditTxn.GetCreditStatsByUserID(ctx, userID)
}

func (s *creditService) GetCreditHistory(ctx context.Context, userID uuid.UUID) ([]*models.CreditTxn, error) {
	return s.models.CreditTxn.GetAllByUserID(ctx, userID)
}

func (s *creditService) AdminGrantCredits(ctx context.Context, userID uuid.UUID, amount int, reason string) error {
	return s.models.CreditTxn.AdminGrantCredits(ctx, userID, amount, reason)
}

func (s *creditService) GetAvailablePackages(ctx context.Context) ([]*models.CreditPackage, error) {
	return s.models.CreditPackages.GetAll(ctx, false)
}
