# kb audit

Full corpus review for drift, accuracy, and coverage gaps.

```
/kb audit
```

No arguments — you audit the KB for the repo you are haunting.

## Overview

Audit compares the current state of the repo against the full knowledge base to produce a change plan. The output is consumed by `/kb build`.

## Phases

### Phase 1 — Study (you, serial)

1. Read the current state of the repo — file tree, README, recent commit history
2. Read all existing KB articles, noting `last_updated` dates and `sources` fields
3. Generate a codeword and create the workspace under `/tmp/<codeword>/`

### Phase 2 — Drift Detection (3 agents, parallel)

Spawn 3 agents, one per lens. Each agent receives:
- The full set of existing KB articles
- The topic placement table
- Access to the repo source

Each agent reviews the entire KB against the source, looking for:

| Category | What to look for |
|---|---|
| **Drift** | Source has changed but articles still describe old behavior. Check `sources` paths — have those files been modified or moved? |
| **Inaccuracy** | Articles that were wrong or are now wrong — incorrect descriptions, broken source references, outdated file paths |
| **Gaps** | Concepts in the source that have no article coverage |
| **Redundancy** | Articles that overlap significantly or could be consolidated |
| **Removals** | Articles documenting code that no longer exists |
| **Quality** | Weak source grounding, missing diagrams, vague descriptions |

Findings must be specific and actionable. Reference exact articles and exact source locations.

### Phase 3 — Edit (2 agents, parallel)

Two editor agents. Each must:
1. Study the repo source and existing KB **first**
2. Read all 3 audit reports from `/tmp` **second**
3. Consolidate — validate findings, merge duplicates, dispute anything inaccurate. Produce a single unified list of findings.

### Phase 4 — Change Plan (you, serial)

1. Review both consolidated reports
2. Synthesize into a change plan — each entry specifies:
   - The article (existing path or proposed new path)
   - The action (`create`, `update`, or `delete`)
   - The scope (what needs to change and why)
   - For updates: what should not change
3. Write the change plan to `/tmp/<codeword>/plan.md`
4. Present the change plan to the user for review
