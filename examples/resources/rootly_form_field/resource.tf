resource "rootly_form_field" "regions_affected" {
  name       = "Regions affected"
  kind       = "custom"
  input_kind = "multi_select"
  shown      = ["web_new_incident_form", "web_update_incident_form"]
  required   = ["web_new_incident_form", "web_update_incident_form"]
}

resource "rootly_form_field_option" "asia" {
  form_field_id = rootly_form_field.regions_affected.id
  value         = "Asia"
}

resource "rootly_form_field_option" "europe" {
  form_field_id = rootly_form_field.regions_affected.id
  value         = "Europe"
}

resource "rootly_form_field_option" "north_america" {
  form_field_id = rootly_form_field.regions_affected.id
  value         = "North America"
}

resource "rootly_form_field" "problem_status" {
  name          = "Problem status"
  kind          = "custom"
  input_kind    = "select"
  resource_type = "problem"
}

resource "rootly_form_field_option" "problem_status_investigating" {
  form_field_id = rootly_form_field.problem_status.id
  value         = "Investigating"
}

resource "rootly_form_field_option" "problem_status_resolved" {
  form_field_id = rootly_form_field.problem_status.id
  value         = "Resolved"
}
