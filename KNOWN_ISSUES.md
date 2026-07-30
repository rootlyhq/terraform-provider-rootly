# Known Test Issues

Skipped acceptance tests and their root causes. Grouped by issue type.

## API returns wrong type for field

| Test | Field | Expected | API returns | Fix |
|---|---|---|---|---|
| `CreateMotionTask` | `labels` | `array` | `string` | Change swagger type to `string` or fix API to return `[]` |
| `UpdateMotionTask` | `labels` | `array` | `string` | Same as above |

## API returns empty objects `{}` for unset map fields

Terraform sees plan drift because unset `TypeMap` fields come back as `{}` instead of `null`.

| Test | Drifting fields | Fix |
|---|---|---|
| `CreateJiraIssue` | `integration`, `priority`, `status` | API should return `null` for unset object fields |
| `CreateJiraSubtask` | `integration`, `priority`, `status` | Same |
| `PagePagerdutyOnCallResponders` | `priority` | Same |

## API requires field not marked required in swagger

| Test | Required field | Fix |
|---|---|---|
| `UpdateCodaPage` | `doc` | Add to `required` in swagger |
| `UpdateDatadogNotebook` | `kind` (empty string rejected) | Add to `required` or add `default` |
| `UpdateAttachedAlerts` | `name` (task name not auto-assigned) | Fix API to auto-assign task name |
| `UpdateGitlabIssue` | `issue_type` (DB NOT NULL) | Add to `required` in swagger (same fix as `CreateGitlabIssue` in PR #20916) |
| `SendDashboardReport` | `from` | Fixed — swagger has `default`, codegen now reads it |
| `SendEmail` | `from` | Fixed — same as above |

## API requires real integration / escalation target

| Test | Issue | Fix |
|---|---|---|
| `InviteToMicrosoftTeamsChannelRootly` | Needs `oneOf` escalation target | Add `oneOf` with escalation targets in swagger |
| `InviteToSlackChannelRootly` | Needs `oneOf` escalation target | Same |
| `PageRootlyOnCallResponders` | Needs escalation target + `alert_urgency` | Add `oneOf` + add `alert_urgency` to `required` |
| `InviteToMicrosoftTeamsChannel` | Team object validation fails with dummy data | Investigate |
| `CreateMicrosoftTeamsChat` | Needs real members with `email`, ≥2 members, `topic` | Change swagger items to `{id, name}` + add `topic` to required |
| `CreatePagertreeAlert` | "Taskable must exist" — needs real integration | Configure PagerTree on test org |
| `UpdatePagertreeAlert` | Same | Same |

## API returns 500 (server errors)

| Test | Notes |
|---|---|
| `SendMicrosoftTeamsBlocks` | 500 on workflow task creation — all 9 retries fail |
| `UpdateConfluencePage` | 500 timeout after 9 retries |

## CI environment limitations

| Test | Issue | Fix |
|---|---|---|
| `ScheduleRotationMembersScheduleType` | CI API key is user-less, schedule creation requires `owner_user` | Use a user-scoped API key or skip |
| `WorkflowActionItemFormFieldCondition` | CI API key returns 403 | Enable feature on test org |

## Swagger/codegen improvements made

These were fixed during this PR:

- **`genResourceTestFile` wired up** — was dead code, now auto-generates tests for all workflow task types
- **`integer` type support** — `genTaskSchemaPropertyType` and `genTestParams` now handle `integer` (was falling through to `TypeString`)
- **`number` type** — maps to `TypeString` (API returns floats as strings like `"0.0"`)
- **`anyOf`/`oneOf` support** — `genTestParams` picks up first conditional option's required fields
- **String defaults from swagger** — optional string fields with swagger `default` now set `Default` in Go schema
- **Optional integers with `minimum`/`default` > 0** — included in test params with valid values
- **Conditional imports** — `reflect`/`encoding/json` only included when JSON diff suppressor is needed
- **Escalation path rule ordering** — `matchRulesOrder` reorders API response to match config
- **Escalation path time_blocks ordering** — `matchTimeBlocksOrder` same pattern
- **`AlertsSource` test** — API now auto-converts `alert_template_attributes` (flipper globally enabled)
