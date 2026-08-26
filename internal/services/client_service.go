package services

import (
	"context"

	"github.com/google/uuid"
	"raven.go.invoice-builder/internal/models"
)

type ClientService interface {
	GetAllByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]*models.Client, int64, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.Client, error)
	CreateClient(ctx context.Context, client *models.Client) error
	UpdateClient(ctx context.Context, client *models.Client) error
	DeleteClient(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type clientService struct {
	models models.Models
}

func NewClientService(models models.Models) ClientService {
	return &clientService{models: models}
}

func (s *clientService) GetAllByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]*models.Client, int64, error) {
	return s.models.Clients.GetAllByUserID(ctx, userID, page, limit)
}

func (s *clientService) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.Client, error) {
	return s.models.Clients.GetByID(ctx, id, userID)
}

func (s *clientService) CreateClient(ctx context.Context, client *models.Client) error {
	if client.ID == uuid.Nil {
		client.ID = uuid.New()
	}
	return s.models.Clients.Insert(ctx, client)
}

func (s *clientService) UpdateClient(ctx context.Context, client *models.Client) error {
	return s.models.Clients.Update(ctx, client)
}

func (s *clientService) DeleteClient(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return s.models.Clients.Delete(ctx, id, userID)
}
