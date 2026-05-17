package pipeline

import (
	"context"
	"encoding/json"
	"log"

	"aggregation-dashboard/internal/models"
	"aggregation-dashboard/internal/repository"
	"aggregation-dashboard/internal/utils"
)

type Processor struct {
	rawRepo        *repository.RawDataRepository
	processedRepo  *repository.ProcessedDataRepository
	validationRepo *repository.ValidationErrorsRepository
	validator      *Validator
	normalizer     *Normalizer
}

func NewProcessor(
	rawRepo *repository.RawDataRepository,
	processedRepo *repository.ProcessedDataRepository,
	validationRepo *repository.ValidationErrorsRepository,
	validator *Validator,
	normalizer *Normalizer,
) *Processor {
	return &Processor{
		rawRepo:        rawRepo,
		processedRepo:  processedRepo,
		validationRepo: validationRepo,
		validator:      validator,
		normalizer:     normalizer,
	}
}

func (p *Processor) ProcessPendingData(
	ctx context.Context,
) (int, int, error) {

	rawDataList, err := p.rawRepo.GetPending(ctx)
	if err != nil {
		return 0, 0, err
	}

	log.Printf(
		"processing %d pending records",
		len(rawDataList),
	)

	processed := 0
	errorsCount := 0

	for _, rawData := range rawDataList {

		err := p.processSingle(ctx, rawData)

		if err != nil {
			errorsCount++

			log.Printf(
				"failed processing raw_data %s: %v",
				rawData.ID,
				err,
			)

			continue
		}

		processed++
	}

	return processed, errorsCount, nil
}
func (p *Processor) processSingle(
	ctx context.Context,
	rawData models.RawData,
) error {
	var payload map[string]any

	err := json.Unmarshal(
		rawData.RawPayload,
		&payload,
	)
	if err != nil {
		return err
	}

	// Validation
	err = p.validator.Validate(payload)
	if err != nil {

		validationError := models.ValidationError{
			RawDataID:    rawData.ID,
			SourceID:     rawData.SourceID,
			ErrorMessage: err.Error(),
			RawPayload:   rawData.RawPayload,
		}

		errInsert := p.validationRepo.Insert(
			ctx,
			validationError,
		)

		if errInsert != nil {
			return errInsert
		}

		return p.rawRepo.UpdateStatus(
			ctx,
			rawData.ID,
			models.StatusInvalid,
		)
	}

	// Normalize
	normalizedPayload := p.normalizer.Normalize(
		payload,
	)

	normalizedJSON, err := json.Marshal(
		normalizedPayload,
	)
	if err != nil {
		return err
	}

	// ETL idempotency key
	idempotencyKey := utils.SHA256(
		rawData.SourceID + string(normalizedJSON),
	)

	processedData := models.ProcessedData{
		RawDataID:         rawData.ID,
		SourceID:          rawData.SourceID,
		NormalizedPayload: normalizedJSON,
		IdempotencyKey:    idempotencyKey,
	}

	err = p.processedRepo.Insert(
		ctx,
		processedData,
	)
	if err != nil {
		return err
	}

	return p.rawRepo.UpdateStatus(
		ctx,
		rawData.ID,
		models.StatusProcessed,
	)
}