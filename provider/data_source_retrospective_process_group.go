package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

func dataSourceRetrospectiveProcessGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRetrospectiveProcessGroupRead,
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
		},
	}
}

func dataSourceRetrospectiveProcessGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListRetrospectiveProcessGroupsParams)
	page_size := 1
	params.PageSize = &page_size

	if value, ok := d.GetOkExists("sub_status_id"); ok {
		sub_status_id := value.(string)
		params.FilterSubStatusId = &sub_status_id
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

	retrospective_process_id := d.Get("retrospective_process_id").(string)
	items, err := c.ListRetrospectiveProcessGroups(retrospective_process_id, params)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(items) == 0 {
		return diag.Errorf("retrospective_process_group not found")
	}
	item, _ := items[0].(*client.RetrospectiveProcessGroup)

	d.SetId(item.ID)

	return nil
}