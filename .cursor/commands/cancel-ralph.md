Cancel the active Ralph Loop.

1. Check if `.sisyphus/.ralph-active` exists
2. If it exists: delete the flag file
3. Report current todo state: how many completed, how many remaining
4. The `stop` hook will no longer auto-continue since the flag file is removed

This immediately stops the ralph-loop/ulw-loop continuation cycle. Any in-progress work completes its current step but does not auto-continue to the next.
