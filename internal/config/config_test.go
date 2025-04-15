package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadDefaultConfig(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir) // Override home directory for testing
	defer os.Unsetenv("HOME")

	cfg, err := Load("")
	assert.NoError(t, err)
	assert.Contains(t, cfg.EnabledNotifiers, "audio")
	assert.Contains(t, cfg.EnabledNotifiers, "dialog")
}

func TestSaveAndLoadConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test_config.yaml")

	cfg := Config{
		EnabledNotifiers: []string{"audio"},
		DialogSettings:   map[string]string{"title": "Test"},
	}

	err := Save(cfg, configPath)
	assert.NoError(t, err)

	loadedCfg, err := Load(configPath)
	assert.NoError(t, err)
	assert.Equal(t, cfg, loadedCfg)
}

func TestListConfigFiles(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir) // Override home directory for testing
	defer os.Unsetenv("HOME")

	configDir := GetConfigDir()
	os.MkdirAll(configDir, 0755)
	os.WriteFile(filepath.Join(configDir, "config1.yaml"), []byte{}, 0644)
	os.WriteFile(filepath.Join(configDir, "config2.yaml"), []byte{}, 0644)

	files, err := ListConfigFiles()
	assert.NoError(t, err)
	assert.Len(t, files, 2)
	assert.Contains(t, files, "config1.yaml")
	assert.Contains(t, files, "config2.yaml")
}
