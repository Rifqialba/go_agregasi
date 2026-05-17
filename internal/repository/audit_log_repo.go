package repository

import (
	"context"
	"strconv"

	"aggregation-dashboard/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditLogRepository struct {
	db *pgxpool.Pool
}

func NewAuditLogRepository(
	db *pgxpool.Pool,
) *AuditLogRepository {
	return &AuditLogRepository{
		db: db,
	}
}

func (r *AuditLogRepository) Insert(
	ctx context.Context,
	auditLog models.AuditLog,
) error {

	query := `
		INSERT INTO audit_log (
			action,
			entity_id,
			details
		)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		auditLog.Action,
		auditLog.EntityID,
		auditLog.Details,
	)

	return err
}

func (r *AuditLogRepository) Find(
	ctx context.Context,
	action string,
	startDate string,
	endDate string,
) ([]models.AuditLog, error) {

	query := `
		SELECT
			id,
			action,
			entity_id,
			details,
			created_at
		FROM audit_log
		WHERE 1=1
	`

	args := []any{}
	argPos := 1

	if action != "" {
		query += `
			AND action = $` + strconv.Itoa(argPos)

		args = append(args, action)

		argPos++
	}

	if startDate != "" {
		query += `
			AND created_at >= $` + strconv.Itoa(argPos)

		args = append(args, startDate)

		argPos++
	}

	if endDate != "" {
		query += `
			AND created_at <= $` + strconv.Itoa(argPos)

		args = append(args, endDate)

		argPos++
	}

	query += `
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []models.AuditLog

	for rows.Next() {
		var logEntry models.AuditLog

		err := rows.Scan(
			&logEntry.ID,
			&logEntry.Action,
			&logEntry.EntityID,
			&logEntry.Details,
			&logEntry.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, logEntry)
	}

	return result, nil
}