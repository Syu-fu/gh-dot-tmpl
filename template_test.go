package main

import (
	"os"
	"path/filepath"
	"testing"
)

// Helper function to create a temporary template file for testing.
func createTempTemplateFile(t *testing.T, dir, content string) (string, func()) {
	filePath := filepath.Join(dir, "template.tpl")

	if err := os.WriteFile(filePath, []byte(content), 0o600); err != nil {
		t.Fatalf("Failed to write temp template file: %v", err)
	}

	return filePath, func() { os.RemoveAll(filePath) }
}

// Helper function to create a temporary config file for testing.
func createTempConfigFile2(t *testing.T) (string, func()) {
	dir, err := os.MkdirTemp("", "testconfig")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	filePath := filepath.Join(dir, "config.yaml")

	if err := os.WriteFile(filePath, []byte(""), 0o600); err != nil {
		t.Fatalf("Failed to write temp config file: %v", err)
	}

	return dir, func() { os.RemoveAll(dir) }
}

func TestGenerateFileFromTemplate(t *testing.T) {
	templateContent := `User: {{.Username}}, Repo: {{.Repository}}`
	tempDir, cleanup := createTempConfigFile2(t)

	defer cleanup()

	templatePath, cleanupTemplate := createTempTemplateFile(t, tempDir, templateContent)
	defer cleanupTemplate()

	outputDir, err := os.MkdirTemp("", "testoutput")
	if err != nil {
		t.Fatalf("Failed to create temp output directory: %v", err)
	}
	defer os.RemoveAll(outputDir)

	outputPath := filepath.Join(outputDir, "output.txt")
	user := "testuser"
	repo := "testrepo"

	if err := GenerateFileFromTemplate(templatePath, outputPath, user, repo); err != nil {
		t.Fatalf("Failed to generate file from template: %v", err)
	}

	outputContent, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expectedContent := "User: testuser, Repo: testrepo"
	if string(outputContent) != expectedContent {
		t.Errorf("Expected output to be %s, got %s", expectedContent, string(outputContent))
	}
}

func TestGenerateFileFromTemplate_ParseError(t *testing.T) {
	// invalid template content
	invalidTemplateContent := `{{.user} {{.repo}}`
	tempDir, cleanup := createTempConfigFile2(t)

	defer cleanup()

	templatePath, cleanupTemplate := createTempTemplateFile(t, tempDir, invalidTemplateContent)
	defer cleanupTemplate()

	outputDir, err := os.MkdirTemp("", "testoutput")
	if err != nil {
		t.Fatalf("Failed to create temp output directory: %v", err)
	}
	defer os.RemoveAll(outputDir)

	outputPath := filepath.Join(outputDir, "output.txt")
	user := "testuser"
	repo := "testrepo"

	// err = GenerateFileFromTemplate(templatePath, outputPath, user, repo)
	if err = GenerateFileFromTemplate(templatePath, outputPath, user, repo); err == nil {
		t.Fatalf("Expected parse error, got %v", err)
	}
}

func TestGenerateFileFromTemplate_ExecuteError(t *testing.T) {
	validTemplateContent := `User: {{.Usernamea}}, Repo: {{.Repository}}`
	tempDir, cleanup := createTempConfigFile2(t)

	defer cleanup()

	templatePath, cleanupTemplate := createTempTemplateFile(t, tempDir, validTemplateContent)
	defer cleanupTemplate()

	outputDir, err := os.MkdirTemp("", "testoutput")
	if err != nil {
		t.Fatalf("Failed to create temp output directory: %v", err)
	}
	defer os.RemoveAll(outputDir)

	outputPath := filepath.Join(outputDir, "output.txt")

	// Missing username and repository fields should cause execute error
	if err = GenerateFileFromTemplate(templatePath, outputPath, "", ""); err == nil {
		t.Fatalf("Expected execute error, got %v", err)
	}
}

func TestGenerateFileFromTemplate_WriteError(t *testing.T) {
	validTemplateContent := `User: {{.Username}}, Repo: {{.Repository}}`
	tempDir, cleanup := createTempConfigFile2(t)

	defer cleanup()

	templatePath, cleanupTemplate := createTempTemplateFile(t, tempDir, validTemplateContent)
	defer cleanupTemplate()

	// Make the output directory read-only to cause write error
	outputDir, err := os.MkdirTemp("", "testoutput")
	if err != nil {
		t.Fatalf("Failed to create temp output directory: %v", err)
	}
	defer os.RemoveAll(outputDir)

	// Change permission to read-only
	if err := os.Chmod(outputDir, 0o400); err != nil {
		t.Fatalf("Failed to change permission: %v", err)
	}

	outputPath := filepath.Join(outputDir, "output.txt")
	user := "testuser"
	repo := "testrepo"

	if err := GenerateFileFromTemplate(templatePath, outputPath, user, repo); err == nil {
		t.Fatalf("Expected write permission error, got %v", err)
	}
}

func TestGetTemplatePath(t *testing.T) {
	templateName := "template"
	expected := "template.tpl"

	config := &Config{
		Templates: map[string]TemplateConfig{
			templateName: {
				TemplateFile: expected,
				OutputFile:   "~/output.txt",
			},
		},
	}

	got, err := GetTemplatePath(config, templateName)
	if err != nil {
		t.Fatalf("Failed to get template path: %v", err)
	}

	if got != expected {
		t.Errorf("Expected template path to be %s, got %s", expected, got)
	}
}
