package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// Helper function to create a temporary git repository for testing.
func setupTempGitRepo(t *testing.T) (string, func()) {
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

func TestIsGitRepository(t *testing.T) {
	dir, cleanup := setupTempGitRepo(t)
	defer cleanup()

	// Change to the temporary directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("Failed to change back to original directory: %v", err)
		}
	}()

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	if !IsGitRepository() {
		t.Errorf("Expected directory to be a git repository")
	}

	// Change to a non-git directory
	nonGitDir, err := os.MkdirTemp("", "nongit")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	defer os.RemoveAll(nonGitDir)

	if err := os.Chdir(nonGitDir); err != nil {
		t.Fatalf("Failed to change to non-git directory: %v", err)
	}

	if IsGitRepository() {
		t.Errorf("Expected directory to not be a git repository")
	}
}

func TestGetGitRoot(t *testing.T) {
	dir, cleanup := setupTempGitRepo(t)
	defer cleanup()

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("Failed to change back to original directory: %v", err)
		}
	}()

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	root, err := GetGitRoot()
	if err != nil {
		t.Errorf("Failed to get git root: %v", err)
	}

	// Evaluate symlinks to get the actual path
	expectedRoot, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatalf("Failed to evaluate symlinks: %v", err)
	}

	if root != expectedRoot {
		t.Errorf("Expected git root to be %v, got %v", expectedRoot, root)
	}
}

func TestGetGitRoot_Error(t *testing.T) {
	nonGitDir, err := os.MkdirTemp("", "nongit")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(nonGitDir)

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("Failed to change back to original directory: %v", err)
		}
	}()

	if err := os.Chdir(nonGitDir); err != nil {
		t.Fatalf("Failed to change to non-git directory: %v", err)
	}

	_, err = GetGitRoot()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	if err.Error() != "not a git repository" {
		t.Errorf("Expected error message 'not a git repository', got %v", err.Error())
	}
}

func TestGetGithubUserRepo(t *testing.T) {
	dir, cleanup := setupTempGitRepo(t)
	defer cleanup()

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("Failed to change back to original directory: %v", err)
		}
	}()

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Set a mock remote origin URL
	cmd := exec.Command("git", "remote", "add", "origin", "https://github.com/testuser/testrepo.git")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to add remote origin: %v", err)
	}

	user, repo, err := GetGithubUserRepo()
	if err != nil {
		t.Errorf("Failed to get GitHub user and repo: %v", err)
	}

	// nolint: goconst
	expectedUser := "testuser"
	// nolint: goconst
	expectedRepo := "testrepo"

	if user != expectedUser {
		t.Errorf("Expected user to be %v, got %v", expectedUser, user)
	}

	if repo != expectedRepo {
		t.Errorf("Expected repo to be %v, got %v", expectedRepo, repo)
	}
}

func TestGetGithubUserRepo_Error(t *testing.T) {
	nonGitDir, err := os.MkdirTemp("", "nongit")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(nonGitDir)

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("Failed to change back to original directory: %v", err)
		}
	}()

	if err := os.Chdir(nonGitDir); err != nil {
		t.Fatalf("Failed to change to non-git directory: %v", err)
	}

	_, _, err = GetGithubUserRepo()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	if err.Error() != "unable to get remote origin URL" {
		t.Errorf("Expected error message 'unable to get remote origin URL', got %v", err.Error())
	}
}

func TestGetGithubUserRepo_InvalidURL(t *testing.T) {
	dir, cleanup := setupTempGitRepo(t)
	defer cleanup()

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("Failed to change back to original directory: %v", err)
		}
	}()

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Set a mock invalid remote origin URL
	cmd := exec.Command("git", "remote", "add", "origin", "invalid_url")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to add remote origin: %v", err)
	}

	_, _, err = GetGithubUserRepo()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	if err.Error() != "invalid GitHub URL" {
		t.Errorf("Expected error message 'invalid GitHub URL', got %v", err.Error())
	}
}
