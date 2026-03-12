#!/bin/bash
# post-edit-lint.sh — Auto-lint after file edits
# Triggered by afterFileEdit hook
# Detects file type and runs the appropriate linter

FILE="$1"
EXT="${FILE##*.}"

case "$EXT" in
  ts|tsx|js|jsx)
    # TypeScript/JavaScript — run ESLint if available
    if command -v npx &> /dev/null && [ -f "node_modules/.bin/eslint" ]; then
      npx eslint --no-error-on-unmatched-pattern --format compact "$FILE" 2>/dev/null
    fi
    ;;
  go)
    # Go — run go vet
    if command -v go &> /dev/null; then
      go vet "$FILE" 2>&1
    fi
    ;;
  css|scss)
    # CSS — run stylelint if available
    if command -v npx &> /dev/null && [ -f "node_modules/.bin/stylelint" ]; then
      npx stylelint "$FILE" 2>/dev/null
    fi
    ;;
esac

# Always exit 0 — linting is advisory, not blocking
exit 0
