package notifier

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/sergiorivas/notify/internal/formatter"
)

// AudioNotifier for audio notifications (using 'say' on macOS)
type AudioNotifier struct{}

func (a *AudioNotifier) Name() string {
	return "Audio (say)"
}

func (a *AudioNotifier) ID() string {
	return "audio"
}

func (a *AudioNotifier) Notify(message string, notificationType string, title string) error {
	formattedMessage := formatter.FormatMessage(message, notificationType, false)
	cmd := exec.Command("say", formattedMessage)
	return cmd.Run()
}

func (a *AudioNotifier) Diagnose() DiagnosticResult {
	// Verify that the 'say' command exists
	path, err := exec.LookPath("say")
	if err != nil {
		return DiagnosticResult{
			Available: false,
			Message:   fmt.Sprintf("The 'say' command is not available: %v", err),
		}
	}

	// Try running a test command
	cmd := exec.Command("say", "-v", "?")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return DiagnosticResult{
			Available: false,
			Message:   fmt.Sprintf("Error testing the 'say' command: %v", err),
		}
	}

	return DiagnosticResult{
		Available: true,
		Message:   fmt.Sprintf("'say' command available at: %s. Available voices: %d", path, strings.Count(string(output), "\n")),
	}
}
