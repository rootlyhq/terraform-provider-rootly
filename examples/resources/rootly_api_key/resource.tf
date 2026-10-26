resource "rootly_api_key" "ci_pipeline" {
  name       = "CI Pipeline"
  expires_at = "2027-01-01T00:00:00Z"
}

resource "rootly_team" "payments" {
  name = "Payments"
}

resource "rootly_api_key" "payments_team" {
  name       = "Payments Team"
  expires_at = "2027-01-01T00:00:00Z"
  kind       = "team"
  group_id   = rootly_team.payments.id
}

# Rotate a key by changing rotation_trigger. The API only lets a key rotate itself,
# so manage the key with a provider authenticated by that same key (for example,
# import the key your pipeline already uses). The new token is only returned once;
# write it to a secret store and point the provider at it before the next run.
provider "rootly" {
  alias     = "deploy_bot"
  api_token = var.deploy_bot_api_token
}

resource "rootly_api_key" "deploy_bot" {
  provider = rootly.deploy_bot

  name       = "Deploy Bot"
  expires_at = "2027-01-01T00:00:00Z"

  rotation_trigger              = "2026-10"
  rotation_grace_period_minutes = 60
}
