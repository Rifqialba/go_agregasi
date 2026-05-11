package service

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"strings"

	"aggregation-dashboard/internal/models"
	"aggregation-dashboard/internal/repository"
	"aggregation-dashboard/internal/utils"
)

type UploadService struct {
	repo *repository.RawDataRepository
}

func NewUploadService(
	repo *repository.RawDataRepository,
) *UploadService {
	return &UploadService{
		repo: repo,
	}
}

func (s *UploadService) ProcessFile(
	sourceID string,
	content []byte,
) error {
	format := detectFormat(content)

	switch format {
	case "CSV":
		return s.processCSV(sourceID, content)

	case "JSON":
		return s.processJSON(sourceID, content)

	default:
		return errors.New("unsupported format")
	}
}

func detectFormat(content []byte) string {
	text := strings.TrimSpace(string(content))

	if strings.HasPrefix(text, "[") {
		return "JSON"
	}

	firstLine := strings.Split(text, "\n")[0]

	if strings.Contains(firstLine, ",") {
		return "CSV"
	}

	return "UNKNOWN"
}
func (s *UploadService) processCSV(
	sourceID string,
	content []byte,
) error {
	reader := csv.NewReader(strings.NewReader(string(content)))

	rows, err := reader.ReadAll()
	if err != nil {
		return err
	}

	if len(rows) < 2 {
		return errors.New("empty csv")
	}

	headers := rows[0]

	var rawDataList []models.RawData

	for _, row := range rows[1:] {
		record := map[string]any{}

		for i, value := range row {
			record[headers[i]] = value
		}

		payload, err := json.Marshal(record)
		if err != nil {
			return err
		}

		idempotencyKey := utils.SHA256(
			sourceID + string(payload),
		)

		rawDataList = append(rawDataList, models.RawData{
			SourceID:       sourceID,
			SourceType:     "FILE_UPLOAD",
			RawPayload:     payload,
			Status:         "PENDING",
			IdempotencyKey: idempotencyKey,
		})
	}

	return s.repo.InsertBatch(
		context.Background(),
		rawDataList,
	)
}
func (s *UploadService) processJSON(
	sourceID string,
	content []byte,
) error {
	var records []map[string]any

	err := json.Unmarshal(content, &records)
	if err != nil {
		return err
	}

	var rawDataList []models.RawData

	for _, record := range records {
		payload, err := json.Marshal(record)
		if err != nil {
			return err
		}

		idempotencyKey := utils.SHA256(
			sourceID + string(payload),
		)

		rawDataList = append(rawDataList, models.RawData{
			SourceID:       sourceID,
			SourceType:     models.SourceTypeFileUpload,
			RawPayload:     payload,
			Status:         models.StatusPending,
			IdempotencyKey: idempotencyKey,
		})
	}

	return s.repo.InsertBatch(
		context.Background(),
		rawDataList,
	)
}