# kb article delete

Remove a single knowledge base article.

```
/kb article delete <article>
```

- `article` — path to the article in the Mycroft content directory

## Overview

Deletion is not just removing a file. Content is interconnected — an article may be referenced by others, may contain the only documentation for a concept, or may need to be relocated rather than removed. Verify the blast radius before you cut, and verify quality after.

## Phases

### Phase 1 — Before Checklist (parallel where possible)

Run these checks before removing anything:

| Check | How |
|---|---|
| **Cross-references** | Search all other articles for links or references to this article. List every article that depends on it. |
| **Unique coverage** | Determine if this article is the only place a concept is documented. If so, does that concept still need coverage elsewhere? |
| **Relocation vs removal** | Is this content misplaced rather than unnecessary? If it belongs in a different topic or article, it should move, not disappear. |
| **Read order impact** | Check the article's position in the topic's read order. Removing it may break the narrative flow — does the topic's guide article reference it? |
| **Downstream consumers** | Check if any other topics list this article in their cross-references or depend on it as prerequisite reading. |

Cross-references, unique coverage, and downstream consumers can run in parallel. Relocation and read order depend on those results.

If any check surfaces a reason to relocate rather than delete, stop and use `article update` or `article create` to move the content first.

### Phase 2 — Remove

1. Delete the article file
2. Renumber remaining articles in the topic to maintain read order continuity

### Phase 3 — After Checklist (parallel where possible)

Verify quality was not compromised:

| Check | How |
|---|---|
| **Broken references** | Search all articles for references that now point to nothing. Fix or remove them. |
| **Narrative gaps** | Read the topic's guide article (first in read order). Does the walkthrough still make sense without the removed article? If not, update it. |
| **Coverage gaps** | Review the topic's question set. Are all questions still answered by the remaining articles? If not, flag the gap. |
| **Cross-reference cleanup** | Remove or redirect any cross-references from other topics that pointed to the deleted article. |

Broken references, coverage gaps, and cross-reference cleanup can run in parallel. Narrative gap assessment depends on whether the guide article was affected.
