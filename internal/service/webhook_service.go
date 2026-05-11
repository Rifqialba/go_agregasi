package service

import (
	"context"
	"log"

	"aggregation-dashboard/internal/models"
	"aggregation-dashboard/internal/repository"
	"aggregation-dashboard/internal/utils"
)

const workerCount = 4

type WebhookJob struct {
	SourceID string
	Payload  []byte
}

type WebhookService struct {
	repo *repository.RawDataRepository
	jobs chan WebhookJob
}

func NewWebhookService(repo *repository.RawDataRepository) *WebhookService {
	service := &WebhookService{
		repo: repo,
		jobs: make(chan WebhookJob, 100),
	}

	service.startWorkers()

	return service
}

func (s *WebhookService) Enqueue(job WebhookJob) {
	s.jobs <- job
}

func (s *WebhookService) startWorkers() {
	for i := 0; i < workerCount; i++ {
		go s.worker(i)
	}
}

func (s *WebhookService) worker(workerID int) {
	log.Printf("webhook worker %d started", workerID)

	for job := range s.jobs {
		s.process(job)
	}
}

func (s *WebhookService) process(job WebhookJob) {
	idempotencyKey := utils.SHA256(
		job.SourceID + string(job.Payload),
	)

	rawData := models.RawData{
		SourceID:       job.SourceID,
		SourceType:     models.SourceTypeWebhook,
		RawPayload:     job.Payload,
		Status:         models.StatusPending,
		IdempotencyKey: idempotencyKey,
	}

	err := s.repo.Insert(context.Background(), rawData)
	if err != nil {
		log.Printf("failed insert raw data: %v", err)
	}
}