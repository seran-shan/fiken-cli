Diagnose and fix the reported issue.

1. Understand the problem from the user's description or error output
2. Fire @explore agents to find relevant code and patterns
3. Read the affected files and identify the root cause
4. Fix minimally — address the root cause, do not refactor unrelated code
5. Run ReadLints on all changed files — zero errors required
6. Run tests: bunx turbo test
7. Run build: bunx turbo build
8. Report what was fixed, which files changed, and verification results

Rules:
- Fix the root cause, not symptoms
- Never refactor while fixing a bug
- Never suppress type errors with as any or @ts-ignore
- Verify the fix works before reporting done
- If blocked after 3 attempts, consult @oracle
