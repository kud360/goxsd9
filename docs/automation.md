# Automation

Codex scheduled tasks invoke repository-local skills. Test each skill manually
before scheduling it. Prefer isolated worktrees for implementation runs and the
main checkout for read-only planning runs.

## Suggested tasks

| Task | Suggested cadence | Prompt |
| --- | --- | --- |
| Plan work | Daily | `Use $plan-work. Reconcile the roadmap and keep a small ready queue.` |
| Implement | Every few hours | `Use $implement-issue. Deliver one ready issue completely.` |
| Conformance | Daily after the harness exists | `Use $ratchet-conformance. Raise one bounded conformance cluster.` |
| Project review | Weekly | `Use $review-project. File only actionable, non-duplicate findings.` |

`$implement-issue` and `$ratchet-conformance` run `$harvest-followups` before
delivery. A scheduled run with no safe ready work exits without inventing work.

## Independent roles

- The spec researcher answers bounded questions from pinned primary sources and
  does not edit implementation code.
- The implementor receives the issue and research artifacts.
- The evaluator starts with fresh context and receives only the issue, diff,
  raw checks, and source references.
- Architecture and user reviewers are read-only and file concise findings.

The orchestrating agent may resolve mechanical feedback, but a material change
requires a fresh evaluator verdict.

## Automatic merge

Configure GitHub to protect `main`, require the `test` check from the `CI`
workflow, allow squash merges, and prevent force pushes. After required CI
succeeds and the read-only
evaluator returns an unqualified pass for that head, the orchestrator applies
`agent:automerge`. The workflow performs one SHA-bound merge. It creates no
persistent auto-merge request. Any new commit or reopen revokes the label and
requires a fresh evaluation.

Recommended labels:

- `status:ready`, `status:blocked`, `status:needs-decision`
- `kind:work`, `kind:conformance`, `kind:architecture`, `kind:docs`
- `area:parser`, `area:datatype`, `area:validator`, `area:xpath`, `area:codegen`
- `agent:automerge`

Autonomous agents may create and edit issues, branches, commits, pull requests,
milestones, roadmap links, and merge passing work. They must not weaken branch
protection, required checks, the conformance baseline, or these authority rules
to make a run succeed.
