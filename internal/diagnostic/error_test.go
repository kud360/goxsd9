package diagnostic

import (
	"errors"
	"testing"
)

var errMalformedValue = errors.New("malformed value")

func TestNestedErrorsRetainCauseAndStructuredContext(t *testing.T) {
	inner := Wrap("decode", "attribute value", errMalformedValue)
	outer := Wrap("parse", "schema document", inner)

	if !errors.Is(outer, errMalformedValue) {
		t.Fatalf("errors.Is(%v, errMalformedValue) = false", outer)
	}
	var contextual *Error
	if !errors.As(outer, &contextual) {
		t.Fatalf("errors.As(%T, *Error) = false", outer)
	}
	if contextual.Operation() != "parse" || contextual.Subject() != "schema document" {
		t.Fatalf("outer context = (%q, %q), want parse schema document", contextual.Operation(), contextual.Subject())
	}
	if got, want := outer.Error(), "parse schema document: decode attribute value: malformed value"; got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func TestWrapItemRetainsItemIdentityAndCause(t *testing.T) {
	err := WrapItem("compile declaration", 3, "{urn:example}widget", errMalformedValue)
	if !errors.Is(err, errMalformedValue) {
		t.Fatalf("errors.Is(%v, errMalformedValue) = false", err)
	}
	if got, want := err.Error(), `compile declaration item[3] "{urn:example}widget": malformed value`; got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func TestWrapItemRejectsNegativeIndexAndRetainsCause(t *testing.T) {
	err := WrapItem("compile declaration", -1, "ignored", errMalformedValue)
	if !errors.Is(err, errNegativeItemIndex) {
		t.Fatalf("errors.Is(%v, errNegativeItemIndex) = false", err)
	}
	if !errors.Is(err, errMalformedValue) {
		t.Fatalf("errors.Is(%v, errMalformedValue) = false", err)
	}
	if got, want := err.Error(), "compile declaration loop item: loop item index must not be negative: -1; malformed value"; got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func TestViolationFormattingAndAccessors(t *testing.T) {
	span, err := NewSpan("file:///instance.xml", 5, 9)
	if err != nil {
		t.Fatalf("NewSpan() error = %v", err)
	}
	// Datatype Valid, XSD 1.1 Part 2:
	// https://www.w3.org/TR/2012/REC-xmlschema11-2-20120405/#cvc-datatype-valid
	const specification = "https://www.w3.org/TR/2012/REC-xmlschema11-2-20120405/#cvc-datatype-valid"
	rule, err := NewRuleRef("GOXSD-VALIDATION-0002", specification)
	if err != nil {
		t.Fatalf("NewRuleRef() error = %v", err)
	}

	violation := Violation("validate", "element value", "value is not valid for its datatype", span, rule, errMalformedValue)
	var contextual *Error
	if !errors.As(violation, &contextual) {
		t.Fatalf("errors.As(%T, *Error) = false", violation)
	}
	gotRule, ok := contextual.Rule()
	if !ok || gotRule != rule {
		t.Fatalf("Error.Rule() = (%v, %t), want (%v, true)", gotRule, ok, rule)
	}
	if contextual.Span() != span {
		t.Fatalf("Error.Span() = %v, want %v", contextual.Span(), span)
	}
	want := "validate element value at file:///instance.xml:5:9 " +
		"[GOXSD-VALIDATION-0002 (W3C cvc-datatype-valid: " + specification + ")]: " +
		"value is not valid for its datatype: malformed value"
	if got := violation.Error(); got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
	if !errors.Is(violation, errMalformedValue) {
		t.Fatalf("errors.Is(%v, errMalformedValue) = false", violation)
	}
}

func TestNilAndMissingCauseEdges(t *testing.T) {
	if err := Wrap("parse", "schema", nil); err != nil {
		t.Fatalf("Wrap(..., nil) = %v, want nil", err)
	}
	if err := WrapAt("parse", "schema", Span{}, nil); err != nil {
		t.Fatalf("WrapAt(..., nil) = %v, want nil", err)
	}
	if err := WrapItem("parse", 0, "first", nil); err != nil {
		t.Fatalf("WrapItem(..., nil) = %v, want nil", err)
	}

	err := New("acquire", "schema", "document is unavailable")
	if errors.Unwrap(err) != nil {
		t.Fatalf("errors.Unwrap(%v) = %v, want nil", err, errors.Unwrap(err))
	}
	if got, want := err.Error(), "acquire schema: document is unavailable"; got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func TestViolationRejectsMissingRuleAndRetainsCause(t *testing.T) {
	span, spanErr := NewSpan("file:///instance.xml", 1, 1)
	if spanErr != nil {
		t.Fatalf("NewSpan() error = %v", spanErr)
	}
	err := Violation("validate", "value", "invalid", span, RuleRef{}, errMalformedValue)
	if err == nil {
		t.Fatal("Violation() = nil, want missing rule error")
	}
	if !errors.Is(err, errMissingRule) {
		t.Fatalf("errors.Is(%v, errMissingRule) = false", err)
	}
	if !errors.Is(err, errMalformedValue) {
		t.Fatalf("errors.Is(%v, errMalformedValue) = false", err)
	}
}

func TestViolationRejectsMissingInputLocation(t *testing.T) {
	rule, ruleErr := NewRuleRef(
		"GOXSD-VALIDATION-0002",
		"https://www.w3.org/TR/2012/REC-xmlschema11-2-20120405/#cvc-datatype-valid",
	)
	if ruleErr != nil {
		t.Fatalf("NewRuleRef() error = %v", ruleErr)
	}
	err := Violation("validate", "value", "invalid", Span{}, rule, nil)
	if !errors.Is(err, errMissingSpan) {
		t.Fatalf("errors.Is(%v, errMissingSpan) = false", err)
	}
}

func TestNilErrorMethodsDoNotPanic(t *testing.T) {
	var err *Error
	if got := err.Error(); got != "<nil>" {
		t.Fatalf("nil Error.Error() = %q, want <nil>", got)
	}
	if err.Unwrap() != nil || err.Operation() != "" || err.Subject() != "" || !err.Span().IsZero() {
		t.Fatal("nil Error accessors returned non-zero values")
	}
	if _, ok := err.Rule(); ok {
		t.Fatal("nil Error.Rule() reports a rule")
	}
}
