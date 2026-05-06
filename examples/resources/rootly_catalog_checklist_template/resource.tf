resource "rootly_catalog_checklist_template" "production_readiness" {
  catalog_id = rootly_catalog.customer_tier.id
  name       = "Production Readiness"
}
