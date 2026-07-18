---
name: review-project
description: Independently review goxsd architecture, Go style and error handling, tests and conformance, and user-facing documentation or CLI behavior. Use for scheduled project health reviews or pre-milestone audits; file only evidenced, non-duplicate findings and do not edit product code.
---

# Review Project

Run independent read-only perspectives and turn their evidence into concise
work items.

## Workflow

1. Identify the changed or milestone-relevant surface. Read existing issues so
   reviewers can avoid duplicates without seeing one another's conclusions.
2. In parallel, spawn:
   - `architecture-reviewer` for illegal states, duplicate state, phase
     boundaries, cycle handling, long functions, `else`, and enforceable checks;
   - `user-reviewer` to use only the README, examples, and CLI help as a new
     user and identify missing or confusing behavior;
   - `evaluator` for correctness, tests, error decoration—especially loop
     boundaries—and specification traceability.
3. Require file/line or command evidence, impact, and the smallest credible
   remedy. Reviewers must not edit code.
4. Reproduce material claims when practical. Deduplicate across reports and
   existing issues.
5. Invoke `$harvest-followups` with the raw findings. Prefer programmatic
   enforcement when a focused standard-library check can reliably prevent a
   repeated style or architecture defect.
6. Report filed links and rejected findings with brief reasons.

Do not file taste-only issues, verbose documentation requests, or broad rewrites
without a concrete failure mode.
