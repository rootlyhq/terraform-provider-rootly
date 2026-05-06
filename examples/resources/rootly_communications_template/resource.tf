resource "rootly_communications_template" "status_update" {
  name = "Incident Status Update"
  body = "**Incident:** {{ incident.title }}\n**Status:** {{ incident.status }}\n**Summary:** {{ incident.summary }}"
}
