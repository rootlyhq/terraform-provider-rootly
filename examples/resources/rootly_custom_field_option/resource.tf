resource "rootly_custom_field_option" "example" {
  custom_field_id = rootly_custom_field.region.id
  value           = "US West"
  color           = "#00FF00"
}
