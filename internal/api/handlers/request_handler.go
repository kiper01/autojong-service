package handler

import (
	"context"
	"encoding/json"
	"main/internal/api/dto"
	"main/internal/domain/models"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type IUsecase interface {
	PostRequest(ctx context.Context, name, phone, email, carInfo string) error
	GetRequest(ctx context.Context, id uuid.UUID) (models.Request, error)
	ListOfRequests(ctx context.Context, page, pageSize int) ([]models.Request, int, int, int, error)
}

type RequestHandler struct {
	uc IUsecase
}

func NewRequestHandler(usecase IUsecase) *RequestHandler {
	return &RequestHandler{
		uc: usecase,
	}
}

// @Summary Добавить заявку
// @Description Добавить заявку
// @Tags Requests
// @Accept json
// @Produce json
// @Param page body dto.PostRequest true "Заявка"
// @Success 201 {object} string "Successful response"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/request/post [post]
func (h *RequestHandler) PostRequest(w http.ResponseWriter, r *http.Request) {
	var req dto.PostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.uc.PostRequest(r.Context(), req.Name, req.Phone, req.Email, req.CarInfo)
	if err != nil {
		http.Error(w, err.Error(), errorCode(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Request created successfully",
	})
}

// @Summary Получить заявку
// @Description Получить заявку по id
// @Tags Requests
// @Accept json
// @Produce json
// @Param id path string true "ID Request"
// @Success 200 {object} dto.Request
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/request/get/ [get]
// @Security ApiKeyAuth
func (h *RequestHandler) GetRequest(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	request, err := h.uc.GetRequest(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), errorCode(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.RequestToDTO(request))
}

// @Summary Получить список всех заявок
// @Description Получить список заявок с пагинацией
// @Tags Requests
// @Accept json
// @Produce json
// @Param page query int true "Page number"
// @Param page_size query int true "Page size"
// @Success 200 {object} dto.ListResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/request/list [get]
// @Security ApiKeyAuth
func (h *RequestHandler) ListOfRequests(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil || pageSize < 1 {
		http.Error(w, "Invalid page size", http.StatusBadRequest)
		return
	}

	safs, page, pageSize, totalPages, err := h.uc.ListOfRequests(r.Context(), page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), errorCode(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.ListToDTOlist(safs, page, pageSize, totalPages))
}
