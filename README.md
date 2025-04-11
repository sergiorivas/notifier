# Notify - CLI Notification System

Notify is a command-line tool that sends notifications through various means.

## Features

- Support for multiple notification types (audio, dialogs)
- Customizable title and message
- Different notification severities (success, error, info, warning)
- Multiple configuration profiles
- Diagnostic tools for notifier availability

## Installation

go install github.com/yourusername/notify@latest

## Usage

### Send a notification

notify --type success "Build completed"
notify --type error --title "Error" "Build failed"

### Create default configuration with only dialog notifier

notify init
notify init --config-file dialog-only.yaml
notify init --force # Overwrite existing configuration

### Use alternative configurations

notify --config-file work.yaml "Work notification"
notify --config-file silent.yaml "Silent notification"


### List available notifier types

notify list-notifiers
notify list-notifiers --verbose # Show detailed information

### List available configurations

notify list-configs

### Run diagnostics

notify diagnose

### Options

- `--type`: Notification type (success, error, info, warning)
- `--title`: Custom title for the notification
- `--config-file`: Use a specific configuration file from ~/.config/notify/
- `--list-configs`: List all available configuration files

## Configuration

Configurations are stored in `~/.config/notify/` directory. The default configuration file is `config.yaml`.

You can create multiple configuration files for different scenarios:

- `work.yaml` - Configuration for work-related notifications
- `silent.yaml` - Configuration with just visual notifications, no audio
- `urgent.yaml` - Configuration for critical notifications

Each configuration file follows this format:

```yaml
enabledNotifiers:
  - audio
  - dialog
dialogSettings:
  title: Notification
awk

## Add: Installation Script (`install.sh`)

```bash
#!/bin/bash

# Create configuration directory
mkdir -p ~/.config/notify

# Copy sample configurations
cp -n configs/*.yaml ~/.config/notify/

# Build and install
go build -o build/notify ./cmd/notify
sudo cp build/notify /usr/local/bin/

echo "Notify installed successfully!"
echo "Configuration files are in ~/.config/notify/"
echo "Run 'notify --list-configs' to see available configurations"
