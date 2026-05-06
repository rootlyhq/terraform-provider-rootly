resource "rootly_catalog_checklist_template" "readiness" {
  catalog_id = rootly_catalog.services.id
  name       = "Production Readiness"
}
