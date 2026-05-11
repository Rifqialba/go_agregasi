package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebhookSignatureValid(t *testing.T) {
	payload := []byte(`{"event":"data_update","count":42}`)

	secret := "webhook-secret-2024"

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)

	signature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	valid := ValidateHMACSignature(
		payload,
		secret,
		signature,
	)

	assert.True(t, valid)
}
func TestWebhookSignatureInvalid(t *testing.T) {
	payload := []byte(`{"event":"data_update","count":42}`)

	secret := "webhook-secret-2024"

	invalidSignature := "sha256=invalidsignature"

	valid := ValidateHMACSignature(
		payload,
		secret,
		invalidSignature,
	)

	assert.False(t, valid)
}