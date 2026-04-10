# kb maintain

Review existing knowledge base content for drift, accuracy, and quality.

```
/kb maintain <repo> [topics...]
```

- `repo` — GitHub URL or `owner/repo` of the documented repository
- `topics` — specific topic folders to review (e.g. `api data`). If omitted, review the full corpus.

## Overview

Maintain does not modify content. It produces a structured change request identifying what has drifted, what is inaccurate, and what is low quality.

## Phases

### Phase 1 — Study (you, serial)

1. Clone the target repo and [Mycroft](https://github.com/alexander-thorwaldson/mycroft)
2. Read the current state of the target repo — file tree, README, recent commits since `last_updated` dates in existing articles
3. Read the existing KB content for the specified topics (or all topics if none specified)
4. Generate a codeword and create the workspace under `/tmp/<codeword>/`

### Phase 2 — Change Discovery (parallel across topics)

For each topic under review, spawn a convergent funnel (see SKILL.md). This funnel produces change requests, not content.

Each funnel agent receives:
- The topic's existing KB articles
- The topic's question set
- The topic placement table
- Access to the target repo source

Each agent must produce a structured change request covering:

| Category | What to look for |
|---|---|
| **Drift** | Source code has changed but articles still describe the old behavior |
| **Inaccuracy** | Articles that were wrong or are now wrong — incorrect descriptions, broken source references, outdated file paths |
| **Gaps** | Concepts that exist in the source but have no article coverage |
| **Redundancy** | Articles that overlap significantly or could be consolidated |
| **Quality** | Weak source grounding, missing diagrams where they'd help, vague descriptions that should be specific |

Change requests must reference specific articles and specific source locations. "The API docs are outdated" is not actionable. "api/public-interfaces.md describes endpoint X at src/routes.ts:40 but that function was moved to src/api/handlers.ts:22" is.

### Phase 3 — Funnel Convergence

The funnel follows the process defined in SKILL.md, but agents are refining change requests rather than articles.

Pass 2 agents study the repo source and existing KB first, forming their own understanding, then read prior change requests from `/tmp`. They produce a consolidated change request that validates, merges, or disputes prior findings.

### Phase 4 — Report (you, serial)

1. Consolidate the final 2 change requests into one authoritative change request
2. Write the final change request to `/tmp/<codeword>/change-request.md`
