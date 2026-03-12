Create a detailed context summary for continuing work in a new session.

Load the `handoff` skill and execute it.

1. **Gather State** (parallel):
   - git status, git diff --stat HEAD~5..HEAD, git log --oneline -10
   - Current todo list state
   - Active plans from `.sisyphus/plans/*.md`
   - Boulder state from `.sisyphus/boulder.json`

2. **Generate Handoff Document**: Write to `.sisyphus/handoffs/{YYYY-MM-DD}-{topic}.md` with:
   - What was done (completed work items)
   - What remains (pending tasks with checkboxes)
   - Key decisions made (with rationale)
   - Files modified (with descriptions)
   - Open questions and blockers
   - Active plan progress
   - Context for next session (2-3 sentence summary)
   - Git diff summary

3. **Update Boulder State**: Write lastHandoff path to `.sisyphus/boulder.json`

4. **Report**: Where the handoff was saved. The sessionStart hook will auto-inject it in the next session.
