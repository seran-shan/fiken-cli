Start ultrawork loop — maximum intensity execution until completion.

Combines ralph-loop (persistent execution) with ultrawork mode (maximum parallelism and delegation).

1. Activate ultrawork mode: fire parallel @explore agents aggressively (5+ concurrent), delegate everything
2. Create flag file: `.sisyphus/.ralph-active` with `{"active": true, "started": "{timestamp}", "mode": "ultrawork"}`
3. Build comprehensive todo list with granular steps (15-30+ items)
4. Execute relentlessly with maximum parallelism:
   - Fire multiple @hephaestus delegations simultaneously for independent tasks
   - Use @explore/@librarian in background for any research needs
   - Self-verify at every step: ReadLints, tests, build
5. When blocked: decompose, try alternative, consult @oracle, skip and return
6. Before marking complete: consult @oracle for final review of all changes
7. Final verification: full test suite, build, ReadLints on ALL modified files
8. Remove flag file and report

**Oracle verification gate**: Before declaring done, spawn @oracle to review the full set of changes. If Oracle raises concerns, address them before completing.

The `stop` hook auto-continues. No pauses. No check-ins. Maximum throughput until done.
