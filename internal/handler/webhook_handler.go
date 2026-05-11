package handler

import (
	"io"
	"net/http"

	"aggregation-dashboard/internal/service"
	"aggregation-dashboard/internal/utils"

	"github.com/go-chi/chi/v5"
)

type WebhookHandler struct {
	service *service.WebhookService
	secret  string
}

func NewWebhookHandler(
	service *service.WebhookService,
	secret string,
) *WebhookHandler {
	return &WebhookHandler{
		service: service,
		secret:  secret,
	}
}

func (h *WebhookHandler) HandleWebhook(
	w http.ResponseWriter,
	r *http.Request,
) {
	sourceID := chi.URLParam(r, "source_id")

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	signature := r.Header.Get("X-Signature")

	valid := utils.ValidateHMACSignature(
		payload,
		h.secret,
		signature,
	)

	if !valid {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	h.service.Enqueue(service.WebhookJob{
		SourceID: sourceID,
		Payload:  payload,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("accepted"))
}