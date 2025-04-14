package notifier

import (
	internal "github.com/sergiorivas/notify/internal/notifier"
)

func Notify(message string, notificationType string) {
	internal.Notify(message, notificationType)
}
