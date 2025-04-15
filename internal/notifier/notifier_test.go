package notifier

import (
	"testing"

	"github.com/sergiorivas/notify/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestGetAllNotifiers(t *testing.T) {
	notifiers := GetAllNotifiers()
	assert.NotEmpty(t, notifiers)
	assert.Equal(t, 2, len(notifiers)) // DialogNotifier and AudioNotifier
}

func TestGetEnabledNotifiers(t *testing.T) {
	cfg := config.Config{
		EnabledNotifiers: []string{"audio"},
	}

	enabledNotifiers := GetEnabledNotifiers(cfg)
	assert.Len(t, enabledNotifiers, 1)
	assert.Equal(t, "audio", enabledNotifiers[0].ID())
}
