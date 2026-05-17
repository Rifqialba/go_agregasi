package pipeline

type ValidationSchema struct {
	RequiredFields []string          `json:"required_fields"`
	FieldTypes     map[string]string `json:"field_types"`
}