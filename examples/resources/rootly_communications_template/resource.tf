resource "rootly_communications_template" "status_update" {
  name = "Status Update Template"
  body = "**Incident:** {{ incident.title }}\n**Status:** {{ incident.status }}\n**Summary:** {{ incident.summary }}"
}
