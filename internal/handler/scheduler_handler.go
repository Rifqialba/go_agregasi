package handler

import (
	"net/http"

	"aggregation-dashboard/internal/scheduler"
	"aggregation-dashboard/internal/utils"
)

type SchedulerHandler struct {
	scheduler *scheduler.Scheduler
}

func NewSchedulerHandler(
	scheduler *scheduler.Scheduler,
) *SchedulerHandler {
	return &SchedulerHandler{
		scheduler: scheduler,
	}
}

func (h *SchedulerHandler) GetStatus(
	w http.ResponseWriter,
	r *http.Request,
) {
	status := h.scheduler.GetStatus()

	utils.JSON(
		w,
		http.StatusOK,
		status,
	)
}