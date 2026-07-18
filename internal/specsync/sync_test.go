package specsync

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (function roundTripFunc) Do(request *http.Request) (*http.Response, error) {
	return function(request)
}

func TestSyncSelectionUsesManifestOrderAndRawDigest(t *testing.T) {
	t.Parallel()
	one := []byte(`<html xmlns="http://www.w3.org/1999/xhtml"><body><h2 id="one">One</h2></body></html>`)
	two := []byte(`<html xmlns="http://www.w3.org/1999/xhtml"><body><h2 id="two">Two</h2></body></html>`)
	manifest := Manifest{Documents: []Document{
		{ID: "one", URL: fixtureURL("one.html"), SHA256: digest(one), Output: "one.md"},
		{ID: "two", URL: fixtureURL("two.html"), SHA256: digest(two), Output: "two.md"},
	}}
	var requests []string
	client := roundTripFunc(func(request *http.Request) (*http.Response, error) {
		requests = append(requests, request.URL.String())
		content := map[string][]byte{manifest.Documents[0].URL: one, manifest.Documents[1].URL: two}[request.URL.String()]
		return response(http.StatusOK, content, nil), nil
	})
	root := t.TempDir()
	options := Options{CacheDir: filepath.Join(root, "cache"), OutputDir: filepath.Join(root, "out"), IDs: []string{"two", "one"}, Client: client}
	if err := Sync(context.Background(), manifest, options); err != nil {
		t.Fatalf("Sync() error = %v", err)
	}
	if got, want := strings.Join(requests, ","), manifest.Documents[0].URL+","+manifest.Documents[1].URL; got != want {
		t.Fatalf("request order = %q, want %q", got, want)
	}
	cached, err := os.ReadFile(filepath.Join(root, "cache", "one.xhtml"))
	if err != nil {
		t.Fatalf("read cached document: %v", err)
	}
	if !bytes.Equal(cached, one) {
		t.Fatalf("cached bytes = %q, want exact raw bytes %q", cached, one)
	}
	selectedMarkdown, err := os.ReadFile(filepath.Join(root, "out", "one.md"))
	if err != nil {
		t.Fatalf("read selected Markdown: %v", err)
	}

	requests = nil
	allRoot := t.TempDir()
	allOptions := Options{CacheDir: filepath.Join(allRoot, "cache"), OutputDir: filepath.Join(allRoot, "out"), Client: client}
	if err := Sync(context.Background(), manifest, allOptions); err != nil {
		t.Fatalf("Sync() all documents error = %v", err)
	}
	if got, want := strings.Join(requests, ","), manifest.Documents[0].URL+","+manifest.Documents[1].URL; got != want {
		t.Fatalf("all-document request order = %q, want %q", got, want)
	}
	allMarkdown, err := os.ReadFile(filepath.Join(allRoot, "out", "one.md"))
	if err != nil {
		t.Fatalf("read all-document Markdown: %v", err)
	}
	if !bytes.Equal(selectedMarkdown, allMarkdown) {
		t.Fatal("selected and all-document Markdown differ")
	}
}

func TestSyncErrorsCarryPerDocumentContext(t *testing.T) {
	t.Parallel()
	body := []byte(`<html xmlns="http://www.w3.org/1999/xhtml"><body/></html>`)
	document := Document{ID: "fixture", URL: fixtureURL("fixture.html"), SHA256: digest(body), Output: "fixture.md"}
	tests := []struct {
		name   string
		client HTTPClient
		want   string
	}{
		{"transport", roundTripFunc(func(*http.Request) (*http.Response, error) { return nil, errors.New("offline") }), "perform GET request: offline"},
		{"nil response", roundTripFunc(func(*http.Request) (*http.Response, error) { return nil, nil }), "nil response"},
		{"HTTP status", roundTripFunc(func(*http.Request) (*http.Response, error) { return response(http.StatusNotFound, nil, nil), nil }), "unexpected HTTP status"},
		{"media type", roundTripFunc(func(*http.Request) (*http.Response, error) {
			result := response(http.StatusOK, body, nil)
			result.Header.Set("Content-Type", "application/xml")
			return result, nil
		}), "unsupported media type"},
		{"malformed media type", roundTripFunc(func(*http.Request) (*http.Response, error) {
			result := response(http.StatusOK, body, nil)
			result.Header.Set("Content-Type", `text/html; charset="`)
			return result, nil
		}), "parse Content-Type"},
		{"read", roundTripFunc(func(*http.Request) (*http.Response, error) { return response(http.StatusOK, nil, failingReader{}), nil }), "read GET response body: fixture read"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			err := Sync(context.Background(), Manifest{Documents: []Document{document}}, Options{CacheDir: root, OutputDir: root, Client: test.client})
			if err == nil || !strings.Contains(err.Error(), `document[0] "fixture"`) || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("Sync() error = %v, want document context and %q", err, test.want)
			}
		})
	}
}

