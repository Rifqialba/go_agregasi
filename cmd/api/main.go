package main

import (
	"log"
	"net/http"
	"time"

	"aggregation-dashboard/internal/audit"
	"aggregation-dashboard/internal/collector"
	"aggregation-dashboard/internal/config"
	"aggregation-dashboard/internal/database"
	"aggregation-dashboard/internal/handler"
	"aggregation-dashboard/internal/pipeline"
	"aggregation-dashboard/internal/repository"
	"aggregation-dashboard/internal/scheduler"
	"aggregation-dashboard/internal/service"
	"aggregation-dashboard/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Load application config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Load scheduler config
	schedulerConfig, err := scheduler.LoadConfig(
		"config.yaml",
	)
	if err != nil {
		log.Fatal(err)
	}

	// Load validation schema
	schema, err := pipeline.LoadValidationSchema(
		"validation_schema.json",
	)
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

	processedRepo := repository.NewProcessedDataRepository(db)

	validationRepo := repository.NewValidationErrorsRepository(db)

	auditRepo := repository.NewAuditLogRepository(db)

	// Initialize audit logger
	auditLogger := audit.NewAuditLogger(
		auditRepo,
	)

	// Initialize pipeline components
	validator := pipeline.NewValidator(schema)

	normalizer := pipeline.NewNormalizer()

	processor := pipeline.NewProcessor(
		rawDataRepo,
		processedRepo,
		validationRepo,
		validator,
		normalizer,
	)

	pipelineRunner := pipeline.NewPipelineRunner(
		processor,
		auditLogger,
	)

	// Initialize collector
	restClient := collector.NewRestClient()

	// Initialize services
	webhookService := service.NewWebhookService(
		rawDataRepo,
	)

	uploadService := service.NewUploadService(
		rawDataRepo,
	)

	pollingService := service.NewPollingService(
		restClient,
		rawDataRepo,
	)

	// Initialize scheduler
	schedulerService := scheduler.NewScheduler(
		pollingService,
	)

	err = schedulerService.RegisterJobs(
		schedulerConfig,
	)
	if err != nil {
		log.Fatal(err)
	}

	schedulerService.Start()

	// Initialize handlers
	webhookHandler := handler.NewWebhookHandler(
		webhookService,
		cfg.WebhookSecret,
	)

	uploadHandler := handler.NewUploadHandler(
		uploadService,
	)

	pipelineHandler := handler.NewPipelineHandler(
		pipelineRunner,
	)

	schedulerHandler := handler.NewSchedulerHandler(
		schedulerService,
	)

	configHandler := handler.NewConfigHandler(
		schedulerService,
		auditLogger,
	)

	auditHandler := handler.NewAuditHandler(
		auditRepo,
	)

	// Initialize router
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Health check endpoint
	r.Get("/health", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		utils.JSON(
			w,
			http.StatusOK,
			map[string]string{
				"status": "ok",
			},
		)
	})

	// Audit log endpoints
	r.Get(
		"/audit-log",
		auditHandler.GetAuditLogs,
	)

	// Webhook endpoints
	r.Post(
		"/webhooks/{source_id}",
		webhookHandler.HandleWebhook,
	)

	// Upload endpoints
	r.Post(
		"/upload",
		uploadHandler.HandleUpload,
	)

	// Pipeline endpoints
	r.Post(
		"/pipeline/run",
		pipelineHandler.RunPipeline,
	)

	r.Get(
		"/pipeline/status",
		pipelineHandler.GetStatus,
	)

	// Scheduler endpoints
	r.Get(
		"/scheduler/status",
		schedulerHandler.GetStatus,
	)

	// Config endpoints
	r.Post(
		"/config/reload",
		configHandler.ReloadConfig,
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

	log.Printf(
		"server running on port %s",
		cfg.AppPort,
	)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}