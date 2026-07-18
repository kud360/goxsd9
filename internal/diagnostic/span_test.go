package diagnostic

import "testing"

func TestSpanFormatting(t *testing.T) {
	tests := []struct {
		name   string
		uri    string
		line   int
		column int
		want   string
	}{
		{name: "absent", want: ""},
		{name: "URI only", uri: "file:///schema.xsd", want: "file:///schema.xsd"},
		{name: "URI and line", uri: "file:///schema.xsd", line: 12, want: "file:///schema.xsd:12"},
		{name: "URI line and column", uri: "file:///schema.xsd", line: 12, column: 7, want: "file:///schema.xsd:12:7"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			span, err := NewSpan(test.uri, test.line, test.column)
			if err != nil {
				t.Fatalf("NewSpan() error = %v", err)
			}
			if got := span.String(); got != test.want {
				t.Fatalf("Span.String() = %q, want %q", got, test.want)
			}
			if span.URI() != test.uri {
				t.Fatalf("Span.URI() = %q, want %q", span.URI(), test.uri)
			}
		})
	}
}

func TestNewSpanRejectsInvalidLocations(t *testing.T) {
	tests := []struct {
		name   string
		uri    string
		line   int
		column int
	}{
		{name: "negative line", uri: "schema.xsd", line: -1},
		{name: "negative column", uri: "schema.xsd", column: -1},
		{name: "line without URI", line: 1},
		{name: "column without line", uri: "schema.xsd", column: 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := NewSpan(test.uri, test.line, test.column); err == nil {
				t.Fatal("NewSpan() error = nil, want invalid location error")
			}
		})
	}
}
