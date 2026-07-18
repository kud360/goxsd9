package main

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func TestReadManifestPreservesReadAndCloseFailures(t *testing.T) {
	t.Parallel()
	readCause := errors.New("fixture read failure")
	closeCause := errors.New("fixture close failure")
	opener := func(string) (io.ReadCloser, error) {
		return &failingReadCloser{readErr: readCause, closeErr: closeCause}, nil
	}

	_, err := readManifest("fixture.json", opener)
	if !errors.Is(err, readCause) {
		t.Fatalf("readManifest() error = %v, want read cause", err)
	}
	if !errors.Is(err, closeCause) {
		t.Fatalf("readManifest() error = %v, want close cause", err)
	}
	for _, want := range []string{`read manifest "fixture.json"`, `close manifest "fixture.json"`} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("readManifest() error = %v, want context %q", err, want)
		}
	}
}

type failingReadCloser struct {
	readErr  error
	closeErr error
}

func (reader *failingReadCloser) Read([]byte) (int, error) {
	return 0, reader.readErr
}

func (reader *failingReadCloser) Close() error {
	return reader.closeErr
}
