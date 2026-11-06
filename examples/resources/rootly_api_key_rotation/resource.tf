# Rotate the API key the provider is authenticated with every 30 days.
resource "time_rotating" "rootly_api_key" {
  rotation_days = 30
}

resource "rootly_api_key_rotation" "this" {
  # Must be the ID of the API key the provider is authenticated with.
  api_key_id = "00000000-0000-0000-0000-000000000000"

  # Keep the previous token valid for an hour, so the new token can be rolled out.
  grace_period_minutes = 60

  triggers = {
    rotation = time_rotating.rootly_api_key.id
  }
}

# Store the new token, then point the provider at it.
resource "aws_secretsmanager_secret_version" "rootly_api_token" {
  secret_id     = "rootly-api-token"
  secret_string = rootly_api_key_rotation.this.token
}
