resource "rootly_catalog_property" "tier" {
  catalog_id = rootly_catalog.services.id
  name       = "Service Tier"
  kind       = "select"
}
