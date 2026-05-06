resource "rootly_cause" "bad_deploy" {
  name        = "Bad Deploy"
  description = "Incident caused by a faulty deployment"
}

resource "rootly_cause" "infrastructure_failure" {
  name        = "Infrastructure Failure"
  description = "Hardware or cloud provider issue"
}
