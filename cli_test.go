package main

import (
	"bytes"
	"os"
	"testing"
)

func TestParseArgs_Help(t *testing.T) {
	tests := []struct {
		args     []string
		expected CliArgs
	}{
		{args: []string{"gh-dot-tmpl", "-h"}, expected: CliArgs{ShowHelp: true}},
		{args: []string{"gh-dot-tmpl", "--help"}, expected: CliArgs{ShowHelp: true}},
	}

	for _, test := range tests {
		cliArgs, err := ParseArgs(test.args[1:]) // Skip the program name
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if cliArgs.ShowHelp != test.expected.ShowHelp {
			t.Errorf("Expected ShowHelp %v, got %v", test.expected.ShowHelp, cliArgs.ShowHelp)
		}
	}
}

func TestParseArgs_Version(t *testing.T) {
	tests := []struct {
		args     []string
		expected CliArgs
	}{
		{args: []string{"gh-dot-tmpl", "-v"}, expected: CliArgs{ShowVersion: true}},
		{args: []string{"gh-dot-tmpl", "--version"}, expected: CliArgs{ShowVersion: true}},
	}

	for _, test := range tests {
		cliArgs, err := ParseArgs(test.args[1:]) // Skip the program name
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if cliArgs.ShowVersion != test.expected.ShowVersion {
			t.Errorf("Expected ShowVersion %v, got %v", test.expected.ShowVersion, cliArgs.ShowVersion)
		}
	}
}

func TestParseArgs_Templates(t *testing.T) {
	tests := []struct {
		args     []string
		expected CliArgs
	}{
		{args: []string{"gh-dot-tmpl", "template1", "template2"}, expected: CliArgs{Templates: []string{"template1", "template2"}}},
	}

	for _, test := range tests {
		cliArgs, err := ParseArgs(test.args[1:]) // Skip the program name
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(cliArgs.Templates) != len(test.expected.Templates) {
			t.Errorf("Expected %d templates, got %d", len(test.expected.Templates), len(cliArgs.Templates))
		}

		for i, template := range cliArgs.Templates {
			if template != test.expected.Templates[i] {
				t.Errorf("Expected template %q, got %q", test.expected.Templates[i], template)
			}
		}
	}
}

func TestCli_Run_Help(t *testing.T) {
	cli := &Cli{}

	// Capture the output
	out := new(bytes.Buffer)
	cli.OutStream = out

	// Simulate the command-line arguments
	oldArgs := os.Args
	os.Args = []string{"gh-dot-tmpl", "-h"}

	defer func() { os.Args = oldArgs }()

	// Run the CLI
	exitCode := cli.Run()

	if exitCode != 0 {
		t.Fatalf("Expected exit code 0, got %d", exitCode)
	}

	expected := `Usage: gh-dot-tmpl [options] [template_name...]

Options:
  -h, --help       Show help message
  -v, --version    Show version

Arguments:
  template_name...  Names of the templates to process
`
	if out.String() != expected {
		t.Errorf("Expected output %q, got %q", expected, out.String())
	}
}

func TestCli_Run_Version(t *testing.T) {
	cli := &Cli{}

	// Capture the output
	out := new(bytes.Buffer)
	cli.OutStream = out

	// Simulate the command-line arguments
	oldArgs := os.Args
	os.Args = []string{"gh-dot-tmpl", "-v"}

	defer func() { os.Args = oldArgs }()

	// Run the CLI
	exitCode := cli.Run()

	if exitCode != 0 {
		t.Fatalf("Expected exit code 0, got %d", exitCode)
	}

	expected := "gh-dot-tmpl version " + version + "\n"
	if out.String() != expected {
		t.Errorf("Expected output %q, got %q", expected, out.String())
	}
}

func TestCli_Run_NoArgs(t *testing.T) {
	cli := &Cli{}

	// Capture the output
	out := new(bytes.Buffer)
	errOut := new(bytes.Buffer)
	cli.OutStream = out
	cli.ErrStream = errOut

	// Simulate the command-line arguments
	oldArgs := os.Args
	os.Args = []string{"gh-dot-tmpl"}

	defer func() { os.Args = oldArgs }()

	// Run the CLI
	exitCode := cli.Run()

	if exitCode != 1 {
		t.Fatalf("Expected exit code 1, got %d", exitCode)
	}

	expectedError := "Error: No template names provided\n"
	if errOut.String() != expectedError {
		t.Errorf("Expected error output %q, got %q", expectedError, errOut.String())
	}

	expectedUsage := `Usage: gh-dot-tmpl [options] [template_name...]

Options:
  -h, --help       Show help message
  -v, --version    Show version

Arguments:
  template_name...  Names of the templates to process
`
	if out.String() != expectedUsage {
		t.Errorf("Expected output %q, got %q", expectedUsage, out.String())
	}
}

func TestCli_Run_InvalidFlag(t *testing.T) {
	cli := &Cli{}

	// Capture the output
	errOut := new(bytes.Buffer)
	cli.ErrStream = errOut

	// Simulate the command-line arguments
	oldArgs := os.Args
	os.Args = []string{"gh-dot-tmpl", "--invalidflag"}

	defer func() { os.Args = oldArgs }()

	// Run the CLI
	exitCode := cli.Run()

	if exitCode != 1 {
		t.Fatalf("Expected exit code 1, got %d", exitCode)
	}

	expectedError := "flag provided but not defined: -invalidflag\n"
	if errOut.String() != expectedError {
		t.Errorf("Expected error output %q, got %q", expectedError, errOut.String())
	}
}
