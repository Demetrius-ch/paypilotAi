package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/paypilot/backend/internal/models"
	"github.com/paypilot/backend/internal/repository"
	"github.com/paypilot/backend/internal/services/paypal"
	"github.com/paypilot/backend/pkg/response"
)

type ActionHandler struct {
	actionRepo *repository.ActionRepository
	orderRepo  *repository.OrderRepository
	paypal     *paypal.Client
	frontendURL string
}

func NewActionHandler(
	actionRepo *repository.ActionRepository,
	orderRepo *repository.OrderRepository,
	paypalClient *paypal.Client,
	frontendURL string,
) *ActionHandler {
	return &ActionHandler{
		actionRepo: actionRepo,
		orderRepo:  orderRepo,
		paypal:     paypalClient,
		frontendURL: frontendURL,
	}
}

func (h *ActionHandler) GetPending(w http.ResponseWriter, r *http.Request) {
	actions, err := h.actionRepo.GetPending()
	if err != nil {
		response.InternalError(w, "Failed to fetch pending actions")
		return
	}

	if actions == nil {
		actions = []models.PendingAction{}
	}

	response.Success(w, actions)
}

func (h *ActionHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, "Action ID is required")
		return
	}

	action, err := h.actionRepo.GetByID(id)
	if err != nil {
		response.NotFound(w, "Action not found")
		return
	}

	if action.Status != string(models.ActionStatusPending) {
		response.BadRequest(w, "Action already processed")
		return
	}

	// Execute the action
	result, err := h.executeAction(action)
	if err != nil {
		response.InternalError(w, "Failed to execute action: "+err.Error())
		return
	}

	// Mark as approved
	if err := h.actionRepo.Approve(id); err != nil {
		response.InternalError(w, "Failed to update action status")
		return
	}

	response.Success(w, result)
}

func (h *ActionHandler) Reject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, "Action ID is required")
		return
	}

	action, err := h.actionRepo.GetByID(id)
	if err != nil {
		response.NotFound(w, "Action not found")
		return
	}

	if action.Status != string(models.ActionStatusPending) {
		response.BadRequest(w, "Action already processed")
		return
	}

	if err := h.actionRepo.Reject(id); err != nil {
		response.InternalError(w, "Failed to reject action")
		return
	}

	response.Success(w, models.ActionResult{
		Success: true,
		Message: "Action annulée",
	})
}

func (h *ActionHandler) executeAction(action *models.PendingAction) (*models.ActionResult, error) {
	var params map[string]interface{}
	if err := json.Unmarshal(action.Params, &params); err != nil {
		return nil, err
	}

	switch action.Type {
	case string(models.ActionTypeCreatePaymentLink):
		return h.executeCreatePaymentLink(params)
	default:
		return nil, fmt.Errorf("unknown action type: %s", action.Type)
	}
}

func (h *ActionHandler) executeCreatePaymentLink(params map[string]interface{}) (*models.ActionResult, error) {
	orderID, ok := params["order_id"].(string)
	if !ok {
		return nil, fmt.Errorf("order_id missing")
	}

	order, err := h.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, err
	}

	returnURL := h.frontendURL + "/payment/success"
	cancelURL := h.frontendURL + "/payment/cancel"

	paypalOrder, err := h.paypal.CreateOrder(
		order.ID,
		order.Description,
		order.Currency,
		fmt.Sprintf("%.2f", order.Amount),
		returnURL,
		cancelURL,
	)
	if err != nil {
		return nil, err
	}

	approveURL := paypalOrder.GetApproveURL()

	if err := h.orderRepo.UpdatePayPalOrder(order.ID, paypalOrder.ID, approveURL); err != nil {
		return nil, err
	}

	return &models.ActionResult{
		Success: true,
		Message: "Lien de paiement créé avec succès",
		Data: map[string]string{
			"payment_link":    approveURL,
			"paypal_order_id": paypalOrder.ID,
		},
	}, nil
}
