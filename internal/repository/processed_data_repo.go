package repository

import (
	"context"

	"aggregation-dashboard/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProcessedDataRepository struct {
	db *pgxpool.Pool
}

func NewProcessedDataRepository(
	db *pgxpool.Pool,
) *ProcessedDataRepository {
	return &ProcessedDataRepository{
		db: db,
	}
}

func (r *ProcessedDataRepository) Insert(
	ctx context.Context,
	data models.ProcessedData,
) error {
	query := `
		INSERT INTO processed_data (
			raw_data_id,
			source_id,
			normalized_payload,
			idempotency_key
		)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (idempotency_key) DO NOTHING
	`

	_, err := r.db.Exec(
		ctx,
		query,
		data.RawDataID,
		data.SourceID,
		data.NormalizedPayload,
		data.IdempotencyKey,
	)

	return err
}