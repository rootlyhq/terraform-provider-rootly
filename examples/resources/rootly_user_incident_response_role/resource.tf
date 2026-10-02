data "rootly_user" "john" {
  email = "john@acme.com"
}

data "rootly_role" "user" {
  slug = "user"
}

resource "rootly_user_incident_response_role" "john" {
  user_id = data.rootly_user.john.id
  role_id = data.rootly_role.user.id
}
