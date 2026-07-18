package specsync

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const maxDocumentBytes = 64 << 20

// HTTPClient is the subset of http.Client used to acquire specifications.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// Options controls one specification synchronization.
type Options struct {
	CacheDir  string
	OutputDir string
	IDs       []string
	Client    HTTPClient
}

// Sync downloads, verifies, caches, and converts selected manifest documents.
// Documents are always processed in manifest order.
func Sync(ctx context.Context, manifest Manifest, options Options) error {
	if ctx == nil {
		return fmt.Errorf("synchronize specifications: context must not be nil")
	}
	if options.Client == nil {
		return fmt.Errorf("synchronize specifications: HTTP client must not be nil")
	}
	documents, err := selectDocuments(manifest.Documents, options.IDs)
	if err != nil {
		return fmt.Errorf("select specification documents: %w", err)
	}
	for index, document := range documents {
		if err := syncDocument(ctx, document, options); err != nil {
			return fmt.Errorf("synchronize specification document[%d] %q: %w", index, document.ID, err)
		}
	}
	return nil
}

func selectDocuments(documents []Document, ids []string) ([]Document, error) {
	if len(ids) == 0 {
		return slices.Clone(documents), nil
	}
	wanted := make(map[string]struct{}, len(ids))
	for index, id := range ids {
		if id == "" {
			return nil, fmt.Errorf("selection[%d]: id must not be empty", index)
		}
		if _, exists := wanted[id]; exists {
			return nil, fmt.Errorf("selection[%d] %q: duplicate id", index, id)
		}
		wanted[id] = struct{}{}
	}
	selected := make([]Document, 0, len(ids))
	for _, document := range documents {
		if _, exists := wanted[document.ID]; !exists {
			continue
		}
		selected = append(selected, document)
		delete(wanted, document.ID)
	}
	if len(wanted) != 0 {
		unknown := make([]string, 0, len(wanted))
		for id := range wanted {
			unknown = append(unknown, id)
		}
		slices.Sort(unknown)
		return nil, fmt.Errorf("unknown ids: %s", strings.Join(unknown, ", "))
	}
	return selected, nil
}

func syncDocument(ctx context.Context, document Document, options Options) error {
	cacheName, err := cacheFilename(document.ID)
	if err != nil {
		return err
	}
	raw, err := acquire(ctx, options.Client, document.URL)
	if err != nil {
		return fmt.Errorf("acquire %s: %w", document.URL, err)
	}
	actual := sha256.Sum256(raw)
	if hex.EncodeToString(actual[:]) != document.SHA256 {
		return fmt.Errorf("verify SHA-256 for %s: expected %s, got %x", document.URL, document.SHA256, actual)
	}
	markdown, err := Convert(bytes.NewReader(raw), document.URL, document.SHA256)
	if err != nil {
		return fmt.Errorf("convert %s: %w", document.URL, err)
	}
	if err := writeFile(options.CacheDir, cacheName, raw); err != nil {
		return fmt.Errorf("write cache for %s: %w", document.URL, err)
	}
	if err := writeFile(options.OutputDir, document.Output, markdown); err != nil {
		return fmt.Errorf("write Markdown for %s: %w", document.URL, err)
	}
	return nil
}

func cacheFilename(id string) (string, error) {
	if !validDocumentID(id) {
		return "", fmt.Errorf("construct cache filename for document %q: invalid document id", id)
	}
	return id + ".xhtml", nil
}

func acquire(ctx context.Context, client HTTPClient, sourceURL string) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create GET request: %w", err)
	}
	request.Header.Set("Accept", "application/xhtml+xml, text/html;q=0.9")
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("perform GET request: %w", err)
	}
	if response == nil {
		return nil, fmt.Errorf("perform GET request: HTTP client returned a nil response")
	}
	if response.Body == nil {
		return nil, fmt.Errorf("read GET response: HTTP client returned a nil body")
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		statusErr := fmt.Errorf("check GET response: unexpected HTTP status %s", response.Status)
		if err := response.Body.Close(); err != nil {
			return nil, errors.Join(statusErr, fmt.Errorf("close GET response body: %w", err))
		}
		return nil, statusErr
	}
	mediaType, _, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if err != nil {
		mediaErr := fmt.Errorf("check GET response: parse Content-Type: %w", err)
		if closeErr := response.Body.Close(); closeErr != nil {
			return nil, errors.Join(mediaErr, fmt.Errorf("close GET response body: %w", closeErr))
		}
		return nil, mediaErr
	}
	mediaType = strings.ToLower(mediaType)
	if mediaType != "application/xhtml+xml" && mediaType != "text/html" {
		mediaErr := fmt.Errorf("check GET response: unsupported media type %q", mediaType)
		if err := response.Body.Close(); err != nil {
			return nil, errors.Join(mediaErr, fmt.Errorf("close GET response body: %w", err))
		}
		return nil, mediaErr
	}
	limited := io.LimitReader(response.Body, maxDocumentBytes+1)
	raw, readErr := io.ReadAll(limited)
	closeErr := response.Body.Close()
	if readErr != nil || closeErr != nil {
		var contextual []error
		if readErr != nil {
			contextual = append(contextual, fmt.Errorf("read GET response body: %w", readErr))
		}
		if closeErr != nil {
			contextual = append(contextual, fmt.Errorf("close GET response body: %w", closeErr))
		}
		return nil, errors.Join(contextual...)
	}
	if len(raw) > maxDocumentBytes {
		return nil, fmt.Errorf("read GET response body: document exceeds %d bytes", maxDocumentBytes)
	}
	return raw, nil
}

func writeFile(directory, name string, content []byte) (returnErr error) {
	if directory == "" {
		return fmt.Errorf("destination directory must not be empty")
	}
	if name == "" || filepath.Base(name) != name {
		return fmt.Errorf("destination name %q must be a basename", name)
	}
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return fmt.Errorf("create directory %q: %w", directory, err)
	}
	temporary, err := os.CreateTemp(directory, ".specsync-*")
	if err != nil {
		return fmt.Errorf("create temporary file in %q: %w", directory, err)
	}
	temporaryName := temporary.Name()
	remove := true
	defer func() {
		if remove {
			if err := os.Remove(temporaryName); err != nil && !errors.Is(err, os.ErrNotExist) {
				returnErr = errors.Join(returnErr, fmt.Errorf("remove temporary file %q: %w", temporaryName, err))
			}
		}
	}()
	if _, err := temporary.Write(content); err != nil {
		writeErr := fmt.Errorf("write temporary file %q: %w", temporaryName, err)
		if closeErr := temporary.Close(); closeErr != nil {
			return errors.Join(writeErr, fmt.Errorf("close temporary file %q: %w", temporaryName, closeErr))
		}
		return writeErr
	}
	if err := temporary.Chmod(0o644); err != nil {
		modeErr := fmt.Errorf("set mode on temporary file %q: %w", temporaryName, err)
		if closeErr := temporary.Close(); closeErr != nil {
			return errors.Join(modeErr, fmt.Errorf("close temporary file %q: %w", temporaryName, closeErr))
		}
		return modeErr
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close temporary file %q: %w", temporaryName, err)
	}
	destination := filepath.Join(directory, name)
	if err := os.Rename(temporaryName, destination); err != nil {
		return fmt.Errorf("replace %q: %w", destination, err)
	}
	remove = false
	return nil
}

// NewHTTPClient returns a client that rejects redirects so the manifest URL is
// the exact representation whose bytes are verified.
func NewHTTPClient() *http.Client {
	return &http.Client{CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	}}
}
