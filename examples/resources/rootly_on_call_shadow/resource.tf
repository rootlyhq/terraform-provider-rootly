resource "rootly_on_call_shadow" "new_hire_shadow" {
  schedule_id     = rootly_schedule.primary_oncall.id
  shadow_user_id  = data.rootly_user.new_hire.id
  shadowable_id   = data.rootly_user.senior_engineer.id
  shadowable_type = "User"
  starts_at       = "2026-06-01T00:00:00Z"
  ends_at         = "2026-06-14T00:00:00Z"
}
