package diagnostic

import (
	"fmt"
	"net/url"
	"strings"
)

// RuleRef relates a stable goxsd validation code to a directly anchored W3C
// specification rule. The project code and W3C anchor have separate identities.
type RuleRef struct {
	code          string
	specification string
}

// NewRuleRef constructs a validation rule reference. specification must be a
// direct HTTPS URL into the W3C site, including a fragment identifier.
func NewRuleRef(code, specification string) (RuleRef, error) {
	if code == "" {
		return RuleRef{}, fmt.Errorf("construct validation rule reference: project code is empty")
	}
	if strings.TrimSpace(code) != code || strings.ContainsAny(code, "\t\r\n ") {
		return RuleRef{}, fmt.Errorf("construct validation rule reference for %q: project code contains whitespace", code)
	}

	reference, err := url.Parse(specification)
	if err != nil {
		return RuleRef{}, fmt.Errorf("construct validation rule reference for %q: parse specification URL %q: %w", code, specification, err)
	}
	if reference.Scheme != "https" || reference.Hostname() != "www.w3.org" || reference.Fragment == "" {
		return RuleRef{}, fmt.Errorf("construct validation rule reference for %q: specification URL must be a direct HTTPS W3C anchor: %q", code, specification)
	}

	return RuleRef{code: code, specification: specification}, nil
}

// Code returns the stable project-defined validation code. It is not a W3C
// constraint name.
func (r RuleRef) Code() string {
	return r.code
}

// Specification returns the direct W3C specification URL.
func (r RuleRef) Specification() string {
	return r.specification
}

// Anchor returns the W3C rule name or clause anchor derived from Specification.
func (r RuleRef) Anchor() string {
	_, anchor, _ := strings.Cut(r.specification, "#")
	return anchor
}

func (r RuleRef) isZero() bool {
	return r.code == ""
}

func (r RuleRef) String() string {
	if r.isZero() {
		return ""
	}
	return fmt.Sprintf("%s (W3C %s: %s)", r.code, r.Anchor(), r.specification)
}
