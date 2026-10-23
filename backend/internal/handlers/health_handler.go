package handlers

import (
	"net/http"

	"github.com/paypilot/backend/pkg/response"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{
		"status":  "ok",
		"service": "paypilot-api",
	})
}
