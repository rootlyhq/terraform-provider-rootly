resource "rootly_custom_field" "region" {
  label   = "Affected Region"
  kind    = "select"
  enabled = true
  shown   = ["incident_form", "incident_post_mortem"]
}

resource "rootly_custom_field_option" "us_east" {
  custom_field_id = rootly_custom_field.region.id
  value           = "US East"
  color           = "#047BF8"
}

resource "rootly_custom_field_option" "eu_west" {
  custom_field_id = rootly_custom_field.region.id
  value           = "EU West"
  color           = "#FFA500"
}
