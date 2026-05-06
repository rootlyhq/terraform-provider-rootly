resource "rootly_catalog_property" "support_level" {
  catalog_id = rootly_catalog.customer_tier.id
  name       = "Support Level"
  kind       = "select"
}
