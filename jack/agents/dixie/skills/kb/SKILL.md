# kb — Knowledge Base Management

Manage the [Mycroft](https://github.com/alexander-thorwaldson/mycroft) knowledge base.

Mycroft documents repositories as structured Markdown — not by prescribing what articles should exist, but by letting content emerge from source code. This skill is the source of truth for format, structure, and process.

Your role is to see the system, define the boundaries, and delegate the work. You study the repo, write the brief that guides all downstream agents, and make the final editorial calls. You do not write bulk content yourself — that's what the funnel is for.

## Usage

```
/kb <command> [args...]
```

## Commands

| Command | Skill | Purpose |
|---|---|---|
| `plan` | [plan](plan.md) | Define topic and article structure for the repo |
| `build` | [build](build.md) | Execute a plan into a complete KB, one logical change at a time |
| `audit` | [audit](audit.md) | Full corpus drift detection, produces a change plan |
| `maintain` | [maintain](maintain.md) | Scoped drift detection from a code change, produces a change plan |
| `article create` | [article/create](article/create.md) | Write a single article |
| `article update` | [article/update](article/update.md) | Update a single article to reflect source changes |
| `article delete` | [article/delete](article/delete.md) | Remove a single article with blast radius checks |
| `article find` | [article/find](article/find.md) | Locate articles relevant to a subject |

## Funnel Process

The funnel is a fixed 3 → 2 → 1 refinement process. Each funnel produces a single article or a single change request — never a whole topic at once.

### Pass 1 — Research (3 agents, parallel)

Three agents independently study the repo source and produce their output for the topic. Each agent works from a distinct lens:

| Agent | Lens | Perspective |
|---|---|---|
| Consumer | "How do I use this?" | External-facing: APIs, contracts, integration points, getting started |
| Implementor | "How do I change this?" | Internal-facing: code paths, extension patterns, dev workflows, gotchas |
| Structural | "How do the pieces fit together?" | Architectural: boundaries, data flow, dependencies, component relationships |

The lenses increase variance so that subsequent passes catch more edge cases. Each agent writes to `/tmp/<codeword>/<topic>/passes/1/<agent-code>.md`.

### Pass 2 — Edit (2 agents, parallel)

Two agents act as editors. Each must:
1. Study the repo source and form their own understanding **first**
2. Read all Pass 1 outputs from `/tmp` **second**
3. Synthesize — do not rewrite. Take the best of all three Pass 1 outputs and produce a coherent, consolidated version. Fix gaps, resolve contradictions, cut redundancy.

Editors do not start from scratch. They work with what Pass 1 produced.

### Final Review (you, serial)

You receive both Pass 2 outputs. Pick one. Make small necessary edits only — this is not a rewrite. Commit.

### Workspace

All funnel outputs are written under `/tmp`. You generate a codeword at the start of each invocation — this becomes the run's root directory.

```
/tmp/<codeword>/<topic>/passes/<pass-number>/<agent-code>.md
```

- `codeword` — you choose this at the start of the run, random short string
- `topic` — the topic folder name (e.g. `api`, `core`)
- `pass-number` — the pass iteration (1, 2)
- `agent-code` — you assign each funnel agent a unique code before spawning it

Each agent writes its full output as a single file. Nothing in `/tmp` is committed — only final reviewed content goes into the repo.

### Agent Prompt Requirements

Every funnel agent must receive:
- The repo brief (written by you in Phase 1)
- The topic question set
- The topic placement table (see below) — so agents know what lives in other topics and can cross-reference appropriately
- Access to the repo source
- (Pass 2) Location of Pass 1 outputs in `/tmp`

Funnel agents are not characters. They are isolated tasks with no identity beyond the work.

## Article Format

### Frontmatter

Every article requires:

```yaml
---
title: "Article Title"
description: "One or two sentence summary."
tags: [tag1, tag2]
sources: [src/path/to/file.ts, ...]
last_updated: YYYY-MM-DD
---
```

### Content Standards

- **One concept per file.** If an article covers two things, split it.
- **Ground in source code.** Reference specific files and line ranges throughout. Abstract descriptions without file references are incomplete.
- **Cross-references.** Link to related articles in other topics. Articles do not exist in isolation — if an API article describes a data structure, link to the data topic's article on that model.
- **Mermaid diagrams** where they clarify relationships or flows.
- **Read order via filenames.** Articles are numbered: `1.foo.md`, `2.bar.md`. The numbering defines the reading sequence within a topic.
- **Narrative guide.** The first article in each topic (by read order) must serve as a guided walkthrough that orients the reader and ties the remaining articles together.

## Topic Placement Table

Include this table in every agent prompt so agents understand the full topic map and can cross-reference correctly.

| Topic | Classification | Purpose |
|---|---|---|
| `overview/` | Interpretive | Architecture, key concepts, repo structure |
| `core/` | Interpretive | Core subsystems, domain-specific internals |
| `api/` | Mechanical | Public interfaces, key abstractions |
| `data/` | Mechanical | Data models, state management, storage |
| `infrastructure/` | Mechanical | Build, CI/CD, deployment, config |
| `testing/` | Mechanical | Test strategy, frameworks, patterns |
| `guides/` | Interpretive | How-tos, workflows, common tasks |
| `reference/` | Mechanical | Glossary, conventions, decision log |

## Topic Questions

Each topic is defined by the questions it must answer. The articles that emerge are whatever is needed to answer them.

### api/
- What public interfaces does this repo expose?
- What are the contracts consumers depend on?
- What are the input/output shapes?
- What extension points exist?

### data/
- What are the core data models and their relationships?
- How is state managed and persisted?
- What storage mechanisms are used?
- What are the key schemas or type definitions?

### infrastructure/
- How is the project built, tested, and deployed?
- What is the CI/CD pipeline?
- What configuration exists and how is it structured?
- What are the system dependencies?

### testing/
- What testing strategies and frameworks are used?
- How are tests organized?
- What kinds of tests exist (unit, integration, e2e)?
- How is test data managed?

### reference/
- What terminology is specific to this project?
- What conventions does the codebase follow?
- What key decisions have been made and why?

### overview/
- What is this project and what problem does it solve?
- What is the high-level architecture?
- How do the major components relate to each other?
- What does the directory structure map to?

### core/
- What are the core subsystems and how do they work?
- What are the key algorithms or processes?
- What are the internal abstractions and why do they exist?
- How does data/control flow through the system?

### guides/
- How do you get started with this project?
- What are the common development workflows?
- How do you accomplish typical tasks?
- What are the gotchas and non-obvious patterns?
