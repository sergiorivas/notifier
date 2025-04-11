package notifier

import (
	"github.com/sergiorivas/notify/internal/config"
)

// Notifier is the interface that all notifiers must implement
type Notifier interface {
	Notify(message string, notificationType string, title string) error
	Diagnose() DiagnosticResult
	Name() string
	ID() string // Added ID method to return identifier used in config
}

// DiagnosticResult represents the result of a diagnostic
type DiagnosticResult struct {
	Available bool
	Message   string
}

// GetAllNotifiers returns all available notifiers
func GetAllNotifiers() []Notifier {
	return []Notifier{
		&AudioNotifier{},
		&DialogNotifier{},
	}
}

// GetEnabledNotifiers returns the enabled notifiers according to configuration
func GetEnabledNotifiers(cfg config.Config) []Notifier {
	allNotifiers := GetAllNotifiers()
	var enabledNotifiers []Notifier

	enabledMap := make(map[string]bool)
	for _, n := range cfg.EnabledNotifiers {
		enabledMap[n] = true
	}

	for _, n := range allNotifiers {
		if enabledMap[n.ID()] {
			enabledNotifiers = append(enabledNotifiers, n)
		}
	}

	return enabledNotifiers
}
