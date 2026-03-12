#!/bin/bash
# todo-continuation.sh — Stop hook for auto-continuation
#
# When the agent stops with "completed" status, this hook checks if there are
# incomplete todos. If so, it returns a followup_message to auto-continue.
#
# Also supports ralph-loop mode: when .sisyphus/.ralph-active exists,
# the hook is more aggressive about continuation.

input=$(cat)

status=$(echo "$input" | jq -r '.status // "unknown"')
loop_count=$(echo "$input" | jq -r '.loop_count // 0')
transcript_path=$(echo "$input" | jq -r '.transcript_path // ""')

# Only auto-continue on "completed" status (not aborted or error)
if [ "$status" != "completed" ]; then
  echo '{}'
  exit 0
fi

# Check if ralph-loop mode is active
ralph_active=false
if [ -f ".sisyphus/.ralph-active" ]; then
  ralph_active=true
fi

# Default loop limit is 5 (configured in hooks.json)
# In ralph-loop mode, we rely on the hook's loop_limit config
# but add our own safety check at 25 iterations
if [ "$ralph_active" = true ] && [ "$loop_count" -ge 25 ]; then
  # Safety valve: even ralph-loop stops after 25 iterations
  rm -f ".sisyphus/.ralph-active"
  echo '{}'
  exit 0
fi

# Check for incomplete todos by scanning the transcript
# The transcript contains TodoWrite calls with status fields
has_incomplete_todos=false

if [ -n "$transcript_path" ] && [ -f "$transcript_path" ]; then
  # Look for the most recent TodoWrite in the transcript
  # Check if any todo has status "pending" or "in_progress"
  last_todos=$(grep -o '"status":"pending"\|"status":"in_progress"' "$transcript_path" 2>/dev/null | tail -5)
  if [ -n "$last_todos" ]; then
    has_incomplete_todos=true
  fi
fi

if [ "$has_incomplete_todos" = true ]; then
  if [ "$ralph_active" = true ]; then
    cat << 'EOF'
{
  "followup_message": "Ralph-loop active. You have incomplete todos. Continue working on the next pending item immediately. Do not stop until all todos are complete."
}
EOF
  else
    cat << 'EOF'
{
  "followup_message": "You have incomplete todos. Continue working on the next pending item."
}
EOF
  fi
else
  echo '{}'
fi

exit 0
