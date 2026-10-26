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

# Rotate the key by changing rotation_trigger, which replaces it with a new key.
# The token is only returned on create; write it to a secret store so consumers
# pick up the new value. Names must be unique, so include the trigger in the name
# if you use create_before_destroy.
resource "rootly_api_key" "rotated" {
  name       = "Deploy Bot v${local.deploy_bot_key_version}"
  expires_at = "2027-01-01T00:00:00Z"
  kind       = "team"
  group_id   = rootly_team.payments.id

  rotation_trigger = local.deploy_bot_key_version

  lifecycle {
    create_before_destroy = true
  }
}

locals {
  deploy_bot_key_version = "1"
}
