resource "rootly_status_page_component_group" "infrastructure" {
  status_page_id       = rootly_status_page.public.id
  name                 = "Infrastructure"
  description          = "Core infrastructure services"
  collapsed_by_default = false
}
