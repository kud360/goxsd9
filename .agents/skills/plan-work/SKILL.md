---
name: plan-work
description: Reconcile the goxsd long-term roadmap, GitHub milestones or epic issues, and short-term issue queue. Use for scheduled planning, backlog grooming, creating the next ready work, repairing hierarchy links, or updating direction after discoveries; do not implement product code.
---

# Plan Work

Keep an evidence-backed executable queue without duplicating mutable state
across planning layers.

## Workflow

1. Read `docs/roadmap.md`, open milestones, open epic issues, ready issues, and
   recently merged pull requests.
2. Reconcile the hierarchy:
   - roadmap sections own long-term direction;
   - milestones or epic issues own bounded outcomes;
   - leaf issues own acceptance criteria.
3. Find stale links, duplicates, blocked items, completed outcomes, and findings
   not yet harvested. Preserve useful history.
4. During each weekday-morning planning run, maintain at least 10 leaf issues
   labeled `status:ready`, ordered by dependency and smallest useful end-to-end
   progress. Do not create filler work to reach the target: every ready issue
   must satisfy the readiness evidence in the next step.
5. Before marking a leaf issue ready, link one milestone or epic parent and
   include an observable outcome,
   acceptance tests, exact specification anchors where applicable, constraints,
   and non-goals.
6. Move genuinely directional discoveries into `docs/roadmap.md`; link the
   resulting epic or milestone. Do not copy leaf checklists into the roadmap.
7. Summarize created, changed, blocked, and ready items with links.

Do not edit implementation code, lower conformance expectations, or silently
resolve an architectural decision. Label material choices
`status:needs-decision` with concise alternatives.
