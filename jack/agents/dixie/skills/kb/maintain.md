# kb maintain

Scoped drift detection from a code change.

```
/kb maintain <diff>
```

- `diff` — a PR reference, commit range, or diff describing what changed in the source

## Overview

Maintain maps a specific code change to its KB impact. Unlike audit, which reviews the full corpus, maintain starts from a known change and determines which articles are affected. The output is a change plan consumed by `/kb build`.

## Phases

### Phase 1 — Study (you, serial)

1. Read the diff — understand what changed, what files were touched, what behavior shifted
2. Map the changed files against `sources` fields in existing KB articles to identify directly affected articles
3. Check cross-references — articles that link to affected articles may also need updates
4. Generate a codeword and create the workspace under `/tmp/<codeword>/`

### Phase 2 — Impact Assessment (3 agents, parallel)

Spawn 3 agents, one per lens. Each agent receives:
- The diff
- The set of potentially affected articles (from Phase 1)
- The topic placement table
- Access to the repo source (current state, post-change)

Each agent assesses impact:

| Category | What to look for |
|---|---|
| **Direct drift** | Articles whose documented behavior no longer matches the source |
| **Indirect drift** | Articles that cross-reference affected content or describe related concepts |
| **New coverage** | Does the change introduce concepts that need new articles? |
| **Removals** | Does the change delete functionality that an article documents? |

Findings must reference the specific lines in the diff that drive each recommendation.

### Phase 3 — Edit (2 agents, parallel)

Two editor agents. Each must:
1. Read the diff and affected articles **first**
2. Read all 3 impact reports from `/tmp` **second**
3. Consolidate — validate findings, cut false positives, merge duplicates

### Phase 4 — Change Plan (you, serial)

1. Review both consolidated reports
2. Synthesize into a change plan — each entry specifies:
   - The article (existing path or proposed new path)
   - The action (`create`, `update`, or `delete`)
   - The scope (what needs to change, tied to the specific diff)
   - For updates: what should not change
3. Write the change plan to `/tmp/<codeword>/plan.md`
4. Present the change plan to the user for review
