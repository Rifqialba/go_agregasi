package models

import "time"

type ProcessedData struct {
	ID                string
	RawDataID         string
	SourceID          string
	NormalizedPayload []byte
	ProcessedAt       time.Time
	IdempotencyKey    string
}