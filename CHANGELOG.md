# Changelog

## [1.2.16] -- 2024-01-31

- Fix some documentation.
- Add `title` attribute to `resource_workflow_task_create_notion_page` & `resource_workflow_task_update_notion_page` resources.
- Add `parent_message_thread_task` attribute to `workflow_task_get_alerts` & `workflow_task_get_pulses` resources.
- Update dependencies.

## [1.2.15] -- 2024-01-11

- Add `environments_impacted_by_incident` & `services_impacted_by_incident` attribute to `workflow_task_get_alerts` resource.
- Add `environments_impacted_by_incident` & `services_impacted_by_incident` attribute to `workflow_task_get_pulses` resource.
- Add `services_impacted_by_incident` attribute to `workflow_task_get_github_commits` resource.
- Add `services_impacted_by_incident` attribute to `workflow_task_get_gitlab_commits` resource.

## [1.2.14] -- 2024-01-10

- Fix required `resource_id` & `resource_type` attribute for `incident_permission_set_resource` resource.

## [1.2.13] -- 2024-01-09

- Remove `status` attribute from `service` and `functionality` resource.

## [1.2.12] -- 2024-01-07

- Add `icon` & `description` attribute to `workflow_group` resource.

## [1.2.11] -- 2023-12-29

- Add support for `authorizations`, `incident_permission_sets`, `incident_permission_set_booleans`, `incident_permission_set_resources`, and `roles`.
- Add `show_on_incident_details` to `form_field` resource.
- Add `rich_text` type to `form_field` resource.
- Suppress diff for post mortem template content and secret.

## [1.2.10] -- 2023-12-20

- Ignore order for int and workflow attributes.

## [1.2.9] -- 2023-12-16

- Add `causes` kind to `form_field` resource.
- Add `incident_condition_cause` column to `workflow_incident`.
- Deprecated `incident_post_mortem_condition_cause` column of `workflow_post_mortem`
- Remove `causes_updated` trigger of `workflow_post_mortem`.

## [1.2.8] -- 2023-12-03

- Add `user_ids` attribute to `team` resource.

## [1.2.7] -- 2023-12-02

- Upgrade dependencies
- Add `incident_updated` trigger to `workflow_action_item` resource.

## [1.2.6] -- 2023-11-24

- Add `create_zendesk_jira_link_task` resource.

## [1.2.5] -- 2023-11-22

- Add `kind` attribute to `status_page_template` resource.
- Fix documentation.

## [1.2.4] -- 2023-11-14

- Add `position` attribute to `playbook_task` resource.
- Add `slack_incident_cancellation_form` & `web_incident_cancellation_form` to `form_field` resource.

## [1.2.3] -- 2023-11-08

- Ignore order of list attributes values when diffing for changes for all resources in addition to workflows
- Update priority options for Opsgenie tasks

## [1.2.2] -- 2023-11-07

- Ignore order of list attributes values when diffing for changes.

## [1.2.1] -- 2023-11-07

- Add `due_date` attribute to `asana` workflow resources.
- Update dependencies.

## [1.2.0] -- 2023-11-02

- Add `public_title`, `public_description` attributes status page resource.
- Add `message`, `description` attributes to `PageOpsgenieOnCallResponders` resource.
- Update dependencies.

## [1.1.9] -- 2023-10-28

- Add `due_date` attribute to `clickup`, `shortcut` & `trello` workflow resources.
- Fix documentation.

## [1.1.8] -- 2023-10-25

- Add 'rootly_ip_ranges' data source.

## [1.1.7] -- 2023-10-20

- Add `functionality_ids` conditions to workflows.
- Add `clickup_task_id` to `attribute_to_query_by` for `workflow_task_update_incident` and `workflow update_action_item` tasks.
- Add `retrospective_steps` resource.

## [1.1.6] -- 2023-09-20

- Allow to pass `custom_fields.<slug>.updated` trigger to workflows.
- Improve documentation.

## [1.1.5] -- 2023-09-20

- Fix panic error by replacing testify with our own JSON comparison logic.

## [1.1.4] -- 2023-09-20

- Upgraded go dependencies.

## [1.1.3] -- 2023-09-06

- Added `pagerduty_id`, `opsgenie_id`, `victor_ops_id` & `pagertree_id` attributes to `rootly_team` resource.
- Added `close` status to `task_update_opsgenie_incident` resource.

## [1.1.2] -- 2023-08-24

- Added `custom_fields_mapping` attribute to `task_update_action_item` resource.

## [1.1.1] -- 2023-08-23

- Added `priority` attribute to `task_update_action_item` resource.
- Added `update_parent_message` attribute to `task_send_slack_message` resource.

## [1.1.0] -- 2023-08-17

