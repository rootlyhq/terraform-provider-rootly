package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

func dataSourceIncidentSubStatus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIncidentSubStatusRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"sub_status_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"assigned_at": &schema.Schema{
				Type:        schema.TypeMap,
				Description: "Filter by date range using 'lt' and 'gt'.",
				Optional:    true,
			},
		},
	}
}

func dataSourceIncidentSubStatusRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListIncidentSubStatusesParams)
	page_size := 1
	params.PageSize = &page_size

	if value, ok := d.GetOkExists("sub_status_id"); ok {
		sub_status_id := value.(string)
		params.FilterSubStatusId = &sub_status_id
	}

	assigned_at_gt := d.Get("assigned_at").(map[string]interface{})
	if value, exists := assigned_at_gt["gt"]; exists {
		v := value.(string)
		params.FilterAssignedAtGt = &v
	}

	assigned_at_lt := d.Get("assigned_at").(map[string]interface{})
	if value, exists := assigned_at_lt["lt"]; exists {
		v := value.(string)
		params.FilterAssignedAtLt = &v
	}

	incident_id := d.Get("incident_id").(string)
	items, err := c.ListIncidentSubStatuses(incident_id, params)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(items) == 0 {
		return diag.Errorf("incident_sub_status not found")
	}
	item, _ := items[0].(*client.IncidentSubStatus)

	d.SetId(item.ID)

	return nil
}