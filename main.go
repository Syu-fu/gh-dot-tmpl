package main

import (
	"os"
)

func main() {
	cli := &Cli{os.Stdout, os.Stderr}
	os.Exit(cli.Run())
}