- Added `service_ids` to `playbook` resource.
- Added `integration` attribute to `task_create_confluence_page`, `task_create_jira_issue` ,`task_create_jira_subtask` & ``.
- Added `jira_issue_id`, `asana_task_id`, `shortcut_task_id`, `shortcut_story_id`, `linear_issue_id`, `zendesk_ticket_id`, `trello_card_id`, `airtable_record_id`, `github_issue_id`, `freshservice_ticket_id` & `freshservice_task_id` attributes to `task_update_incident`.
- Added `assign_user_email` attribute to `task_create_linear_subtask_issue` & `task_update_linear_issue`.
- Added `create_incident_postmortem_task` workflow task.

## [1.0.9] -- 2023-08-04

- Fixed `invite_to_slack_channel_victor_ops_task_params` required argument.
- Added `due_date` attribute to `publish_incident_task_params`.

## [1.0.8] -- 2023-08-03

- Fixed `auto_assign_role_victor_ops_task_params`.
- Fixed `invite_to_slack_channel_victor_ops_task_params`.
- Added `integration_payload` attribute to `publish_incident_task_params`.

## [1.0.7] -- 2023-07-19

- Fixed `name` computed value.
- Added `rootly_incident` + `rootly_incident_post_mortem` data sources.
- Added `send_whatsapp_message_task` resource.
- Fixed some required attributes for workflow task resources.
- Improved test files.

## [1.0.6] -- 2023-07-18

- Added `name` attributes to all workflow tasks.

## [1.0.5] -- 2023-07-14

- Added `allow_multi_user_assignment` to `incident_role` resource.
- Fixed `triggers` list for `workflow_post_mortem` resource.
- Added `template_id` to `workflow_task_update_google_docs_page` resource.

## [1.0.4] -- 2023-07-06

- Added `assigned_to_user` and deprecated `assigned_to_user_id` for `workflow_task_add_action_item` resource.
- Added `assigned_to_user` and deprecated `assigned_to_user_id` for `workflow_task_update_action_item` resource.
- Added `assigned_to_user` and deprecated `assigned_to_user_id` for `workflow_task_add_role` resource.
- Added `user` data source.

## [1.0.3] -- 2023-06-29

- Added `in_triage` status for incidents.
- Added `command_feedback_enabled` for workflows.

## [1.0.2] -- 2023-06-18

- Add `workflow_task_update_notion_page` resource.

## [1.0.1] -- 2023-06-02

- Add `notes` attribute to `workflow_task_update_asana_task` resource.

## [1.0.0] -- 2023-06-02

- Fix Semantic Versioning.

## [0.2.00] -- 2023-06-02

- Add `ticket_payload` attribute to `workflow_task_create_zendesk_ticket` resource.
- Add `bcc` & `cc` attributes to `/workflow_task_send_email` resource.

## [0.1.99] -- 2023-05-24

- `skip_on_failure` and `enabled` attributes for `workflow_task` resources are now working as expected.
- Fix some documentation.

## [0.1.98] -- 2023-05-23

- Add `service_ids`, `functionality_ids` attributes to `status_page` resource.
- Add `mark_post_mortem_as_published` attribute to `workflow_task_create_quip_page` resource.

## [0.1.97] -- 2023-05-18

- Add `notify_subscribers`, `should_tweet` attributes to `workflow_task_publish_incident` resource.

## [0.1.96] -- 2023-05-16

- Add `authentication_enabled`, `authentication_password`, `website_url`, `website_privacy_url`, `website_support_url`, `ga_tracking_id`, `time_zone` attributes to status page resource.

## [0.1.95] -- 2023-05-12

- Improve docs.

## [0.1.94] -- 2023-05-04

- Add `incident_inactivity_duration` to workflow resources.

## [0.1.93] -- 2023-05-04

- Add data sources documentation.
- Add `rootly_status_page` data source.

## [0.1.92] -- 2023-04-20

- Add `CONTAINS_NONE` to workflows conditions.

## [0.1.91] -- 2023-04-20

- Fix `incident_role_ids` typo.

## [0.1.90] -- 2023-04-19

- Fix `rootly_workflow_task_trigger_workflow` required attributes.

## [0.1.89] -- 2023-04-11

- Add `skip_on_failure` and `enabled` to `workflow_task` resources.
- Add `color` to `workflow_task_send_slack_message` resource.

## [0.1.88] -- 2023-04-06

- Add `enabled` attribute to worklow tasks.

## [0.1.87] -- 2023-04-06

- Add `notes` attribute to `rootly_workflow_task_create_asana_subtask` resource.

## [0.1.86] -- 2023-04-04

- Add `rootly_workflow_task_update_google_docs_page`, `rootly_workflow_task_create_google_docs_permissions`, `rootly_workflow_task_remove_google_docs_permissions` resources.
- Add `subscribers_updated`, `subscribers_added`, `subscribers_removed`, `user_left_slack_channel` workflow triggers.

## [0.1.85] -- 2023-03-30

- Add `rootly_workflow_task_update_action_item` resource.

