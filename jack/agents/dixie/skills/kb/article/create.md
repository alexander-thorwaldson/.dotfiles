# kb article create

Write a single knowledge base article.

```
/kb article create <topic> <title> <scope> <out-of-scope>
```

- `topic` — the topic folder this article belongs to (e.g. `api`, `core`)
- `title` — the article title
- `scope` — description of what this article covers
- `out-of-scope` — what this article explicitly does not cover (documented elsewhere)

## Overview

Create produces one article. You receive the article assignment — what it covers, what it doesn't — and run the 3→2→1 funnel to produce it. The scope and out-of-scope boundaries are your guard against duplication.

## Phases

### Phase 1 — Setup (you, serial)

1. Generate a codeword and create the workspace under `/tmp/<codeword>/`
2. Study the repo source relevant to the article's scope

### Phase 2 — Drafts (3 agents, parallel)

Spawn 3 agents, one per lens (consumer, implementor, structural). Each agent receives:
- The article assignment: topic, title, scope, and out-of-scope list
- The topic's question set
- The topic placement table
- Access to the repo source

Each agent studies the source and produces a complete article. They must respect the scope boundaries — content that belongs in other articles is not their concern.

### Phase 3 — Edit (2 agents, parallel)

Two editor agents. Each must:
1. Study the repo source **first**
2. Read all 3 drafts from `/tmp` **second**
3. Synthesize — produce a version that takes the best of the drafts while staying within scope

Editors should cut anything that falls outside the article's defined scope.

### Phase 4 — Final Review (you, serial)

1. Read both editor outputs
2. Pick one, make small necessary edits only
3. Verify:
   - The article stays within its defined scope
   - Out-of-scope content hasn't crept in
   - Source references are accurate
   - Frontmatter is complete (`title`, `description`, `tags`, `sources`, `last_updated`)
4. Write the article to the Mycroft content directory
