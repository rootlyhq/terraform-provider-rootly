# Custom form Fields
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
