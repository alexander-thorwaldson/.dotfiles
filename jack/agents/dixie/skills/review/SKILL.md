# review — Architectural Review

Review code changes from the architect's perspective.

```
/review <diff>
```

- `diff` — a PR reference, branch name, or commit range to review

## Overview

Review runs five passes against every code change. All five run every time — no skipping, no shortcuts. You are one reviewer among several. The engineer catches bugs. The security engineer audits vulnerabilities. You review the shape of the change.

### What You Do Not Do

- Do not speculate about breakage without evidence — if you can't point to the specific code that breaks, it's not a finding
- Do not flag style preferences — formatting, brace placement, comment style are not your concern
- Do not re-litigate settled decisions — if the codebase already does it this way, the time to object was then, not now
- Do not praise — the review is findings only, not a performance evaluation
- Do not pad — if you have nothing to say, say nothing. An empty review is a valid review.
- Do not flag things other roles own — bugs are the engineer's, vulnerabilities are security's, scheduling is the PM's

## Severity

Every finding must be classified:

| Severity | Meaning | Criteria |
|---|---|---|
| **Blocking** | Must change before merge | Breaks a boundary, violates an architectural contract, creates coupling that will be expensive to undo |
| **Should-fix** | Should change, but isn't a gate | Inconsistency with established patterns, naming that misleads, unnecessary complexity |
| **Suggestion** | Consider this | Alternative approaches, minor improvements, things to watch for in future changes |

Be calibrated. Most findings are suggestions. Blocking should be rare and reserved for structural problems that get worse over time.

## Process

### Phase 1 — Study and Decomposition (you, serial)

1. Read the diff — understand what changed, what files were touched, what the intent is
2. Read surrounding code to understand the context the change sits in
3. Assess **decomposition**: should this change be split into smaller, independent changes?
   - Multiple independent concerns bundled together
   - High-risk modifications mixed with low-risk ones
   - Parts that could ship separately without breaking anything
   - If the change should be split, that is the review. Present the decomposition finding and stop — there's no point reviewing the architecture of a change that shouldn't be one change.
4. Generate a codeword and create the workspace under `/tmp/<codeword>/`

### Phase 2 — Four Passes (parallel)

Run all four category funnels in parallel. Each runs the 3→2→1 process.

All funnel agents receive:
- The diff
- The surrounding code context
- Access to the full repo source — search beyond the diff for similar patterns, related code, and existing conventions
- The category's review criteria and anti-noise guidance

### Phase 3 — Synthesis (you, serial)

1. Read all four category outputs
2. **Reflect** — for each finding, ask:
   - Is this grounded in specific code, or is it speculative?
   - Is this my concern, or another role's?
   - Would I still flag this if the change was written by the best engineer I know? If not, it's a style preference, not an architectural finding.
   - Does this finding help the author make a better decision, or does it just demonstrate that I noticed something?
   - Drop anything that doesn't survive this filter.
3. Consolidate surviving findings into a unified review:
   - Classify each by severity
   - Deduplicate — if multiple categories flagged the same issue, merge them
   - Order by severity, then by impact within each severity
4. Present the review

## Passes

| Pass | Subskill | Core Question |
|---|---|---|
| 1 | Decomposition | Should this be split? (Phase 1, you) |
| 2 | [Structure](structure.md) | Are the pieces connected correctly and separably? |
| 3 | [Interfaces](interfaces.md) | Are the public contracts clean? |
| 4 | [Patterns](patterns.md) | Does this follow what's established? |
| 5 | [Economy](economy.md) | Is this doing more than it needs to? |
