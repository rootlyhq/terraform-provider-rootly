<!--
Use a Conventional Commit PR title: <type>(<optional scope>): <description>

Scope the title to the resource type that changed:
  resource/<name>    resource only
  datasource/<name>  data source only
  <name>             both resource and data source

Examples:
  feat(resource/rootly_api_key_rotation): add API key rotation resource [PRF-3498]
  fix(datasource/rootly_alert_route): return all matching alert routes [TER-227]
  fix(rootly_alert_route): allow clearing optional property and condition fields [TER-227]
  docs: clarify provider configuration
  chore(deps): update dependencies

Use feat for new features, fix for bug fixes, and ! before the colon for
breaking changes, e.g. feat!: remove a deprecated resource.

Release Please uses commit messages for release notes and version bumps.
When squash-merging, ensure the final commit title follows this format too.
Keep relevant Linear issue IDs as a title suffix, as shown above, and link
the ticket in the description. The title must start with the Conventional
Commit type, not the ticket ID (e.g. avoid [PRF-3498] feat(...): ...).
-->

## Description

<!--- What changes are being made in this PR? --->

## Motivation

<!--- Why are we making this change? What problem does it solve? Link any related issues. --->

## Testing

- [ ] Added new tests for this functionality
- [ ] Existing tests cover this change
- [ ] No tests needed (explain why below)

## Risk Assessment

- [ ] Added risk level label (low/medium/high risk)
