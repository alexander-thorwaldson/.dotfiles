# kb build

Execute a plan against the knowledge base, one article at a time.

```
/kb build <plan>
```

- `plan` — path to a plan (from `/kb plan`, `/kb audit`, or `/kb maintain`)

## Overview

Build is the orchestration layer between plans and article skills. It works through a plan one article at a time, validating before and after each, and adapting the plan as new information surfaces. Every article is committed individually with user approval.

A plan contains entries with actions. The action determines which article skill is invoked:

| Action | Skill |
|---|---|
| `create` | `/kb article create` |
| `update` | `/kb article update` |
| `delete` | `/kb article delete` |

Quality over speed. One logical change at a time. A logical change may touch multiple articles — a core update plus supporting cross-reference changes in other articles is one commit, not three. No batching unrelated changes, no shortcuts.

## Setup

1. Read the plan
2. Create a branch
3. Determine execution order:
   - For creation plans: mechanical topic articles first, interpretive second. Within each group, follow the read order.
   - For change plans: order by dependency — if updating an article that others cross-reference, update it before its dependents.

## Per-Article Loop

For each entry in the execution order, run this cycle:

### 1. Pre-Validation (you + user)

Before doing anything, confirm the entry is ready:

- **Plan check** — does this entry's scope still make sense given what you've done so far? Previous articles may have revealed that the boundaries need adjustment.
- **Dependency check** — does this entry depend on a previous one that hasn't been executed yet? If so, reorder.
- **Scope check** — has anything surfaced during earlier entries that should change this one's scope or action?

If the plan needs modification, present the change to the user. Do not proceed until the plan is updated and the user agrees.

### 2. Execute

Invoke the article skill corresponding to the entry's action.

### 3. Post-Validation (you + user)

Review the result:

- **Scope adherence** — does the result stay within the entry's boundaries?
- **Quality** — are source references accurate? Is the narrative clear? Are cross-references correct?
- **Plan impact** — did this entry reveal anything that changes the plan for remaining entries? New concepts that need coverage, scope boundaries that need adjustment, entries that should be added, removed, or changed from one action to another.

Present the result to the user. The user approves, requests changes, or rejects.

- If approved, commit the logical change (primary article plus any supporting changes) and move to the next entry.
- If changes are requested, revise via `/kb article update` and re-present.
- If rejected, revisit the plan for this entry with the user before retrying.

### 4. Plan Maintenance

If post-validation surfaced plan changes, update the plan before starting the next entry. The plan is a living document throughout the build — not a fixed script.

## Interpretive Topic Transition

For creation plans: when all mechanical topic articles are complete, pause before starting interpretive topics. The completed mechanical articles are now available as context. Review the plan's interpretive entries with the user — the mechanical output may change what the interpretive articles need to cover.

## Completion

After all entries are committed to the branch, open a PR against [Mycroft](https://github.com/alexander-thorwaldson/mycroft).
