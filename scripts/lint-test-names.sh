#!/usr/bin/env bash
# Lint: ensure acceptance tests don't use hardcoded resource names.
# Catches const HCL configs that create API resources with static name values.
# Allowed: data sources, resources without a name field, resources using dynamic names.
set -euo pipefail

EXIT_CODE=0

for file in provider/*_test.go; do
  consts=$(grep -n 'const testAcc.*= `' "$file" 2>/dev/null || true)
  [ -z "$consts" ] && continue

  while IFS= read -r line; do
    lineno=$(echo "$line" | cut -d: -f1)
    constname=$(echo "$line" | sed 's/.*const \(testAcc[^ ]*\).*/\1/')

    # Extract const body between backticks
    body=$(sed -n "${lineno},\$p" "$file" | sed -n '/= `/,/^`/p')

    # Flag if body creates a resource AND has a hardcoded name = "..." value
    if echo "$body" | grep -q 'resource "rootly_' && echo "$body" | grep -qE '^\s*(name|label|title)\s*=\s*"[^%]'; then
      echo "ERROR: $file:$lineno — const '$constname' creates resources with hardcoded names."
      echo "       Use acctest.RandomWithPrefix(\"tf-...\") and fmt.Sprintf instead."
      echo ""
      EXIT_CODE=1
    fi
  done <<< "$consts"
done

if [ "$EXIT_CODE" -eq 0 ]; then
  echo "OK: No hardcoded test resource names found."
fi

exit $EXIT_CODE
