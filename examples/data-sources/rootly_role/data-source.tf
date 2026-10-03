# Built-in roles examples
data "rootly_role" "owner" {
  slug = "owner"
}

data "rootly_role" "admin" {
  slug = "admin"
}

data "rootly_role" "user" {
  slug = "user"
}

data "rootly_role" "observer" {
  slug = "observer"
}

# This is the "None" role in the UI
data "rootly_role" "no_access" {
  slug = "no_access"
}
