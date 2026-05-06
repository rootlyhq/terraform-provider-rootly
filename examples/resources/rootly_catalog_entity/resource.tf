resource "rootly_catalog_entity" "payments_api" {
  catalog_id  = rootly_catalog.services.id
  name        = "Payments API"
  description = "Handles payment processing"
}
