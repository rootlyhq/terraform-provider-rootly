resource "rootly_custom_field" "affected_region" {
  label   = "Affected Region"
  kind    = "select"
  enabled = true
  shown   = ["incident_form", "incident_post_mortem"]
}

resource "rootly_custom_field_option" "us_east_1" {
  custom_field_id = rootly_custom_field.affected_region.id
  value           = "us-east-1"
  color           = "#047BF8"
}

resource "rootly_custom_field_option" "eu_west_1" {
  custom_field_id = rootly_custom_field.affected_region.id
  value           = "eu-west-1"
  color           = "#FFA500"
}
