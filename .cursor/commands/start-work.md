Start a Sisyphus work session from a Prometheus plan.

Resume existing work or start fresh from the latest plan.

1. **Check for active boulder**: Read `.sisyphus/boulder.json`
   - If exists and has active plan: RESUME — report progress and continue from where we left off
   - If not exists: FRESH START — proceed to step 2

2. **Find latest plan**: Read the most recent `.sisyphus/plans/*.md` file
   - If no plans exist: tell user to run `/plan` first

3. **Initialize boulder state**: Write to `.sisyphus/boulder.json`:
   - activePlan: path to the plan file
   - startedAt: ISO timestamp
   - status: "in_progress"

4. **Switch to @atlas**: Delegate the full plan execution to @atlas with:
   - Plan file path
   - Boulder state for tracking
   - Any accumulated wisdom from `.sisyphus/notepads/`

5. @atlas handles: task analysis, notepad initialization, delegation to @hephaestus, verification, and Final Verification Wave

When complete, @atlas updates boulder.json with status "completed".
