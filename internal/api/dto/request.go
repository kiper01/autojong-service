package dto

import (
	"main/internal/domain/models"
	"time"

	"github.com/google/uuid"
)

type PostRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email,omitempty"`
	CarInfo string `json:"car_info,omitempty"`
}

type Request struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Phone   string    `json:"phone"`
	Email   string    `json:"email,omitempty"`
	CarInfo string    `json:"car_info,omitempty"`
	Date    time.Time `json:"date"`
}

type ListResponse struct {
	Requests   []Request `json:"requests" `
	Page       int       `json:"page"`
	PageSize   int       `json:"page_size"`
	TotalPages int       `json:"total_pages" `
}

func ListToDTOlist(requests []models.Request, page, pageSize, totalPages int) ListResponse {
	dto := ListResponse{
		Requests:   make([]Request, len(requests)),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	for i, request := range requests {
		dto.Requests[i] = RequestToDTO(request)
	}

	return dto
}

func RequestToDTO(request models.Request) Request {
	dto := Request{
		ID:    request.ID,
		Name:  request.Name,
		Phone: request.Phone,
		Email: request.Email,
		Date:  request.Date,
	}

	return dto
}
