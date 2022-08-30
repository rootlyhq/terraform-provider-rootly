package provider

import (
	"context"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

func dataSourceWorkflowGroup() *schema.Resource{
	return &schema.Resource{
		ReadContext: dataSourceWorkflowGroupRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			
			"kind": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: false,
			},
			

			"name": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: false,
			},
			

			"expanded": &schema.Schema{
				Type: schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			

			"position": &schema.Schema{
				Type: schema.TypeInt,
				Computed: true,
				Optional: false,
			},
			
		},
	}
}

func dataSourceWorkflowGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListWorkflowGroupsParams)
	page_size := 1
	params.PageSize = &page_size

	
				name := d.Get("name").(string)
				params.FilterName = &name
			

				expanded := d.Get("expanded").(bool)
				params.FilterExpanded = &expanded
			

				position := d.Get("position").(int)
				params.FilterPosition = &position
			

	items, err := c.ListWorkflowGroups(params)
	if err != nil {
		return diag.FromErr(err)
	}

	item, _ := items[0].(*client.WorkflowGroup)

	d.SetId(item.ID)

	return nil
}
