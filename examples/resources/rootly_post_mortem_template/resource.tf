resource "rootly_post_mortem_template" "default" {
  name    = "Standard Post-Mortem"
  default = true
  format  = "markdown"
  content = <<-EOT
## Incident Summary
{{ incident.title }}

## Timeline
{{ incident.timeline }}

## Root Cause
_To be filled in_

## Action Items
{{ incident.action_items }}
EOT
}
