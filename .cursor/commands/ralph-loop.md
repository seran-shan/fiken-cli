Start self-referential development loop until completion.

Load the `ralph-loop` skill and execute it against the current task.

This activates persistent execution mode — the agent works continuously until ALL work is 100% complete, with no pauses or check-ins.

1. Create flag file: `.sisyphus/.ralph-active` with `{"active": true, "started": "{timestamp}"}`
2. Build comprehensive todo list covering the ENTIRE scope (10-20+ items for large tasks)
3. Execute every item relentlessly: mark in_progress, execute, verify (ReadLints), mark completed, next
4. Self-verify after each item: all todos done? lints clean? build passes? tests pass?
5. When blocked: try alternative, decompose, consult @oracle, skip and return later
6. When ALL done: final verification (full test suite, build, re-read original request)
7. Remove flag file and report results

The `stop` hook auto-continues when incomplete todos exist. Default loop limit: 5 iterations.
