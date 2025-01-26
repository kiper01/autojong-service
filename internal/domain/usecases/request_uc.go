package usecases

import (
	"context"
	"main/internal/domain/models"
	"time"

	"github.com/google/uuid"
)

type IService interface {
	PostRequest(ctx context.Context, request models.Request) error
	GetRequest(ctx context.Context, id uuid.UUID) (models.Request, error)
	ListOfRequests(ctx context.Context, page, pageSize int) ([]models.Request, int, error)
}

type RequestUC struct {
	requestS IService
}

func NewRequestUC(requestS IService) *RequestUC {
	return &RequestUC{requestS: requestS}
}

func (uc *RequestUC) PostRequest(ctx context.Context, name, phone, email, carInfo string) error {
	request := models.Request{
		ID:      uuid.New(),
		Name:    name,
		Phone:   phone,
		Email:   email,
		CarInfo: carInfo,
		Date:    time.Now(),
	}
	return uc.requestS.PostRequest(ctx, request)
}

func (uc *RequestUC) GetRequest(ctx context.Context, id uuid.UUID) (models.Request, error) {
	return uc.requestS.GetRequest(ctx, id)
}

func (uc *RequestUC) ListOfRequests(ctx context.Context, page, pageSize int) ([]models.Request, int, int, int, error) {
	if page == 0 {
		page = 1 // по умолчанию
	}
	if pageSize == 0 {
		pageSize = 20 // по умолчанию
	}

	resp, totalPages, err := uc.requestS.ListOfRequests(ctx, page, pageSize)
	return resp, page, pageSize, totalPages, err
}
