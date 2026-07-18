package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

var errWriterFailed = errors.New("writer failed")

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, errWriterFailed
}

func TestRunHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := run([]string{"help"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run help error = %v, want nil", err)
	}
	if code != 0 {
		t.Fatalf("run help exit code = %d, want 0", code)
	}
	if !strings.Contains(stdout.String(), "Usage: goxsd") {
		t.Fatalf("run help output = %q, want usage", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("run help stderr = %q, want empty", stderr.String())
	}
}

func TestRunUnknownCommand(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := run([]string{"nope"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run unknown error = %v, want nil", err)
	}
	if code != 2 {
		t.Fatalf("run unknown exit code = %d, want 2", code)
	}
	if !strings.Contains(stderr.String(), `unknown command "nope"`) {
		t.Fatalf("run unknown stderr = %q, want command context", stderr.String())
	}
}

func TestRunRejectsUnexpectedArguments(t *testing.T) {
	for _, command := range []string{"help", "version", "--help"} {
		t.Run(command, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code, err := run([]string{command, "extra"}, &stdout, &stderr)
			if err != nil {
				t.Fatalf("run %s extra error = %v, want nil", command, err)
			}
			if code != 2 {
				t.Fatalf("run %s extra exit code = %d, want 2", command, code)
			}
			if !strings.Contains(stderr.String(), `unexpected arguments: ["extra"]`) {
				t.Fatalf("run %s extra stderr = %q, want argument context", command, stderr.String())
			}
		})
	}
}

func TestRunDecoratesWriterErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "help", args: []string{"help"}, want: "write help"},
		{name: "version", args: []string{"version"}, want: "write version"},
		{name: "usage", args: []string{"nope"}, want: "write usage error for unknown command"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			code, err := run(test.args, failingWriter{}, failingWriter{})
			if code != 1 {
				t.Fatalf("run exit code = %d, want 1", code)
			}
			if !errors.Is(err, errWriterFailed) {
				t.Fatalf("run error = %v, want writer failure cause", err)
			}
			if !strings.Contains(err.Error(), test.want) {
				t.Fatalf("run error = %q, want context %q", err, test.want)
			}
		})
	}
}
