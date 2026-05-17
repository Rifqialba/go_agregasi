package handler

import (
	"net/http"

	"aggregation-dashboard/internal/repository"
	"aggregation-dashboard/internal/utils"
)

type AuditHandler struct {
	repo *repository.AuditLogRepository
}

func NewAuditHandler(
	repo *repository.AuditLogRepository,
) *AuditHandler {
	return &AuditHandler{
		repo: repo,
	}
}

func (h *AuditHandler) GetAuditLogs(
	w http.ResponseWriter,
	r *http.Request,
) {

	action := r.URL.Query().Get("action")

	startDate := r.URL.Query().Get(
		"start_date",
	)

	endDate := r.URL.Query().Get(
		"end_date",
	)

	logs, err := h.repo.Find(
		r.Context(),
		action,
		startDate,
		endDate,
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

	utils.JSON(
		w,
		http.StatusOK,
		logs,
	)
}