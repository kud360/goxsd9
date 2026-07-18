---
name: harvest-followups
description: Find and route actionable follow-up work from goxsd diffs, tests, TODOs, review findings, specification research, and conformance results. Use before pull-request delivery, after reviews, or during scheduled backlog hygiene to prevent discoveries from being lost or duplicated.
---

# Harvest Follow-ups

Turn evidence into linked work at the smallest correct planning layer.

## Workflow

1. Inspect the current issue, diff, test output, TODO/FIXME comments, skipped
   tests, evaluator findings, research notes, and nearby open issues.
2. Discard observations that are not actionable. Merge duplicates into the
   existing issue with new evidence instead of opening another.
3. Classify each remaining finding:
   - required for current acceptance: keep it in the current issue;
   - bounded follow-up: create a linked leaf issue;
   - bounded multi-issue outcome: link or create an epic/milestone;
   - long-term direction: update `docs/roadmap.md` and link an outcome;
   - material unresolved choice: file `status:needs-decision`.
4. Give new issues an observable outcome, acceptance tests, source anchors,
   constraints, non-goals, and parent link. Add dependency links when order
   matters.
5. Remove obsolete TODO comments only when their issue link or completed work
   makes the comment redundant.
6. Return a compact list of created, updated, deduplicated, and retained items.

Do not create speculative backlogs, copy the same checklist between layers, or
hide unfinished acceptance criteria in a follow-up.
