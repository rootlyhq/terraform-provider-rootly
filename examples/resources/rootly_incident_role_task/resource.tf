resource "rootly_incident_role_task" "open_war_room" {
  incident_role_id = rootly_incident_role.incident_commander.id
  task             = "Open a war room and invite responders"
  priority         = "high"
}

resource "rootly_incident_role_task" "send_status_update" {
  incident_role_id = rootly_incident_role.communications_lead.id
  task             = "Send initial status update to stakeholders"
  priority         = "high"
}
