resource "rootly_form_field_placement_condition_group" "teams" {
  form_field_placement_id = rootly_form_field_placement.example.id
  group_type              = "required"
  group_operator          = "or"
}

resource "rootly_form_field_placement_condition" "team_a" {
  form_field_placement_id = rootly_form_field_placement.example.id
  condition_group_id      = rootly_form_field_placement_condition_group.teams.id
  conditioned             = "required"
  form_field_id           = rootly_form_field.teams.id
  comparison              = "equal"
  values                  = [rootly_team.a.id]
}

resource "rootly_form_field_placement_condition" "team_b" {
  form_field_placement_id = rootly_form_field_placement.example.id
  condition_group_id      = rootly_form_field_placement_condition_group.teams.id
  conditioned             = "required"
  form_field_id           = rootly_form_field.teams.id
  comparison              = "equal"
  values                  = [rootly_team.b.id]
}
