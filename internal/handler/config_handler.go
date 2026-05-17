package handler

import (
	"net/http"

	"aggregation-dashboard/internal/audit"
	"aggregation-dashboard/internal/models"
	"aggregation-dashboard/internal/scheduler"
	"aggregation-dashboard/internal/utils"
)

type ConfigHandler struct {
	scheduler   *scheduler.Scheduler
	auditLogger *audit.AuditLogger
}

func NewConfigHandler(
	scheduler *scheduler.Scheduler,
	auditLogger *audit.AuditLogger,
) *ConfigHandler {
	return &ConfigHandler{
		scheduler:   scheduler,
		auditLogger: auditLogger,
	}
}

func (h *ConfigHandler) ReloadConfig(
	w http.ResponseWriter,
	r *http.Request,
) {

	err := h.scheduler.ReloadConfig(
		"config.yaml",
	)

	if err != nil {
		utils.JSON(
			w,
			http.StatusInternalServerError,
			map[string]string{
				"error": err.Error(),
			},
		)

		return
	}

	h.auditLogger.Log(
		r.Context(),
		models.AuditActionConfigReload,
		"",
		map[string]any{
			"message": "scheduler config reloaded",
		},
	)

	utils.JSON(
		w,
		http.StatusOK,
		map[string]string{
			"status": "config reloaded",
		},
	)
}