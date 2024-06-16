package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config struct represents the configuration file structure.
type Config struct {
	Templates map[string]TemplateConfig `yaml:"templates"`
}

// TemplateConfig represents the mapping of template files to generated files.
type TemplateConfig struct {
	TemplateFile string `yaml:"template_file"`
	OutputFile   string `yaml:"output_file"`
}

// LoadConfig reads the configuration file and unmarshals it into a Config struct.
func LoadConfig(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("unable to open config file: %w", err)
	}
	defer file.Close()

	var config Config

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config file: %w", err)
	}

	return &config, nil
}

// GetConfigPath returns the path to the configuration file.
func GetConfigPath() string {
	configDir := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "gh-dot-tmpl")
	if os.Getenv("XDG_CONFIG_HOME") == "" {
		configDir = filepath.Join(os.Getenv("HOME"), ".config", "gh-dot-tmpl")
	}

	return filepath.Join(configDir, "config.yaml")
}
