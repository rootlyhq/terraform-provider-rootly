package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

func dataSourceWorkflowTask() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWorkflowTaskRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"workflow_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"slug": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func dataSourceWorkflowTaskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListWorkflowTasksParams)
	page_size := 1
	params.PageSize = &page_size

	if value, ok := d.GetOkExists("name"); ok {
		name := value.(string)
		params.FilterName = &name
	}

	if value, ok := d.GetOkExists("slug"); ok {
		slug := value.(string)
		params.FilterSlug = &slug
	}

	workflow_id := d.Get("workflow_id").(string)

	items, err := c.ListWorkflowTasks(workflow_id, params)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(items) == 0 {
		return diag.Errorf("workflow task not found")
	}
	item, _ := items[0].(*client.WorkflowTask)

	d.SetId(item.ID)

	return nil
}