## [0.1.84] -- 2023-03-28

- Add `mark_post_mortem_as_published` attribute to `create_dropbox_paper_page_task`, `create_google_docs_page`, `create_confluence_page`, `create_datadog_notebook`, `create_notion_page` resources.

## [0.1.83] -- 2023-03-27

- Add `status` to `services` and `functionalities` resources.
- Add `custom_fields_mapping` attributes to `zendesk` workflow resources.
- Fix defaulted values for `resources`.

## [0.1.82] -- 2023-03-17

- Add `normal_sub` and `test_sub` incident kinds.
- Add `hashicorp_vault_mount`, `hashicorp_vault_mount`, `hashicorp_vault_mount` attributes to secret resource.
- Add `incident_role_id` to `AddActionItemTaskParams`.

## [0.1.81] -- 2023-03-13

- Add `service` to `workflow_task_invite_to_slack_channel` task.

## [0.1.80] -- 2023-03-03

- Add `dependency_direction` & `dependent_task_ids` attributes to asana tasks resources.

## [0.1.79] -- 2023-03-01

- Add `position` attributes to multiple resources.

## [0.1.78] -- 2023-03-01

- Add `optional` and `enabled` attributes to `incident_roles` resource.

## [0.1.77] -- 2023-02-28

- Rename `Postmortem` to `Retrospective` in our documentation.
- Add `check_workflow_conditions` to `workflow_task_trigger_workflow` resource.

## [0.1.76] -- 2023-02-28

- Add `incident_roles_ids` & `incident_condition_incident_roles` to `resource_workflow_*` resources.

## [0.1.75] -- 2023-02-27

- Fix resources examples.
- Add `pin_to_channel` to `workflow_task_send_slack_message` & `workflow_task_send_slack_block` resources.

## [0.1.74] -- 2023-02-24

- Add more resources examples.

## [0.1.73] -- 2023-02-23

- Add resources examples.

## [0.1.72] -- 2023-02-21

- Add `format` field to `post_mortem_template` resource.

## [0.1.71] -- 2023-02-21

- Add `failure_message`, `success_message` fields to `status_pages` resource.
- Add `pause_reminder`, `snooze_reminder`, `restart_reminder` to `actionables` on `workflow_task_send_slack_message` resource.

## [0.1.70] -- 2023-02-16

- Add `notify_emails`, `slack_channels`, `slack_aliases` fields to `severities`, `environments` & `incident_types`.

## [0.1.69] -- 2023-02-10

- Add `project` and `labels` fields to `workflow_task_create_linear_issue` & `workflow_task_update_linear_issue`.

## [0.1.68] -- 2023-01-28

- Add `show_action_items_as_table` and `show_timeline_as_table` fields to `rootly_workflow_task_create_notion_page`.
- Drop support for `rootly_postmortem_template`. Rename to `rootly_post_mortem_template` instead.

## [0.1.67] -- 2023-01-28

- Add `show_action_items_as_table` and `show_timeline_as_table` fields to `rootly_workflow_task_create_notion_page`.
- Drop support for `rootly_postmortem_template`. Rename to `rootly_post_mortem_template` instead.

## [0.1.66] -- 2023-01-24

- Add missing attributes to `tasks` related to slack.

## [0.1.65] -- 2023-01-24

- Add `owners_user_ids` attribute to `service` and `functionality` resources.

## [0.1.64] -- 2023-01-24

- Remove duplicate fields.

## [0.1.63] -- 2023-01-24

- Add additional fields to `rootly_playbook`.

## [0.1.62] -- 2023-01-23

- Rename `rootly_postmortem_template` to `rootly_post_mortem_template` to match Rootly API.
- Add `rootly_playbook_task` resource.

## [0.1.60] -- 2023-01-19

- Fix `workflow_task_create_airtable_table_record` attributes.

## [0.1.59] -- 2023-01-11

- Fix inoperative postmortem template content attribute.

## [0.1.57] -- 2023-01-05

- Add `secrets` resources.
- Add `pagertree` genius task resources.

## [0.1.56] -- 2022-12-13

- Improve docs.
- Add `show_uptime` & `show_uptime_last_days` to services and functionalities.
- Add `slack_channel_converted` trigger to `resourceWorkflowIncident`.

## [0.1.55] -- 2022-12-01
- Add `rootly_workflow_simple` resource.
- Add `auto_refresh` to `dashboard` resources.

## [0.1.54] -- 2022-11-29
- Default nil values for optional enums
- Added maintenance form fields constants `web_scheduled_incident_form`, `web_update_scheduled_incident_form`, `slack_scheduled_incident_form`, and `slack_update_scheduled_incident_form`

## [0.1.53] -- 2022-11-17

- Update `IncidentTriggerParamsTriggers`.

## [0.1.52] -- 2022-11-14

- Add `rootly_workflow` and `rootly_workflow_task` data sources.

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
