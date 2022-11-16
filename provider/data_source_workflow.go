package provider

import (
	"context"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

func dataSourceWorkflow() *schema.Resource{
	return &schema.Resource{
		ReadContext: dataSourceWorkflowRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			
			"name": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

			"slug": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			
		},
	}
}

func dataSourceWorkflowRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListWorkflowsParams)
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
			

	items, err := c.ListWorkflows(params)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(items) == 0 {
		return diag.Errorf("workflow not found")
	}
	item, _ := items[0].(*client.Workflow)

	d.SetId(item.ID)

	return nil
}
