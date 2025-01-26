package services

import (
	"context"
	"main/internal/domain/models"

	"github.com/google/uuid"
)

type IRepository interface {
	Create(ctx context.Context, request models.Request) error
	GetByID(ctx context.Context, id uuid.UUID) (models.Request, error)
	List(ctx context.Context, page, pageSize int) ([]models.Request, int, error)
}

type RequestService struct {
	Repo IRepository
}

func NewRequestService(repo IRepository) *RequestService {
	return &RequestService{Repo: repo}
}

func (s *RequestService) PostRequest(ctx context.Context, request models.Request) error {
	return s.Repo.Create(ctx, request)
}

func (s *RequestService) GetRequest(ctx context.Context, id uuid.UUID) (models.Request, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *RequestService) ListOfRequests(ctx context.Context, page, pageSize int) ([]models.Request, int, error) {
	return s.Repo.List(ctx, page, pageSize)
}
