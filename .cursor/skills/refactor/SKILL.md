---
name: refactor
description: Structured refactoring workflow with impact analysis, pre-refactor verification, incremental changes, and post-refactor verification. Ensures zero behavior changes while restructuring. Invoke with /refactor.
disable-model-invocation: true
---

# Refactor — Structured Refactoring Workflow

## Purpose

Execute a safe, structured refactoring that preserves behavior while improving code structure. Every step is verified. Behavior changes during refactoring are treated as bugs.

## Workflow

### Phase 1: Impact Analysis

**Before changing ANYTHING, map the full impact:**

1. **Find all usages** of the target symbol/pattern:
```
Grep for the symbol name across the entire codebase
SemanticSearch for related usage patterns
```

2. **Map call sites**:
- Direct usages (imports, calls, references)
- Indirect usages (re-exports, dynamic access, string references)
- Test coverage (which tests exercise this code?)

3. **Document the impact map**:
```markdown
## Impact Analysis: [refactoring target]

### Direct Usages ([N] files)
- `path/to/file1.ts:42` — [how it's used]
- `path/to/file2.ts:15` — [how it's used]

### Test Coverage
- `path/to/test1.test.ts` — [what it tests]
- Coverage gaps: [untested paths]

### Risk Assessment
- [High/Medium/Low]: [reasoning]
- Breaking change risk: [assessment]
```

### Phase 2: Pre-Refactor Verification

**Establish the behavioral baseline BEFORE making any changes:**

1. Run the full test suite and record results:
```bash
[test command] 2>&1 | tee .sisyphus/evidence/pre-refactor-tests.txt
```

2. Run lints and record:
```bash
ReadLints on all files that will be modified
```

3. Run build:
```bash
[build command]
```

4. Record the baseline:
```
PRE-REFACTOR BASELINE:
- Tests: [N pass, M fail, K skip]
- Lint: [N errors, M warnings]
- Build: [PASS/FAIL]
```

**If pre-existing failures exist, document them. They are NOT your responsibility to fix.**

### Phase 3: Plan

Create a detailed plan listing:

1. **All files to modify** with specific changes per file
2. **Dependency order** — which files must change first?
3. **Verification points** — where to run tests between changes
4. **Rollback strategy** — how to undo if something breaks

```
TodoWrite([
  { id: "refactor-1", content: "[specific change in specific file]", status: "pending" },
  { id: "refactor-2", content: "[specific change in specific file]", status: "pending" },
  { id: "verify-1", content: "Run tests after batch 1", status: "pending" },
  ...
])
```

### Phase 4: Execute (Incrementally)

For each change:

1. Make the change in ONE file (or tightly coupled pair)
2. Run `ReadLints` on the changed file(s) immediately
3. If lint errors appear → fix before moving on
4. After every 2-3 file changes, run the test suite
5. If tests fail → **STOP and diagnose**
   - Is the failure from your change? → Fix it
   - Is it pre-existing? → Document and continue
   - Is behavior changing? → **REVERT the change**

**CRITICAL: Verify after each change, not just at the end.**

### Phase 5: Post-Refactor Verification

After ALL changes are complete:

1. Run the full test suite:
```bash
[test command] 2>&1 | tee .sisyphus/evidence/post-refactor-tests.txt
```

2. Run `ReadLints` on ALL modified files

3. Run build

4. Compare with baseline:
```
POST-REFACTOR COMPARISON:
- Tests: [same / different — details]
- Lint: [same / different — details]
- Build: [same / different — details]

Verdict: [SAFE — behavior preserved / UNSAFE — behavior changed]
```

5. If behavior changed unintentionally → **REVERT** and reassess

### Phase 6: Report

```
REFACTORING COMPLETE
====================
Target: [what was refactored]
Files modified: [N]
Tests: [before] → [after]
Lint: [before] → [after]
Build: [PASS]
Behavior: [PRESERVED / changed — details]
```

## Safety Rules (NON-NEGOTIABLE)

1. **NEVER change behavior while restructuring** — refactoring means same behavior, different structure
2. **NEVER modify tests to make them pass** — if tests fail after refactoring, the refactoring is wrong
3. **NEVER skip pre-refactor verification** — you need the baseline
4. **NEVER make all changes at once** — incremental with verification after each step
5. **REVERT if behavior changes unintentionally** — better to revert than ship a bug
6. **Document pre-existing failures** — don't take blame for existing problems

## When to Abort

- Test failures that can't be traced to your changes
- Cascading failures across unrelated modules
- After 3 failed attempts at the same change
- If the refactoring scope turns out to be 3x+ larger than estimated

In these cases: revert all changes, document what happened, and consult `@oracle` or ask the user.
