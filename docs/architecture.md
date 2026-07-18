# Architecture

## Pipeline

Schemas move through explicit phases:

1. Acquire documents and retain source locations.
2. Parse XML syntax into a lossless schema-document representation.
3. Index declarations and establish construction order.
4. Resolve names and references.
5. Construct typed schema components.
6. Check schema representation and component constraints.
7. Freeze an immutable compiled schema for assessment and code generation.

Each phase consumes the previous phase's type and produces a distinct type.
Partially resolved values must not inhabit the compiled model. Handle recursion
at the construction boundary where the specification permits it; do not add
cycle detection to unrelated lookups.

Parsing, compilation, validation, and generation are synchronous. They do not
start goroutines or hide shared mutable caches. Callers may run independent
operations concurrently. A compiled schema must eventually be safe for
concurrent read-only use.

## Built-in datatypes

Generate declarative built-in metadata from the illustrative definitions in
XSD 1.1 Part 2 Appendix C, including the HFP facet and fundamental-property
annotations. Keep the pinned source and generator version in the generated
artifact.

The HFP document is informative and says primitive definitions are intrinsic.
Primitive lexical mappings, value operations, equality, ordering, and canonical
mappings therefore remain reviewed Go implementations. Generated metadata must
not pretend to supply those algorithms.

A built-in backend owns value representation and primitive operations. Planned
backends are:

- strict: preserves XSD value semantics and precision;
- Go-friendly: uses ergonomic Go values with documented conversions;
- user supplied: implements the same capability-oriented interfaces.

Backend capabilities should be expressed by interfaces and concrete types, not
parallel booleans. Exact decimal support is foundational rather than an
afterthought.

## Assessment and errors

Validation produces an immutable assessment result; it never annotates or
mutates an input tree. Streaming validation should remain possible even if a
tree-backed adapter is added later.

Operational errors gain context at every boundary. Validation violations add a
stable project code, the W3C constraint name, a direct specification anchor,
and source location. Multiple violations may be reported, but their ordering
must be deterministic. `slog` is an optional diagnostic sink, never the error
transport.

## Code generation

Generate explicit choice types and type switches rather than bags of pointer
fields. XML and JSON codecs come first; BER follows. Generated codecs should
minimize allocation while retaining enough path and state information to make
failures debuggable by humans and agents.
