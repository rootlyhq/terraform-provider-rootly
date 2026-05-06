resource "rootly_catalog_entity" "enterprise_tier" {
  catalog_id  = rootly_catalog.customer_tier.id
  name        = "Enterprise"
  description = "Enterprise customer tier"
}

resource "rootly_catalog_entity" "growth_tier" {
  catalog_id  = rootly_catalog.customer_tier.id
  name        = "Growth"
  description = "Growth customer tier"
}
