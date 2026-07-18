// Package specsync acquires pinned XHTML specifications and converts them to
// deterministic Markdown for local research.
package specsync

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

// Manifest describes the documents owned by the specification sync tool.
type Manifest struct {
	Documents []Document `json:"documents"`
}

// Document is one immutable XHTML representation and its Markdown filename.
type Document struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	SHA256 string `json:"sha256"`
	Output string `json:"output"`
}

// ReadManifest reads and validates a manifest.
func ReadManifest(r io.Reader) (Manifest, error) {
	if r == nil {
		return Manifest{}, fmt.Errorf("decode specification manifest: reader must not be nil")
	}
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()
	var manifest Manifest
	if err := decoder.Decode(&manifest); err != nil {
		return Manifest{}, fmt.Errorf("decode specification manifest: %w", err)
	}
	if err := requireJSONEnd(decoder); err != nil {
		return Manifest{}, err
	}
	if len(manifest.Documents) == 0 {
		return Manifest{}, fmt.Errorf("validate specification manifest: documents must not be empty")
	}
	ids := make(map[string]struct{}, len(manifest.Documents))
	outputs := make(map[string]struct{}, len(manifest.Documents))
	for index, document := range manifest.Documents {
		if err := validateDocument(document, ids, outputs); err != nil {
			return Manifest{}, fmt.Errorf("validate specification manifest document[%d] %q: %w", index, document.ID, err)
		}
		ids[document.ID] = struct{}{}
		outputs[document.Output] = struct{}{}
	}
	return manifest, nil
}

func requireJSONEnd(decoder *json.Decoder) error {
	var trailing any
	if err := decoder.Decode(&trailing); err == io.EOF {
		return nil
	} else if err != nil {
		return fmt.Errorf("decode trailing specification manifest data: %w", err)
	}
	return fmt.Errorf("decode specification manifest: unexpected trailing JSON value")
}

func validateDocument(document Document, ids, outputs map[string]struct{}) error {
	if !validDocumentID(document.ID) {
		return fmt.Errorf("id must use lowercase ASCII letters, digits, and single interior hyphens")
	}
	if _, exists := ids[document.ID]; exists {
		return fmt.Errorf("duplicate id")
	}
	parsed, err := url.Parse(document.URL)
	if err != nil {
		return fmt.Errorf("parse representation URL %q: %w", document.URL, err)
	}
	if parsed.Scheme != "https" || parsed.Hostname() != "www.w3.org" || parsed.Port() != "" || parsed.User != nil || parsed.Fragment != "" || parsed.RawQuery != "" || parsed.RawPath != "" {
		return fmt.Errorf("representation URL %q must be an official HTTPS www.w3.org URL without credentials, port, query, or fragment", document.URL)
	}
	year, publicationDate, ok := datedW3CRepresentationPath(parsed.Path)
	if !ok {
		return fmt.Errorf("representation URL %q must identify a dated W3C REC or NOTE XHTML representation", document.URL)
	}
	date, err := time.Parse("20060102", publicationDate)
	if err != nil || date.Format("2006") != year {
		return fmt.Errorf("representation URL %q must contain a valid publication date matching its year directory", document.URL)
	}
	digest, err := hex.DecodeString(document.SHA256)
	if err != nil || len(digest) != 32 || strings.ToLower(document.SHA256) != document.SHA256 {
		return fmt.Errorf("sha256 must be exactly 64 lowercase hexadecimal characters")
	}
	if filepath.Base(document.Output) != document.Output || !strings.HasSuffix(document.Output, ".md") {
		return fmt.Errorf("output %q must be a Markdown basename", document.Output)
	}
	if _, exists := outputs[document.Output]; exists {
		return fmt.Errorf("duplicate output %q", document.Output)
	}
	return nil
}

func datedW3CRepresentationPath(path string) (string, string, bool) {
	parts := strings.Split(path, "/")
	if len(parts) != 5 || parts[0] != "" || parts[1] != "TR" || len(parts[2]) != 4 {
		return "", "", false
	}
	if !asciiDigits(parts[2]) || !strings.HasSuffix(parts[4], ".html") || len(parts[4]) == len(".html") {
		return "", "", false
	}
	publication := parts[3]
	prefix := ""
	if strings.HasPrefix(publication, "REC-") {
		prefix = "REC-"
	}
	if strings.HasPrefix(publication, "NOTE-") {
		prefix = "NOTE-"
	}
	dateSeparator := strings.LastIndexByte(publication, '-')
	if prefix == "" || dateSeparator <= len(prefix) {
		return "", "", false
	}
	name := publication[len(prefix):dateSeparator]
	date := publication[dateSeparator+1:]
	if !portablePublicationName(name) || len(date) != 8 || !asciiDigits(date) {
		return "", "", false
	}
	return parts[2], date, true
}

func portablePublicationName(name string) bool {
	for _, character := range name {
		if character >= 'a' && character <= 'z' || character >= 'A' && character <= 'Z' || character >= '0' && character <= '9' {
			continue
		}
		if character != '-' && character != '_' && character != '.' {
			return false
		}
	}
	return name != ""
}

func asciiDigits(value string) bool {
	for _, character := range value {
		if character < '0' || character > '9' {
			return false
		}
	}
	return value != ""
}

func validDocumentID(id string) bool {
	if id == "" || id[0] == '-' || id[len(id)-1] == '-' {
		return false
	}
	previousHyphen := false
	for _, character := range id {
		if character == '-' {
			if previousHyphen {
				return false
			}
			previousHyphen = true
			continue
		}
		if character < 'a' || character > 'z' {
			if character < '0' || character > '9' {
				return false
			}
		}
		previousHyphen = false
	}
	return true
}
