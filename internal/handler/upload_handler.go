package handler

import (
	"io"
	"net/http"

	"aggregation-dashboard/internal/service"
)

const maxUploadSize = 10 << 20 // 10MB

type UploadHandler struct {
	service *service.UploadService
}

func NewUploadHandler(
	service *service.UploadService,
) *UploadHandler {
	return &UploadHandler{
		service: service,
	}
}

func (h *UploadHandler) HandleUpload(
	w http.ResponseWriter,
	r *http.Request,
) {
	r.Body = http.MaxBytesReader(
		w,
		r.Body,
		maxUploadSize,
	)

	err := r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		http.Error(
			w,
			"file too large",
			http.StatusRequestEntityTooLarge,
		)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(
			w,
			"file is required",
			http.StatusBadRequest,
		)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(
			w,
			"failed read file",
			http.StatusInternalServerError,
		)
		return
	}

	err = h.service.ProcessFile(
		"manual-upload",
		content,
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("upload success"))
}
