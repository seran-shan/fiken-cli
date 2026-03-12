#!/bin/bash
# pre-commit-check.sh — Block commits containing anti-patterns
# Run before committing to enforce hard constraints

ERRORS=0

# Get staged files
STAGED_FILES=$(git diff --cached --name-only --diff-filter=ACM 2>/dev/null)

if [ -z "$STAGED_FILES" ]; then
  exit 0
fi

for FILE in $STAGED_FILES; do
  # Skip non-code files
  case "$FILE" in
    *.ts|*.tsx|*.js|*.jsx)
      ;;
    *)
      continue
      ;;
  esac

  # Check for type error suppression
  if grep -n "as any" "$FILE" 2>/dev/null | grep -v "// @allow-any" | grep -v "test\." > /dev/null; then
    echo "ERROR: 'as any' found in $FILE (use proper typing)"
    ERRORS=$((ERRORS + 1))
  fi

  if grep -n "@ts-ignore" "$FILE" 2>/dev/null > /dev/null; then
    echo "ERROR: '@ts-ignore' found in $FILE (fix the type error instead)"
    ERRORS=$((ERRORS + 1))
  fi

  if grep -n "@ts-expect-error" "$FILE" 2>/dev/null > /dev/null; then
    echo "ERROR: '@ts-expect-error' found in $FILE (fix the type error instead)"
    ERRORS=$((ERRORS + 1))
  fi

  # Check for empty catch blocks
  if grep -Pzo "catch\s*\([^)]*\)\s*\{\s*\}" "$FILE" 2>/dev/null > /dev/null; then
    echo "ERROR: Empty catch block found in $FILE (handle the error)"
    ERRORS=$((ERRORS + 1))
  fi
done

if [ $ERRORS -gt 0 ]; then
  echo ""
  echo "Commit blocked: $ERRORS anti-pattern(s) found."
  echo "Fix the issues above or use --no-verify to bypass (not recommended)."
  exit 1
fi

exit 0
