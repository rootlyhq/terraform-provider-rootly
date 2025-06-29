// DO NOT MODIFY: This file is generated by tools/generate.js. Any changes will be overwritten during the next build.

package provider

import (
	"context"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

func dataSourceFormFieldPosition() *schema.Resource {
	return &schema.Resource {
		ReadContext: dataSourceFormFieldPositionRead,
		Schema: map[string]*schema.Schema {
			"id": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
			},
			
			"form": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
		ValidateFunc: validation.StringInSlice([]string{"web_new_incident_form", "web_update_incident_form", "web_incident_post_mortem_form", "web_incident_mitigation_form", "web_incident_resolution_form", "web_incident_cancellation_form", "web_scheduled_incident_form", "web_update_scheduled_incident_form", "incident_post_mortem", "slack_new_incident_form", "slack_update_incident_form", "slack_update_incident_status_form", "slack_incident_mitigation_form", "slack_incident_resolution_form", "slack_incident_cancellation_form", "slack_scheduled_incident_form", "slack_update_scheduled_incident_form"}, false),
			},
			
		},
	}
}

func dataSourceFormFieldPositionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListFormFieldPositionsParams)
	page_size := 1
	params.PageSize = &page_size

	
				if value, ok := d.GetOkExists("form"); ok {
					form := value.(string)
					params.FilterForm = &form
				}
			

	form_field_id := d.Get("form_field_id").(string)
			items, err := c.ListFormFieldPositions(form_field_id, params)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(items) == 0 {
		return diag.Errorf("form_field_position not found")
	}
	item, _ := items[0].(*client.FormFieldPosition)

	d.SetId(item.ID)

	return nil
}
