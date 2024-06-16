package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// version is the current version of the CLI.
var version string

// Cli represents the command-line interface.
type Cli struct {
	OutStream, ErrStream io.Writer
}

// CliArgs holds the parsed command-line arguments.
type CliArgs struct {
	ShowHelp    bool
	ShowVersion bool
	Templates   []string
}

// ParseArgs parses command-line arguments.
func ParseArgs(args []string) (CliArgs, error) {
	var cliArgs CliArgs

	flags := flag.NewFlagSet("gh-dot-tmpl", flag.ContinueOnError)
	flags.BoolVar(&cliArgs.ShowHelp, "h", false, "Show help message")
	flags.BoolVar(&cliArgs.ShowHelp, "help", false, "Show help message")
	flags.BoolVar(&cliArgs.ShowVersion, "v", false, "Show version")
	flags.BoolVar(&cliArgs.ShowVersion, "version", false, "Show version")

	if err := flags.Parse(args); err != nil {
		// nolint: wrapcheck
		return cliArgs, err
	}

	cliArgs.Templates = flags.Args()

	return cliArgs, nil
}

// Run parses command-line arguments and executes the appropriate action.
func (cli *Cli) Run() int {
	cliArgs, err := ParseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintln(cli.ErrStream, err)
		return 1
	}

	if cliArgs.ShowHelp {
		cli.usage()
		return 0
	}

	if cliArgs.ShowVersion {
		fmt.Fprintf(cli.OutStream, "gh-dot-tmpl version %s\n", version)
		return 0
	}

	if len(cliArgs.Templates) == 0 {
		fmt.Fprintf(cli.ErrStream, "Error: No template names provided\n")
		cli.usage()

		return 1
	}

	templateName := cliArgs.Templates

	err = Generate(templateName)
	if err != nil {
		fmt.Fprintf(cli.ErrStream, "Error: %s\n", err)
		return 1
	}

	return 0
}

// usage prints the help message.
func (cli *Cli) usage() {
	fmt.Fprintf(cli.OutStream, `Usage: gh-dot-tmpl [options] [template_name...]

Options:
  -h, --help       Show help message
  -v, --version    Show version

Arguments:
  template_name...  Names of the templates to process
`)
}
