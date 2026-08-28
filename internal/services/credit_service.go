package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"raven.go.invoice-builder/internal/models"
)

type CreditService interface {
	GetBalance(ctx context.Context, userID uuid.UUID) (int, error)
	GetCreditStats(ctx context.Context, userID uuid.UUID) (balance, totalPurchased, totalUsed int, err error)
	GetCreditHistory(ctx context.Context, userID uuid.UUID) ([]*models.CreditTxn, error)
	AdminGrantCredits(ctx context.Context, userID uuid.UUID, amount int, reason string) error
	GetAvailablePackages(ctx context.Context) ([]*models.CreditPackage, error)
	DeductCredits(ctx context.Context, userID uuid.UUID, amount int, reason string) error
	FulfillTopup(ctx context.Context, userID uuid.UUID, packageID uuid.UUID, reference string) error
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

func (s *creditService) DeductCredits(ctx context.Context, userID uuid.UUID, amount int, reason string) error {
	balance, err := s.GetBalance(ctx, userID)
	if err != nil {
		return err
	}
	if balance < amount {
		return errors.New("insufficient credits to perform this action")
	}

	txn := &models.CreditTxn{
		UserID:      userID,
		Amount:      -amount,
		Type:        models.CreditTxnUsage,
		Description: reason,
	}

	return s.models.CreditTxn.Insert(ctx, txn)
}

func (s *creditService) FulfillTopup(ctx context.Context, userID uuid.UUID, packageID uuid.UUID, reference string) error {
	exists, err := s.models.CreditTxn.CheckReferenceExists(ctx, reference)
	if err != nil {
		return err
	}
	if exists {
		return models.ErrTransactionAlreadyProcessed
	}

	pkg, err := s.models.CreditPackages.GetByID(ctx, packageID)
	if err != nil {
		return err
	}

	txn := &models.CreditTxn{
		UserID:      userID,
		Amount:      pkg.CreditsGranted,
		Type:        models.CreditTxnPurchase,
		Description: "Topup Ref: " + reference,
		CreatedAt:   time.Now(),
	}

	return s.models.CreditTxn.Insert(ctx, txn)
}
