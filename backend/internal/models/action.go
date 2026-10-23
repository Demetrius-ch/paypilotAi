package models

import (
	"encoding/json"
	"time"
)

type ActionType string

const (
	ActionTypeCreatePaymentLink ActionType = "create_payment_link"
	ActionTypeSendReminder      ActionType = "send_reminder"
)

type ActionStatus string

const (
	ActionStatusPending  ActionStatus = "pending"
	ActionStatusApproved ActionStatus = "approved"
	ActionStatusRejected ActionStatus = "rejected"
)

type PendingAction struct {
	ID          string          `json:"id" db:"id"`
	Type        string          `json:"type" db:"type"`
	Description string          `json:"description" db:"description"`
	Params      json.RawMessage `json:"params" db:"params"`
	OrderID     string          `json:"order_id,omitempty" db:"order_id"`
	Status      string          `json:"status" db:"status"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	ResolvedAt  *time.Time      `json:"resolved_at,omitempty" db:"resolved_at"`
}

type CreateActionRequest struct {
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Params      interface{} `json:"params"`
	OrderID     string      `json:"order_id,omitempty"`
}

type ActionResult struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
