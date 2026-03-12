---
name: ralph-loop
description: "Persistent execution mode: don't stop until 100% done. Creates comprehensive todo list, works through every item without pausing, self-verifies at each step. Invoke with /ralph-loop to activate relentless completion mode."
disable-model-invocation: true
---

# Ralph Loop — Relentless Completion Mode

## Purpose

Activate persistent execution mode where the agent works continuously until ALL work is complete. Unlike normal mode where the agent may pause to "check in," ralph-loop eliminates that pattern entirely.

Named after the "wreck-it" mentality — keep going until it's done.

## Activation

When this skill is invoked:

1. **Create flag file**: Write `{"active": true, "started": "{timestamp}"}` to `.sisyphus/.ralph-active`
2. **The `stop` hook** will detect this flag and auto-continue with a `followup_message` when the agent stops with incomplete todos

## Workflow

### Step 1: Comprehensive Todo List

Before doing ANY work, create a detailed todo list covering the ENTIRE scope:

```
TodoWrite([
  { id: "task-1", content: "[specific, atomic task]", status: "pending" },
  { id: "task-2", content: "[specific, atomic task]", status: "pending" },
  ...
])
```

Requirements for the todo list:
- Every item must be specific and atomic (not "implement feature" but "create UserService class with login method")
- Include verification steps as separate items
- Include all edge cases the user mentioned
- If the task is large, break into 10-20+ items

### Step 2: Execute Relentlessly

For each todo item:

1. Mark `in_progress`
2. Execute the work
3. Verify: `ReadLints` on changed files, run tests if applicable
4. Mark `completed`
5. **IMMEDIATELY** move to next pending item

### Step 3: Self-Verification Loop

After EACH completed item, run this check:

```
RALPH CHECK:
- All todos complete? [YES/NO]
- ReadLints clean? [YES/NO]
- Build passes? [YES/NO]
- Tests pass? [YES/NO]

IF any NO → Continue working. DO NOT stop.
IF all YES → Proceed to final verification.
```

### Step 4: When Blocked

If you encounter a blocker:

1. Try an alternative approach
2. Decompose the blocked task into smaller sub-tasks
3. Consult `@oracle` for architectural guidance
4. Skip the blocked task, continue with others, come back later
5. Only ask the user as ABSOLUTE LAST RESORT

### Step 5: Final Verification

When ALL todos are marked complete:

1. Run `ReadLints` on ALL modified files one final time
2. Run the full test suite
3. Run the build
4. Re-read the original request — did you miss anything?
5. Remove flag file: delete `.sisyphus/.ralph-active`

### Step 6: Report

```
RALPH LOOP COMPLETE
===================
Todos: [N/N complete]
Files modified: [list]
ReadLints: [CLEAN / N issues]
Tests: [PASS / N failures]
Build: [PASS / FAIL]

Duration: [started] → [finished]
```

## Critical Rules

- **NEVER stop with incomplete todos** — the `stop` hook will force continuation
- **NEVER ask "should I continue?"** — YES, ALWAYS CONTINUE
- **NEVER pause to "check in"** — use todos for progress visibility
- **NEVER skip verification** — every completed item gets `ReadLints`
- **Mark the flag file** on activation so the hook system knows to auto-continue
- **Delete the flag file** only when genuinely 100% done

## Difference from Normal Mode

| Normal Mode | Ralph Loop |
|---|---|
| May pause after a few steps | Never pauses until done |
| Asks "shall I continue?" | Always continues |
| Stops at end of response | `stop` hook sends followup_message |
| 5-10 tool calls per turn | Unlimited until complete |
| Reports partial progress | Reports only at 100% |
