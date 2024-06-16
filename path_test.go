package main

import (
	"fmt"
	"os"
	"os/user"
	"strings"
	"testing"
)

// Mock user lookup function.
var mockUserLookup = func(username string) (*user.User, error) {
	if username == "testuser" {
		return &user.User{
			Username: "testuser",
			HomeDir:  "/home/testuser",
		}, nil
	}
	return nil, fmt.Errorf("user not found")
}

func TestExpandTilde(t *testing.T) {
	// Save the original lookupUser function and HOME environment variable.
	originalLookupUser := lookupUser
	originalHome := os.Getenv("HOME")

	// Mock lookupUser function and set HOME environment variable.
	lookupUser = mockUserLookup

	os.Setenv("HOME", "/mock/home")

	// Restore the original lookupUser function and HOME environment variable after tests.
	defer func() {
		lookupUser = originalLookupUser

		os.Setenv("HOME", originalHome)
	}()

	testCases := []struct {
		input    string
		expected string
		errMsg   string
	}{
		{"", "", "No Path provided"},
		{"~", "/mock/home", ""},
		{"~/test/path", "/mock/home/test/path", ""},
		{"~testuser/test/path", "/home/testuser/test/path", ""},
		{"~nosuchuser/test/path", "", "user not found"},
		{"/no/tilde/path", "/no/tilde/path", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			output, err := ExpandTilde(tc.input)
			if err != nil {
				if tc.errMsg == "" || !strings.Contains(err.Error(), tc.errMsg) {
					t.Errorf("expected error %v, got %v", tc.errMsg, err)
				}
			} else {
				if output != tc.expected {
					t.Errorf("expected %v, got %v", tc.expected, output)
				}
			}
		})
	}
}
