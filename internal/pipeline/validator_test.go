package pipeline

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatorValidPayload(t *testing.T) {
	schema := &ValidationSchema{
		RequiredFields: []string{
			"id",
			"name",
			"value",
			"date",
		},
		FieldTypes: map[string]string{
			"id":    "string",
			"name":  "string",
			"value": "number",
			"date":  "date:2006-01-02",
		},
	}

	validator := NewValidator(schema)

	payload := map[string]any{
		"id":    "001",
		"name":  "Alpha",
		"value": 100,
		"date":  "2024-03-15",
	}

	err := validator.Validate(payload)

	assert.NoError(t, err)
}

func TestValidatorMissingField(t *testing.T) {
	schema := &ValidationSchema{
		RequiredFields: []string{
			"id",
			"name",
		},
	}

	validator := NewValidator(schema)

	payload := map[string]any{
		"id": "001",
	}

	err := validator.Validate(payload)

	assert.Error(t, err)
}

func TestValidatorInvalidDate(t *testing.T) {
	schema := &ValidationSchema{
		FieldTypes: map[string]string{
			"date": "date:2006-01-02",
		},
	}

	validator := NewValidator(schema)

	payload := map[string]any{
		"date": "15-03-2024",
	}

	err := validator.Validate(payload)

	assert.Error(t, err)
}