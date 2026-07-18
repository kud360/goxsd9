package main

import (
	"fmt"
	"io"
	"os"
)

const version = "dev"

func main() {
	code, err := run(os.Args[1:], os.Stdout, os.Stderr)
	if err != nil {
		if _, writeErr := fmt.Fprintf(os.Stderr, "goxsd: %v\n", err); writeErr != nil {
			code = 1
		}
	}
	os.Exit(code)
}

func run(args []string, stdout, stderr io.Writer) (int, error) {
	if len(args) == 0 {
		if err := printHelp(stdout); err != nil {
			return 1, err
		}
		return 0, nil
	}

	switch args[0] {
	case "help", "-h", "--help":
		if len(args) != 1 {
			return unexpectedArguments(args[0], args[1:], stderr)
		}
		if err := printHelp(stdout); err != nil {
			return 1, err
		}
		return 0, nil
	case "version":
		if len(args) != 1 {
			return unexpectedArguments(args[0], args[1:], stderr)
		}
		if _, err := fmt.Fprintln(stdout, version); err != nil {
			return 1, fmt.Errorf("write version: %w", err)
		}
		return 0, nil
	default:
		return usageError(stderr, "unknown command %q", args[0])
	}
}

func unexpectedArguments(command string, args []string, stderr io.Writer) (int, error) {
	return usageError(stderr, "%s: unexpected arguments: %q", command, args)
}

func usageError(w io.Writer, format string, args ...any) (int, error) {
	message := fmt.Sprintf(format, args...)
	if _, err := fmt.Fprintf(w, "goxsd: %s\nRun 'goxsd help' for usage.\n", message); err != nil {
		return 1, fmt.Errorf("write usage error for %s: %w", message, err)
	}
	return 2, nil
}

func printHelp(w io.Writer) error {
	if _, err := fmt.Fprintln(w, `Usage: goxsd <command>

Commands:
  help       Show this help
  version    Print the build version`); err != nil {
		return fmt.Errorf("write help: %w", err)
	}
	return nil
}
