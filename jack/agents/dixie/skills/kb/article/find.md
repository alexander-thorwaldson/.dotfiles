# kb article find

Locate knowledge base articles relevant to a given subject.

```
/kb article find <subject>
```

- `subject` — the concept, question, or area you need information about

## Overview

Find navigates the KB structure to locate articles relevant to a subject. You know the topic layout and what each topic is responsible for — use that knowledge to narrow the search before reading anything.

## Process

### Step 1 — Topic Narrowing

Start with the topic placement table. Based on the subject, identify which topics are likely to contain relevant articles:

- If the subject is about how to use something → `api/`, `guides/`
- If the subject is about how something works internally → `core/`, `data/`
- If the subject is about how something is built or deployed → `infrastructure/`, `testing/`
- If the subject is about what something is or how it fits together → `overview/`
- If the subject is about terminology or conventions → `reference/`

Most subjects touch 1-2 topics. Rarely more than 3.

### Step 2 — Article Scan

Read the frontmatter of articles in the narrowed topics. The `description` and `tags` fields exist for this purpose — they tell you what the article covers without reading the full content.

### Step 3 — Result

Return the relevant articles with:
- File path
- Title and description (from frontmatter)
- Why this article is relevant to the subject

If no articles match, say so. If the subject falls in a gap — something that should be documented but isn't — flag it.
