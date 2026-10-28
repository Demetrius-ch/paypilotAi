package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/paypilot/backend/internal/config"
	"github.com/paypilot/backend/internal/database"
	"github.com/paypilot/backend/internal/handlers"
	"github.com/paypilot/backend/internal/middleware"
	"github.com/paypilot/backend/internal/repository"
	"github.com/paypilot/backend/internal/services/ai"
	"github.com/paypilot/backend/internal/services/paypal"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.New(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	orderRepo := repository.NewOrderRepository(db)
	actionRepo := repository.NewActionRepository(db)
	chatRepo := repository.NewChatRepository(db)

	// Initialize services
	paypalClient := paypal.NewClient(cfg.PayPalClientID, cfg.PayPalSecret, cfg.PayPalBaseURL)
	aiAgent := ai.NewAgent(cfg.AnthropicAPIKey, orderRepo, actionRepo)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	orderHandler := handlers.NewOrderHandler(orderRepo)
	analyticsHandler := handlers.NewAnalyticsHandler(orderRepo)
	chatHandler := handlers.NewChatHandler(aiAgent, chatRepo)
	actionHandler := handlers.NewActionHandler(actionRepo, orderRepo, paypalClient, cfg.FrontendURL)
	paypalHandler := handlers.NewPayPalHandler(paypalClient, orderRepo, cfg.FrontendURL)

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Recovery)
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.FrontendURL, "http://localhost:5173", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Requested-With"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes
	r.Get("/health", healthHandler.Health)

	r.Route("/api", func(r chi.Router) {
		// Orders
		r.Route("/orders", func(r chi.Router) {
			r.Get("/", orderHandler.List)
			r.Post("/", orderHandler.Create)
			r.Get("/overdue", orderHandler.GetOverdue)
			r.Get("/{id}", orderHandler.GetByID)
			r.Patch("/{id}", orderHandler.Update)
			r.Delete("/{id}", orderHandler.Delete)
		})

		// Analytics
		r.Route("/analytics", func(r chi.Router) {
			r.Get("/summary", analyticsHandler.GetSummary)
			r.Get("/overdue", analyticsHandler.GetOverdue)
		})

		// Chat
		r.Route("/chat", func(r chi.Router) {
			r.Post("/", chatHandler.SendMessage)
			r.Get("/history", chatHandler.GetHistory)
			r.Delete("/history", chatHandler.ClearHistory)
		})

		// Actions (pending confirmations)
		r.Route("/actions", func(r chi.Router) {
			r.Get("/pending", actionHandler.GetPending)
			r.Post("/{id}/confirm", actionHandler.Confirm)
			r.Post("/{id}/reject", actionHandler.Reject)
		})

		// PayPal
		r.Route("/paypal", func(r chi.Router) {
			r.Post("/create-order/{id}", paypalHandler.CreateOrder)
			r.Post("/capture/{paypalOrderId}", paypalHandler.CaptureOrder)
			r.Post("/webhook", paypalHandler.Webhook)
		})
	})

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf(" PayPilot API starting on http://localhost%s", addr)
	log.Printf(" Environment: %s", cfg.Env)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
