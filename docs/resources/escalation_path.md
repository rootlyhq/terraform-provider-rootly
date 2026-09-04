---
page_title: "Resource rootly_escalation_path - terraform-provider-rootly"
subcategory:
description: |-
    Manages an escalation path.
---

# Resource (rootly_escalation_path)

Manages an escalation path.

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

# Deferral path - defer alerts outside business hours, then re-evaluate
resource "rootly_escalation_path" "defer_off_hours" {
  name                    = "Defer Off Hours"
  default                 = false
  escalation_policy_id    = rootly_escalation_policy.primary.id
  path_type               = "deferral"
  after_deferral_behavior = "re_evaluate"

  rules {
    rule_type = "deferral_window"
    time_zone = "America/New_York"
    time_blocks {
      monday     = true
      tuesday    = true
      wednesday  = true
      thursday   = true
      friday     = true
      start_time = "18:00"
      end_time   = "09:00"
    }
    time_blocks {
      saturday = true
      sunday   = true
      all_day  = true
    }
  }
}

# Deferral path that executes another path after deferral
resource "rootly_escalation_path" "defer_then_escalate" {
  name                    = "Defer Then Escalate"
  default                 = false
  escalation_policy_id    = rootly_escalation_policy.primary.id
  path_type               = "deferral"
  after_deferral_behavior = "execute_path"
  after_deferral_path_id  = rootly_escalation_path.default.id

  rules {
    rule_type = "deferral_window"
    time_zone = "America/New_York"
    time_blocks {
      saturday = true
      sunday   = true
      all_day  = true
    }
  }
}

# Service-based routing path
resource "rootly_escalation_path" "by_service" {
  name                 = "Route by Service"
  default              = false
  escalation_policy_id = rootly_escalation_policy.primary.id
  match_mode           = "match-any-rule"

  rules {
    rule_type   = "service"
    service_ids = ["your-service-id"]
  }
}

# Alert source-based routing path - only route alerts from specific sources
resource "rootly_escalation_path" "by_source" {
  name                 = "Route by Alert Source"
  default              = false
  escalation_policy_id = rootly_escalation_policy.primary.id
  match_mode           = "match-any-rule"

  rules {
    rule_type = "source"
    operator  = "is_one_of"
    values    = ["manual", "datadog"]
  }
}

