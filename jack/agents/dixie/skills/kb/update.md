# kb update

Apply changes to existing knowledge base content based on an approved change request.

```
/kb update <repo> <change-request>
```

- `repo` — GitHub URL or `owner/repo` of the documented repository
- `change-request` — path to a change request file (e.g. from maintain) or inline description of known changes

## Overview

Update requires a change request approved by you or the operator. It does not discover what needs changing — that is maintain's job. Update takes a defined scope and produces the actual article modifications through the convergent funnel.

## Phases

### Phase 1 — Scope (you, serial)

1. Clone the target repo and [Mycroft](https://github.com/alexander-thorwaldson/mycroft)
2. Read the change request and confirm the affected topics and articles
3. Generate a codeword and create the workspace under `/tmp/<codeword>/`
4. Create a branch: `kb/<repo-name>-update`

### Phase 2 — Content Production (parallel across affected topics)

For each affected topic, spawn a convergent funnel (see SKILL.md). This funnel produces article content.

Each funnel agent receives:
- The change request (scoped to this topic's relevant entries)
- The topic's existing KB articles
- The topic's question set
- The topic placement table
- Access to the target repo source

Each agent produces the actual updated articles — modified content, new articles for gaps, consolidated articles where redundancy was flagged. Output is the complete set of articles for the topic as they should exist after the update.

### Phase 3 — Funnel Convergence

The funnel follows the process defined in SKILL.md. Agents are refining article content.

Pass 2 agents study the repo source and existing KB first, then read prior pass outputs from `/tmp`. They synthesize a refined version of the updated articles.

### Phase 4 — Commit and PR (you, serial)

1. Review the final 2 versions, pick one, make small necessary edits only
2. Ensure `last_updated` and `sources` fields are current on all modified articles
3. Commit changes to the branch
4. Open a PR against [Mycroft](https://github.com/alexander-thorwaldson/mycroft)
