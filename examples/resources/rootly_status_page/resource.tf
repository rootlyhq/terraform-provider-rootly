resource "rootly_status_page" "public" {
  title                 = "Acme Status"
  public_title          = "Acme System Status"
  description           = "Current status of Acme services"
  success_message       = "All Systems Operational"
  failure_message       = "Degraded Performance"
  show_uptime           = true
  show_uptime_last_days = 90
  service_ids           = [rootly_service.api.id, rootly_service.web.id]
  functionality_ids     = [rootly_functionality.payments.id]
}
