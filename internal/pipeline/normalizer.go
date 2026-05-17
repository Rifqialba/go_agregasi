package pipeline

import (
	"strings"
	"time"
)

type Normalizer struct {
	fieldMappings map[string]string
	dateFields    map[string]bool
}

func NewNormalizer() *Normalizer {
	return &Normalizer{
		fieldMappings: map[string]string{
			"nm": "name",
		},

		dateFields: map[string]bool{
			"date": true,
		},
	}
}

func (n *Normalizer) Normalize(
	payload map[string]any,
) map[string]any {
	normalized := map[string]any{}

	for key, value := range payload {

		// Rename field
		if mappedField, exists := n.fieldMappings[key]; exists {
			key = mappedField
		}

		switch v := value.(type) {

		case string:
			trimmed := strings.TrimSpace(v)

			// Normalize known date fields
			if n.dateFields[key] {
				parsedDate, err := time.Parse(
					"2006-01-02",
					trimmed,
				)

				if err == nil {
					normalized[key] = parsedDate.UTC().Format(
						time.RFC3339,
					)

					continue
				}
			}

			normalized[key] = trimmed

		default:
			normalized[key] = value
		}
	}

	return normalized
}