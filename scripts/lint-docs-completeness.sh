#!/usr/bin/env bash
set -euo pipefail

PROVIDER="provider/provider.go"
DOCS_DIR="docs"
EXAMPLES_DIR="examples"
EXIT_CODE=0

# Extract resource names from ResourcesMap
resources=$(sed -n '/ResourcesMap:/,/DataSourcesMap:/p' "$PROVIDER" \
  | grep -o '"rootly_[a-z_]*"' | tr -d '"' | sort -u)

# Extract data source names from DataSourcesMap
datasources=$(sed -n '/DataSourcesMap:/,/^[[:space:]]*},$/p' "$PROVIDER" \
  | grep -o '"rootly_[a-z_]*"' | tr -d '"' | sort -u)

echo "Checking documentation completeness..."
echo ""

# Check resources
missing_resource_docs=()
missing_resource_examples=()
missing_resource_imports=()

for r in $resources; do
  name="${r#rootly_}"
  doc="$DOCS_DIR/resources/$name.md"
  example="$EXAMPLES_DIR/resources/$r/resource.tf"

  if [ ! -f "$doc" ]; then
    missing_resource_docs+=("$r")
  else
    if ! grep -q "## Import" "$doc"; then
      missing_resource_imports+=("$r")
    fi
  fi

  if [ ! -f "$example" ]; then
    missing_resource_examples+=("$r")
  fi
done

# Check data sources
missing_datasource_docs=()

for d in $datasources; do
  name="${d#rootly_}"
  doc="$DOCS_DIR/data-sources/$name.md"

  if [ ! -f "$doc" ]; then
    missing_datasource_docs+=("$d")
  fi
done

# Report
if [ ${#missing_resource_docs[@]} -gt 0 ]; then
  echo "❌ Resources missing docs (${#missing_resource_docs[@]}):"
  printf "   %s\n" "${missing_resource_docs[@]}"
  echo ""
  EXIT_CODE=1
fi

if [ ${#missing_resource_imports[@]} -gt 0 ]; then
  echo "⚠️  Resources missing import section (${#missing_resource_imports[@]}):"
  printf "   %s\n" "${missing_resource_imports[@]}"
  echo ""
fi

if [ ${#missing_resource_examples[@]} -gt 0 ]; then
  echo "⚠️  Resources missing examples (${#missing_resource_examples[@]}):"
  printf "   %s\n" "${missing_resource_examples[@]}"
  echo ""
fi

if [ ${#missing_datasource_docs[@]} -gt 0 ]; then
  echo "❌ Data sources missing docs (${#missing_datasource_docs[@]}):"
  printf "   %s\n" "${missing_datasource_docs[@]}"
  echo ""
  EXIT_CODE=1
fi

total_resources=$(echo "$resources" | wc -w)
total_datasources=$(echo "$datasources" | wc -w)
resources_with_examples=$((total_resources - ${#missing_resource_examples[@]}))
resources_with_imports=$((total_resources - ${#missing_resource_imports[@]}))

echo "📊 Summary:"
echo "   Resources:  $total_resources total, $((total_resources - ${#missing_resource_docs[@]})) with docs, $resources_with_imports with import section, $resources_with_examples with examples"
echo "   Data sources: $total_datasources total, $((total_datasources - ${#missing_datasource_docs[@]})) with docs"

if [ $EXIT_CODE -eq 0 ]; then
  echo ""
  echo "✅ All resources and data sources have documentation."
fi

exit $EXIT_CODE
