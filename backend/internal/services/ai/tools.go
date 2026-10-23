package ai

type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"input_schema"`
}

var AvailableTools = []Tool{
	{
		Name:        "list_orders",
		Description: "Liste les commandes avec filtres optionnels. Utilise cet outil pour afficher les commandes de l'utilisateur.",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"status": map[string]interface{}{
					"type":        "string",
					"enum":        []string{"pending", "paid", "cancelled"},
					"description": "Filtrer par statut de commande",
				},
				"period": map[string]interface{}{
					"type":        "string",
					"enum":        []string{"today", "this_week", "this_month"},
					"description": "Filtrer par période",
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"default":    10,
					"description": "Nombre maximum de commandes à retourner",
				},
			},
		},
	},
	{
		Name:        "get_order_details",
		Description: "Récupère les détails complets d'une commande spécifique par son ID.",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"order_id": map[string]interface{}{
					"type":        "string",
					"description": "L'identifiant unique de la commande",
				},
			},
			"required": []string{"order_id"},
		},
	},
	{
		Name:        "create_payment_link",
		Description: "Crée un lien de paiement PayPal pour une commande. IMPORTANT: Cette action nécessite une confirmation de l'utilisateur avant exécution.",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"order_id": map[string]interface{}{
					"type":        "string",
					"description": "L'identifiant de la commande pour laquelle créer le lien",
				},
			},
			"required": []string{"order_id"},
		},
	},
	{
		Name:        "get_revenue_summary",
		Description: "Calcule un résumé des revenus sur une période donnée: total, nombre de commandes payées, en attente, etc.",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"period": map[string]interface{}{
					"type":        "string",
					"enum":        []string{"today", "this_week", "this_month", "last_month"},
					"default":     "this_month",
					"description": "La période pour le calcul des revenus",
				},
			},
		},
	},
	{
		Name:        "draft_reminder_message",
		Description: "Prépare un message de relance personnalisé pour un client avec une commande en attente.",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"order_id": map[string]interface{}{
					"type":        "string",
					"description": "L'identifiant de la commande concernée",
				},
				"tone": map[string]interface{}{
					"type":        "string",
					"enum":        []string{"friendly", "professional", "urgent"},
					"default":     "friendly",
					"description": "Le ton du message de relance",
				},
			},
			"required": []string{"order_id"},
		},
	},
}

// Tools that require human confirmation before execution
var ToolsRequiringConfirmation = map[string]bool{
	"create_payment_link": true,
	"send_reminder":       true,
}

func RequiresConfirmation(toolName string) bool {
	return ToolsRequiringConfirmation[toolName]
}
