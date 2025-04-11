package formatter

import (
	"fmt"
)

// FormatMessage formats the message with the corresponding emoji
func FormatMessage(message string, notificationType string, isDialog bool) string {
	var emoji string

	// Fixed emojis for each type
	switch notificationType {
	case "success":
		emoji = "✅"
	case "error":
		emoji = "❌"
	case "info":
		emoji = "ℹ️"
	case "warning":
		emoji = "⚠️"
	default:
		return message
	}

	if isDialog {
		return fmt.Sprintf("%s %s", emoji, message)
	}
	return fmt.Sprintf("%s, %s", emoji, message)
}
