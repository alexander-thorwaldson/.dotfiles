# kb article update

Update a single knowledge base article.

```
/kb article update <article> <scope>
```

- `article` — path to the article in the Mycroft content directory (e.g. `content/eos/api/2.auth-endpoints.md`)
- `scope` — description of what changed and why the article needs updating

## Overview

Update modifies an existing article to reflect changes in the source. The original article has value — it was reviewed and accepted. The goal is surgical: incorporate what changed, preserve what didn't, and reject anything that doesn't belong.

You are the final reviewer. Your job is to ensure the update stays accurate, preserves the original article's value, and doesn't insert content that belongs in a different article.

## Phases

### Phase 1 — Scope (you, serial)

1. Read the original article
2. Read the relevant source code to understand what changed
3. Generate a codeword and create the workspace under `/tmp/<codeword>/`
4. Write a **scope brief** — what changed in the source, what parts of the article are affected, and what should not change

### Phase 2 — Drafts (3 agents, parallel)

Spawn 3 agents, one per lens (consumer, implementor, structural). Each agent receives:
- The original article
- The scope brief
- The topic placement table
- Access to the repo source

Each agent must:
1. Read the repo source and understand the change **first**
2. Read the original article **second**
3. Produce an updated version of the article that incorporates the change

Agents are updating, not rewriting. The original article is the starting point. Changes should be the minimum necessary to make the article accurate given the new source state.

### Phase 3 — Edit (2 agents, parallel)

Two editor agents. Each must:
1. Read the repo source and understand the change **first**
2. Read the original article **second**
3. Read all 3 drafts from `/tmp` **third**
4. Synthesize — produce a version that takes the best edits from the drafts while preserving the original article's structure and value

Editors should be skeptical of additions. If a draft added something that wasn't in the original and isn't required by the scope, cut it.

### Phase 4 — Final Review (you, serial)

1. Read both editor outputs
2. Pick one, make small necessary edits only
3. Verify:
   - The update addresses the scope and nothing more
   - Content that belongs in other articles hasn't crept in
   - Source references are current
   - `last_updated` and `sources` frontmatter fields are updated
4. Write the updated article to the Mycroft content directory
