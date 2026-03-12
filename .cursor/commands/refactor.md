Intelligent refactoring with verification at every step.

Load the `refactor` skill and execute it.

Arguments:
- `--scope=<file|module|project>` — refactoring scope (default: inferred from context)
- `--strategy=<safe|aggressive>` — safe preserves all behavior, aggressive allows minor API changes

1. **Impact Analysis**: Find ALL usages of the target via Grep, SemanticSearch, and ReadLints. Map call sites, test coverage, and risk.
2. **Pre-Refactor Baseline**: Run tests, lints, build. Record results as the behavioral baseline.
3. **Plan**: Create detailed todo list with specific changes per file, dependency order, and verification points.
4. **Execute Incrementally**: One file at a time. ReadLints after each change. Run tests every 2-3 files. STOP if behavior changes.
5. **Post-Refactor Verification**: Run full test suite, lints, build. Compare against baseline.
6. **Report**: Files modified, before/after comparison, behavior preservation verdict.

Safety rules:
- NEVER change behavior while restructuring
- NEVER modify tests to make them pass — if tests fail, the refactoring is wrong
- REVERT if behavior changes unintentionally
- After 3 failed attempts at same change, consult @oracle or abort
