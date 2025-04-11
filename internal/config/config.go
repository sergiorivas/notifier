package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the user configuration
type Config struct {
	EnabledNotifiers []string          `yaml:"enabledNotifiers"`
	DialogSettings   map[string]string `yaml:"dialogSettings"`
}

// GetConfigDir returns the directory where all configurations are stored
func GetConfigDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "notify"
	}
	return filepath.Join(homeDir, ".config", "notify")
}

// GetConfigPath returns the path to the configuration file
func GetConfigPath(configFile string) string {
	configDir := GetConfigDir()

	if configFile != "" {
		// If a specific file is provided, look for it in the config directory
		// unless it's already an absolute path
		if filepath.IsAbs(configFile) {
			return configFile
		}
		// Use the file within the configuration directory
		return filepath.Join(configDir, configFile)
	}

	// Default configuration file is config.yaml in config directory
	return filepath.Join(configDir, "config.yaml")
}

// Load loads the configuration from the file
func Load(configFile string) (Config, error) {
	var config Config
	configPath := GetConfigPath(configFile)
	configDir := filepath.Dir(configPath)

	// Create default configuration if it doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config = Config{
			EnabledNotifiers: []string{"audio", "dialog"},
			DialogSettings: map[string]string{
				"title": "Notification",
			},
		}

		// Make sure the config directory exists
		err := os.MkdirAll(configDir, 0755)
		if err != nil {
			return config, err
		}

		// Save default configuration
		return config, Save(config, configPath)
	}

	// Read existing configuration
	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	return config, err
}

// Save saves the configuration to the specified file
func Save(config Config, configPath string) error {
	configYAML, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, configYAML, 0644)
}

// ListConfigFiles returns a list of all config files in the config directory
func ListConfigFiles() ([]string, error) {
	configDir := GetConfigDir()

	// Ensure directory exists
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return nil, err
		}
		return []string{}, nil
	}

	entries, err := os.ReadDir(configDir)
	if err != nil {
		return nil, err
	}

	var configFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && (filepath.Ext(entry.Name()) == ".yaml" || filepath.Ext(entry.Name()) == ".yml") {
			configFiles = append(configFiles, entry.Name())
		}
	}

	return configFiles, nil
}
