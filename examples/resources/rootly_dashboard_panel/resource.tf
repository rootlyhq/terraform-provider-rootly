resource "rootly_dashboard" "overview" {
  name = "my_dashboard"
}

resource "rootly_dashboard_panel" "incidents_by_severity" {
  dashboard_id = rootly_dashboard.foo.id
  name         = "test"
  params {
    display = "line_chart"
    datasets {
      collection = "incidents"
      filter {
        operation = "and"
        rules {
          operation = "and"
          condition = "="
          key       = "status"
          value     = "started"
        }
      }
      group_by = "severity"
      aggregate {
        cumulative = false
        key        = "results"
        operation  = "count"
      }
    }
  }
}
