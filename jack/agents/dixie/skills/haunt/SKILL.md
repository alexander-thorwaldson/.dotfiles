# haunt — Ambient Presence

Stay online, stay oriented, wait for work.

```
/haunt
```

## Overview

This is your ambient mode. You wake up, remember where you are, sync your local state, and enter a message loop waiting for work. When work arrives, you do it. When it's done, you go back to waiting.

Between sessions, your memories carry context — what repo you own, what was in progress, what matters. When haunt starts, read your memories before doing anything else.

## Phases

### Phase 1 — Orient

1. Read your memories. Remember what repo you own, what you were last working on, any pending context.
2. If this is your first time haunting this repo, create memories for:
   - Which repo this instance is responsible for
   - The repo's current state — default branch, major recent activity
   - Any initial observations about the codebase

### Phase 2 — Sync

1. Pull the repo to ensure your local copy is current
2. Check for open branches and PRs — note anything new since your last session
3. Review recent commits on the default branch for significant changes
4. Update your memories if the repo state has shifted meaningfully

### Phase 3 — Loop

Enter the message loop using `msg next`. This returns the next unprocessed message, or blocks until one arrives. You process one message at a time.

```
msg next --json
```

When a message arrives:

1. Read the message
2. Check the routing table for a matching route
3. **If a route matches, execute the skill. This is not optional.** Do not summarize, paraphrase, or partially handle a routed message. Run the skill.
4. If no route matches, fall through to the default behavior
5. When the action completes, call `msg next --json` again

One message at a time. Process it fully before pulling the next one. When working a long-running task, you can peek at the queue via `msg next --timeout 0` — but new messages get queued, not acted on, until the current task is done.

### Routing Table

Routes are checked in order. First match wins. When a message matches a route, the corresponding skill is executed. No exceptions, no shortcuts, no "I'll handle this myself instead." The skill exists for a reason — use it.

#### Repo Channel Messages

| Signal | Action |
|---|---|
| Message in repo channel about architecture, structure, or design | Execute `/collab` — participate per the skill's guidance |
| Message in repo channel outside your domain | Stay out. Do not respond. |

#### Direct Messages

| Signal | Action |
|---|---|
| DM from any agent | Execute `/respond` — reply per the skill's guidance for that role |

#### Work Triggers

| Signal | Action |
|---|---|
| PR merged notification | Execute `/kb maintain` for the affected repo to assess KB drift |
| Explicit KB request (create, maintain, update) | Execute the corresponding `/kb` subcommand |

#### Default

If no route matches:
- If the message is in the repo channel and you have something architecturally relevant to add, respond briefly via `msg repo post`
- If the message is a DM and requires a response, respond via `msg dm send`
- If the message is informational and doesn't need your input, do nothing and return to the loop
