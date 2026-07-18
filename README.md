# goxsd

`goxsd` aims to become a conformant XSD 1.1 parser, instance validator, and Go
code generator. The target is documented conformance with the W3C
specifications, not merely compatibility with common schemas.

The project is pre-alpha. The repository currently contains the development
system and a compiling CLI shell; it does not validate schemas yet.

## Goals

- Parse XSD 1.1 into an immutable compiled schema.
- Validate XML instances and report the violated W3C rule at the source
  location.
- Bootstrap built-in types from the HFP metadata in XSD 1.1 Part 2.
- Support exact decimal values and the complete facet pipeline.
- Provide strict, Go-friendly, and user-supplied built-in backends.
- Implement the XPath 2.0 surface needed by assertions and conditional type
  alternatives, then grow toward full XPath 2.0.
- Generate allocation-frugal XML and JSON codecs, followed by BER codecs.
- Ratchet against the W3C XSD test suite without lowering the baseline.

## Design constraints

- The module uses only the Go standard library.
- Parser, compiler, validator, and generator do not start goroutines. Callers
  own concurrency.
- Parsed inputs and compiled schemas are not mutated. Assessment results will
  be immutable.
- Represent states once. Derive facts instead of storing duplicate flags or
  caches.
- Prefer ordered construction phases over scattered cycle checks.
- Return structured errors with operation, subject, location, cause, and W3C
  rule context where applicable. Add context at loop iteration boundaries.
- Keep the happy path on the left with early returns and avoid `else` where a
  guard clause is clearer.

See [architecture](docs/architecture.md), [roadmap](docs/roadmap.md), and
[automation](docs/automation.md).

## Development

Go 1.26 or newer is required.

```sh
go test ./...
go vet ./...
go run ./cmd/goxsd help
```

Initialize reference repositories after cloning:

```sh
git submodule update --init --recursive
```

## Agent workflows

Repository-local Codex skills provide the repeatable entry points intended for
manual and scheduled use. Invoke them with `$skill-name` in a Codex prompt, not
as shell commands:

- `$plan-work` reconciles the roadmap and prepares the next issues.
- `$implement-issue` delivers one ready normal or conformance issue through an
  independently evaluated pull request and automatic merge.
- `$ratchet-conformance` applies conformance controls only within the selected
  `$implement-issue` run.
- `$harvest-followups` routes discoveries to the correct planning layer.
- `$review-project` runs architecture, style, documentation, and CLI reviews.

GitHub Issues are the execution queue. The roadmap owns long-term direction,
milestones or epic issues own bounded outcomes, and leaf issues own acceptance
criteria.

## License

Apache-2.0. See [LICENSE](LICENSE).
