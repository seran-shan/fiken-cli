#!/bin/bash
# session-init.sh — Session start context injection
#
# Runs at session start to inject project-specific context:
# - Reads .sisyphus/boulder.json (current active task state)
# - Reads most recent handoff file
# - Reads most recent active plan
# Returns { "additional_context": "..." } with assembled state

context_parts=()

# Check for boulder.json (active task state)
if [ -f ".sisyphus/boulder.json" ]; then
  boulder_content=$(cat ".sisyphus/boulder.json" 2>/dev/null)
  if [ -n "$boulder_content" ] && [ "$boulder_content" != "{}" ]; then
    context_parts+=("## Active Task State (.sisyphus/boulder.json)")
    context_parts+=("$boulder_content")
    context_parts+=("")
  fi
fi

# Check for most recent handoff
if [ -d ".sisyphus/handoffs" ]; then
  latest_handoff=$(ls -t .sisyphus/handoffs/*.md 2>/dev/null | head -1)
  if [ -n "$latest_handoff" ]; then
    handoff_content=$(head -30 "$latest_handoff" 2>/dev/null)
    if [ -n "$handoff_content" ]; then
      context_parts+=("## Latest Handoff ($latest_handoff)")
      context_parts+=("$handoff_content")
      context_parts+=("...")
      context_parts+=("")
    fi
  fi
fi

# Check for most recent plan
if [ -d ".sisyphus/plans" ]; then
  latest_plan=$(ls -t .sisyphus/plans/*.md 2>/dev/null | head -1)
  if [ -n "$latest_plan" ]; then
    # Count incomplete tasks
    total_tasks=$(grep -c '^\- \[' "$latest_plan" 2>/dev/null || echo "0")
    incomplete_tasks=$(grep -c '^\- \[ \]' "$latest_plan" 2>/dev/null || echo "0")
    complete_tasks=$(grep -c '^\- \[x\]' "$latest_plan" 2>/dev/null || echo "0")

    # Read the TL;DR section (first ~15 lines after the title)
    plan_summary=$(head -20 "$latest_plan" 2>/dev/null)

    context_parts+=("## Active Plan ($latest_plan)")
    context_parts+=("Progress: $complete_tasks/$total_tasks tasks complete, $incomplete_tasks remaining")
    context_parts+=("$plan_summary")
    context_parts+=("")
  fi
fi

# Check for ralph-loop active flag
if [ -f ".sisyphus/.ralph-active" ]; then
  context_parts+=("## Ralph Loop ACTIVE")
  context_parts+=("Persistent execution mode is enabled. Do not stop until all todos are complete.")
  context_parts+=("")
fi

# Read stdin (required by hook protocol)
cat > /dev/null

# Build response
if [ ${#context_parts[@]} -eq 0 ]; then
  echo '{}'
  exit 0
fi

# Join context parts with newlines and escape for JSON
joined_context=""
for part in "${context_parts[@]}"; do
  if [ -n "$joined_context" ]; then
    joined_context="$joined_context\n$part"
  else
    joined_context="$part"
  fi
done

# Use jq to properly escape the string for JSON output
echo "$joined_context" | jq -Rs '{ "additional_context": . }'

exit 0
