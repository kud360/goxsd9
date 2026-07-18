package specsync

import (
	"strings"
	"testing"
)

func TestReadManifest(t *testing.T) {
	t.Parallel()
	manifest, err := ReadManifest(strings.NewReader(`{
  "documents": [{
    "id": "structures",
    "url": "https://www.w3.org/TR/2012/REC-fixture-20120405/structures.html",
    "sha256": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
    "output": "structures.md"
  }]
}`))
	if err != nil {
		t.Fatalf("ReadManifest() error = %v", err)
	}
	if got := manifest.Documents[0].ID; got != "structures" {
		t.Fatalf("document ID = %q, want structures", got)
	}
}

func TestReadManifestRejectsNilReader(t *testing.T) {
	t.Parallel()
	_, err := ReadManifest(nil)
	if err == nil || !strings.Contains(err.Error(), "decode specification manifest: reader must not be nil") {
		t.Fatalf("ReadManifest(nil) error = %v, want contextual nil-reader error", err)
	}
}

func TestReadManifestRejectsInvalidDocumentsWithIndex(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		document string
		want     string
	}{
		{"mutable URL", `{"id":"x","url":"https://www.w3.org/TR/xmlschema11-1/structures.html","sha256":"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef","output":"x.md"}`, "dated W3C REC or NOTE"},
		{"traversal ID", `{"id":"../escaped","url":"https://www.w3.org/TR/2012/REC-fixture-20120405/x.html","sha256":"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef","output":"x.md"}`, "lowercase ASCII letters"},
		{"uppercase digest", `{"id":"x","url":"https://www.w3.org/TR/2012/REC-fixture-20120405/x.html","sha256":"0123456789ABCDEF0123456789abcdef0123456789abcdef0123456789abcdef","output":"x.md"}`, "lowercase hexadecimal"},
		{"nested output", `{"id":"x","url":"https://www.w3.org/TR/2012/REC-fixture-20120405/x.html","sha256":"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef","output":"sub/x.md"}`, "Markdown basename"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			_, err := ReadManifest(strings.NewReader(`{"documents":[` + test.document + `]}`))
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("ReadManifest() error = %v, want substring %q", err, test.want)
			}
		})
	}
}

func TestReadManifestRejectsDuplicateIDAndOutput(t *testing.T) {
	t.Parallel()
	input := `{"documents":[
{"id":"x","url":"https://www.w3.org/TR/2012/REC-fixture-20120405/x.html","sha256":"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef","output":"x.md"},
{"id":"x","url":"https://www.w3.org/TR/2012/REC-fixture-20120405/y.html","sha256":"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef","output":"x.md"}]}`
	_, err := ReadManifest(strings.NewReader(input))
	if err == nil || !strings.Contains(err.Error(), `document[1] "x": duplicate id`) {
		t.Fatalf("ReadManifest() error = %v, want duplicate id with index", err)
	}
}
