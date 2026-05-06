resource "rootly_cause" "deployment" {
  name        = "Deployment"
  description = "Caused by a recent deployment"
}

resource "rootly_cause" "configuration_change" {
  name        = "Configuration Change"
  description = "Caused by a configuration change"
}
