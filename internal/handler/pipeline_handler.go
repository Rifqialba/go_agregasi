package handler

import (
	"net/http"

	"aggregation-dashboard/internal/pipeline"
	"aggregation-dashboard/internal/utils"
)

type PipelineHandler struct {
	runner *pipeline.PipelineRunner
}

func NewPipelineHandler(
	runner *pipeline.PipelineRunner,
) *PipelineHandler {
	return &PipelineHandler{
		runner: runner,
	}
}

func (h *PipelineHandler) RunPipeline(
	w http.ResponseWriter,
	r *http.Request,
) {
	h.runner.Run()

	utils.JSON(
		w,
		http.StatusAccepted,
		map[string]string{
			"status": "pipeline started",
		},
	)
}

func (h *PipelineHandler) GetStatus(
	w http.ResponseWriter,
	r *http.Request,
) {
	status := h.runner.GetStatus()

	utils.JSON(
		w,
		http.StatusOK,
		status,
	)
}