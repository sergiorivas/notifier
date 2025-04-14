# Notify - CLI Notification System

Notify is a versatile command-line tool designed to streamline the process of sending notifications. It supports multiple notification methods, including audio alerts and dialog boxes, making it suitable for a wide range of use cases. Whether you need to deliver success messages, error alerts, or informational updates, Notify provides a customizable and efficient solution.

## Features

- Support for multiple notification types (audio, dialogs)
- Customizable title and message
- Different notification severities (success, error, info, warning)
- Multiple configuration profiles
- Diagnostic tools for notifier availability

## Installation

### Using Homebrew

```bash
brew install sergiorivas/tap/notify
```

### Using Go

```bash
go install github.com/sergiorivas/cmd/notify@latest
```

## Usage

### Send a notification

```bash
notify send --type success "Build completed"
notify send --type error --title "Error" "Build failed"
```

### Manage configurations

#### List available configurations

```bash
notify config list
```

#### Create a default configuration with only dialog notifier

```bash
notify config init
notify config init --config dialog-only.yaml
notify config init --force # Overwrite existing configuration
```

### List available notifier types

```bash
notify notifiers
```

### Run diagnostics

```bash
notify diagnose
```

### Options

- `--type`: Notification type (success, error, info, warning)
- `--title`: Custom title for the notification
- `--config`: Use a specific configuration file from ~/.config/notify/

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
```
