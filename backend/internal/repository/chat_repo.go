package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/paypilot/backend/internal/models"
)

type ChatRepository struct {
	db *sqlx.DB
}

func NewChatRepository(db *sqlx.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) Create(role, content string) (*models.ChatMessage, error) {
	msg := &models.ChatMessage{
		ID:        uuid.New().String(),
		Role:      role,
		Content:   content,
		CreatedAt: time.Now(),
	}

	query := `INSERT INTO chat_messages (id, role, content, created_at) VALUES (?, ?, ?, ?)`
	_, err := r.db.Exec(query, msg.ID, msg.Role, msg.Content, msg.CreatedAt)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (r *ChatRepository) GetHistory(limit int) ([]models.ChatMessage, error) {
	if limit <= 0 {
		limit = 50
	}

	var messages []models.ChatMessage
	err := r.db.Select(&messages,
		"SELECT * FROM chat_messages ORDER BY created_at DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (r *ChatRepository) Clear() error {
	_, err := r.db.Exec("DELETE FROM chat_messages")
	return err
}
