Stop all continuation mechanisms for this session.

Halts ralph-loop, todo-continuation, and boulder state tracking.

1. Delete `.sisyphus/.ralph-active` if it exists (stops ralph-loop)
2. Update `.sisyphus/boulder.json` to set status to "paused" (stops boulder continuation)
3. Report what was stopped and current state:
   - Ralph loop: was it active? How many todos completed?
   - Boulder: was there an active plan? Progress so far?
   - Todo state: how many pending/in-progress/completed?

After this command, the agent will stop at the end of its current response without auto-continuing.

To resume later, use `/start-work` to pick up from the boulder state.
