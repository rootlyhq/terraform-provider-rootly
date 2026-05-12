resource "rootly_custom_field_option" "us_west_2" {
  custom_field_id = rootly_custom_field.affected_region.id
  value           = "us-west-2"
  color           = "#00AA00"
}
