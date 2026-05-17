package repository

import (
	"context"

	"aggregation-dashboard/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ValidationErrorsRepository struct {
	db *pgxpool.Pool
}

func NewValidationErrorsRepository(
	db *pgxpool.Pool,
) *ValidationErrorsRepository {
	return &ValidationErrorsRepository{
		db: db,
	}
}

func (r *ValidationErrorsRepository) Insert(
	ctx context.Context,
	data models.ValidationError,
) error {
	query := `
		INSERT INTO validation_errors (
			raw_data_id,
			source_id,
			error_message,
			raw_payload
		)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		data.RawDataID,
		data.SourceID,
		data.ErrorMessage,
		data.RawPayload,
	)

	return err
}