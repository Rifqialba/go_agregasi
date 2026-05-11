package models

import "time"

type RawData struct {
	ID             string
	SourceID       string
	SourceType     string
	RawPayload     []byte
	ReceivedAt     time.Time
	Status         string
	IdempotencyKey string
}