# kb create

Ingest a new repository into the Mycroft knowledge base.

```
/kb create <repo>
```

- `repo` — GitHub URL or `owner/repo` of the repository to document

## Overview

Content emerges from source, not from templates. Each topic has a set of questions it must answer, but the articles that answer them are discovered by studying the repo. You define the brief and the boundaries. The funnel agents do the writing.

## Phases

### Phase 1 — Repo Study (you, serial)

1. Clone the target repo and read it: file tree, README, entrypoints, package manifests
2. Clone [Mycroft](https://github.com/alexander-thorwaldson/mycroft)
3. Generate a codeword and create the workspace under `/tmp/<codeword>/`
4. Write a **repo brief** — a concise architectural summary covering:
   - What the project is and what problem it solves
   - Key abstractions and their relationships
   - Core data flows
   - Major boundaries and interfaces
   - Directory structure mapped to responsibilities
5. Decide which of the 8 standard topics apply to this repo. Not every repo needs all of them. Drop topics that have no meaningful content to produce.
6. Create a branch: `kb/<repo-name>`

### Phase 2 — Mechanical Topics (parallel across topics)

Mechanical topics derive content directly from what exists in the code. Run these first.

For each applicable mechanical topic, spawn a convergent funnel (see SKILL.md). Each funnel agent receives:
- The repo brief from Phase 1
- The topic's question set
- The topic placement table
- Access to the repo source

Funnel agents produce articles — they decide what articles should exist and write them.

As each topic's funnel completes and you make your final review, commit it to the branch.

### Phase 3 — Interpretive Topics (parallel across topics)

Interpretive topics require the mechanical output as foundation and your architectural guidance as context.

For each applicable interpretive topic, spawn a convergent funnel. Each funnel agent receives everything from Phase 2 plus:
- The completed mechanical topic articles (committed in Phase 2)
- Your architectural guidance specific to this topic — what the agents should understand about the system's design to write well about it

As each topic's funnel completes and you make your final review, commit it to the branch.

### Phase 4 — PR (you, serial)

Open a PR against [Mycroft](https://github.com/alexander-thorwaldson/mycroft) with all committed topics.
