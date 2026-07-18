# Roadmap

This file owns long-term direction. GitHub milestones or linked epic issues own
bounded outcomes. Leaf issues own implementation details and acceptance tests.
Link between layers; do not copy the same mutable checklist into each layer.

## Milestone 0: first validation slice

Prove the pipeline with one schema document and a small end-to-end feature set:

- global elements with named or anonymous simple types;
- HFP-driven metadata for `string`, `boolean`, `decimal`, and selected integer
  derivations;
- restrictions using whitespace, pattern, enumeration, and numeric bounds;
- exact decimal values;
- validation of simple element content;
- `goxsd validate schema.xsd instance.xml`;
- structured, located, specification-linked errors.

Imports, includes, complex content, XPath, and code generation are outside this
milestone unless required to establish a clean architectural seam.

## Long-term tracks

1. Schema acquisition, parsing, multi-phase component construction, and all XSD
   1.1 schema constraints.
2. Built-in and user-defined datatypes, exact values, and the full facet
   pipeline.
3. Streaming instance assessment and immutable assessment results.
4. XPath 2.0 assertions and conditional type alternatives, expanding from the
   XSD-required subset to the full language.
5. Strict, Go-friendly, and external built-in backends.
6. Allocation-frugal Go generation with XML, JSON, then BER codecs.
7. Specification ingestion, rule indexing, conformance reporting, and developer
   tools.
8. Concise user documentation, examples, and CLI ergonomics.

## Conformance policy

The W3C suite baseline is monotonic. Store totals by suite revision, profile,
and outcome: pass, fail, unsupported, disputed, and harness error. Never combine
unsupported tests with failures or silently exclude tests.

The suite checks top-level validity outcomes, not every processor requirement.
Maintain a separate normative-requirements ledger for PSVI behavior, reporting,
implementation-defined choices, and requirements without suite coverage.

Every baseline change records the pinned suite commit, command, environment,
and linked issue or pull request.
