package diagnostic

import (
	"fmt"
	"strconv"
)

// Span identifies a source URI and, when known, a one-based line and column.
// Its zero value represents an unknown source location.
type Span struct {
	uri    string
	line   int
	column int
}

// NewSpan constructs a source span. Line and column are one-based; zero means
// that part of the location is unknown. A column requires a line, and a line
// requires a URI.
func NewSpan(uri string, line, column int) (Span, error) {
	if line < 0 {
		return Span{}, fmt.Errorf("construct source span for %q: line must not be negative: %d", uri, line)
	}
	if column < 0 {
		return Span{}, fmt.Errorf("construct source span for %q: column must not be negative: %d", uri, column)
	}
	if uri == "" && line != 0 {
		return Span{}, fmt.Errorf("construct source span at line %d: line requires a URI", line)
	}
	if line == 0 && column != 0 {
		return Span{}, fmt.Errorf("construct source span for %q at column %d: column requires a line", uri, column)
	}

	return Span{uri: uri, line: line, column: column}, nil
}

// URI returns the source URI, or an empty string when it is unknown.
func (s Span) URI() string {
	return s.uri
}

// Line returns the one-based line and whether it is known.
func (s Span) Line() (int, bool) {
	return s.line, s.line != 0
}

// Column returns the one-based column and whether it is known.
func (s Span) Column() (int, bool) {
	return s.column, s.column != 0
}

// IsZero reports whether the span contains no source location.
func (s Span) IsZero() bool {
	return s.uri == ""
}

func (s Span) String() string {
	if s.uri == "" {
		return ""
	}
	if s.line == 0 {
		return s.uri
	}

	location := s.uri + ":" + strconv.Itoa(s.line)
	if s.column == 0 {
		return location
	}
	return location + ":" + strconv.Itoa(s.column)
}
