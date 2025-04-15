package formatter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatMessage(t *testing.T) {
	tests := []struct {
		message          string
		notificationType string
		isDialog         bool
		expected         string
	}{
		{"Test message", "success", true, "✅ Test message"},
		{"Test message", "error", false, "error, Test message"},
		{"Test message", "info", true, "ℹ️ Test message"},
		{"Test message", "warning", false, "alert, Test message"},
	}

	for _, tt := range tests {
		result := FormatMessage(tt.message, tt.notificationType, tt.isDialog)
		assert.Equal(t, tt.expected, result)
	}
}
