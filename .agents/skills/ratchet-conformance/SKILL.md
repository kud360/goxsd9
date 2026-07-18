---
name: ratchet-conformance
description: Provide the conformance procedure within $implement-issue for one ready goxsd conformance issue: raise the pinned W3C XSD baseline by one coherent cluster while preserving all prior passes. Use only after the conformance harness exists.
---

# Ratchet Conformance

Use this procedure only within the selected `kind:conformance` issue in the
unified `$implement-issue` workflow; do not schedule it as separate delivery
automation. Improve one bounded cluster and preserve auditable test
classifications.

## Workflow

1. Verify the test-suite submodule commit and current baseline. Run the harness
   before editing and retain its raw report.
2. Select the smallest coherent failing or unsupported cluster that unlocks
   useful behavior. Do not cherry-pick tests merely because they are easy.
3. Spawn `spec-researcher` with test metadata and exact normative questions.
   Treat suite expectations as evidence, not a replacement for the spec and
   errata.
4. Spawn `implementer` with the cluster, baseline, and research report. Require
   focused regression tests plus the full suite command.
5. Classify every changed outcome as pass, fail, unsupported, disputed, or
   harness error. Never turn a failure into an exclusion to raise the score.
6. Update the machine-readable baseline with suite commit, environment,
   command, counts, and changed test IDs. The pass set must be monotonic.
7. Invoke `$harvest-followups`.
8. Spawn a fresh `evaluator` with the issue, diff, before/after raw reports, and
   sources. Require confirmation that no pass regressed and classifications are
   evidenced.
9. Deliver through a pull request and automatic merge using the same gates as
   `$implement-issue`.

If the harness does not exist, create or refine a linked harness issue and stop.
Do not fabricate a baseline.
