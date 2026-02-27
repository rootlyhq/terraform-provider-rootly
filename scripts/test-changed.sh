#!/usr/bin/env bash
# Run acceptance tests for resources/data sources changed in this branch.
#
# Usage:
#   scripts/test-changed.sh [BASE_REF] [extra go test flags]
#
# BASE_REF defaults to origin/master. Extra flags are passed to go test.
#
# Examples:
#   scripts/test-changed.sh
#   scripts/test-changed.sh origin/main -timeout 20m
#   scripts/test-changed.sh HEAD~1

set -euo pipefail

REPO_ROOT="$(git -C "$(dirname "$0")" rev-parse --show-toplevel)"
BASE_REF="${1:-origin/master}"
shift 2>/dev/null || true

# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

to_pascal() {
  # Split on underscores, capitalize the first letter of each segment, rejoin with no separator.
  # e.g. alert_routing_rule -> AlertRoutingRule
  echo "$1" | awk 'BEGIN{FS="_";OFS=""} {for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) substr($i,2); print}'
}

# Naive singularization of the last underscore-segment:
#   workflow_form_field_conditions -> workflow_form_field_condition
#   functionalities                -> functionality
#   status_pages                   -> status_page
singularize_last() {
  local name="$1"
  local prefix last

  if [[ "$name" == *_* ]]; then
    prefix="${name%_*}_"
    last="${name##*_}"
  else
    prefix=""
    last="$name"
  fi

  if [[ "$last" == *ies ]]; then
    last="${last%ies}y"
  elif [[ "$last" =~ (s|x|z|ch|sh)es$ ]]; then
    last="${last%es}"
  elif [[ "$last" == *[^s]s ]]; then
    last="${last%s}"
  fi

  echo "${prefix}${last}"
}

PATTERNS=()

add_pattern() {
  local p="$1"
  local existing
  for existing in "${PATTERNS[@]+"${PATTERNS[@]}"}"; do
    [[ "$existing" == "$p" ]] && return
  done
  PATTERNS+=("$p")
}

# Derive a test prefix from a resource/data-source name and add it.
# Also checks whether a corresponding test file actually exists so we don't
# emit patterns for resources that have no tests yet.
add_resource_pattern() {
  local kind="$1"   # "Resource" or "DataSource"
  local name="$2"   # snake_case resource name
  local prefix="TestAcc${kind}$(to_pascal "$name")"
  local file_prefix

  if [[ "$kind" == "Resource" ]]; then
    file_prefix="${REPO_ROOT}/provider/resource_${name}_test.go"
  else
    file_prefix="${REPO_ROOT}/provider/data_source_${name}_test.go"
  fi

  if [[ -f "$file_prefix" ]]; then
    add_pattern "$prefix"
  fi
}

# ---------------------------------------------------------------------------
# Get changed files
# ---------------------------------------------------------------------------

CHANGED=$(git diff --name-only "${BASE_REF}...HEAD" 2>/dev/null \
  || git diff --name-only HEAD~1 2>/dev/null \
  || true)

if [[ -z "$CHANGED" ]]; then
  echo "No changed files found relative to ${BASE_REF}."
  exit 0
fi

echo "Changed files vs ${BASE_REF}:"
echo "$CHANGED" | sed 's/^/  /'
echo ""

while IFS= read -r file; do
  # ── provider/resource_<name>_test.go or data_source_<name>_test.go ───────
  # Derive prefix directly from the test filename — covers both generated and
  # hand-rolled tests since they all follow the same naming convention.
  if [[ "$file" =~ ^provider/resource_([a-z0-9_]+)_test\.go$ ]]; then
    add_pattern "TestAccResource$(to_pascal "${BASH_REMATCH[1]}")"
    continue
  fi
  if [[ "$file" =~ ^provider/data_source_([a-z0-9_]+)_test\.go$ ]]; then
    add_pattern "TestAccDataSource$(to_pascal "${BASH_REMATCH[1]}")"
    continue
  fi

  # ── provider/resource_<name>.go ──────────────────────────────────────────
  if [[ "$file" =~ ^provider/resource_([a-z0-9_]+)\.go$ ]]; then
    add_resource_pattern "Resource" "${BASH_REMATCH[1]}"
    continue
  fi

  # ── provider/data_source_<name>.go ───────────────────────────────────────
  if [[ "$file" =~ ^provider/data_source_([a-z0-9_]+)\.go$ ]]; then
    add_resource_pattern "DataSource" "${BASH_REMATCH[1]}"
    continue
  fi

  # ── client/<plural>.go — singularize and check for test files ────────────
  if [[ "$file" =~ ^client/([a-z0-9_]+)\.go$ ]]; then
    plural="${BASH_REMATCH[1]}"
    singular=$(singularize_last "$plural")
    for candidate in "$singular" "$plural"; do
      add_resource_pattern "Resource"   "$candidate"
      add_resource_pattern "DataSource" "$candidate"
    done
    continue
  fi
done <<< "$CHANGED"

# ---------------------------------------------------------------------------
# Run
# ---------------------------------------------------------------------------

if [[ ${#PATTERNS[@]} -eq 0 ]]; then
  echo "No test patterns resolved from changed files — nothing to run."
  exit 0
fi

RUN_PATTERN=$(
  IFS='|'
  echo "${PATTERNS[*]}"
)

echo "Test pattern: ${RUN_PATTERN}"
echo ""

exec go test -v ./provider -run "${RUN_PATTERN}" -timeout 10m "$@"
