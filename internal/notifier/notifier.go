package notifier

import (
	"fmt"

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

// Notify sends a notification using the available notifiers.
// If no configuration file is found, it defaults to dialog and audio notifiers.
func Notify(message string, notificationType string) {
	// Attempt to load configuration
	cfg, err := config.Load("")
	var notifiers []Notifier

	if err != nil {
		// Default to dialog and audio if configuration is not found
		notifiers = []Notifier{
			&DialogNotifier{},
			&AudioNotifier{},
		}
	} else {
		// Get enabled notifiers from configuration
		notifiers = GetEnabledNotifiers(cfg)
		if len(notifiers) == 0 {
			// Default to dialog and audio if no notifiers are enabled
			notifiers = []Notifier{
				&DialogNotifier{},
				&AudioNotifier{},
			}
		}
	}

	// Send notification using all available notifiers
	for _, n := range notifiers {
		if err := n.Notify(message, notificationType, ""); err != nil {
			// Log error (can be replaced with proper logging)
			fmt.Printf("Error notifying with %s: %v\n", n.Name(), err)
		}
	}
}
