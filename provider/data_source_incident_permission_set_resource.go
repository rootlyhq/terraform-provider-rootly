package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

func dataSourceIncidentPermissionSetResource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIncidentPermissionSetResourceRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"kind": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"created_at": &schema.Schema{
				Type:        schema.TypeMap,
				Description: "Filter by date range using 'lt' and 'gt'.",
				Optional:    true,
			},
		},
	}
}

func dataSourceIncidentPermissionSetResourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListIncidentPermissionSetResourcesParams)
	page_size := 1
	params.PageSize = &page_size

	if value, ok := d.GetOkExists("kind"); ok {
		kind := value.(string)
		params.FilterKind = &kind
	}

	created_at_gt := d.Get("created_at").(map[string]interface{})
	if value, exists := created_at_gt["gt"]; exists {
		v := value.(string)
		params.FilterCreatedAtGt = &v
	}

	created_at_lt := d.Get("created_at").(map[string]interface{})
	if value, exists := created_at_lt["lt"]; exists {
		v := value.(string)
		params.FilterCreatedAtLt = &v
	}

	incident_permission_set_id := d.Get("incident_permission_set_id").(string)
	items, err := c.ListIncidentPermissionSetResources(incident_permission_set_id, params)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(items) == 0 {
		return diag.Errorf("incident_permission_set_resource not found")
	}
	item, _ := items[0].(*client.IncidentPermissionSetResource)

	d.SetId(item.ID)

	return nil
}
