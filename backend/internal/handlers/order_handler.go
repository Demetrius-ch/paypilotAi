package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/paypilot/backend/internal/models"
	"github.com/paypilot/backend/internal/repository"
	"github.com/paypilot/backend/pkg/response"
)

type OrderHandler struct {
	repo *repository.OrderRepository
}

func NewOrderHandler(repo *repository.OrderRepository) *OrderHandler {
	return &OrderHandler{repo: repo}
}

func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := models.OrderFilter{
		Limit: 50,
	}

	if status := r.URL.Query().Get("status"); status != "" {
		filter.Status = &status
	}
	if period := r.URL.Query().Get("period"); period != "" {
		filter.Period = &period
	}

	orders, err := h.repo.List(filter)
	if err != nil {
		response.InternalError(w, "Failed to fetch orders")
		return
	}

	if orders == nil {
		orders = []models.Order{}
	}

	response.Success(w, orders)
}

func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, "Order ID is required")
		return
	}

	order, err := h.repo.GetByID(id)
	if err != nil {
		response.NotFound(w, "Order not found")
		return
	}

	response.Success(w, order)
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	if req.ClientName == "" {
		response.BadRequest(w, "Client name is required")
		return
	}
	if req.Amount <= 0 {
		response.BadRequest(w, "Amount must be greater than 0")
		return
	}

	order, err := h.repo.Create(req)
	if err != nil {
		response.InternalError(w, "Failed to create order")
		return
	}

	response.Created(w, order)
}

func (h *OrderHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, "Order ID is required")
		return
	}

	var req models.UpdateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	order, err := h.repo.Update(id, req)
	if err != nil {
		response.NotFound(w, "Order not found")
		return
	}

	response.Success(w, order)
}

func (h *OrderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, "Order ID is required")
		return
	}

	if err := h.repo.Delete(id); err != nil {
		response.NotFound(w, "Order not found")
		return
	}

	response.NoContent(w)
}

func (h *OrderHandler) GetOverdue(w http.ResponseWriter, r *http.Request) {
	days := 7

	orders, err := h.repo.GetOverdue(days)
	if err != nil {
		response.InternalError(w, "Failed to fetch overdue orders")
		return
	}

	if orders == nil {
		orders = []models.Order{}
	}

	response.Success(w, orders)
}
