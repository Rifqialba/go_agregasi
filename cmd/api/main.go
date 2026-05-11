package main

import (
	"log"
	"net/http"
	"time"

	"aggregation-dashboard/internal/collector"
	"aggregation-dashboard/internal/config"
	"aggregation-dashboard/internal/database"
	"aggregation-dashboard/internal/handler"
	"aggregation-dashboard/internal/repository"
	"aggregation-dashboard/internal/service"
	"aggregation-dashboard/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize database connection
	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize repositories
	rawDataRepo := repository.NewRawDataRepository(db)

	// Initialize services
	webhookService := service.NewWebhookService(rawDataRepo)

	uploadService := service.NewUploadService(rawDataRepo)

	restClient := collector.NewRestClient()

	_ = service.NewPollingService(
		restClient,
		rawDataRepo,
	)

	// Initialize handlers
	webhookHandler := handler.NewWebhookHandler(
		webhookService,
		cfg.WebhookSecret,
	)

	uploadHandler := handler.NewUploadHandler(
		uploadService,
	)

	// Initialize router
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		utils.JSON(w, http.StatusOK, map[string]string{
	"status": "ok",
})
	})

	// Webhook endpoint
	r.Post(
		"/webhooks/{source_id}",
		webhookHandler.HandleWebhook,
	)

	// File upload endpoint
	r.Post(
		"/upload",
		uploadHandler.HandleUpload,
	)

	// HTTP server configuration
	server := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("server running on port %s", cfg.AppPort)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}