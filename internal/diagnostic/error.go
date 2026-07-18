package diagnostic

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

var (
	errMissingRule       = errors.New("validation violation requires a rule reference")
	errMissingSpan       = errors.New("validation violation requires an input location")
	errNegativeItemIndex = errors.New("loop item index must not be negative")
)

type combinedError struct {
	causes []error
}

func combineErrors(causes ...error) error {
	nonNil := make([]error, 0, len(causes))
	for _, cause := range causes {
		if cause != nil {
			nonNil = append(nonNil, cause)
		}
	}
	if len(nonNil) == 0 {
		return nil
	}
	if len(nonNil) == 1 {
		return nonNil[0]
	}
	return combinedError{causes: nonNil}
}

func (e combinedError) Error() string {
	messages := make([]string, len(e.causes))
	for index, cause := range e.causes {
		messages[index] = cause.Error()
	}
	return strings.Join(messages, "; ")
}

func (e combinedError) Unwrap() []error {
	return slices.Clone(e.causes)
}

// Error adds operation, subject, source, and validation-rule context to a
// failure. It preserves an underlying cause when one exists.
type Error struct {
	operation string
	subject   string
	message   string
	span      Span
	rule      RuleRef
	cause     error
}

// New creates a contextual error without an underlying cause.
func New(operation, subject, message string) error {
	return &Error{operation: operation, subject: subject, message: message}
}

// NewAt creates a contextual error at a source location without an underlying
// cause.
func NewAt(operation, subject, message string, span Span) error {
	return &Error{operation: operation, subject: subject, message: message, span: span}
}

// Wrap adds operation and subject context. It returns nil when cause is nil.
func Wrap(operation, subject string, cause error) error {
	if cause == nil {
		return nil
	}
	return &Error{operation: operation, subject: subject, cause: cause}
}

// WrapAt adds operation, subject, and source context. It returns nil when cause
// is nil.
func WrapAt(operation, subject string, span Span, cause error) error {
	if cause == nil {
		return nil
	}
	return &Error{operation: operation, subject: subject, span: span, cause: cause}
}

// WrapItem adds the zero-based index and optional identity of a failing loop
// item before returning its child error. It returns nil when cause is nil.
func WrapItem(operation string, index int, identity string, cause error) error {
	if cause == nil {
		return nil
	}
	if index < 0 {
		indexError := fmt.Errorf("%w: %d", errNegativeItemIndex, index)
		return &Error{
			operation: operation,
			subject:   "loop item",
			cause:     combineErrors(indexError, cause),
		}
	}

	subject := fmt.Sprintf("item[%d]", index)
	if identity != "" {
		subject += " " + fmt.Sprintf("%q", identity)
	}
	return Wrap(operation, subject, cause)
}

// Violation constructs a located validation violation. rule must have been
// created by NewRuleRef. An underlying cause is optional.
func Violation(operation, subject, message string, span Span, rule RuleRef, cause error) error {
	var invariantErrors []error
	if span.IsZero() {
		invariantErrors = append(invariantErrors, errMissingSpan)
	}
	if rule.isZero() {
		invariantErrors = append(invariantErrors, errMissingRule)
	}
	if len(invariantErrors) != 0 {
		invariantErrors = append(invariantErrors, cause)
		cause = combineErrors(invariantErrors...)
	}

	return &Error{
		operation: operation,
		subject:   subject,
		message:   message,
		span:      span,
		rule:      rule,
		cause:     cause,
	}
}

// Operation returns the failed operation.
func (e *Error) Operation() string {
	if e == nil {
		return ""
	}
	return e.operation
}

// Subject returns the operation's subject.
func (e *Error) Subject() string {
	if e == nil {
		return ""
	}
	return e.subject
}

// Span returns the source span. Its zero value means no location is known.
func (e *Error) Span() Span {
	if e == nil {
		return Span{}
	}
	return e.span
}

// Rule returns the validation rule and whether one is present.
func (e *Error) Rule() (RuleRef, bool) {
	if e == nil || e.rule.isZero() {
		return RuleRef{}, false
	}
	return e.rule, true
}

// Unwrap returns the underlying cause, if any.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.cause
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	var context strings.Builder
	context.WriteString(e.operation)
	if e.subject != "" {
		if context.Len() != 0 {
			context.WriteByte(' ')
		}
		context.WriteString(e.subject)
	}
	if !e.span.IsZero() {
		if context.Len() != 0 {
			context.WriteString(" at ")
		}
		context.WriteString(e.span.String())
	}
	if !e.rule.isZero() {
		if context.Len() != 0 {
			context.WriteByte(' ')
		}
		context.WriteByte('[')
		context.WriteString(e.rule.String())
		context.WriteByte(']')
	}

	detail := e.message
	if e.cause != nil {
		if detail != "" {
			detail += ": "
		}
		detail += e.cause.Error()
	}
	if detail == "" {
		detail = "failed"
	}
	if context.Len() == 0 {
		return detail
	}
	return context.String() + ": " + detail
}
