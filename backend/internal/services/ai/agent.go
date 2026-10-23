package ai

import (
	"encoding/json"
	"fmt"

	"github.com/paypilot/backend/internal/models"
	"github.com/paypilot/backend/internal/repository"
)

type Agent struct {
	client     *Client
	orderRepo  *repository.OrderRepository
	actionRepo *repository.ActionRepository
}

func NewAgent(apiKey string, orderRepo *repository.OrderRepository, actionRepo *repository.ActionRepository) *Agent {
	return &Agent{
		client:     NewClient(apiKey),
		orderRepo:  orderRepo,
		actionRepo: actionRepo,
	}
}

type AgentResponse struct {
	Message       string                `json:"message"`
	PendingAction *models.PendingAction `json:"pending_action,omitempty"`
}

func (a *Agent) ProcessMessage(userMessage string, history []models.ChatMessage) (*AgentResponse, error) {
	messages := make([]Message, 0, len(history)+1)

	for _, msg := range history {
		messages = append(messages, Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	messages = append(messages, Message{
		Role:    "user",
		Content: userMessage,
	})

	resp, err := a.client.SendMessage(messages)
	if err != nil {
		return nil, err
	}

	if resp.HasToolUse() {
		toolUse := resp.GetToolUse()
		if toolUse != nil {
			return a.handleToolUse(toolUse, messages)
		}
	}

	return &AgentResponse{
		Message: resp.GetTextContent(),
	}, nil
}

func (a *Agent) handleToolUse(toolUse *ContentBlock, messages []Message) (*AgentResponse, error) {
	if RequiresConfirmation(toolUse.Name) {
		return a.createPendingAction(toolUse)
	}

	result, err := a.executeTool(toolUse.Name, toolUse.Input)
	if err != nil {
		return &AgentResponse{
			Message: fmt.Sprintf("Erreur lors de l'exécution de l'outil: %v", err),
		}, nil
	}

	toolResultMessages := append(messages, Message{
		Role:    "assistant",
		Content: fmt.Sprintf("[Outil utilisé: %s]", toolUse.Name),
	}, Message{
		Role:    "user",
		Content: fmt.Sprintf("[Résultat de l'outil %s]: %s", toolUse.Name, result),
	})

	resp, err := a.client.SendMessage(toolResultMessages)
	if err != nil {
		return nil, err
	}

	return &AgentResponse{
		Message: resp.GetTextContent(),
	}, nil
}

func (a *Agent) createPendingAction(toolUse *ContentBlock) (*AgentResponse, error) {
	var description string
	var orderID string

	switch toolUse.Name {
	case "create_payment_link":
		if id, ok := toolUse.Input["order_id"].(string); ok {
			orderID = id
			order, err := a.orderRepo.GetByID(orderID)
			if err != nil {
				return &AgentResponse{
					Message: fmt.Sprintf("Commande introuvable: %s", orderID),
				}, nil
			}
			description = fmt.Sprintf("Créer un lien de paiement de %.2f %s pour %s", order.Amount, order.Currency, order.ClientName)
		}
	case "send_reminder":
		if id, ok := toolUse.Input["order_id"].(string); ok {
			orderID = id
			order, err := a.orderRepo.GetByID(orderID)
			if err != nil {
				return &AgentResponse{
					Message: fmt.Sprintf("Commande introuvable: %s", orderID),
				}, nil
			}
			description = fmt.Sprintf("Envoyer une relance à %s pour la commande de %.2f %s", order.ClientName, order.Amount, order.Currency)
		}
	}

	action, err := a.actionRepo.Create(models.CreateActionRequest{
		Type:        toolUse.Name,
		Description: description,
		Params:      toolUse.Input,
		OrderID:     orderID,
	})

	if err != nil {
		return nil, err
	}

	return &AgentResponse{
		Message:       fmt.Sprintf("⚠️ Action en attente de confirmation:\n\n**%s**\n\nVeuillez confirmer ou annuler cette action.", description),
		PendingAction: action,
	}, nil
}

func (a *Agent) executeTool(name string, input map[string]interface{}) (string, error) {
	switch name {
	case "list_orders":
		return a.executeListOrders(input)
	case "get_order_details":
		return a.executeGetOrderDetails(input)
	case "get_revenue_summary":
		return a.executeGetRevenueSummary(input)
	case "draft_reminder_message":
		return a.executeDraftReminderMessage(input)
	default:
		return "", fmt.Errorf("outil inconnu: %s", name)
	}
}

func (a *Agent) executeListOrders(input map[string]interface{}) (string, error) {
	filter := models.OrderFilter{Limit: 10}

	if status, ok := input["status"].(string); ok {
		filter.Status = &status
	}
	if period, ok := input["period"].(string); ok {
		filter.Period = &period
	}
	if limit, ok := input["limit"].(float64); ok {
		filter.Limit = int(limit)
	}

	orders, err := a.orderRepo.List(filter)
	if err != nil {
		return "", err
	}

	if len(orders) == 0 {
		return "Aucune commande trouvée avec ces critères.", nil
	}

	result, _ := json.MarshalIndent(orders, "", "  ")
	return string(result), nil
}

func (a *Agent) executeGetOrderDetails(input map[string]interface{}) (string, error) {
	orderID, ok := input["order_id"].(string)
	if !ok {
		return "", fmt.Errorf("order_id manquant")
	}

	order, err := a.orderRepo.GetByID(orderID)
	if err != nil {
		return "", err
	}

	result, _ := json.MarshalIndent(order, "", "  ")
	return string(result), nil
}

func (a *Agent) executeGetRevenueSummary(input map[string]interface{}) (string, error) {
	period := "this_month"
	if p, ok := input["period"].(string); ok {
		period = p
	}

	summary, err := a.orderRepo.GetSummary(period)
	if err != nil {
		return "", err
	}

	result, _ := json.MarshalIndent(summary, "", "  ")
	return string(result), nil
}

func (a *Agent) executeDraftReminderMessage(input map[string]interface{}) (string, error) {
	orderID, ok := input["order_id"].(string)
	if !ok {
		return "", fmt.Errorf("order_id manquant")
	}

	tone := "friendly"
	if t, ok := input["tone"].(string); ok {
		tone = t
	}

	order, err := a.orderRepo.GetByID(orderID)
	if err != nil {
		return "", err
	}

	var message string
	switch tone {
	case "friendly":
		message = fmt.Sprintf(`Bonjour %s,

J'espère que vous allez bien ! Je me permets de vous contacter concernant la commande "%s" d'un montant de %.2f %s.

N'hésitez pas à me faire signe si vous avez des questions.

Bien cordialement`, order.ClientName, order.Description, order.Amount, order.Currency)
	case "professional":
		message = fmt.Sprintf(`Madame, Monsieur %s,

Je vous contacte au sujet de la facture en attente pour "%s" (%.2f %s).

Je reste à votre disposition pour tout renseignement complémentaire.

Cordialement`, order.ClientName, order.Description, order.Amount, order.Currency)
	case "urgent":
		message = fmt.Sprintf(`Bonjour %s,

Je me permets de vous relancer concernant le règlement de la commande "%s" (%.2f %s) qui est en attente depuis plusieurs jours.

Merci de procéder au paiement dans les meilleurs délais.

Cordialement`, order.ClientName, order.Description, order.Amount, order.Currency)
	}

	return message, nil
}
