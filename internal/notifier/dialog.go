package notifier

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/sergiorivas/notify/internal/formatter"
)

// DialogNotifier for dialog notifications (using 'osascript' on macOS)
type DialogNotifier struct{}

func (d *DialogNotifier) Name() string {
	return "Dialog (osascript)"
}

func (d *DialogNotifier) ID() string {
	return "dialog"
}

func (d *DialogNotifier) Notify(message string, notificationType string, title string) error {
	formattedMessage := formatter.FormatMessage(message, notificationType, true)
	iconType := "stop" // Default value for error

	switch notificationType {
	case "success":
		iconType = "note"
	case "info":
		iconType = "note"
	case "warning":
		iconType = "caution"
	}

	script := fmt.Sprintf(
		`display dialog "%s" buttons {"OK"} default button "OK" with icon %s with title "%s"`,
		formattedMessage, iconType, title,
	)
	cmd := exec.Command("osascript", "-e", script)
	return cmd.Run()
}

func (d *DialogNotifier) Diagnose() DiagnosticResult {
	// Verify that the 'osascript' command exists
	path, err := exec.LookPath("osascript")
	if err != nil {
		return DiagnosticResult{
			Available: false,
			Message:   fmt.Sprintf("The 'osascript' command is not available: %v", err),
		}
	}

	// Check macOS version
	cmd := exec.Command("sw_vers", "-productVersion")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return DiagnosticResult{
			Available: true,
			Message:   fmt.Sprintf("'osascript' command available at: %s. Could not determine macOS version.", path),
		}
	}

	return DiagnosticResult{
		Available: true,
		Message:   fmt.Sprintf("'osascript' command available at: %s. macOS version: %s", path, strings.TrimSpace(string(output))),
	}
}
