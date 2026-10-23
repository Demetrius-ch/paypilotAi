package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/paypilot/backend/internal/models"
	"github.com/paypilot/backend/internal/repository"
	"github.com/paypilot/backend/internal/services/ai"
	"github.com/paypilot/backend/pkg/response"
)

type ChatHandler struct {
	agent    *ai.Agent
	chatRepo *repository.ChatRepository
}

func NewChatHandler(agent *ai.Agent, chatRepo *repository.ChatRepository) *ChatHandler {
	return &ChatHandler{
		agent:    agent,
		chatRepo: chatRepo,
	}
}

func (h *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var req models.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	if req.Message == "" {
		response.BadRequest(w, "Message is required")
		return
	}

	// Save user message
	_, err := h.chatRepo.Create(string(models.ChatRoleUser), req.Message)
	if err != nil {
		response.InternalError(w, "Failed to save message")
		return
	}

	// Get history for context
	history, err := h.chatRepo.GetHistory(20)
	if err != nil {
		history = []models.ChatMessage{}
	}

	// Process with AI agent
	agentResp, err := h.agent.ProcessMessage(req.Message, history)
	if err != nil {
		response.InternalError(w, "Failed to process message: "+err.Error())
		return
	}

	// Save assistant response
	_, err = h.chatRepo.Create(string(models.ChatRoleAssistant), agentResp.Message)
	if err != nil {
		response.InternalError(w, "Failed to save response")
		return
	}

	response.Success(w, models.ChatResponse{
		Message:       agentResp.Message,
		PendingAction: agentResp.PendingAction,
	})
}

func (h *ChatHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	messages, err := h.chatRepo.GetHistory(50)
	if err != nil {
		response.InternalError(w, "Failed to fetch history")
		return
	}

	if messages == nil {
		messages = []models.ChatMessage{}
	}

	response.Success(w, messages)
}

func (h *ChatHandler) ClearHistory(w http.ResponseWriter, r *http.Request) {
	if err := h.chatRepo.Clear(); err != nil {
		response.InternalError(w, "Failed to clear history")
		return
	}

	response.NoContent(w)
}
