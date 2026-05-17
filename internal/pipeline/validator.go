package pipeline

import (
	"fmt"
	"time"
)

type Validator struct {
	schema *ValidationSchema
}

func NewValidator(
	schema *ValidationSchema,
) *Validator {
	return &Validator{
		schema: schema,
	}
}

func (v *Validator) Validate(
	payload map[string]any,
) error {
	// Validate required fields
	for _, field := range v.schema.RequiredFields {
		value, exists := payload[field]

		if !exists {
			return fmt.Errorf(
				"missing required field: %s",
				field,
			)
		}

		if value == nil {
			return fmt.Errorf(
				"field is nil: %s",
				field,
			)
		}
	}

	// Validate field types
	for field, expectedType := range v.schema.FieldTypes {
		value, exists := payload[field]

		if !exists {
			continue
		}

		err := validateType(
			field,
			value,
			expectedType,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func validateType(
	field string,
	value any,
	expectedType string,
) error {
	switch expectedType {

	case "string":
		_, ok := value.(string)

		if !ok {
			return fmt.Errorf(
				"field %s must be string",
				field,
			)
		}

	case "number":
		switch value.(type) {
		case float64, int, int64:
			// valid
		default:
			return fmt.Errorf(
				"field %s must be number",
				field,
			)
		}

	default:
		// Date validation
		if len(expectedType) > 5 &&
			expectedType[:5] == "date:" {

			dateFormat := expectedType[5:]

			strValue, ok := value.(string)
			if !ok {
				return fmt.Errorf(
					"field %s must be string date",
					field,
				)
			}

			_, err := time.Parse(
				dateFormat,
				strValue,
			)

			if err != nil {
				return fmt.Errorf(
					"field %s invalid date format",
					field,
				)
			}
		}
	}

	return nil
}