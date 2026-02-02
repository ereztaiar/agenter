package version

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestVersionDefaults(t *testing.T) {
    // Test that default values are set
    assert.NotEmpty(t, Version)
    assert.NotEmpty(t, GitCommit)
    assert.NotEmpty(t, BuildDate)
}