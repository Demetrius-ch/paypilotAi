package models

import (
	"encoding/json"
	"time"
)

type ChatRole string

const (
	ChatRoleUser      ChatRole = "user"
	ChatRoleAssistant ChatRole = "assistant"
)

type ChatMessage struct {
	ID        string          `json:"id" db:"id"`
	Role      string          `json:"role" db:"role"`
	Content   string          `json:"content" db:"content"`
	ToolCalls json.RawMessage `json:"tool_calls,omitempty" db:"tool_calls"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
}

type ChatRequest struct {
	Message string `json:"message"`
}

type ChatResponse struct {
	Message       string         `json:"message"`
	PendingAction *PendingAction `json:"pending_action,omitempty"`
}
