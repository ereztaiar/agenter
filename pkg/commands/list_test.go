package commands

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListCmd_Structure(t *testing.T) {
	setupTestEnv(t)

	assert.NotNil(t, listCmd)
	assert.Equal(t, "list", listCmd.Use)
	assert.Equal(t, "List all running agents", listCmd.Short)
	assert.NotEmpty(t, listCmd.Long)
}
