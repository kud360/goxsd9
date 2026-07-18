---
name: implement-issue
description: Autonomously deliver one ready goxsd GitHub issue through specification research, implementation, independent evaluation, CI, pull request, follow-up harvesting, and merge to main. Use for manual or scheduled development runs; exit cleanly when no safe ready issue exists.
---

# Implement Issue

Act as an orchestrator. Keep implementation and evaluation in independent
contexts.

## Select and claim

1. Inspect `status:ready` issues in dependency order. Exclude issues with an
   open implementation pull request or unresolved decision.
2. Choose one bounded issue. If none is safe, report why and stop without
   inventing work.
3. Mark it in progress, create a short-lived `codex/` branch in an isolated
   worktree, and record the branch on the issue.

## Research and implement

1. Spawn the `spec-researcher` agent with the raw issue and bounded questions.
   Require primary-source URLs, anchors, rule names, implementation-defined
   choices, and test implications. It must not edit code.
2. Spawn the `implementer` agent with the issue plus the researcher's factual
   report. Require focused tests, full error-context handling, standard-library
   only code, and the commands in `AGENTS.md`.
3. Preserve raw command output for evaluation. Inspect the diff for unrelated
   changes and remove none of the user's pre-existing work.
4. Invoke `$harvest-followups`; link every resulting issue or roadmap update.

## Evaluate independently

1. Spawn a fresh `evaluator` agent. Give it only the original issue, final diff,
   raw test output, and primary sources. Do not provide implementation reasoning
   or the implementor's conclusions.
2. Require an unqualified `PASS` or actionable `FAIL` with file and line
   evidence. A material fix requires another fresh evaluation.
3. On failure, send only the findings and evidence to the implementer, rerun all
   checks, then repeat evaluation with a new context. Stop and mark blocked only
   for a genuine external or product decision.

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
