package main

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"
)

// Helper function to create a temporary Git repository.
func setupTempGitRepoGenerate(t *testing.T) (string, func()) {
	dir, err := os.MkdirTemp("", "testrepo")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = dir

	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	return dir, func() { os.RemoveAll(dir) }
}

// Helper function to create a temporary config file.
func createTempConfigFileGenerate(t *testing.T, dir, content string) {
	configDir := filepath.Join(dir, ".config", "gh-dot-tmpl")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	configPath := filepath.Join(configDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(content), 0o600); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}
}

// Helper function to create a temporary template file.
func createTempTemplateFileGenerate(t *testing.T, dir, template, content string) string {
	templatePath := filepath.Join(dir, template)
	if err := os.WriteFile(templatePath, []byte(content), 0o600); err != nil {
		t.Fatalf("Failed to write template file: %v", err)
	}

	return templatePath
}

func TestGenerate(t *testing.T) {
	dir, cleanup := setupTempGitRepoGenerate(t)
	defer cleanup()

	// Set up the environment for the test
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// nolint: errcheck
	defer os.Chdir(originalDir)

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Set up a mock GitHub remote
	cmd := exec.Command("git", "remote", "add", "origin", "https://github.com/testuser/testrepo.git")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to add remote origin: %v", err)
	}

	// Set up the configuration file
	templatePath := path.Join(dir, "template1.tpl")
	configContent := `
templates:
  template1:
    template_file: ` + templatePath + `
    output_file: "output1.txt"
`
	createTempConfigFileGenerate(t, dir, configContent)

	// Set the environment variables
	os.Setenv("HOME", dir)
	os.Setenv("XDG_CONFIG_HOME", "")

	// Create a template file
	templateContent := "Template: template1.tpl, User: {{.Username}}, Repo: {{.Repository}}"
	createTempTemplateFileGenerate(t, dir, "template1.tpl", templateContent)

	// Run the Generate function
	templates := []string{"template1"}

	err = Generate(templates)
	if err != nil {
		t.Fatalf("Generate function failed: %v", err)
	}

	// Verify the generated file
	outputPath := filepath.Join(dir, "output1.txt")

	outputContent, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	expectedContent := "Template: template1.tpl, User: testuser, Repo: testrepo"
	if string(outputContent) != expectedContent {
		t.Errorf("Expected generated file content to be %s, got %s", expectedContent, string(outputContent))
	}
}

func TestGenerateNotGitRepo(t *testing.T) {
	// Test when not a git repository
	dir, err := os.MkdirTemp("", "testrepo")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// nolint: errcheck
	defer os.Chdir(originalDir)

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	err = Generate([]string{"template1"})
	if err == nil || err.Error() != "not a git repository" {
		t.Errorf("Expected 'not a git repository' error, got %v", err)
		t.Errorf(GetGitRoot())
	}
}

func TestGenerateErrors(t *testing.T) {
	// Test when failing to get git root
	dir, cleanup := setupTempGitRepoGenerate(t)
	defer cleanup()

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// nolint: errcheck
	defer os.Chdir(originalDir)

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Test when failing to get GitHub user/repo
	cmd := exec.Command("git", "remote", "add", "origin", "invalid-url")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to add remote origin: %v", err)
	}

	err = Generate([]string{"template1"})
	if err == nil || err.Error() == "" {
		t.Errorf("Expected error due to invalid GitHub URL, got %v", err)
	}
}

// func TestGenerate(t *testing.T) {
// 	dir, cleanup := setupTempGitRepoGenerate(t)
// 	defer cleanup()
//
// 	originalDir, err := os.Getwd()
// 	if err != nil {
// 		t.Fatalf("Failed to get current directory: %v", err)
// 	}
// 	defer os.Chdir(originalDir)
//
// 	if err := os.Chdir(dir); err != nil {
// 		t.Fatalf("Failed to change directory: %v", err)
// 	}
//
// Set up the correct remote URL
// }

func TestGenerateConfigErrors(t *testing.T) {
	dir, cleanup := setupTempGitRepoGenerate(t)
	defer cleanup()

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	// nolint: errcheck
	defer os.Chdir(originalDir)

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	os.Setenv("XDG_CONFIG_HOME", dir)

	// Test when config loading fails
	invalidConfigContent := `
templates:
  template1:
    template_file: 
    output_file: 
`
	createTempConfigFileGenerate(t, dir, invalidConfigContent)

	err = Generate([]string{"template1"})
	if err == nil || err.Error() == "" {
		t.Errorf("Expected error due to invalid config, got %v", err)
	}
}

func TestGenerateTemplateErrors(t *testing.T) {
	dir, cleanup := setupTempGitRepoGenerate(t)
	defer cleanup()

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	// nolint: errcheck
	defer os.Chdir(originalDir)

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	os.Setenv("XDG_CONFIG_HOME", path.Join(dir, ".config"))

	cmd := exec.Command("git", "remote", "add", "origin", "https://github.com/testuser/testrepo.git")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to set remote URL: %v", err)
	}

	// Test when template not found in config
	validConfigContent := `
templates:
  template2:
    template_file: template2.tpl
    output_file: output2.txt
`
	createTempConfigFileGenerate(t, dir, validConfigContent)

	err = Generate([]string{"template1"})
	if err == nil || err.Error() != "No Path provided" {
		t.Errorf("Expected 'No Path provided' error, got %v", err)
	}
}
