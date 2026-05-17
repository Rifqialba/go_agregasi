package audit

import (
	"context"
	"encoding/json"
	"log"

	"aggregation-dashboard/internal/models"
	"aggregation-dashboard/internal/repository"
)

type AuditLogger struct {
	repo *repository.AuditLogRepository
}

func NewAuditLogger(
	repo *repository.AuditLogRepository,
) *AuditLogger {
	return &AuditLogger{
		repo: repo,
	}
}

func (a *AuditLogger) Log(
	ctx context.Context,
	action string,
	entityID string,
	details any,
) {
	payload, err := json.Marshal(details)
	if err != nil {
		log.Printf(
			"failed marshal audit log: %v",
			err,
		)

		return
	}

	auditLog := models.AuditLog{
		Action:   action,
		EntityID: entityID,
		Details:  payload,
	}

	err = a.repo.Insert(ctx, auditLog)
	if err != nil {
		log.Printf(
			"failed insert audit log: %v",
			err,
		)
	}
}