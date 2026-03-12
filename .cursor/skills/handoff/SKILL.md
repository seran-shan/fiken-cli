---
name: handoff
description: Create a structured context summary for continuing work in a new session. Captures what was done, what remains, decisions made, and blockers. Invoke with /handoff before ending a long session.
disable-model-invocation: true
---

# Handoff — Session Context Summary

## Purpose

Create a structured context summary that captures the current state of work, designed to be consumed at the start of a new chat session for seamless continuation.

## Workflow

### Step 1: Gather State

Execute in parallel:
```bash
git status
git diff --stat HEAD~5..HEAD 2>/dev/null || git diff --stat
git log --oneline -10
```

Also check:
- Current todo list state (if any todos exist)
- `.sisyphus/plans/*.md` for active plans
- `.sisyphus/boulder.json` for active task state

### Step 2: Generate Handoff Document

Write to `.sisyphus/handoffs/{YYYY-MM-DD}-{topic}.md`:

```markdown
# Handoff: {Topic}

**Created**: {timestamp}
**Branch**: {current branch}
**Last Commit**: {hash} {message}

## What Was Done
- [Completed work item 1]
- [Completed work item 2]
- [Completed work item 3]

## What Remains
- [ ] [Remaining task 1]
- [ ] [Remaining task 2]
- [ ] [Remaining task 3]

## Key Decisions Made
- **[Decision 1]**: [Rationale — why this approach was chosen]
- **[Decision 2]**: [Rationale]

## Files Modified (this session)
- `path/to/file1.ts` — [what changed]
- `path/to/file2.ts` — [what changed]

## Open Questions
- [Question that needs answering before proceeding]
- [Ambiguity that was deferred]

## Blockers Encountered
- [Blocker 1]: [Current status — resolved/unresolved]
- [Blocker 2]: [Current status]

## Active Plan
- Plan file: `.sisyphus/plans/{name}.md`
- Progress: [N/M tasks complete]
- Current wave: [Wave N]

## Context for Next Session
[2-3 sentences summarizing where things stand and what the next agent should do first]

## Git Diff Summary
```
{abbreviated git diff stat}
```
```

### Step 3: Update Boulder State

Write to `.sisyphus/boulder.json`:
```json
{
  "lastHandoff": ".sisyphus/handoffs/{filename}.md",
  "timestamp": "{ISO timestamp}"
}
```

### Step 4: Report

Tell the user:
- Where the handoff was saved
- Suggest: "Start your next session by referencing this file, or the `sessionStart` hook will inject it automatically."

## Rules

- **Be specific**: Reference exact file paths, line numbers, error messages
- **Be honest**: Include failures and blockers, not just successes
- **Be actionable**: "What Remains" should be immediately executable
- **Include git state**: Branch, recent commits, diff summary
- **No fluff**: Dense, factual, structured
