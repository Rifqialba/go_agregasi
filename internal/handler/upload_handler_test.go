package handler

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"aggregation-dashboard/internal/service"

	"github.com/stretchr/testify/assert"
)

func TestFileUploadSizeLimit(t *testing.T) {
	var body bytes.Buffer

	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile(
		"file",
		"large.txt",
	)
	assert.NoError(t, err)

	largeContent := make([]byte, 11<<20)

	_, err = part.Write(largeContent)
	assert.NoError(t, err)

	writer.Close()

	uploadService := &service.UploadService{}

	handler := NewUploadHandler(uploadService)

	req := httptest.NewRequest(
		http.MethodPost,
		"/upload",
		&body,
	)

	req.Header.Set(
		"Content-Type",
		writer.FormDataContentType(),
	)

	rr := httptest.NewRecorder()

	handler.HandleUpload(rr, req)

	assert.Equal(
		t,
		http.StatusRequestEntityTooLarge,
		rr.Code,
	)
}