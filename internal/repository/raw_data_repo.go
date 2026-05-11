package repository

import (
	"context"

	"aggregation-dashboard/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jackc/pgx/v5"
)

type RawDataRepository struct {
	db *pgxpool.Pool
}

func NewRawDataRepository(db *pgxpool.Pool) *RawDataRepository {
	return &RawDataRepository{
		db: db,
	}
}

func (r *RawDataRepository) Insert(ctx context.Context, data models.RawData) error {
	query := `
		INSERT INTO raw_data (
			source_id,
			source_type,
			raw_payload,
			status,
			idempotency_key
		)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (idempotency_key) DO NOTHING
	`

	_, err := r.db.Exec(
		ctx,
		query,
		data.SourceID,
		data.SourceType,
		data.RawPayload,
		data.Status,
		data.IdempotencyKey,
	)

	return err
}
func (r *RawDataRepository) InsertBatch(
	ctx context.Context,
	dataList []models.RawData,
) error {
	query := `
		INSERT INTO raw_data (
			source_id,
			source_type,
			raw_payload,
			status,
			idempotency_key
		)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (idempotency_key) DO NOTHING
	`

	batch := &pgx.Batch{}

	for _, data := range dataList {
		batch.Queue(
			query,
			data.SourceID,
			data.SourceType,
			data.RawPayload,
			data.Status,
			data.IdempotencyKey,
		)
	}

	results := r.db.SendBatch(ctx, batch)

	defer results.Close()

	for range dataList {
		_, err := results.Exec()
		if err != nil {
			return err
		}
	}

	return nil
}