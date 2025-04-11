package diagnose

import (
	"fmt"

	"github.com/sergiorivas/notify/internal/config"
	"github.com/sergiorivas/notify/internal/notifier"
)

// RunDiagnostic runs diagnostics on all notifiers
func RunDiagnostic(cfg config.Config) {
	allNotifiers := notifier.GetAllNotifiers()
	fmt.Println("Running notifier diagnostics...")
	fmt.Println()

	for _, n := range allNotifiers {
		result := n.Diagnose()

		// Check if it's enabled in the configuration
		var enabled bool
		for _, name := range cfg.EnabledNotifiers {
			if name == n.ID() {
				enabled = true
				break
			}
		}

		// Show diagnostic result
		statusStr := "NOT AVAILABLE"
		if result.Available {
			statusStr = "AVAILABLE"
		}

		enabledStr := "DISABLED"
		if enabled {
			enabledStr = "ENABLED"
		}

		fmt.Printf("Notifier: %s (%s)\n", n.Name(), n.ID())
		fmt.Printf("Status: %s\n", statusStr)
		fmt.Printf("Configuration: %s\n", enabledStr)
		fmt.Printf("Details: %s\n", result.Message)
		fmt.Println()
	}
}
