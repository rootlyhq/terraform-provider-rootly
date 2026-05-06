resource "rootly_sub_status" "investigating" {
  name          = "Investigating"
  parent_status = "started"
  description   = "Team is investigating the issue"
}

resource "rootly_sub_status" "fix_deployed" {
  name          = "Fix Deployed"
  parent_status = "resolved"
  description   = "A fix has been deployed"
}
