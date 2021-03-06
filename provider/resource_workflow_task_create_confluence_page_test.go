package provider

// This file was auto-generated by tools/gen_tasks.js

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskCreateConfluencePage(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskCreateConfluencePage,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "name", "test-workflow"),
				),
			},
			{
				Config: testAccResourceWorkflowTaskCreateConfluencePageUpdate,
			},
		},
	})
}

const testAccResourceWorkflowTaskCreateConfluencePage = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_create_confluence_page" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		space = {
					id = "foo"
					name = "bar"
				}
title = "test"
	}
}
`

const testAccResourceWorkflowTaskCreateConfluencePageUpdate = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_create_confluence_page" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		space = {
					id = "foo"
					name = "bar"
				}
title = "test"
	}
}
`
