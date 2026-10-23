package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/paypilot/backend/internal/repository"
	"github.com/paypilot/backend/internal/services/paypal"
	"github.com/paypilot/backend/pkg/response"
)

type PayPalHandler struct {
	client      *paypal.Client
	orderRepo   *repository.OrderRepository
	frontendURL string
}

func NewPayPalHandler(client *paypal.Client, orderRepo *repository.OrderRepository, frontendURL string) *PayPalHandler {
	return &PayPalHandler{
		client:      client,
		orderRepo:   orderRepo,
		frontendURL: frontendURL,
	}
}

type CreatePayPalOrderResponse struct {
	PayPalOrderID string `json:"paypal_order_id"`
	ApproveURL    string `json:"approve_url"`
}

func (h *PayPalHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "id")
	if orderID == "" {
		response.BadRequest(w, "Order ID is required")
		return
	}

	order, err := h.orderRepo.GetByID(orderID)
	if err != nil {
		response.NotFound(w, "Order not found")
		return
	}

	if order.PayPalOrderID != "" {
		response.Success(w, CreatePayPalOrderResponse{
			PayPalOrderID: order.PayPalOrderID,
			ApproveURL:    order.PaymentLink,
		})
		return
	}

	returnURL := h.frontendURL + "/payment/success?order_id=" + order.ID
	cancelURL := h.frontendURL + "/payment/cancel?order_id=" + order.ID

	paypalOrder, err := h.client.CreateOrder(
		order.ID,
		order.Description,
		order.Currency,
		fmt.Sprintf("%.2f", order.Amount),
		returnURL,
		cancelURL,
	)
	if err != nil {
		response.InternalError(w, "Failed to create PayPal order: "+err.Error())
		return
	}

	approveURL := paypalOrder.GetApproveURL()

	if err := h.orderRepo.UpdatePayPalOrder(order.ID, paypalOrder.ID, approveURL); err != nil {
		response.InternalError(w, "Failed to update order")
		return
	}

	response.Success(w, CreatePayPalOrderResponse{
		PayPalOrderID: paypalOrder.ID,
		ApproveURL:    approveURL,
	})
}

func (h *PayPalHandler) CaptureOrder(w http.ResponseWriter, r *http.Request) {
	paypalOrderID := chi.URLParam(r, "paypalOrderId")
	if paypalOrderID == "" {
		response.BadRequest(w, "PayPal Order ID is required")
		return
	}

	captureResult, err := h.client.CaptureOrder(paypalOrderID)
	if err != nil {
		response.InternalError(w, "Failed to capture payment: "+err.Error())
		return
	}

	// Find and update our order
	orders, _ := h.orderRepo.List(struct {
		Status *string
		Period *string
		Limit  int
		Offset int
	}{Limit: 1000})

	for _, order := range orders {
		if order.PayPalOrderID == paypalOrderID {
			if captureResult.Status == "COMPLETED" {
				h.orderRepo.MarkAsPaid(order.ID)
			}
			break
		}
	}

	response.Success(w, map[string]interface{}{
		"status":  captureResult.Status,
		"message": "Payment captured successfully",
	})
}

type WebhookEvent struct {
	EventType string          `json:"event_type"`
	Resource  json.RawMessage `json:"resource"`
}

func (h *PayPalHandler) Webhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.BadRequest(w, "Failed to read body")
		return
	}

	var event WebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		response.BadRequest(w, "Invalid webhook payload")
		return
	}

	switch event.EventType {
	case "CHECKOUT.ORDER.APPROVED", "PAYMENT.CAPTURE.COMPLETED":
		var resource struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(event.Resource, &resource); err == nil {
			// Find order by PayPal order ID and mark as paid
			orders, _ := h.orderRepo.List(struct {
				Status *string
				Period *string
				Limit  int
				Offset int
			}{Limit: 1000})

			for _, order := range orders {
				if order.PayPalOrderID == resource.ID {
					h.orderRepo.MarkAsPaid(order.ID)
					break
				}
			}
		}
	}

	response.Success(w, map[string]string{"status": "received"})
}
