package repository

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/paypilot/backend/internal/models"
)

type ActionRepository struct {
	db *sqlx.DB
}

func NewActionRepository(db *sqlx.DB) *ActionRepository {
	return &ActionRepository{db: db}
}

func (r *ActionRepository) Create(req models.CreateActionRequest) (*models.PendingAction, error) {
	paramsJSON, err := json.Marshal(req.Params)
	if err != nil {
		return nil, err
	}

	action := &models.PendingAction{
		ID:          uuid.New().String(),
		Type:        req.Type,
		Description: req.Description,
		Params:      paramsJSON,
		OrderID:     req.OrderID,
		Status:      string(models.ActionStatusPending),
		CreatedAt:   time.Now(),
	}

	query := `INSERT INTO pending_actions (id, type, description, params, order_id, status, created_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err = r.db.Exec(query,
		action.ID,
		action.Type,
		action.Description,
		string(action.Params),
		action.OrderID,
		action.Status,
		action.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return action, nil
}

func (r *ActionRepository) GetByID(id string) (*models.PendingAction, error) {
	var action models.PendingAction
	err := r.db.Get(&action, "SELECT * FROM pending_actions WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &action, nil
}

func (r *ActionRepository) GetPending() ([]models.PendingAction, error) {
	var actions []models.PendingAction
	err := r.db.Select(&actions, "SELECT * FROM pending_actions WHERE status = 'pending' ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	return actions, nil
}

func (r *ActionRepository) UpdateStatus(id string, status models.ActionStatus) error {
	now := time.Now()
	query := `UPDATE pending_actions SET status = ?, resolved_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, string(status), now, id)
	return err
}

func (r *ActionRepository) Approve(id string) error {
	return r.UpdateStatus(id, models.ActionStatusApproved)
}

func (r *ActionRepository) Reject(id string) error {
	return r.UpdateStatus(id, models.ActionStatusRejected)
}

func (r *ActionRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM pending_actions WHERE id = ?", id)
	return err
}
