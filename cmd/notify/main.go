package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/sergiorivas/notify/pkg/diagnose"

	"github.com/sergiorivas/notify/internal/config"
	"github.com/sergiorivas/notify/internal/notifier"
)

var version = "0.1.0"

func main() {
	// Check if any arguments were provided
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	// Handle subcommands
	switch os.Args[1] {
	case "diagnose":
		// Diagnostic subcommand
		diagnosticCmd := flag.NewFlagSet("diagnose", flag.ExitOnError)
		diagnosticConfigFile := diagnosticCmd.String("config-file", "", "Configuration file to use (from ~/.config/notify/)")
		if err := diagnosticCmd.Parse(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing diagnose command: %v\n", err)
			os.Exit(1)
		}

		cfg, err := config.Load(*diagnosticConfigFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
			os.Exit(1)
		}
		diagnose.RunDiagnostic(cfg)
		return

	case "init":
		// Init subcommand to create default configuration
		initCmd := flag.NewFlagSet("init", flag.ExitOnError)
		initConfigFile := initCmd.String("config-file", "", "Configuration file to create (in ~/.config/notify/)")
		force := initCmd.Bool("force", false, "Overwrite existing configuration if it exists")
		if err := initCmd.Parse(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing init command: %v\n", err)
			os.Exit(1)
		}

		// Create dialog-only configuration
		dialogOnlyConfig := config.Config{
			EnabledNotifiers: []string{"dialog"},
			DialogSettings: map[string]string{
				"title": "Notification",
			},
		}

		configPath := config.GetConfigPath(*initConfigFile)

		// Check if file exists and force is not set
		if _, err := os.Stat(configPath); err == nil && !*force {
			fmt.Printf("Configuration file %s already exists. Use --force to overwrite.\n", configPath)
			return
		}

		// Make sure the directory exists
		configDir := config.GetConfigDir()
		if err := os.MkdirAll(configDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating configuration directory: %v\n", err)
			os.Exit(1)
		}

		// Save the configuration
		if err := config.Save(dialogOnlyConfig, configPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving configuration: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Created dialog-only configuration at %s\n", configPath)
		return

	case "list-notifiers":
		// List available notifier types
		listNotifiersCmd := flag.NewFlagSet("list-notifiers", flag.ExitOnError)
		verbose := listNotifiersCmd.Bool("verbose", false, "Show detailed information about each notifier")
		if err := listNotifiersCmd.Parse(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing list-notifiers command: %v\n", err)
			os.Exit(1)
		}

		allNotifiers := notifier.GetAllNotifiers()

		fmt.Println("Available notifier types:")
		fmt.Println("")

		for _, n := range allNotifiers {
			fmt.Printf("- %s\n", n.Name())

			if *verbose {
				// Run a quick diagnostic to get more info
				diagResult := n.Diagnose()
				statusStr := "NOT AVAILABLE"
				if diagResult.Available {
					statusStr = "AVAILABLE"
				}

				fmt.Printf("  Status: %s\n", statusStr)
				fmt.Printf("  Details: %s\n", diagResult.Message)
				fmt.Println("")
			}
		}

		fmt.Println("")
		fmt.Println("To use these notifiers, add them to your configuration file.")
		fmt.Println("Example config that enables all notifiers:")
		fmt.Println("")
		fmt.Println("```yaml")
		fmt.Println("enabledNotifiers:")
		for _, n := range allNotifiers {
			var id string
			switch n.(type) {
			case *notifier.AudioNotifier:
				id = "audio"
			case *notifier.DialogNotifier:
				id = "dialog"
			}
			fmt.Printf("  - %s\n", id)
		}
		fmt.Println("dialogSettings:")
		fmt.Println("  title: Notification")
		fmt.Println("```")
		return

	case "list-configs":
		// List available configuration files
		listConfigsFlag := true
		if listConfigsFlag {
			configFiles, err := config.ListConfigFiles()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error listing configuration files: %v\n", err)
				os.Exit(1)
			}

			if len(configFiles) == 0 {
				fmt.Println("No configuration files found in", config.GetConfigDir())
				fmt.Println("Run 'notify init' to create a default configuration.")
				return
			}

			fmt.Println("Available configuration files in", config.GetConfigDir()+":")
			for _, file := range configFiles {
				fmt.Println("  -", file)
			}
			return
		}

	case "version":
		// Print version information
		fmt.Println("notify version:", version)
		return
	}

	// Main command for sending notifications
	mainCmd := flag.NewFlagSet("notify", flag.ExitOnError)
	configFlag := mainCmd.Bool("config", false, "Edit configuration")
	configFile := mainCmd.String("config-file", "", "Configuration file to use (from ~/.config/notify/)")
	notificationType := mainCmd.String("type", "info", "Notification type: success, error, info, warning")
	title := mainCmd.String("title", "", "Custom title for notification (optional)")

	// Parse the remaining arguments
	if err := mainCmd.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing main command: %v\n", err)
		os.Exit(1)
	}

	// Check if help was requested
	if mainCmd.NArg() == 0 && !*configFlag && *configFile == "" && *title == "" && *notificationType == "info" {
		printUsage()
		return
	}

	// Load configuration
	cfg, err := config.Load(*configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Configuration mode
	if *configFlag {
		configPath := config.GetConfigPath(*configFile)
		fmt.Println("Configuration mode not implemented in version 0")
		fmt.Println("Manually edit the file:", configPath)
		fmt.Println("Configuration directory:", config.GetConfigDir())
		return
	}

	// Validate notification type
	validTypes := map[string]bool{"success": true, "error": true, "info": true, "warning": true}
	if !validTypes[*notificationType] {
		fmt.Fprintf(os.Stderr, "Invalid notification type: %s\n", *notificationType)
		fmt.Fprintln(os.Stderr, "Valid types: success, error, info, warning")
		os.Exit(1)
	}

	// Get the message
	args := mainCmd.Args()
	message := "Notification received"
	if len(args) > 0 {
		message = args[0]
	}

	// Use the title provided by command line or from configuration
	notificationTitle := cfg.DialogSettings["title"]
	if *title != "" {
		notificationTitle = *title
	}

	// Get notifiers and send notification asynchronously
	notifiers := notifier.GetEnabledNotifiers(cfg)
	if len(notifiers) == 0 {
		fmt.Fprintln(os.Stderr, "No enabled notifiers found in configuration. Run 'notify list-notifiers' to see available options.")
		os.Exit(1)
	}

	var wg sync.WaitGroup

	for _, n := range notifiers {
		wg.Add(1)
		go func(n notifier.Notifier) {
			defer wg.Done()
			err := n.Notify(message, *notificationType, notificationTitle)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error notifying with %s: %v\n", n.Name(), err)
			}
		}(n)
	}

	wg.Wait()
}

