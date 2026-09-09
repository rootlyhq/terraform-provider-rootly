# Ad-hoc component with its own name and description
resource "rootly_status_page_component" "website" {
  status_page_id                 = rootly_status_page.public.id
  status_page_component_group_id = rootly_status_page_component_group.infrastructure.id
  name                           = "Website"
  description                    = "Public marketing website"
}

# Catalog-backed component mirroring a service; name and description
# are derived from the backing service.
resource "rootly_status_page_component" "database" {
  status_page_id = rootly_status_page.public.id
  source_type    = "Service"
  source_id      = rootly_service.elasticsearch_prod.id
}
