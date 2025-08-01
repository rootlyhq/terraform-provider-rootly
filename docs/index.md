# Rootly Provider

[Rootly](https://rootly.com) is the all-in-one AI-native platform for on-call and incident management, including status pages—built for fast-moving engineering teams to detect, manage, learn from, and resolve incidents faster.

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `api_host` (String) The Rootly API host. Defaults to https://api.rootly.com. Can also be sourced from the `ROOTLY_API_URL` environment variable.
- `api_token` (String, Sensitive) The Rootly API Token. Generate it from your account at https://rootly.com/account. It must be provided but can also be sourced from the `ROOTLY_API_TOKEN` environment variable.

## Example Usage

### Provider

```terraform
terraform {
  required_providers {
    rootly = {
      source = "rootlyhq/rootly"
    }
  }
}

provider "rootly" {
  # We recommend using the `ROOTLY_API_TOKEN` env var to set the API Token
  # when interacting with Rootly's API.
  # api_token = var.rootly_api_key
}
```

Data sources

```terraform
# uses output of severity data source as input for workflow
data "rootly_severity" "critical" {
  slug = "sev0"
}

resource "rootly_workflow_incident" "ping_oncall" {
  name        = "Ping on-call when critical incident"
  description = "ping on-call when critical incident happens"
  trigger_params {
    triggers                  = ["incident_created"]
    incident_condition_kind   = "IS"
    incident_kinds            = ["normal"]
    incident_condition_status = "IS"
    incident_statuses         = ["started"]
    severity_ids              = [data.rootly_severity.critical.id]
  }
  enabled = true
}

resource "rootly_workflow_task_send_sms" "sms_oncall" {
  workflow_id = rootly_workflow_incident.ping_oncall.id
  name        = "On-call team"

  task_params {
    phone_numbers = ["+11231231234"]
    content       = "Critical incident started"
  }
}
```

### Workflows

```terraform
# Jira workflow
resource "rootly_workflow_incident" "jira" {
  name        = "Create a Jira Issue"
  description = "Open Jira ticket whenever incident starts"
  trigger_params {
    triggers                  = ["incident_created"]
    incident_condition_kind   = "IS"
    incident_kinds            = ["normal"]
    incident_condition_status = "IS"
    incident_statuses         = ["started"]
  }
  enabled = true
}

resource "rootly_workflow_task_create_jira_issue" "jira" {
  workflow_id = rootly_workflow_incident.jira.id
  task_params {
    title       = "{{ incident.title }}"
    description = "{{ incident.summary }}"
    project_key = "ROOT"
    issue_type = {
      id   = "10001"
      name = "Task"
    }
    status = {
      id   = "10000"
      name = "To Do"
    }
    labels = "{{ incident.environment_slugs | concat: incident.service_slugs | concat: incident.functionality_slugs | concat: incident.group_slugs | join: \",\" }}"
  }
}
```

### Dashboards

```terraform
resource "rootly_dashboard" "overview" {
  name = "my_dashboard"
}

resource "rootly_dashboard_panel" "incidents_by_severity" {
  dashboard_id = rootly_dashboard.foo.id
  name         = "test"
  params {
    display = "line_chart"
    datasets {
      collection = "incidents"
      filter {
        operation = "and"
        rules {
          operation = "and"
          condition = "="
          key       = "status"
          value     = "started"
        }
      }
      group_by = "severity"
      aggregate {
        cumulative = false
        key        = "results"
        operation  = "count"
      }
    }
  }
}
```

### On-Call

```terraform
data "rootly_user" "john" {
  email = "demo@rootly.com"
}

data "rootly_user" "jane" {
  email = "demo1@rootly.com"
}

data "rootly_alert_urgency" "low" {
  name = "Low"
}

resource "rootly_team" "sre" {
  name     = "SREs On-Call"
  user_ids = [data.rootly_user.john.id, data.rootly_user.jane.id]
}

resource "rootly_schedule" "primary" {
  name              = "Primary On-Call Schedule"
  owner_user_id     = data.rootly_user.john.id
  all_time_coverage = false
}

resource "rootly_schedule_rotation" "weekdays" {
  schedule_id     = rootly_schedule.primary.id
  name            = "weekdays"
  active_all_week = false
  active_days = [
    "M",
    "T",
    "W",
    "R",
    "F",
  ]
  active_time_type = "custom"
  position         = 1
  schedule_rotationable_attributes = {
    handoff_time = "10:00"
  }
  schedule_rotationable_type = "ScheduleDailyRotation"
  time_zone                  = "America/Toronto"
}

# Define active days for the weekday rotation
# Monday
resource "rootly_schedule_rotation_active_day" "m1-weekday" {
  schedule_rotation_id = rootly_schedule_rotation.weekdays.id
  day_name             = "M"
  active_time_attributes {
    start_time = "10:00"
    end_time   = "18:00"
  }
}
# Tuesday
resource "rootly_schedule_rotation_active_day" "t1-weekday" {
  schedule_rotation_id = rootly_schedule_rotation.weekdays.id
  day_name             = "T"
  active_time_attributes {
    start_time = "10:00"
    end_time   = "18:00"
  }
}
# Wednesday
resource "rootly_schedule_rotation_active_day" "w1-weekday" {
  schedule_rotation_id = rootly_schedule_rotation.weekdays.id
  day_name             = "W"
  active_time_attributes {
    start_time = "10:00"
    end_time   = "18:00"
  }
}
# Thursday
resource "rootly_schedule_rotation_active_day" "th1-weekday" {
  schedule_rotation_id = rootly_schedule_rotation.weekdays.id
  day_name             = "R"
  active_time_attributes {
    start_time = "10:00"
    end_time   = "18:00"
  }
}
# Friday
resource "rootly_schedule_rotation_active_day" "f1-weekday" {
  schedule_rotation_id = rootly_schedule_rotation.weekdays.id
  day_name             = "F"
  active_time_attributes {
    start_time = "10:00"
    end_time   = "18:00"
  }
}

resource "rootly_schedule_rotation_user" "john" {
  schedule_rotation_id = rootly_schedule_rotation.weekdays.id
  position             = 1
  user_id              = data.rootly_user.john.id
}

resource "rootly_schedule_rotation_user" "jane" {
  schedule_rotation_id = rootly_schedule_rotation.weekdays.id
  position             = 2
  user_id              = data.rootly_user.jane.id
}

resource "rootly_escalation_policy" "primary" {
  name      = "Primary"
  group_ids = [rootly_team.sre.id]
}

resource "rootly_escalation_path" "default" {
  name                 = "Default"
  default              = true
  escalation_policy_id = rootly_escalation_policy.primary.id
}

resource "rootly_escalation_path" "ignore" {
  name                 = "Ignore"
  default              = false
  escalation_policy_id = rootly_escalation_policy.primary.id
  rules {
    rule_type   = "alert_urgency"
    urgency_ids = [data.rootly_alert_urgency.low.id]
  }
}

resource "rootly_escalation_level" "first" {
  escalation_policy_path_id = rootly_escalation_path.default.id
  escalation_policy_id      = rootly_escalation_policy.primary.id
  position                  = 1
  notification_target_params {
    team_members = "all"
    type         = "slack_channel"
    id           = "C06D4QHLAUE"
  }
  notification_target_params {
    type         = "schedule"
    id           = rootly_schedule.primary.id
    team_members = "all"
  }
}

# cycle-based round-robin everyone on the schedule
resource "rootly_escalation_level" "second" {
  escalation_policy_path_id                       = rootly_escalation_path.default.id
  escalation_policy_id                            = rootly_escalation_policy.primary.id
  position                                        = 2
  delay                                           = 5
  paging_strategy_configuration_strategy          = "cycle"
  paging_strategy_configuration_schedule_strategy = "everyone"
  notification_target_params {
    type         = "schedule"
    id           = rootly_schedule.primary.id
    team_members = "all"
  }
}
```
