package models

import (
	"time"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID            string     `json:"id" db:"id"`
	ClientName    string     `json:"client_name" db:"client_name"`
	ClientEmail   string     `json:"client_email,omitempty" db:"client_email"`
	Description   string     `json:"description,omitempty" db:"description"`
	Amount        float64    `json:"amount" db:"amount"`
	Currency      string     `json:"currency" db:"currency"`
	Status        string     `json:"status" db:"status"`
	PayPalOrderID string     `json:"paypal_order_id,omitempty" db:"paypal_order_id"`
	PaymentLink   string     `json:"payment_link,omitempty" db:"payment_link"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	PaidAt        *time.Time `json:"paid_at,omitempty" db:"paid_at"`
}

type CreateOrderRequest struct {
	ClientName  string  `json:"client_name"`
	ClientEmail string  `json:"client_email"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
}

type UpdateOrderRequest struct {
	ClientName  *string `json:"client_name,omitempty"`
	ClientEmail *string `json:"client_email,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
}

type OrderFilter struct {
	Status *string
	Period *string
	Limit  int
	Offset int
}