func TestNewHTTPClientRejectsRedirects(t *testing.T) {
	t.Parallel()
	client := NewHTTPClient()
	request, err := http.NewRequest(http.MethodGet, fixtureURL("redirect.html"), nil)
	if err != nil {
		t.Fatalf("construct fixture request: %v", err)
	}
	err = client.CheckRedirect(request, nil)
	if !errors.Is(err, http.ErrUseLastResponse) {
		t.Fatalf("CheckRedirect() error = %v, want http.ErrUseLastResponse", err)
	}
}

func TestSyncRejectsDigestAndParseFailuresBeforeWriting(t *testing.T) {
	t.Parallel()
	valid := []byte(`<html xmlns="http://www.w3.org/1999/xhtml"><body/></html>`)
	tests := []struct {
		name     string
		body     []byte
		expected string
		want     string
	}{
		{"digest", append(valid, '\n'), digest(valid), "verify SHA-256"},
		{"parse", []byte(`<html xmlns="http://www.w3.org/1999/xhtml"><body>`), digest([]byte(`<html xmlns="http://www.w3.org/1999/xhtml"><body>`)), "parse XHTML"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			document := Document{ID: "x", URL: fixtureURL("x.html"), SHA256: test.expected, Output: "x.md"}
			client := roundTripFunc(func(*http.Request) (*http.Response, error) { return response(http.StatusOK, test.body, nil), nil })
			err := Sync(context.Background(), Manifest{Documents: []Document{document}}, Options{CacheDir: filepath.Join(root, "cache"), OutputDir: filepath.Join(root, "out"), Client: client})
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("Sync() error = %v, want substring %q", err, test.want)
			}
			if _, statErr := os.Stat(filepath.Join(root, "cache", "x.xhtml")); !errors.Is(statErr, os.ErrNotExist) {
				t.Fatalf("cache stat error = %v, want not exist", statErr)
			}
		})
	}
}

func TestSyncRejectsUnknownSelectionDeterministically(t *testing.T) {
	t.Parallel()
	err := Sync(context.Background(), Manifest{}, Options{IDs: []string{"z", "a"}, Client: roundTripFunc(nil)})
	if err == nil || !strings.Contains(err.Error(), "unknown ids: a, z") {
		t.Fatalf("Sync() error = %v, want sorted unknown ids", err)
	}
}

func TestSyncRejectsNilContext(t *testing.T) {
	t.Parallel()
	err := Sync(nil, Manifest{}, Options{Client: roundTripFunc(nil)})
	if err == nil || !strings.Contains(err.Error(), "synchronize specifications: context must not be nil") {
		t.Fatalf("Sync() error = %v, want nil context error", err)
	}
}

func TestSyncRejectsUnsafeProgrammaticDocumentIDBeforeAcquisition(t *testing.T) {
	t.Parallel()
	client := roundTripFunc(func(*http.Request) (*http.Response, error) {
		t.Fatal("HTTP client called for unsafe document ID")
		return nil, nil
	})
	document := Document{
		ID:     "../escaped",
		URL:    fixtureURL("x.html"),
		SHA256: strings.Repeat("0", 64),
		Output: "x.md",
	}
	err := Sync(context.Background(), Manifest{Documents: []Document{document}}, Options{
		CacheDir:  filepath.Join(t.TempDir(), "cache"),
		OutputDir: filepath.Join(t.TempDir(), "output"),
		Client:    client,
	})
	if err == nil || !strings.Contains(err.Error(), `document[0] "../escaped": construct cache filename`) {
		t.Fatalf("Sync() error = %v, want unsafe ID context", err)
	}
}

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) { return 0, errors.New("fixture read") }
func (failingReader) Close() error             { return nil }

func response(status int, body []byte, reader io.ReadCloser) *http.Response {
	if reader == nil {
		reader = io.NopCloser(bytes.NewReader(body))
	}
	return &http.Response{
		StatusCode: status,
		Status:     fmt.Sprintf("%d %s", status, http.StatusText(status)),
		Header:     http.Header{"Content-Type": []string{"application/xhtml+xml; charset=utf-8"}},
		Body:       reader,
	}
}

func digest(body []byte) string {
	value := sha256.Sum256(body)
	return fmt.Sprintf("%x", value)
}

func fixtureURL(name string) string {
	return "https://www.w3.org/TR/2012/REC-fixture-20120405/" + name
}
