# Custom Fields
resource "rootly_custom_field" "regions_affected" {
  name = "Regions affected"
  kind = "multi_select"
  shown = ["incident_form"]
  required = ["incident_form"]
}

resource "rootly_custom_field_option" "asia" {
  custom_field_id = rootly_custom_field.regions_affected.id
  value = "Asia"
}

resource "rootly_custom_field_option" "europe" {
  custom_field_id = rootly_custom_field.regions_affected.id
  value = "Europe"
}

resource "rootly_custom_field_option" "north_america" {
  custom_field_id = rootly_custom_field.regions_affected.id
  value = "North America"
}
