package pipeline

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizerTrimWhitespace(t *testing.T) {
	normalizer := NewNormalizer()

	payload := map[string]any{
		"name": "  Rifqi Alba  ",
	}

	result := normalizer.Normalize(payload)

	assert.Equal(
		t,
		"Rifqi Alba",
		result["name"],
	)
}

func TestNormalizerDateFormat(t *testing.T) {
	normalizer := NewNormalizer()

	payload := map[string]any{
		"date": "2026-03-15",
	}

	result := normalizer.Normalize(payload)

	assert.Equal(
		t,
		"2026-03-15T00:00:00Z",
		result["date"],
	)
}

func TestNormalizerFieldRename(t *testing.T) {
	normalizer := NewNormalizer()

	payload := map[string]any{
		"nm": "Alba",
	}

	result := normalizer.Normalize(payload)

	assert.Equal(
		t,
		"Alba",
		result["name"],
	)
}