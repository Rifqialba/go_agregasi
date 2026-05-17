package pipeline

import (
	"testing"

	"aggregation-dashboard/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestPipelineRerun(t *testing.T) {
	payload := `{"id":"1","name":"Alpha"}`

	keyRun1 := utils.SHA256(
		"source-001" + payload,
	)

	keyRun2 := utils.SHA256(
		"source-001" + payload,
	)

	assert.Equal(
		t,
		keyRun1,
		keyRun2,
	)
}