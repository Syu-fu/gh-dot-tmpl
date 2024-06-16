package main

import (
	"fmt"
	"os"
)

func Generate(templates []string) error {
	if !IsGitRepository() {
		return fmt.Errorf("not a git repository")
	}

	gitRoot, err := GetGitRoot()
	if err != nil {
		return err
	}

	if err := os.Chdir(gitRoot); err != nil {
		return fmt.Errorf("failed to change directory to git root: %w", err)
	}

	user, repo, err := GetGithubUserRepo()
	if err != nil {
		return err
	}

	configPath := GetConfigPath()

	config, err := LoadConfig(configPath)
	if err != nil {
		return err
	}

	for _, template := range templates {
		if err := processTemplate(config, template, user, repo); err != nil {
			return err
		}
	}

	return nil
}

func processTemplate(config *Config, template, user, repo string) error {
	tempPath, err := GetTemplatePath(config, template)
	if err != nil {
		return err
	}

	outputFile := config.Templates[template].OutputFile
	if err := GenerateFileFromTemplate(tempPath, outputFile, user, repo); err != nil {
		return err
	}

	return nil
}
