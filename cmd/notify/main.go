package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/sergiorivas/notify/internal/config"
	"github.com/sergiorivas/notify/internal/notifier"
	"github.com/sergiorivas/notify/pkg/diagnose"
)

var version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: No command specified.")
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "send":
		sendCommand(os.Args[2:])
	case "config":
		configCommand(os.Args[2:])
	case "notifiers":
		notifiersCommand(os.Args[2:])
	case "diagnose":
		diagnoseCommand(os.Args[2:])
	case "version":
		versionCommand()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func defaultCommand() {
	fmt.Println("Default command: send notification")
	printUsage()
}

func sendCommand(args []string) {
	cmd := flag.NewFlagSet("send", flag.ExitOnError)
	notificationType := cmd.String("type", "info", "Notification type: success, error, info, warning")
	title := cmd.String("title", "", "Custom title for notification")
	configFile := cmd.String("config", "", "Configuration file to use")
	if err := cmd.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing send command: %v\n", err)
		os.Exit(1)
	}

	if cmd.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "Error: No message specified.")
		cmd.Usage()
		os.Exit(1)
	}

	cfg, err := config.Load(*configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	message := cmd.Arg(0)

	notifiers := notifier.GetEnabledNotifiers(cfg)
	if len(notifiers) == 0 {
		fmt.Fprintln(os.Stderr, "No enabled notifiers found in configuration.")
		os.Exit(1)
	}

	var wg sync.WaitGroup
	for _, n := range notifiers {
		wg.Add(1)
		go func(n notifier.Notifier) {
			defer wg.Done()
			if err := n.Notify(message, *notificationType, *title); err != nil {
				fmt.Fprintf(os.Stderr, "Error notifying with %s: %v\n", n.Name(), err)
			}
		}(n)
	}
	wg.Wait()
}

func configCommand(args []string) {
	if len(args) < 1 {
		printConfigHelp()
		os.Exit(0)
	}

	switch args[0] {
	case "list":
		configFiles, err := config.ListConfigFiles()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing configuration files: %v\n", err)
			os.Exit(1)
		}
		for _, file := range configFiles {
			fmt.Println(file)
		}
	case "init":
		initCmd := flag.NewFlagSet("init", flag.ExitOnError)
		configFile := initCmd.String("config", "config.yaml", "Configuration file to create")
		force := initCmd.Bool("force", false, "Overwrite existing configuration")
		if err := initCmd.Parse(args[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing init command: %v\n", err)
			os.Exit(1)
		}

		configPath := config.GetConfigPath(*configFile)
		if _, err := os.Stat(configPath); err == nil && !*force {
			fmt.Printf("Configuration file %s already exists. Use --force to overwrite.\n", configPath)
			return
		}

		defaultConfig := config.Config{
			EnabledNotifiers: []string{"audio", "dialog"},
			DialogSettings: map[string]string{
				"title": "Notification",
			},
		}

		if err := config.Save(defaultConfig, configPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving configuration: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Created configuration file at %s\n", configPath)
	default:
		fmt.Fprintf(os.Stderr, "Unknown config subcommand: %s\n", args[0])
		printConfigHelp()
		os.Exit(1)
	}
}

func printConfigHelp() {
	fmt.Println("Usage: notify config [subcommand] [options]")
	fmt.Println("Subcommands:")
	fmt.Println("  list       List available configuration files")
	fmt.Println("  init       Initialize a new configuration file")
	fmt.Println("Options for 'init':")
	fmt.Println("  --config   Specify the configuration file name (default: config.yaml)")
	fmt.Println("  --force    Overwrite existing configuration file if it exists")
}

func notifiersCommand(args []string) {
	if len(args) > 0 && args[0] == "--help" {
		printNotifiersHelp()
		os.Exit(0)
	}

	allNotifiers := notifier.GetAllNotifiers()
	for _, n := range allNotifiers {
		fmt.Println(n.Name())
	}
}

func printNotifiersHelp() {
	fmt.Println("Usage: notify notifiers")
	fmt.Println("Description:")
	fmt.Println("  List all available notifiers")
}

func diagnoseCommand(args []string) {
	cmd := flag.NewFlagSet("diagnose", flag.ExitOnError)
	configFile := cmd.String("config", "", "Configuration file to use")
	if err := cmd.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing diagnose command: %v\n", err)
		os.Exit(1)
	}

	cfg, err := config.Load(*configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}
	diagnose.RunDiagnostic(cfg)
}

func versionCommand() {
	fmt.Println("notify version:", version)
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  notify [command] [options]")
	fmt.Println("Commands:")
	fmt.Println("  send       Send a notification")
	fmt.Println("  config     Manage configuration files")
	fmt.Println("  notifiers  Manage notification providers")
	fmt.Println("  diagnose   Run system diagnostics")
	fmt.Println("  version    Display version information")
}
