---
name: init-deep
description: Generate hierarchical AGENTS.md files throughout the project for per-directory agent context. Invoke with /init-deep to create comprehensive codebase documentation that Cursor auto-injects when working in each directory.
disable-model-invocation: true
---

# Init Deep — Hierarchical AGENTS.md Generator

## Purpose

Analyze the project directory structure and generate `AGENTS.md` files at each significant directory level. Cursor natively supports `AGENTS.md` — these files are automatically injected into agent context when working in that directory.

## Workflow

### Step 1: Analyze Project Structure

```bash
find . -type d -not -path '*/node_modules/*' -not -path '*/.git/*' -not -path '*/dist/*' -not -path '*/build/*' -not -path '*/.next/*' -not -path '*/.cursor/*' -maxdepth 3 | sort
```

Also examine:
- Package manager files (package.json, go.mod, requirements.txt, Cargo.toml)
- Config files (tsconfig.json, .eslintrc, Makefile)
- Existing README.md files for context

### Step 2: Generate Root AGENTS.md

Create `AGENTS.md` at the project root with:

```markdown
# [Project Name]

## Overview
[1-2 sentences: what this project does]

## Tech Stack
- [Language/runtime]
- [Framework]
- [Key dependencies]

## Project Structure
```
[directory tree with purpose annotations]
```

## Conventions
- [Naming patterns]
- [File organization rules]
- [Import conventions]
- [Error handling patterns]

## Build & Test
```bash
[build command]
[test command]
[lint command]
```

## Key Patterns
- [Pattern 1: description + example file]
- [Pattern 2: description + example file]

## Anti-Patterns
- [What NOT to do in this codebase]
```

### Step 3: Generate Subdirectory AGENTS.md Files

For each significant subdirectory (max depth 3), create an `AGENTS.md` covering:

```markdown
# [Directory Name]

## Purpose
[What this directory contains and why]

## Key Files
- `file.ts` — [what it does]
- `other.ts` — [what it does]

## Conventions
- [Directory-specific patterns]
- [Naming conventions here]

## Dependencies
- Depends on: [other directories/modules]
- Depended on by: [who uses this]

## Patterns to Follow
- [Specific patterns in this directory with examples]

## Testing
- Test files: [where tests live]
- Run: `[test command for this module]`
```

### Step 4: Verify

- Read each generated `AGENTS.md` to confirm accuracy
- Ensure no sensitive information is included
- Check that file references actually exist

## Rules

- **DO NOT generate for**: `node_modules`, `.git`, `dist`, `build`, `.next`, `.cursor`, vendor directories
- **DO generate for**: `src/`, `apps/`, `packages/`, `lib/`, `internal/`, `cmd/`, `test/`, `scripts/`, and their significant subdirectories
- **Max depth**: 3 levels by default
- **Skip empty directories**: Only generate for directories with meaningful code
- **Be specific**: Reference actual files and patterns, not generic advice
- **Read before writing**: Examine existing code to understand real patterns, don't guess