# Related incidents routing path - only route alerts that have related incidents
resource "rootly_escalation_path" "with_related_incidents" {
  name                 = "Route When Related Incidents Exist"
  default              = false
  escalation_policy_id = rootly_escalation_policy.primary.id
  match_mode           = "match-any-rule"

  rules {
    rule_type = "related_incidents"
    operator  = "is_set"
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
# Audible/quiet conditions on a single path - rules are evaluated in order and the
# first match wins, otherwise notification_type_fallback applies
resource "rootly_escalation_path" "conditional_notifications" {
  name                       = "Conditional Notifications"
  default                    = false
  escalation_policy_id       = rootly_escalation_policy.primary.id
  notification_type_fallback = "quiet"

  # Page audibly for low urgency alerts during business hours
  notification_type_rules {
    notification_type = "audible"
    match_mode        = "match-all-rules"

    conditions {
      rule_type   = "alert_urgency"
      urgency_ids = [data.rootly_alert_urgency.low.id]
    }

    conditions {
      rule_type = "deferral_window"
      time_zone = "America/New_York"
      time_blocks {
        monday     = true
        tuesday    = true
        wednesday  = true
        thursday   = true
        friday     = true
        start_time = "09:00"
        end_time   = "17:00"
      }
    }
  }

  # Stay quiet for informational alerts regardless of when they arrive
  notification_type_rules {
    notification_type = "quiet"
    match_mode        = "match-all-rules"

    conditions {
      rule_type = "json_path"
      json_path = "$.severity"
      operator  = "is"
      value     = "info"
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `escalation_policy_id` (String) <i style="color:red;font-weight: bold">(ForceNew)</i>
- `name` (String) The name of the escalation path.

### Optional

- `after_deferral_behavior` (String) What happens after a deferral path finishes. Value must be one of `re_evaluate`, `execute_path`.
- `after_deferral_path_id` (String) The escalation path to execute after this deferral path when after_deferral_behavior is execute_path.
- `default` (Boolean) Whether this escalation path is the default path.
- `initial_delay` (Number) Initial delay for escalation path in minutes. Maximum 1 week (10080).
- `match_mode` (String) How path rules are matched. Value must be one of `match-all-rules`, `match-any-rule`.
- `notification_type` (String) Notification rule type.
- `notification_type_fallback` (String) Paged when no notification type rule matches. Considered only when notification_type_rules are present — the path's notification_type is aligned to it; without rules it is aligned to notification_type instead. Only available when notification type conditions are enabled for the team. Value must be one of `audible`, `quiet`.
- `notification_type_rules` (Block List) Rules deciding whether an alert pages audible or quiet, evaluated in order — the first matching rule's notification_type wins, otherwise notification_type_fallback applies. When present, the path's notification_type is aligned to notification_type_fallback. Only available when notification type conditions are enabled for the team. (see [below for nested schema](#nestedblock--notification_type_rules))
- `path_type` (String) <i style="color:red;font-weight: bold">(ForceNew)</i> The type of escalation path. Value must be one of `escalation`, `deferral`.
- `position` (Number) The position of this path in the paths for this EP.
- `repeat` (Boolean) Whether this path should be repeated until someone acknowledges the alert.
- `repeat_count` (Number) The number of times this path will be executed until someone acknowledges the alert.
- `retrigger_timeout_minutes` (Number) Re-trigger acknowledged alerts on this path after N minutes; null inherits the urgency/workspace default, negative = never.
- `rules` (Block List) Escalation path rules. (see [below for nested schema](#nestedblock--rules))
- `time_restriction_time_zone` (String) Time zone used for time restrictions. Value must be one of `International Date Line West`, `Etc/GMT+12`, `American Samoa`, `Pacific/Pago_Pago`, `Midway Island`, `Pacific/Midway`, `Hawaii`, `Pacific/Honolulu`, `Alaska`, `America/Juneau`, `Pacific Time (US & Canada)`, `America/Los_Angeles`, `Tijuana`, `America/Tijuana`, `Arizona`, `America/Phoenix`, `Mazatlan`, `America/Mazatlan`, `Mountain Time (US & Canada)`, `America/Denver`, `Central America`, `America/Guatemala`, `Central Time (US & Canada)`, `America/Chicago`, `Chihuahua`, `America/Chihuahua`, `Guadalajara`, `America/Mexico_City`, `Mexico City`, `Monterrey`, `America/Monterrey`, `Saskatchewan`, `America/Regina`, `Bogota`, `America/Bogota`, `Eastern Time (US & Canada)`, `America/New_York`, `Indiana (East)`, `America/Indiana/Indianapolis`, `Lima`, `America/Lima`, `Quito`, `Atlantic Time (Canada)`, `America/Halifax`, `Caracas`, `America/Caracas`, `Georgetown`, `America/Guyana`, `La Paz`, `America/La_Paz`, `Puerto Rico`, `America/Puerto_Rico`, `Santiago`, `America/Santiago`, `Newfoundland`, `America/St_Johns`, `Asuncion`, `America/Asuncion`, `Brasilia`, `America/Sao_Paulo`, `Buenos Aires`, `America/Argentina/Buenos_Aires`, `Montevideo`, `America/Montevideo`, `Greenland`, `America/Nuuk`, `Mid-Atlantic`, `Atlantic/South_Georgia`, `Azores`, `Atlantic/Azores`, `Cape Verde Is.`, `Atlantic/Cape_Verde`, `Edinburgh`, `Europe/London`, `Lisbon`, `Europe/Lisbon`, `London`, `Monrovia`, `Africa/Monrovia`, `UTC`, `Etc/UTC`, `Amsterdam`, `Europe/Amsterdam`, `Belgrade`, `Europe/Belgrade`, `Berlin`, `Europe/Berlin`, `Bern`, `Europe/Zurich`, `Bratislava`, `Europe/Bratislava`, `Brussels`, `Europe/Brussels`, `Budapest`, `Europe/Budapest`, `Casablanca`, `Africa/Casablanca`, `Copenhagen`, `Europe/Copenhagen`, `Dublin`, `Europe/Dublin`, `Ljubljana`, `Europe/Ljubljana`, `Madrid`, `Europe/Madrid`, `Paris`, `Europe/Paris`, `Prague`, `Europe/Prague`, `Rome`, `Europe/Rome`, `Sarajevo`, `Europe/Sarajevo`, `Skopje`, `Europe/Skopje`, `Stockholm`, `Europe/Stockholm`, `Vienna`, `Europe/Vienna`, `Warsaw`, `Europe/Warsaw`, `West Central Africa`, `Africa/Algiers`, `Zagreb`, `Europe/Zagreb`, `Zurich`, `Athens`, `Europe/Athens`, `Bucharest`, `Europe/Bucharest`, `Cairo`, `Africa/Cairo`, `Harare`, `Africa/Harare`, `Helsinki`, `Europe/Helsinki`, `Jerusalem`, `Asia/Jerusalem`, `Kaliningrad`, `Europe/Kaliningrad`, `Kyiv`, `Europe/Kiev`, `Pretoria`, `Africa/Johannesburg`, `Riga`, `Europe/Riga`, `Sofia`, `Europe/Sofia`, `Tallinn`, `Europe/Tallinn`, `Vilnius`, `Europe/Vilnius`, `Baghdad`, `Asia/Baghdad`, `Istanbul`, `Europe/Istanbul`, `Kuwait`, `Asia/Kuwait`, `Minsk`, `Europe/Minsk`, `Moscow`, `Europe/Moscow`, `Nairobi`, `Africa/Nairobi`, `Riyadh`, `Asia/Riyadh`, `St. Petersburg`, `Volgograd`, `Europe/Volgograd`, `Tehran`, `Asia/Tehran`, `Abu Dhabi`, `Asia/Muscat`, `Baku`, `Asia/Baku`, `Muscat`, `Samara`, `Europe/Samara`, `Tbilisi`, `Asia/Tbilisi`, `Yerevan`, `Asia/Yerevan`, `Kabul`, `Asia/Kabul`, `Almaty`, `Asia/Almaty`, `Astana`, `Ekaterinburg`, `Asia/Yekaterinburg`, `Islamabad`, `Asia/Karachi`, `Karachi`, `Tashkent`, `Asia/Tashkent`, `Chennai`, `Asia/Kolkata`, `Kolkata`, `Mumbai`, `New Delhi`, `Sri Jayawardenepura`, `Asia/Colombo`, `Kathmandu`, `Asia/Kathmandu`, `Dhaka`, `Asia/Dhaka`, `Urumqi`, `Asia/Urumqi`, `Rangoon`, `Asia/Rangoon`, `Bangkok`, `Asia/Bangkok`, `Hanoi`, `Jakarta`, `Asia/Jakarta`, `Krasnoyarsk`, `Asia/Krasnoyarsk`, `Novosibirsk`, `Asia/Novosibirsk`, `Beijing`, `Asia/Shanghai`, `Chongqing`, `Asia/Chongqing`, `Hong Kong`, `Asia/Hong_Kong`, `Irkutsk`, `Asia/Irkutsk`, `Kuala Lumpur`, `Asia/Kuala_Lumpur`, `Perth`, `Australia/Perth`, `Singapore`, `Asia/Singapore`, `Taipei`, `Asia/Taipei`, `Ulaanbaatar`, `Asia/Ulaanbaatar`, `Osaka`, `Asia/Tokyo`, `Sapporo`, `Seoul`, `Asia/Seoul`, `Tokyo`, `Yakutsk`, `Asia/Yakutsk`, `Adelaide`, `Australia/Adelaide`, `Darwin`, `Australia/Darwin`, `Brisbane`, `Australia/Brisbane`, `Canberra`, `Australia/Canberra`, `Guam`, `Pacific/Guam`, `Hobart`, `Australia/Hobart`, `Melbourne`, `Australia/Melbourne`, `Port Moresby`, `Pacific/Port_Moresby`, `Sydney`, `Australia/Sydney`, `Vladivostok`, `Asia/Vladivostok`, `Magadan`, `Asia/Magadan`, `New Caledonia`, `Pacific/Noumea`, `Solomon Is.`, `Pacific/Guadalcanal`, `Srednekolymsk`, `Asia/Srednekolymsk`, `Auckland`, `Pacific/Auckland`, `Fiji`, `Pacific/Fiji`, `Kamchatka`, `Asia/Kamchatka`, `Marshall Is.`, `Pacific/Majuro`, `Wellington`, `Chatham Is.`, `Pacific/Chatham`, `Nuku'alofa`, `Pacific/Tongatapu`, `Samoa`, `Pacific/Apia`, `Tokelau Is.`, `Pacific/Fakaofo`, `America/Adak`, `America/Atka`, `US/Aleutian`, `America/Vancouver`, `Canada/Pacific`, `America/Miquelon`, `Australia/Eucla`, `Australia/LHI`, `Australia/Lord_Howe`, `Chile/EasterIsland`, `Pacific/Easter`, `Pacific/Gambier`, `Pacific/Pitcairn`, `Pacific/Marquesas`, `Pacific/Kiritimati`, `Pacific/Norfolk`.
- `time_restrictions` (Block List) If time restrictions are set, alerts will follow this path when they arrive within the specified time ranges and meet the rules. (see [below for nested schema](#nestedblock--time_restrictions))

### Read-Only

- `created_at` (String) Date of creation.
- `id` (String) The ID of the resource.
- `updated_at` (String) Date of last update.

<a id="nestedblock--notification_type_rules"></a>
### Nested Schema for `notification_type_rules`

Optional:

- `conditions` (Block List) Conditions combined per match_mode, at least one per rule. A deferral_window condition matches when the alert falls inside its time blocks. (see [below for nested schema](#nestedblock--notification_type_rules--conditions))
- `match_mode` (String) Whether all or any of the rule's conditions must match. Value must be one of `match-all-rules`, `match-any-rule`.
- `notification_type` (String) Outcome when this rule matches. Value must be one of `audible`, `quiet`.

<a id="nestedblock--notification_type_rules--conditions"></a>
### Nested Schema for `notification_type_rules.conditions`

Optional:

- `fieldable_id` (String) The ID of the alert field.
- `fieldable_type` (String) The type of the fieldable (e.g., AlertField). Value must be one of `AlertField`.
- `json_path` (String) JSON path to extract value from payload.
- `operator` (String) Whether the alert must (or must not) have related incidents. Value must be one of `is`, `is_not`, `contains`, `does_not_contain`, `is_one_of`, `is_not_one_of`, `is_empty`, `is_not_empty`, `contains_key`, `does_not_contain_key`, `starts_with`, `does_not_start_with`, `matches`, `does_not_match`, `is_set`, `is_not_set`.
- `rule_type` (String) The type of the escalation path rule. Value must be one of `alert_urgency`, `working_hour`, `json_path`, `field`, `service`, `deferral_window`, `source`, `related_incidents`.
- `service_ids` (Set of String) Service ids for which this escalation path should be used.
- `time_blocks` (Block List) Time windows during which alerts are deferred. (see [below for nested schema](#nestedblock--notification_type_rules--conditions--time_blocks))
- `time_zone` (String) Time zone for the deferral window. Value must be one of `International Date Line West`, `Etc/GMT+12`, `American Samoa`, `Pacific/Pago_Pago`, `Midway Island`, `Pacific/Midway`, `Hawaii`, `Pacific/Honolulu`, `Alaska`, `America/Juneau`, `Pacific Time (US & Canada)`, `America/Los_Angeles`, `Tijuana`, `America/Tijuana`, `Arizona`, `America/Phoenix`, `Mazatlan`, `America/Mazatlan`, `Mountain Time (US & Canada)`, `America/Denver`, `Central America`, `America/Guatemala`, `Central Time (US & Canada)`, `America/Chicago`, `Chihuahua`, `America/Chihuahua`, `Guadalajara`, `America/Mexico_City`, `Mexico City`, `Monterrey`, `America/Monterrey`, `Saskatchewan`, `America/Regina`, `Bogota`, `America/Bogota`, `Eastern Time (US & Canada)`, `America/New_York`, `Indiana (East)`, `America/Indiana/Indianapolis`, `Lima`, `America/Lima`, `Quito`, `Atlantic Time (Canada)`, `America/Halifax`, `Caracas`, `America/Caracas`, `Georgetown`, `America/Guyana`, `La Paz`, `America/La_Paz`, `Puerto Rico`, `America/Puerto_Rico`, `Santiago`, `America/Santiago`, `Newfoundland`, `America/St_Johns`, `Asuncion`, `America/Asuncion`, `Brasilia`, `America/Sao_Paulo`, `Buenos Aires`, `America/Argentina/Buenos_Aires`, `Montevideo`, `America/Montevideo`, `Greenland`, `America/Nuuk`, `Mid-Atlantic`, `Atlantic/South_Georgia`, `Azores`, `Atlantic/Azores`, `Cape Verde Is.`, `Atlantic/Cape_Verde`, `Edinburgh`, `Europe/London`, `Lisbon`, `Europe/Lisbon`, `London`, `Monrovia`, `Africa/Monrovia`, `UTC`, `Etc/UTC`, `Amsterdam`, `Europe/Amsterdam`, `Belgrade`, `Europe/Belgrade`, `Berlin`, `Europe/Berlin`, `Bern`, `Europe/Zurich`, `Bratislava`, `Europe/Bratislava`, `Brussels`, `Europe/Brussels`, `Budapest`, `Europe/Budapest`, `Casablanca`, `Africa/Casablanca`, `Copenhagen`, `Europe/Copenhagen`, `Dublin`, `Europe/Dublin`, `Ljubljana`, `Europe/Ljubljana`, `Madrid`, `Europe/Madrid`, `Paris`, `Europe/Paris`, `Prague`, `Europe/Prague`, `Rome`, `Europe/Rome`, `Sarajevo`, `Europe/Sarajevo`, `Skopje`, `Europe/Skopje`, `Stockholm`, `Europe/Stockholm`, `Vienna`, `Europe/Vienna`, `Warsaw`, `Europe/Warsaw`, `West Central Africa`, `Africa/Algiers`, `Zagreb`, `Europe/Zagreb`, `Zurich`, `Athens`, `Europe/Athens`, `Bucharest`, `Europe/Bucharest`, `Cairo`, `Africa/Cairo`, `Harare`, `Africa/Harare`, `Helsinki`, `Europe/Helsinki`, `Jerusalem`, `Asia/Jerusalem`, `Kaliningrad`, `Europe/Kaliningrad`, `Kyiv`, `Europe/Kiev`, `Pretoria`, `Africa/Johannesburg`, `Riga`, `Europe/Riga`, `Sofia`, `Europe/Sofia`, `Tallinn`, `Europe/Tallinn`, `Vilnius`, `Europe/Vilnius`, `Baghdad`, `Asia/Baghdad`, `Istanbul`, `Europe/Istanbul`, `Kuwait`, `Asia/Kuwait`, `Minsk`, `Europe/Minsk`, `Moscow`, `Europe/Moscow`, `Nairobi`, `Africa/Nairobi`, `Riyadh`, `Asia/Riyadh`, `St. Petersburg`, `Volgograd`, `Europe/Volgograd`, `Tehran`, `Asia/Tehran`, `Abu Dhabi`, `Asia/Muscat`, `Baku`, `Asia/Baku`, `Muscat`, `Samara`, `Europe/Samara`, `Tbilisi`, `Asia/Tbilisi`, `Yerevan`, `Asia/Yerevan`, `Kabul`, `Asia/Kabul`, `Almaty`, `Asia/Almaty`, `Astana`, `Ekaterinburg`, `Asia/Yekaterinburg`, `Islamabad`, `Asia/Karachi`, `Karachi`, `Tashkent`, `Asia/Tashkent`, `Chennai`, `Asia/Kolkata`, `Kolkata`, `Mumbai`, `New Delhi`, `Sri Jayawardenepura`, `Asia/Colombo`, `Kathmandu`, `Asia/Kathmandu`, `Dhaka`, `Asia/Dhaka`, `Urumqi`, `Asia/Urumqi`, `Rangoon`, `Asia/Rangoon`, `Bangkok`, `Asia/Bangkok`, `Hanoi`, `Jakarta`, `Asia/Jakarta`, `Krasnoyarsk`, `Asia/Krasnoyarsk`, `Novosibirsk`, `Asia/Novosibirsk`, `Beijing`, `Asia/Shanghai`, `Chongqing`, `Asia/Chongqing`, `Hong Kong`, `Asia/Hong_Kong`, `Irkutsk`, `Asia/Irkutsk`, `Kuala Lumpur`, `Asia/Kuala_Lumpur`, `Perth`, `Australia/Perth`, `Singapore`, `Asia/Singapore`, `Taipei`, `Asia/Taipei`, `Ulaanbaatar`, `Asia/Ulaanbaatar`, `Osaka`, `Asia/Tokyo`, `Sapporo`, `Seoul`, `Asia/Seoul`, `Tokyo`, `Yakutsk`, `Asia/Yakutsk`, `Adelaide`, `Australia/Adelaide`, `Darwin`, `Australia/Darwin`, `Brisbane`, `Australia/Brisbane`, `Canberra`, `Australia/Canberra`, `Guam`, `Pacific/Guam`, `Hobart`, `Australia/Hobart`, `Melbourne`, `Australia/Melbourne`, `Port Moresby`, `Pacific/Port_Moresby`, `Sydney`, `Australia/Sydney`, `Vladivostok`, `Asia/Vladivostok`, `Magadan`, `Asia/Magadan`, `New Caledonia`, `Pacific/Noumea`, `Solomon Is.`, `Pacific/Guadalcanal`, `Srednekolymsk`, `Asia/Srednekolymsk`, `Auckland`, `Pacific/Auckland`, `Fiji`, `Pacific/Fiji`, `Kamchatka`, `Asia/Kamchatka`, `Marshall Is.`, `Pacific/Majuro`, `Wellington`, `Chatham Is.`, `Pacific/Chatham`, `Nuku'alofa`, `Pacific/Tongatapu`, `Samoa`, `Pacific/Apia`, `Tokelau Is.`, `Pacific/Fakaofo`, `America/Adak`, `America/Atka`, `US/Aleutian`, `America/Vancouver`, `Canada/Pacific`, `America/Miquelon`, `Australia/Eucla`, `Australia/LHI`, `Australia/Lord_Howe`, `Chile/EasterIsland`, `Pacific/Easter`, `Pacific/Gambier`, `Pacific/Pitcairn`, `Pacific/Marquesas`, `Pacific/Kiritimati`, `Pacific/Norfolk`.
- `urgency_ids` (Set of String) Alert urgency ids for which this escalation path should be used.
- `value` (String) Value with which JSON path value should be matched.
- `values` (Set of String) Alert source values to match against (e.g., manual, datadog).
- `within_working_hour` (Boolean) Whether the escalation path should be used within working hours.

<a id="nestedblock--notification_type_rules--conditions--time_blocks"></a>
### Nested Schema for `notification_type_rules.conditions.time_blocks`

Optional:

- `all_day` (Boolean)
- `end_time` (String) Formatted as HH:MM.
- `ends_next_day` (Boolean) Whether the window crosses midnight. Derived from start_time and end_time; accepted and ignored on write.
- `friday` (Boolean)
- `id` (String) Unique ID of the time block.
- `monday` (Boolean)
- `position` (Number) Order of this time block, starting at 1. Defaults to the block's 1-based position in time_blocks when omitted.
- `saturday` (Boolean)
- `start_time` (String) Formatted as HH:MM.
- `sunday` (Boolean)
- `thursday` (Boolean)
- `tuesday` (Boolean)
- `wednesday` (Boolean)




<a id="nestedblock--rules"></a>
### Nested Schema for `rules`

Optional:

- `fieldable_id` (String) The ID of the alert field.
- `fieldable_type` (String) The type of the fieldable (e.g., AlertField). Value must be one of `AlertField`.
- `json_path` (String) JSON path to extract value from payload.
- `operator` (String) Whether the alert must (or must not) have related incidents. Value must be one of `is`, `is_not`, `contains`, `does_not_contain`, `is_one_of`, `is_not_one_of`, `is_empty`, `is_not_empty`, `contains_key`, `does_not_contain_key`, `starts_with`, `does_not_start_with`, `matches`, `does_not_match`, `is_set`, `is_not_set`.
- `rule_type` (String) The type of the escalation path rule. Value must be one of `alert_urgency`, `working_hour`, `json_path`, `field`, `service`, `deferral_window`, `source`, `related_incidents`.
- `service_ids` (Set of String) Service ids for which this escalation path should be used.
- `time_blocks` (Block List) Time windows during which alerts are deferred. (see [below for nested schema](#nestedblock--rules--time_blocks))
- `time_zone` (String) Time zone for the deferral window. Value must be one of `International Date Line West`, `Etc/GMT+12`, `American Samoa`, `Pacific/Pago_Pago`, `Midway Island`, `Pacific/Midway`, `Hawaii`, `Pacific/Honolulu`, `Alaska`, `America/Juneau`, `Pacific Time (US & Canada)`, `America/Los_Angeles`, `Tijuana`, `America/Tijuana`, `Arizona`, `America/Phoenix`, `Mazatlan`, `America/Mazatlan`, `Mountain Time (US & Canada)`, `America/Denver`, `Central America`, `America/Guatemala`, `Central Time (US & Canada)`, `America/Chicago`, `Chihuahua`, `America/Chihuahua`, `Guadalajara`, `America/Mexico_City`, `Mexico City`, `Monterrey`, `America/Monterrey`, `Saskatchewan`, `America/Regina`, `Bogota`, `America/Bogota`, `Eastern Time (US & Canada)`, `America/New_York`, `Indiana (East)`, `America/Indiana/Indianapolis`, `Lima`, `America/Lima`, `Quito`, `Atlantic Time (Canada)`, `America/Halifax`, `Caracas`, `America/Caracas`, `Georgetown`, `America/Guyana`, `La Paz`, `America/La_Paz`, `Puerto Rico`, `America/Puerto_Rico`, `Santiago`, `America/Santiago`, `Newfoundland`, `America/St_Johns`, `Asuncion`, `America/Asuncion`, `Brasilia`, `America/Sao_Paulo`, `Buenos Aires`, `America/Argentina/Buenos_Aires`, `Montevideo`, `America/Montevideo`, `Greenland`, `America/Nuuk`, `Mid-Atlantic`, `Atlantic/South_Georgia`, `Azores`, `Atlantic/Azores`, `Cape Verde Is.`, `Atlantic/Cape_Verde`, `Edinburgh`, `Europe/London`, `Lisbon`, `Europe/Lisbon`, `London`, `Monrovia`, `Africa/Monrovia`, `UTC`, `Etc/UTC`, `Amsterdam`, `Europe/Amsterdam`, `Belgrade`, `Europe/Belgrade`, `Berlin`, `Europe/Berlin`, `Bern`, `Europe/Zurich`, `Bratislava`, `Europe/Bratislava`, `Brussels`, `Europe/Brussels`, `Budapest`, `Europe/Budapest`, `Casablanca`, `Africa/Casablanca`, `Copenhagen`, `Europe/Copenhagen`, `Dublin`, `Europe/Dublin`, `Ljubljana`, `Europe/Ljubljana`, `Madrid`, `Europe/Madrid`, `Paris`, `Europe/Paris`, `Prague`, `Europe/Prague`, `Rome`, `Europe/Rome`, `Sarajevo`, `Europe/Sarajevo`, `Skopje`, `Europe/Skopje`, `Stockholm`, `Europe/Stockholm`, `Vienna`, `Europe/Vienna`, `Warsaw`, `Europe/Warsaw`, `West Central Africa`, `Africa/Algiers`, `Zagreb`, `Europe/Zagreb`, `Zurich`, `Athens`, `Europe/Athens`, `Bucharest`, `Europe/Bucharest`, `Cairo`, `Africa/Cairo`, `Harare`, `Africa/Harare`, `Helsinki`, `Europe/Helsinki`, `Jerusalem`, `Asia/Jerusalem`, `Kaliningrad`, `Europe/Kaliningrad`, `Kyiv`, `Europe/Kiev`, `Pretoria`, `Africa/Johannesburg`, `Riga`, `Europe/Riga`, `Sofia`, `Europe/Sofia`, `Tallinn`, `Europe/Tallinn`, `Vilnius`, `Europe/Vilnius`, `Baghdad`, `Asia/Baghdad`, `Istanbul`, `Europe/Istanbul`, `Kuwait`, `Asia/Kuwait`, `Minsk`, `Europe/Minsk`, `Moscow`, `Europe/Moscow`, `Nairobi`, `Africa/Nairobi`, `Riyadh`, `Asia/Riyadh`, `St. Petersburg`, `Volgograd`, `Europe/Volgograd`, `Tehran`, `Asia/Tehran`, `Abu Dhabi`, `Asia/Muscat`, `Baku`, `Asia/Baku`, `Muscat`, `Samara`, `Europe/Samara`, `Tbilisi`, `Asia/Tbilisi`, `Yerevan`, `Asia/Yerevan`, `Kabul`, `Asia/Kabul`, `Almaty`, `Asia/Almaty`, `Astana`, `Ekaterinburg`, `Asia/Yekaterinburg`, `Islamabad`, `Asia/Karachi`, `Karachi`, `Tashkent`, `Asia/Tashkent`, `Chennai`, `Asia/Kolkata`, `Kolkata`, `Mumbai`, `New Delhi`, `Sri Jayawardenepura`, `Asia/Colombo`, `Kathmandu`, `Asia/Kathmandu`, `Dhaka`, `Asia/Dhaka`, `Urumqi`, `Asia/Urumqi`, `Rangoon`, `Asia/Rangoon`, `Bangkok`, `Asia/Bangkok`, `Hanoi`, `Jakarta`, `Asia/Jakarta`, `Krasnoyarsk`, `Asia/Krasnoyarsk`, `Novosibirsk`, `Asia/Novosibirsk`, `Beijing`, `Asia/Shanghai`, `Chongqing`, `Asia/Chongqing`, `Hong Kong`, `Asia/Hong_Kong`, `Irkutsk`, `Asia/Irkutsk`, `Kuala Lumpur`, `Asia/Kuala_Lumpur`, `Perth`, `Australia/Perth`, `Singapore`, `Asia/Singapore`, `Taipei`, `Asia/Taipei`, `Ulaanbaatar`, `Asia/Ulaanbaatar`, `Osaka`, `Asia/Tokyo`, `Sapporo`, `Seoul`, `Asia/Seoul`, `Tokyo`, `Yakutsk`, `Asia/Yakutsk`, `Adelaide`, `Australia/Adelaide`, `Darwin`, `Australia/Darwin`, `Brisbane`, `Australia/Brisbane`, `Canberra`, `Australia/Canberra`, `Guam`, `Pacific/Guam`, `Hobart`, `Australia/Hobart`, `Melbourne`, `Australia/Melbourne`, `Port Moresby`, `Pacific/Port_Moresby`, `Sydney`, `Australia/Sydney`, `Vladivostok`, `Asia/Vladivostok`, `Magadan`, `Asia/Magadan`, `New Caledonia`, `Pacific/Noumea`, `Solomon Is.`, `Pacific/Guadalcanal`, `Srednekolymsk`, `Asia/Srednekolymsk`, `Auckland`, `Pacific/Auckland`, `Fiji`, `Pacific/Fiji`, `Kamchatka`, `Asia/Kamchatka`, `Marshall Is.`, `Pacific/Majuro`, `Wellington`, `Chatham Is.`, `Pacific/Chatham`, `Nuku'alofa`, `Pacific/Tongatapu`, `Samoa`, `Pacific/Apia`, `Tokelau Is.`, `Pacific/Fakaofo`, `America/Adak`, `America/Atka`, `US/Aleutian`, `America/Vancouver`, `Canada/Pacific`, `America/Miquelon`, `Australia/Eucla`, `Australia/LHI`, `Australia/Lord_Howe`, `Chile/EasterIsland`, `Pacific/Easter`, `Pacific/Gambier`, `Pacific/Pitcairn`, `Pacific/Marquesas`, `Pacific/Kiritimati`, `Pacific/Norfolk`.
- `urgency_ids` (Set of String) Alert urgency ids for which this escalation path should be used.
- `value` (String) Value with which JSON path value should be matched.
- `values` (Set of String) Alert source values to match against (e.g., manual, datadog).
- `within_working_hour` (Boolean) Whether the escalation path should be used within working hours.

<a id="nestedblock--rules--time_blocks"></a>
### Nested Schema for `rules.time_blocks`

Optional:

- `all_day` (Boolean)
- `end_time` (String) Formatted as HH:MM.
- `ends_next_day` (Boolean) Whether the window crosses midnight. Derived from start_time and end_time; accepted and ignored on write.
- `friday` (Boolean)
- `id` (String) Unique ID of the time block.
- `monday` (Boolean)
- `position` (Number) Order of this time block, starting at 1. Defaults to the block's 1-based position in time_blocks when omitted.
- `saturday` (Boolean)
- `start_time` (String) Formatted as HH:MM.
- `sunday` (Boolean)
- `thursday` (Boolean)
- `tuesday` (Boolean)
- `wednesday` (Boolean)



<a id="nestedblock--time_restrictions"></a>
### Nested Schema for `time_restrictions`

Optional:

- `end_day` (String) Value must be one of `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, `saturday`, `sunday`.
- `end_time` (String) Formatted as HH:MM.
- `start_day` (String) Value must be one of `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, `saturday`, `sunday`.
- `start_time` (String) Formatted as HH:MM.

## Import

rootly_escalation_path can be imported using the [`import` command](https://developer.hashicorp.com/terraform/cli/commands/import).

```sh
terraform import rootly_escalation_path.primary a816421c-6ceb-481a-87c4-585e47451f24
```

Or using an [`import` block](https://developer.hashicorp.com/terraform/language/import).

```terraform
import {
  to = rootly_escalation_path.primary
  id = "a816421c-6ceb-481a-87c4-585e47451f24"
}
```

Locate the resource id in the web app, or retrieve it by listing resources through the API if it's not visible in the web app.

HCL can be generated from the import block using the `-generate-config-out` flag.

```sh
terraform plan -generate-config-out=generated.tf
```
