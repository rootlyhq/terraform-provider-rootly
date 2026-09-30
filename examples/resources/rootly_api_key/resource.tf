resource "rootly_api_key" "ci_pipeline" {
  name       = "CI Pipeline"
  expires_at = "2027-01-01T00:00:00Z"
}

resource "rootly_api_key" "payments_team" {
  name       = "Payments Team"
  expires_at = "2027-01-01T00:00:00Z"
  kind       = "team"
  group_id   = rootly_team.payments.id
}
