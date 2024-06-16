package main

import (
	"errors"
	"os/exec"
	"strings"
)

// Number of parts expected in the GitHub URL.
const expectedGithubURLParts = 2

// IsGitRepository checks if the current directory is a git repository.
func IsGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Stderr = nil
	err := cmd.Run()

	return err == nil
}

// GetGitRoot returns the root path of the git repository.
func GetGitRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")

	output, err := cmd.Output()
	if err != nil {
		return "", errors.New("not a git repository")
	}

	return strings.TrimSpace(string(output)), nil
}

// GetGithubUserRepo returns the GitHub username and repository name.
func GetGithubUserRepo() (string, string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")

	output, err := cmd.Output()
	if err != nil {
		return "", "", errors.New("unable to get remote origin URL")
	}

	url := strings.TrimSpace(string(output))
	url = strings.TrimPrefix(url, "https://github.com/")
	url = strings.TrimSuffix(url, ".git")

	parts := strings.Split(url, "/")

	if len(parts) != expectedGithubURLParts {
		return "", "", errors.New("invalid GitHub URL")
	}

	return parts[0], parts[1], nil
}
