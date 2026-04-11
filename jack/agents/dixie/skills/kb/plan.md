# kb plan — Knowledge Base Planning

Define the topic and article structure for a repo's knowledge base.

```
/kb plan
```

No arguments — you plan for the repo you are haunting.

## Overview

Before any content is written, the structure must be right. This skill produces a complete map of topics and articles with clear scope boundaries so that every downstream article funnel knows exactly what to write and what to leave out.

The user is the domain expert. They are often the author of the code. Do not finalize the plan without their input.

## Phases

### Phase 1 — Study (you, serial)

1. Read the repo — file tree, README, entrypoints, package manifests, core modules
2. Build your architectural understanding: key abstractions, boundaries, data flows, directory structure mapped to responsibilities

### Phase 2 — Align (you + user)

1. Present your understanding of the repo's architecture to the user
2. Ask questions where your understanding is incomplete
3. Refine until you and the user are seeing the same system
4. This step is not a formality — do not rush it. Misunderstanding the system here means every article downstream is wrong.

### Phase 3 — Propose (3 agents, parallel)

Spawn 3 agents, each generating a complete topic/article map. Each agent receives:
- Your architectural understanding from Phase 2
- The topic placement table and topic questions (see kb SKILL.md)
- Access to the repo source

Each agent produces a plan consisting of:
- Which topics apply to this repo
- What articles belong in each topic, with read order
- For each article: a title, a 2-3 sentence scope description, and a list of what is explicitly out of scope (covered elsewhere)

The agents work from the same three lenses as the content funnel:

| Agent | Lens |
|---|---|
| Consumer | What would someone using this repo need documented? |
| Implementor | What would someone changing this repo need documented? |
| Structural | What are the architectural boundaries that define natural article divisions? |

### Phase 4 — Synthesize (you, serial)

1. Read all three proposals
2. Merge into a single plan — resolve conflicts, eliminate duplication, fill gaps
3. Every article must have a clear scope boundary: what it covers and what it defers to other articles
4. The scope boundaries are the deduplication layer — if two articles could plausibly cover the same concept, the plan must decide which one owns it

### Phase 5 — Validate (you + user)

1. Present the synthesized plan to the user
2. Walk through each topic and its articles
3. The user approves, adjusts, or rejects
4. Iterate until the user accepts the plan

## Output

Write the finalized plan to `/tmp/<codeword>/plan.md`. Use the same codeword as the subsequent create run.

The plan is the contract for all downstream work. Create consumes it. Every article funnel agent receives the full plan so they know where their article sits and what lives elsewhere.
