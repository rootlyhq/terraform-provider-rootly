package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskCreateAirtableTableRecord(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskCreateAirtableTableRecord,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "name", "test-workflow"),
				),
			},
			{
				Config: testAccResourceWorkflowTaskCreateAirtableTableRecordUpdate,
			},
		},
	})
}

const testAccResourceWorkflowTaskCreateAirtableTableRecord = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_create_airtable_table_record" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		base_key = "test"
table_name = "test"
	}
}
`

const testAccResourceWorkflowTaskCreateAirtableTableRecordUpdate = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_create_airtable_table_record" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		base_key = "test"
table_name = "test"
	}
}
`