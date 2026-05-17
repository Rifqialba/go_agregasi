package models

import "time"

type AuditLog struct {
	ID        string
	Action    string
	EntityID  string
	Details   []byte
	CreatedAt time.Time
}