func printUsage() {
	fmt.Println("Notify - CLI Notification System")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  notify [options] <message>    - Send a notification")
	fmt.Println("  notify diagnose [options]     - Run notifier diagnostics")
	fmt.Println("  notify init [options]         - Create a default configuration")
	fmt.Println("  notify list-notifiers         - List all available notifier types")
	fmt.Println("  notify list-configs           - List available configuration files")
	fmt.Println("")
	fmt.Println("Options for sending notifications:")
	fmt.Println("  --type string      Notification type: success, error, info, warning (default \"info\")")
	fmt.Println("  --title string     Custom title for notification (optional)")
	fmt.Println("  --config-file string  Configuration file to use (from ~/.config/notify/)")
	fmt.Println("")
	fmt.Println("Options for init:")
	fmt.Println("  --config-file string  Configuration file to create (default \"config.yaml\")")
	fmt.Println("  --force            Overwrite existing configuration if it exists")
	fmt.Println("")
	fmt.Println("Options for list-notifiers:")
	fmt.Println("  --verbose          Show detailed information about each notifier")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  notify --type success \"Build successful\"")
	fmt.Println("  notify --type error --title \"Error\" \"Build failed\"")
	fmt.Println("  notify --config-file work.yaml \"Using custom config\"")
	fmt.Println("  notify init --config-file custom.yaml")
	fmt.Println("  notify list-notifiers --verbose")
}
