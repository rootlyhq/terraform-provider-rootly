package provider

import (
	"context"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

func dataSourceWorkflowCustomFieldSelection() *schema.Resource{
	return &schema.Resource{
		ReadContext: dataSourceWorkflowCustomFieldSelectionRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			
			"workflow_id": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

			"custom_field_id": &schema.Schema{
				Type: schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			

			"incident_condition": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

				"values": &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Computed: true,
					Optional: true,
				},
				

				"selected_option_ids": &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Computed: true,
					Optional: true,
				},
				
		},
	}
}

func dataSourceWorkflowCustomFieldSelectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListWorkflowCustomFieldSelectionsParams)
	page_size := 1
	params.PageSize = &page_size

	

	workflow_id := d.Get("workflow_id").(string)
			items, err := c.ListWorkflowCustomFieldSelections(workflow_id, params)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(items) == 0 {
		return diag.Errorf("workflow_custom_field_selection not found")
	}
	item, _ := items[0].(*client.WorkflowCustomFieldSelection)

	d.SetId(item.ID)

	return nil
}
