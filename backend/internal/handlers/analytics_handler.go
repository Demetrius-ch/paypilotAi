package handlers

import (
	"net/http"

	"github.com/paypilot/backend/internal/repository"
	"github.com/paypilot/backend/pkg/response"
)

type AnalyticsHandler struct {
	orderRepo *repository.OrderRepository
}

func NewAnalyticsHandler(orderRepo *repository.OrderRepository) *AnalyticsHandler {
	return &AnalyticsHandler{orderRepo: orderRepo}
}

func (h *AnalyticsHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	if period == "" {
		period = "this_month"
	}

	summary, err := h.orderRepo.GetSummary(period)
	if err != nil {
		response.InternalError(w, "Failed to calculate summary")
		return
	}

	response.Success(w, summary)
}

func (h *AnalyticsHandler) GetOverdue(w http.ResponseWriter, r *http.Request) {
	days := 7

	orders, err := h.orderRepo.GetOverdue(days)
	if err != nil {
		response.InternalError(w, "Failed to fetch overdue orders")
		return
	}

	response.Success(w, orders)
}
