# Project guidance

Build `goxsd` toward XSD 1.1 conformance. Read `docs/architecture.md`,
`docs/roadmap.md`, and the issue before changing code.

## Engineering rules

- Use Go 1.26 and the standard library only. Do not add module dependencies.
- Do not start goroutines in parsing, compilation, validation, or generation.
- Prefer immutable values and types that make illegal states unrepresentable.
- Do not store derivable state. In particular, do not pair a discriminator with
  fields whose validity duplicates that discriminator.
- Use ordered construction phases to resolve references and recursion. Do not
  scatter defensive cycle maps through the model.
- Keep the happy path left-aligned with guard clauses. Avoid `else` when an
  early return is clearer.
- Split functions before they exceed roughly 100 lines or become difficult to
  reason about.
- Keep exported documentation concise. Prefer a small example to a long essay.
- Generated artifacts must be deterministic and carry a generated-code header.

## Errors

- Handle every error. Wrap non-sentinel errors with the failed operation and
  useful subject before returning across a package boundary.
- In loops, decorate an error inside the iteration with the item index, QName,
  URI, path, or phase before returning. Never return a bare child error after
  losing which item failed.
- Validation violations must carry a stable code, specification URL and anchor,
  input location, and underlying cause when one exists.
- Avoid panics. Return contextual errors for malformed input, unsupported
  behavior, violated invariants, and operational failures.
- Return structured errors to callers. Use caller-provided `slog` logging only
  for optional diagnostics; logging is not error handling.
- Preserve causes with `Unwrap`/`%w` so `errors.Is` and `errors.As` work.

## Verification

Run before handing work to evaluation:

```sh
gofmt -w .
go vet ./...
go test ./...
```

Run relevant conformance tests when the harness exists. Never lower a checked-in
baseline or hide a regression by reclassifying a test without evidence.

## Delivery

Use GitHub Issues as the executable queue. Keep one issue per change. Cite the
relevant specification rules in the issue, tests, errors, and pull request.

Implementation and evaluation must use independent agent contexts. The
evaluator receives the issue, diff, raw test output, and specifications—not the
implementor's reasoning. Only apply `agent:automerge` after an evaluator returns
an unqualified pass and required CI succeeds for that exact head commit.

Before merging, run `$harvest-followups`. Record bounded work as linked issues;
record longer direction in the roadmap and link an epic or milestone. Do not
leave actionable TODO comments as a substitute for tracking.
