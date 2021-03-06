package provider

// This file was auto-generated by tools/gen_tasks.js

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskSnapshotGrafanaDashboard(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskSnapshotGrafanaDashboard,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "name", "test-workflow"),
				),
			},
			{
				Config: testAccResourceWorkflowTaskSnapshotGrafanaDashboardUpdate,
			},
		},
	})
}

const testAccResourceWorkflowTaskSnapshotGrafanaDashboard = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_snapshot_grafana_dashboard" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		dashboards {
						id = "foo"
						name = "bar"
					}
	}
}
`

const testAccResourceWorkflowTaskSnapshotGrafanaDashboardUpdate = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_snapshot_grafana_dashboard" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		dashboards {
						id = "foo"
						name = "bar"
					}
	}
}
`
