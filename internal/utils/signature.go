package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func ValidateHMACSignature(payload []byte, secret string, signature string) bool {
	if !strings.HasPrefix(signature, "sha256=") {
		return false
	}

	receivedSignature := strings.TrimPrefix(signature, "sha256=")

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)

	expectedMAC := mac.Sum(nil)

	receivedMAC, err := hex.DecodeString(receivedSignature)
	if err != nil {
		return false
	}

	return hmac.Equal(expectedMAC, receivedMAC)
}