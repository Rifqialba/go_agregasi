package models

import "time"

type ValidationError struct {
	ID           string
	RawDataID    string
	SourceID     string
	ErrorMessage string
	RawPayload   []byte
	CreatedAt    time.Time
}