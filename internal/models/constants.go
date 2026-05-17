package models

const (
	SourceTypeWebhook   = "WEBHOOK"
	SourceTypeFileUpload = "FILE_UPLOAD"
	SourceTypeRestAPI   = "REST_API"
	SourceTypePolling   = "POLLING"
	StatusProcessed = "PROCESSED"
	StatusInvalid   = "INVALID"
	StatusPending = "PENDING"
	AuditActionPipelineRun  = "PIPELINE_RUN"
	AuditActionWebhook      = "WEBHOOK_RECEIVED"
	AuditActionUpload       = "FILE_UPLOAD"
	AuditActionConfigReload = "CONFIG_RELOAD"
)