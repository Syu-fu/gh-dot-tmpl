package main

import (
	"os"
	"path/filepath"
	"testing"
)

// Helper function to create a temporary config file for testing.
func createTempConfigFile(t *testing.T, content string) (string, func()) {
	dir, err := os.MkdirTemp("", "testconfig")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	filePath := filepath.Join(dir, "config.yaml")

	if err := os.WriteFile(filePath, []byte(content), 0o600); err != nil {
		t.Fatalf("Failed to write temp config file: %v", err)
	}

	return filePath, func() { os.RemoveAll(dir) }
}

func TestLoadConfig(t *testing.T) {
	configContent := `
templates:
  template1:
    template_file: "template1.tpl"
    output_file: "output1.txt"
  template2:
    template_file: "template2.tpl"
    output_file: "output2.txt"
`
	filePath, cleanup := createTempConfigFile(t, configContent)

	defer cleanup()

	config, err := LoadConfig(filePath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(config.Templates) != 2 {
		t.Errorf("Expected 2 templates, got %d", len(config.Templates))
	}

	if config.Templates["template1"].TemplateFile != "template1.tpl" {
		t.Errorf("Expected template1 file to be 'template1.tpl', got %s", config.Templates["template1"].TemplateFile)
	}

	if config.Templates["template1"].OutputFile != "output1.txt" {
		t.Errorf("Expected template1 output to be 'output1.txt', got %s", config.Templates["template1"].OutputFile)
	}
}

func TestLoadConfigNotExist(t *testing.T) {
	_, err := LoadConfig("nonexistent.yaml")
	if err == nil {
		t.Fatalf("Expected error for nonexistent config file, got nil")
	}
}

func TestLoadConfigError(t *testing.T) {
	invalidConfigContent := `
templates:
     template1:
   template_file: "template1.tpl"
    output_file: "output1.txt
  template2:
    temp_file: "template2.tpl"
    out_file: "output2.txt"
`
	filePath, cleanup := createTempConfigFile(t, invalidConfigContent)

	defer cleanup()

	_, err := LoadConfig(filePath)
	if err == nil {
		t.Fatalf("Expected error for invalid config file, got nil")
	}
}

func TestGetConfigPath(t *testing.T) {
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	os.Setenv("XDG_CONFIG_HOME", "/tmp")

	expected := filepath.Join("/tmp", "gh-dot-tmpl", "config.yaml")
	got := GetConfigPath()

	if got != expected {
		t.Errorf("Expected config path to be %s, got %s", expected, got)
	}
}

func TestGetConfigPathNotExistXdgConfigHome(t *testing.T) {
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	os.Setenv("XDG_CONFIG_HOME", "")

	originalHOME := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHOME)

	os.Setenv("HOME", "/tmp")

	expected := filepath.Join("/tmp", ".config", "gh-dot-tmpl", "config.yaml")
	got := GetConfigPath()

	if got != expected {
		t.Errorf("Expected config path to be %s, got %s", expected, got)
	}
}
