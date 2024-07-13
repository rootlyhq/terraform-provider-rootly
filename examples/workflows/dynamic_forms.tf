data "rootly_form_field" "severity" {
  slug = "severity"
}

data "rootly_form_field" "cause" {
  slug = "causes"
}

data "rootly_severity" "sev0" {
  slug = "sev0"
}

# create a new form set conditioned on sev0 severity
resource "rootly_form_set" "sev0" {
  name = "Resolving sev0 incidents requires cause"
  forms = ["web_incident_resolution_form"]
}

# condition the form set on sev0 severity
resource "rootly_form_set_condition" "sev0" {
  form_set_id = rootly_form_set.sev0.id
  form_field_id = data.rootly_form_field.severity.id
  values = [data.rootly_severity.sev0.id]
}

# place the cause field on the resolution form for this form set
resource "rootly_form_field_placement" "severity" {
  depends_on = [rootly_form_set_condition.sev0]
  form_set_id = rootly_form_set.sev0.id
  form_field_id = data.rootly_form_field.severity.id
  form = "web_incident_resolution_form"
  required = true
  position = 1
}

resource "rootly_form_field_placement" "cause" {
  depends_on = [rootly_form_set_condition.sev0]
  form_set_id = rootly_form_set.sev0.id
  form_field_id = data.rootly_form_field.cause.id
  form = "web_incident_resolution_form"
  required = true
  position = 2
}
