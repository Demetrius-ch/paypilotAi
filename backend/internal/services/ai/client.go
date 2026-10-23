package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		baseURL:    "https://api.anthropic.com",
		model:      "claude-sonnet-4-20250514",
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ToolUse struct {
	Type  string                 `json:"type"`
	ID    string                 `json:"id"`
	Name  string                 `json:"name"`
	Input map[string]interface{} `json:"input"`
}

type ContentBlock struct {
	Type  string                 `json:"type"`
	Text  string                 `json:"text,omitempty"`
	ID    string                 `json:"id,omitempty"`
	Name  string                 `json:"name,omitempty"`
	Input map[string]interface{} `json:"input,omitempty"`
}

type Request struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	System    string    `json:"system,omitempty"`
	Tools     []Tool    `json:"tools,omitempty"`
	Messages  []Message `json:"messages"`
}

type Response struct {
	ID           string         `json:"id"`
	Type         string         `json:"type"`
	Role         string         `json:"role"`
	Content      []ContentBlock `json:"content"`
	Model        string         `json:"model"`
	StopReason   string         `json:"stop_reason"`
	StopSequence string         `json:"stop_sequence,omitempty"`
}

const SystemPrompt = `Tu es PayPilot, un assistant intelligent pour aider les freelances et petites entreprises à gérer leurs paiements et commandes.

Tu es capable de:
- Lister et filtrer les commandes
- Afficher les détails d'une commande
- Créer des liens de paiement PayPal (avec confirmation de l'utilisateur)
- Calculer des résumés de revenus
- Préparer des messages de relance pour les clients

Réponds toujours en français de manière concise et professionnelle.
Quand tu utilises un outil, explique brièvement ce que tu fais.
Si une action nécessite une confirmation (comme créer un lien de paiement), indique clairement que tu attends l'approbation de l'utilisateur.

Format les montants avec le symbole de devise approprié (€ pour EUR, $ pour USD).
Pour les dates, utilise le format français (jour/mois/année).`

func (c *Client) SendMessage(messages []Message) (*Response, error) {
	req := Request{
		Model:     c.model,
		MaxTokens: 1024,
		System:    SystemPrompt,
		Tools:     AvailableTools,
		Messages:  messages,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/v1/messages", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(respBody))
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (r *Response) GetTextContent() string {
	for _, block := range r.Content {
		if block.Type == "text" {
			return block.Text
		}
	}
	return ""
}

func (r *Response) GetToolUse() *ContentBlock {
	for _, block := range r.Content {
		if block.Type == "tool_use" {
			return &block
		}
	}
	return nil
}

func (r *Response) HasToolUse() bool {
	return r.StopReason == "tool_use"
}
