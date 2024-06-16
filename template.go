package main

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

// TemplateData holds the data to be inserted into the template.
type TemplateData struct {
	Username   string
	Repository string
}

const permission = 0o600

// GenerateFileFromTemplate generates a file from a template with the provided data.
func GenerateFileFromTemplate(templatePath, outputPath, username, repository string) error {
	data := TemplateData{
		Username:   username,
		Repository: repository,
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template file: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	if err := os.WriteFile(outputPath, buf.Bytes(), permission); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

// GetTemplatePath returns the full path of a template file based on the config directory.
func GetTemplatePath(config *Config, templateName string) (string, error) {
	templatePath, err := ExpandTilde(config.Templates[templateName].TemplateFile)
	if err != nil {
		return "", err
	}

	return templatePath, nil
}
