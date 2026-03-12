Generate hierarchical AGENTS.md files throughout the project.

Load the `init-deep` skill and execute it.

Arguments:
- `--max-depth=N` (default: 3) — maximum directory depth for AGENTS.md generation
- `--create-new` — only create new AGENTS.md files, don't update existing ones

1. Analyze the full project directory structure (excluding node_modules, .git, dist, build, .next, .cursor)
2. Generate root AGENTS.md with: overview, tech stack, project structure, conventions, build/test commands, key patterns, anti-patterns
3. Generate subdirectory AGENTS.md files for each significant directory (max depth 3)
4. Each subdirectory AGENTS.md includes: purpose, key files, conventions, dependencies, patterns, testing
5. Verify all generated files: read each one, confirm file references exist, ensure no sensitive info
6. Report: number of files created/updated, directories covered
