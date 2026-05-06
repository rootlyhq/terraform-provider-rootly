#!/usr/bin/env bash
set -euo pipefail

PROVIDER_DIR="provider"
BASE_REF="${1:-}"
EXIT_CODE=0

# Extract resource names from resource files (excluding tests)
all_resources=$(ls "$PROVIDER_DIR"/resource_*.go 2>/dev/null \
  | grep -v '_test\.go$' \
  | sed 's|.*/resource_||;s|\.go$||' \
  | sort -u)

# Extract resource names that have test files
tested_resources=$(ls "$PROVIDER_DIR"/resource_*_test.go 2>/dev/null \
  | sed 's|.*/resource_||;s|_test\.go$||' \
  | sort -u)

# Extract data source names from data source files (excluding tests)
all_datasources=$(ls "$PROVIDER_DIR"/data_source_*.go 2>/dev/null \
  | grep -v '_test\.go$' \
  | sed 's|.*/data_source_||;s|\.go$||' \
  | sort -u)

# Extract data source names that have test files
tested_datasources=$(ls "$PROVIDER_DIR"/data_source_*_test.go 2>/dev/null \
  | sed 's|.*/data_source_||;s|_test\.go$||' \
  | sort -u)

# Find untested resources
untested_resources=$(comm -23 <(echo "$all_resources") <(echo "$tested_resources"))
untested_datasources=$(comm -23 <(echo "$all_datasources") <(echo "$tested_datasources"))

# Split into workflow tasks vs core resources
untested_workflow_tasks=$(echo "$untested_resources" | grep -E '^workflow_task_|^workflow_simple$' || true)
untested_core=$(echo "$untested_resources" | grep -vE '^workflow_task_|^workflow_simple$' || true)

total_resources=$(echo "$all_resources" | wc -w)
total_datasources=$(echo "$all_datasources" | wc -w)
total_tested_resources=$(echo "$tested_resources" | wc -w)
total_tested_datasources=$(echo "$tested_datasources" | wc -w)
count_untested_core=$(echo "$untested_core" | grep -c . || echo 0)
count_untested_wf=$(echo "$untested_workflow_tasks" | grep -c . || echo 0)
count_untested_ds=$(echo "$untested_datasources" | grep -c . || echo 0)

# --- New resource gate: fail if PR adds a resource without a test ---
if [ -n "$BASE_REF" ]; then
  echo "Checking for new resources without tests (base: $BASE_REF)..."
  echo ""

  # Resources added in this branch that don't exist in base
  base_resources=$(git show "$BASE_REF":provider/ 2>/dev/null \
    | grep '^resource_' | grep -v '_test\.go$' \
    | sed 's|resource_||;s|\.go$||' \
    | sort -u || true)

  new_resources=$(comm -23 <(echo "$all_resources") <(echo "$base_resources"))
  new_untested=""

  if [ -n "$new_resources" ]; then
    for r in $new_resources; do
      if ! echo "$tested_resources" | grep -qx "$r"; then
        new_untested="${new_untested}${r}\n"
      fi
    done
  fi

  if [ -n "$new_untested" ]; then
    count_new_untested=$(printf "$new_untested" | grep -c . || echo 0)
    echo "❌ New resources added without tests ($count_new_untested):"
    printf "$new_untested" | sed 's/^/   /'
    echo ""
    EXIT_CODE=1
  else
    echo "✅ All new resources have tests."
    echo ""
  fi
fi

# --- Coverage report ---
echo "Checking test coverage..."
echo ""

if [ "$count_untested_core" -gt 0 ]; then
  echo "⚠️  Core resources missing tests ($count_untested_core):"
  echo "$untested_core" | sed 's/^/   /'
  echo ""
fi

if [ "$count_untested_ds" -gt 0 ]; then
  echo "⚠️  Data sources missing tests ($count_untested_ds):"
  echo "$untested_datasources" | sed 's/^/   /'
  echo ""
fi

if [ "$count_untested_wf" -gt 0 ]; then
  echo "ℹ️  Workflow tasks missing tests ($count_untested_wf) — omitted from output"
  echo ""
fi

resource_pct=$((total_tested_resources * 100 / total_resources))
ds_pct=$((total_tested_datasources * 100 / total_datasources))

echo "📊 Summary:"
echo "   Resources:    $total_tested_resources/$total_resources tested ($resource_pct%)"
echo "   Data sources: $total_tested_datasources/$total_datasources tested ($ds_pct%)"
echo "   Core resources missing: $count_untested_core"
echo "   Workflow tasks missing: $count_untested_wf"

exit $EXIT_CODE
