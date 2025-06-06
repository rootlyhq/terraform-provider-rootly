---
page_title: "Resource rootly_escalation_level - terraform-provider-rootly"
subcategory:
description: |-
    
---

# Resource (rootly_escalation_level)



## Example Usage

```shell
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

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import rootly_escalation_level using the `id`. For example:

```terraform
import {
  to = rootly_escalation_level.my-resource
  id = "00000000-0000-0000-0000-000000000000"
}
```

Using `terraform import`, import rootly_escalation_level using the `id`. For example:

```console
% terraform import rootly_escalation_level.my-resource 00000000-0000-0000-0000-000000000000
```

Locate the resource id in the web app, or retrieve it by listing resources through the API if it's not visible in the web app.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `notification_target_params` (Block List, Min: 1) Escalation level's notification targets (see [below for nested schema](#nestedblock--notification_target_params))
- `position` (Number) Position of the escalation policy level

### Optional

- `delay` (Number) Delay before notification targets will be alerted.
- `escalation_policy_id` (String) The ID of the escalation policy
- `escalation_policy_path_id` (String) The ID of the dynamic escalation policy path the level will belong to. If nothing is specified it will add the level to your default path.
- `paging_strategy_configuration_schedule_strategy` (String) Value must be one of `on_call_only`, `everyone`.
- `paging_strategy_configuration_strategy` (String) Value must be one of `default`, `random`, `cycle`, `alert`.

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--notification_target_params"></a>
### Nested Schema for `notification_target_params`

Optional:

- `id` (String) The ID of notification target
- `team_members` (String) For targets with type=team, controls whether to notify admins, all team members, or escalate to team EP.. Value must be one of `all`, `admins`, `escalate`.
- `type` (String) The type of the notification target. Value must be one of `team`, `user`, `schedule`, `slack_channel`, `service`.
