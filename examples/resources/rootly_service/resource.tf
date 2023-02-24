resource "rootly_service" "elasticsearch_prod" {
  name  = "elasticsearch-prod"
  color = "#800080"
}

resource "rootly_service" "customer_postgresql_prod" {
  name  = "customer-postgresql-prod"
  color = "#800080"
}
