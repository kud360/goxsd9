---
name: implement-issue
description: Autonomously deliver one ready normal or conformance goxsd issue through specification research, implementation, raw verification, follow-up harvesting, independent evaluation, CI, pull request, and merge to main. Use for manual or scheduled development runs; exit cleanly when no safe ready issue exists.
---

# Implement Issue

Act as an orchestrator. Keep implementation and evaluation in independent
contexts.

## Select and claim

1. Inspect `status:ready` normal and conformance issues in dependency order.
   Exclude issues with an open implementation pull request or unresolved
   decision.
2. Choose exactly one safe, bounded `status:ready` issue. For a conformance
   issue, scope its coherent test cluster to that issue's acceptance criteria.
   If none is safe, report why and stop without inventing work.
3. Mark it in progress, create a short-lived `codex/` branch in an isolated
   worktree, and record the branch on the issue.

## Establish the conformance baseline

For a selected `kind:conformance` issue, before any implementation work:

1. Record the pinned suite commit and current machine-readable baseline.
2. Run the harness before editing and retain its raw measurement.

## Research and implement

1. Spawn the `spec-researcher` agent with the raw issue and bounded questions.
   Require primary-source URLs, anchors, rule names, implementation-defined
   choices, and test implications. It must not edit code.
2. For a conformance issue, give the implementer the pre-edit measurement and
   baseline with the selected issue and its scoped cluster. Otherwise, give the
   implementer the selected issue. In either case, include the researcher's
   factual report and require focused tests, full error-context handling,
   standard-library-only code, and the commands in `AGENTS.md`.
3. Run the required verification and preserve its raw command output for
   evaluation. Inspect the diff for unrelated changes and remove none of the
   user's pre-existing work.

## Ratchet conformance within the selected issue

For every selected `kind:conformance` issue, after implementation and before
evaluation:

1. Run the harness after implementation and retain its raw measurement.
   Classify every changed outcome as pass, fail, unsupported, disputed, or
   harness error; do not hide a failure by excluding it.
2. Update the machine-readable baseline with the suite commit, environment,
   command, counts, and changed test IDs. Preserve a monotonic pass set.
3. Supply the before-and-after reports to evaluation and require confirmation
   that no prior pass regressed and all classifications are evidenced.

## Harvest follow-ups

After any applicable conformance controls, invoke `$harvest-followups`; link
every resulting issue or roadmap update before evaluation.

## Evaluate independently

1. Spawn a fresh `evaluator` agent. Give it only the original issue, final diff,
   raw test output, and primary sources. For a conformance issue, include the
   before-and-after harness reports and baseline. Do not provide implementation
   reasoning or the implementor's conclusions.
2. Require an unqualified `PASS` or actionable `FAIL` with file and line
   evidence. A material fix requires another fresh evaluation.
3. On failure, send only the findings and evidence to the implementer, then
   rerun all checks. After every material fix, rerun applicable conformance
   controls and `$harvest-followups` before repeating evaluation with a new
   context. Stop and mark blocked only for a genuine external or product
   decision.

## Deliver

1. Commit intentionally, push, and open a pull request using the repository
   template. Link the issue, specification anchors, verification, evaluator
   verdict, and follow-ups.
2. Wait for required CI. Address failures without weakening checks.
3. Apply `agent:automerge` only after both `PASS` and green required CI. The
   SHA-bound workflow merges to `main`; confirm the issue closed and report the
   merge commit.

Never force-push `main`, bypass required checks, merge an evaluator failure, or
add an external Go dependency.
