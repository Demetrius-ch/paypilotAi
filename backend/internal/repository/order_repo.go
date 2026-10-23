package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/paypilot/backend/internal/models"
)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(req models.CreateOrderRequest) (*models.Order, error) {
	order := &models.Order{
		ID:          uuid.New().String(),
		ClientName:  req.ClientName,
		ClientEmail: req.ClientEmail,
		Description: req.Description,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Status:      string(models.OrderStatusPending),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if order.Currency == "" {
		order.Currency = "USD"
	}

	query := `INSERT INTO orders (id, client_name, client_email, description, amount, currency, status, created_at, updated_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		order.ID,
		order.ClientName,
		order.ClientEmail,
		order.Description,
		order.Amount,
		order.Currency,
		order.Status,
		order.CreatedAt,
		order.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *OrderRepository) GetByID(id string) (*models.Order, error) {
	var order models.Order
	err := r.db.Get(&order, "SELECT * FROM orders WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) List(filter models.OrderFilter) ([]models.Order, error) {
	query := "SELECT * FROM orders WHERE 1=1"
	args := []interface{}{}

	if filter.Status != nil && *filter.Status != "" {
		query += " AND status = ?"
		args = append(args, *filter.Status)
	}

	if filter.Period != nil && *filter.Period != "" {
		switch *filter.Period {
		case "today":
			query += " AND date(created_at) = date('now')"
		case "this_week":
			query += " AND created_at >= date('now', '-7 days')"
		case "this_month":
			query += " AND created_at >= date('now', 'start of month')"
		case "last_month":
			query += " AND created_at >= date('now', 'start of month', '-1 month') AND created_at < date('now', 'start of month')"
		}
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, filter.Limit)
	}

	if filter.Offset > 0 {
		query += " OFFSET ?"
		args = append(args, filter.Offset)
	}

	var orders []models.Order
	err := r.db.Select(&orders, query, args...)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepository) Update(id string, req models.UpdateOrderRequest) (*models.Order, error) {
	order, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.ClientName != nil {
		order.ClientName = *req.ClientName
	}
	if req.ClientEmail != nil {
		order.ClientEmail = *req.ClientEmail
	}
	if req.Description != nil {
		order.Description = *req.Description
	}
	if req.Status != nil {
		order.Status = *req.Status
		if *req.Status == string(models.OrderStatusPaid) {
			now := time.Now()
			order.PaidAt = &now
		}
	}
	order.UpdatedAt = time.Now()

	query := `UPDATE orders SET client_name = ?, client_email = ?, description = ?, status = ?, updated_at = ?, paid_at = ? WHERE id = ?`
	_, err = r.db.Exec(query, order.ClientName, order.ClientEmail, order.Description, order.Status, order.UpdatedAt, order.PaidAt, id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *OrderRepository) UpdatePayPalOrder(id, paypalOrderID, paymentLink string) error {
	query := `UPDATE orders SET paypal_order_id = ?, payment_link = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, paypalOrderID, paymentLink, time.Now(), id)
	return err
}

func (r *OrderRepository) MarkAsPaid(id string) error {
	now := time.Now()
	query := `UPDATE orders SET status = ?, paid_at = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, string(models.OrderStatusPaid), now, now, id)
	return err
}

func (r *OrderRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM orders WHERE id = ?", id)
	return err
}

func (r *OrderRepository) GetOverdue(days int) ([]models.Order, error) {
	query := `SELECT * FROM orders
			  WHERE status = 'pending'
			  AND created_at <= date('now', ? || ' days')
			  ORDER BY created_at ASC`

	var orders []models.Order
	err := r.db.Select(&orders, query, -days)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetSummary(period string) (*models.Summary, error) {
	var dateFilter string
	var prevDateFilter string

	switch period {
	case "today":
		dateFilter = "date(created_at) = date('now')"
		prevDateFilter = "date(created_at) = date('now', '-1 day')"
	case "this_week":
		dateFilter = "created_at >= date('now', '-7 days')"
		prevDateFilter = "created_at >= date('now', '-14 days') AND created_at < date('now', '-7 days')"
	case "this_month":
		dateFilter = "created_at >= date('now', 'start of month')"
		prevDateFilter = "created_at >= date('now', 'start of month', '-1 month') AND created_at < date('now', 'start of month')"
	default:
		dateFilter = "created_at >= date('now', 'start of month')"
		prevDateFilter = "created_at >= date('now', 'start of month', '-1 month') AND created_at < date('now', 'start of month')"
	}

	summary := &models.Summary{}

	// Current period stats
	query := `SELECT
		COALESCE(SUM(CASE WHEN status = 'paid' THEN amount ELSE 0 END), 0) as total_revenue,
		COUNT(*) as order_count,
		SUM(CASE WHEN status = 'paid' THEN 1 ELSE 0 END) as paid_count,
		SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END) as pending_count,
		COALESCE(SUM(CASE WHEN status = 'pending' THEN amount ELSE 0 END), 0) as pending_amount
		FROM orders WHERE ` + dateFilter

	err := r.db.QueryRow(query).Scan(
		&summary.TotalRevenue,
		&summary.OrderCount,
		&summary.PaidCount,
		&summary.PendingCount,
		&summary.PendingAmount,
	)
	if err != nil {
		return nil, err
	}

	// Previous period revenue
	prevQuery := `SELECT COALESCE(SUM(CASE WHEN status = 'paid' THEN amount ELSE 0 END), 0) FROM orders WHERE ` + prevDateFilter
	err = r.db.QueryRow(prevQuery).Scan(&summary.PreviousPeriod)
	if err != nil {
		return nil, err
	}

	// Calculate percentage change
	if summary.PreviousPeriod > 0 {
		summary.PercentageChange = ((summary.TotalRevenue - summary.PreviousPeriod) / summary.PreviousPeriod) * 100
	}

	return summary, nil
}
