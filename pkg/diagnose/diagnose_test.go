package diagnose

import (
	"testing"

	"github.com/sergiorivas/notify/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestRunDiagnostic(t *testing.T) {
	cfg := config.Config{
		EnabledNotifiers: []string{"audio", "dialog"},
	}

	// Mocking output for testing
	assert.NotPanics(t, func() {
		RunDiagnostic(cfg)
	})
}
