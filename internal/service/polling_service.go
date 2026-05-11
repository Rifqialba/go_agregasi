package service

import (
	"context"
	"log"
	"time"

	"aggregation-dashboard/internal/collector"
	"aggregation-dashboard/internal/models"
	"aggregation-dashboard/internal/repository"
	"aggregation-dashboard/internal/utils"
)

type PollingService struct {
	client *collector.RestClient
	repo   *repository.RawDataRepository
}

func NewPollingService(
	client *collector.RestClient,
	repo *repository.RawDataRepository,
) *PollingService {
	return &PollingService{
		client: client,
		repo:   repo,
	}
}

func (s *PollingService) Poll(
	ctx context.Context,
	sourceID string,
	url string,
	lastFetchedAt time.Time,
) error {
	backoffs := []time.Duration{
		5 * time.Second,
		15 * time.Second,
		30 * time.Second,
	}

	for attempt := 0; attempt <= len(backoffs); attempt++ {
		body, statusCode, err := s.client.Fetch(
			ctx,
			url,
			lastFetchedAt,
		)

		if err == nil {
			if statusCode == 304 {
				return nil
			}

			idempotencyKey := utils.SHA256(
				sourceID + string(body),
			)

			rawData := models.RawData{
				SourceID:       sourceID,
				SourceType:     models.SourceTypePolling,
				RawPayload:     body,
				Status:         models.StatusPending,
				IdempotencyKey: idempotencyKey,
			}

			return s.repo.Insert(ctx, rawData)
		}

		if statusCode >= 400 && statusCode < 500 {
			log.Printf("client error: %v", err)

			return err
		}

		if attempt < len(backoffs) {
			wait := backoffs[attempt]

			log.Printf(
				"retrying in %v...",
				wait,
			)

			time.Sleep(wait)

			continue
		}

		return err
	}

	return nil
}