package pipeline

import (
	"testing"

	"aggregation-dashboard/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestIdempotencyKey(t *testing.T) {
	payload1 := `{"name":"Alpha"}`
	payload2 := `{"name":"Alpha"}`

	key1 := utils.SHA256(
		"source-001" + payload1,
	)

	key2 := utils.SHA256(
		"source-001" + payload2,
	)

	assert.Equal(t, key1, key2)
}