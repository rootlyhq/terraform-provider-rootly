# Changelog

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
