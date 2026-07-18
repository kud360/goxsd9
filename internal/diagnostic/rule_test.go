package diagnostic

import "testing"

func TestRuleRefUsesStableProjectCodeAndDirectW3CAnchor(t *testing.T) {
	// Element Locally Valid (Type), XSD 1.1 Part 1:
	// https://www.w3.org/TR/2012/REC-xmlschema11-1-20120405/#cvc-type
	const specification = "https://www.w3.org/TR/2012/REC-xmlschema11-1-20120405/#cvc-type"
	rule, err := NewRuleRef("GOXSD-VALIDATION-0001", specification)
	if err != nil {
		t.Fatalf("NewRuleRef() error = %v", err)
	}
	if got := rule.Code(); got != "GOXSD-VALIDATION-0001" {
		t.Fatalf("RuleRef.Code() = %q", got)
	}
	if got := rule.Anchor(); got != "cvc-type" {
		t.Fatalf("RuleRef.Anchor() = %q, want cvc-type", got)
	}
	if got := rule.Specification(); got != specification {
		t.Fatalf("RuleRef.Specification() = %q, want %q", got, specification)
	}
}

func TestNewRuleRefRejectsIndirectReferences(t *testing.T) {
	tests := []struct {
		name string
		code string
		url  string
	}{
		{name: "empty code", url: "https://www.w3.org/TR/example/#rule"},
		{name: "code whitespace", code: "GOXSD 1", url: "https://www.w3.org/TR/example/#rule"},
		{name: "non-W3C host", code: "GOXSD-1", url: "https://example.com/spec#rule"},
		{name: "insecure URL", code: "GOXSD-1", url: "http://www.w3.org/TR/example/#rule"},
		{name: "missing anchor", code: "GOXSD-1", url: "https://www.w3.org/TR/example/"},
		{name: "malformed URL", code: "GOXSD-1", url: "://bad"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := NewRuleRef(test.code, test.url); err == nil {
				t.Fatal("NewRuleRef() error = nil, want invalid reference error")
			}
		})
	}
}
