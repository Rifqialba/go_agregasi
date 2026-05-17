package pipeline

import (
	"encoding/json"
	"os"
)

func LoadValidationSchema(
	path string,
) (*ValidationSchema, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var schema ValidationSchema

	err = json.Unmarshal(file, &schema)
	if err != nil {
		return nil, err
	}

	return &schema, nil
}