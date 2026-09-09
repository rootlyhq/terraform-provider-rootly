data "rootly_user" "john" {
  email = "john@acme.com"
}

resource "rootly_user_on_call_role" "john" {
  user_id         = data.rootly_user.john.id
  on_call_role_id = "54725583-b47a-4397-aa50-c6b1db3c1da0"
}
