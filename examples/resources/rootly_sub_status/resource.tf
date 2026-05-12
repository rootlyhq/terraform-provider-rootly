resource "rootly_sub_status" "investigating" {
  name          = "Investigating"
  parent_status = "started"
  description   = "Responders are investigating"
}

resource "rootly_sub_status" "fix_in_progress" {
  name          = "Fix In Progress"
  parent_status = "started"
  description   = "A fix is being implemented"
}
