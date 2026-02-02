package commands

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestVersionCommand_Structure(t *testing.T) {
	setupTestEnv(t)

	assert.NotNil(t, versionCommand)
	assert.Equal(t, "version", versionCommand.Use)
	assert.Equal(t, "Get the current version", versionCommand.Short)
	assert.NotEmpty(t, versionCommand.Long)
	assert.NotNil(t, versionCommand.Run)
}
