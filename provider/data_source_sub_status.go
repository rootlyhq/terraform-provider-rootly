package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

func dataSourceSubStatus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSubStatusRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
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

			"parent_status": &schema.Schema{
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

func dataSourceSubStatusRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListSubStatusesParams)
	page_size := 1
	params.PageSize = &page_size

	if value, ok := d.GetOkExists("slug"); ok {
		slug := value.(string)
		params.FilterSlug = &slug
	}

	if value, ok := d.GetOkExists("name"); ok {
		name := value.(string)
		params.FilterName = &name
	}

	if value, ok := d.GetOkExists("parent_status"); ok {
		parent_status := value.(string)
		params.FilterParentStatus = &parent_status
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

	items, err := c.ListSubStatuses(params)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(items) == 0 {
		return diag.Errorf("sub_status not found")
	}
	item, _ := items[0].(*client.SubStatus)

	d.SetId(item.ID)

	return nil
}