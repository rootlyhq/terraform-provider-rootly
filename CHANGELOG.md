# Changelog

## [0.1.50] -- 2022-11-04

- Add `CreateQuipPageTaskParams` resource.
- Add `workflow_task_create_quip_page` to `workflow_task_send_slack_message` resource.

## [0.1.49] -- 2022-11-02

- Add support for `number` with `form_field` resources.
- Add `skip_on_failure` to `workflow_task` resources.

## [0.1.48] -- 2022-11-02

- Add `form_field` resources.
- Deprecate `custom_fields` resources.

## [0.1.47] -- 2022-10-25

- Add `slug` attribute to `workflow_groups`.

## [0.1.46] -- 2022-10-25

- Make all data source attributes optional.

## [0.1.45] -- 2022-10-13

- Add `attribute_to_query_by` to `rootly_workflow_task_update_incident`.

## [0.1.44] -- 2022-10-13

- Add enum values in docs for array fields.

## [0.1.43] -- 2022-10-13

- Add `actionables` to `send_slack_message_task`.

## [0.1.42] -- 2022-10-13

- Ignore whitespace and keyorder when diffing `custom_fields_mapping` JSON fields.

## [0.1.41] -- 2022-10-12

- Whitespace and order agnostic diff for `rootly_workflow_task_send_slack_blocks` JSON.

## [0.1.40] -- 2022-10-10

- Convert `rootly_workflow_task_create_slack_channel` private field to `text` and supports `auto` options.

## [0.1.39] -- 2022-10-09

- Add `rootly_post_mortem_template`
- Fix diff for `rootly_workflow_task_send_slack_blocks`

## [0.1.38] -- 2022-10-08

- Improve documentation.

## [0.1.37] -- 2022-10-07

- Add `include_header` and `include_footer` to `workflow_task_send_email`.
- Add `content` attribute to `PostMortemTemplate`.

## [0.1.36] -- 2022-10-07

- Fix missing position for workflow tasks.
- Display API errors for all failing resource requests.

## [0.1.35] -- 2022-10-06

- Fix `workflow_task_send_slack_blocks` required parameters.

## [0.1.34] -- 2022-10-06

- Add enum values and required object properties to attribute description.

## [0.1.33] -- 2022-10-05

- Add `workflow_task_send_slack_blocks` resource.
- Add `incident_role_task` resource.

## [0.1.32] -- 2022-10-04

- Add `slack_channel_created` trigger type for workflows.
- Add `priority` attribute for workflows.
- Add `rootly_workflow_task_send_dashboard_report` resource.
- Add `due_date` attribute for action items.

## [0.1.31] -- 2022-09-22

- Add date range filters to data sources.
- Add `rootly_workflow_task_update_opsgenie_alert` resource.
- Add `rootly_workflow_task_create_slack_channel` resource.

## [0.1.30] -- 2022-09-20

- Add more genius workflow triggers.
- Add `create_linear_issue_comment`.
- Add `template` field for `create_confluence_page`.
- `conference_solution_key` is now optional for `create_google_meeting`.

## [0.1.29] -- 2022-09-15

- Fix `update_payload` type for Jira Task.
- Add `create_opsgenie_alert_task`.
- Add `default` field for custom fields.

## [0.1.28] -- 2022-09-09

- Fix `channels` type for `invite_to_slack_channel_pagerduty`, `invite_to_slack_channel_opsgenie` & `invite_to_slack_channel_victor_ops_task`.
- Add `trigger_workflow_task`.

## [0.1.27] -- 2022-08-31

- Add "slug" to team data source.

## [0.1.26] -- 2022-08-31

- Change data sources to return one resource. Example: `data.rootly_causes.name.causes[0].id` changed to `data.rootly_cause.name.id`
- Add playbook resources/data-source.
- Add status_page resources/data-source.
- Add status_page_template resources/data-source.

## [0.1.25] -- 2022-08-29

- Add `rootly_workflow_custom_field_selection` resource.

## [0.1.24] -- 2022-08-22

- Regenerate docs

## [0.1.23] -- 2022-08-22

- Add `from` for `resource_workflow_task_send_email`.
- Add `ListWorkflowRunsWithResponse`

## [0.1.22] -- 2022-08-22

- Add filtering by backstage_id/pagerduty_id/opsgenie_id for services & functionalities.

## [0.1.21] -- 2022-08-22

- Change enabled default from false to true for incident_roles and custom_fields.

## [0.1.20] -- 2022-08-18

- Fix missing fields.

## [0.1.19] -- 2022-08-17

Add `conference_solution_key` to `resource_workflow_task_create_google_calendar_event.go` &  `resource_workflow_task_update_google_calendar_event.go`
ADD `alert_condition_source_use_regexp`, `alert_condition_label_use_regexp`, `alert_condition_payload_use_regexp` flags
ADD `pulse_condition_source_use_regexp`, `pulse_condition_label_use_regexp`, `pulse_condition_payload_use_regexp` flags

# Changelog

## [0.1.18] -- 2022-08-17

Add CONTAINS_ALL to genius workflows

# Changelog

## [0.1.17] -- 2022-08-12

### Changed

Add categories to docs.

## [0.1.15] -- 2022-08-12

### Added

- Added data sources:
	- rootly_causes
	- rootly_custom_fields
	- rootly_custom_field_options
	- rootly_environments
	- rootly_functionalities
	- rootly_incident_roles
	- rootly_incident_types
	- rootly_severities
	- rootly_services
	- rootly_teams
- Added resources:
	- rootly_dashboard
	- rootly_dashboard_panel

## [0.1.14] -- 2022-08-08

### Added

- Add workflow_group_id field to workflow resources.

### Changed

- Fixed type of incident_action_item_statuses, incident_action_item_kinds, incident_action_item_priorities for workflow_action_item resource.

## [0.1.13] -- 2022-08-05

### Change

- Changed workflow task workflow_id field to ForceNew. Workflow tasks cannot
	change their workflow. They must be deleted and recreated with the correct
	workflow.

## [0.1.12] -- 2022-08-03

### Added

- Added workflow_group resource.

### Changed

- changed workflow_task_snapshot_datadog_dashboard.dashboards  optional

## [0.1.11] -- 2022-08-03

### Changed

- Fixed position field persistence for workflows.

### Added

## [0.1.10] -- 2022-07-26

### Added

- Added position field to workflows and tasks

## [0.1.9] -- 2022-07-23

### Changed

- Updated GPG github action

## [0.1.8] -- 2022-07-23

### Added

- Added additional schema fields.

### Fixed

- Fixed "enabled" fields with "false".

## [0.1.7] -- 2022-07-12

### Added

- Added workflow examples.

## [0.1.6] -- 2022-07-11

### Changed

- Retry failed HTTP requests with exponential backoff.

## [0.1.5] -- 2022-07-07

### Added

- Custom field resource
- Workflow and workflow task resources
- Environments resource

## [0.1.4] -- 2022-05-11

### Changed

- Updated docs and examples

## [0.1.3] -- 2022-05-10

### Fixed

- Fix incorrect dependency URLs.
