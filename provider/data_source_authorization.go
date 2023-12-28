package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

func dataSourceAuthorization() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAuthorizationRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"authorizable_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"authorizable_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"grantee_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"grantee_type": &schema.Schema{
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

func dataSourceAuthorizationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListAuthorizationsParams)
	page_size := 1
	params.PageSize = &page_size

	if value, ok := d.GetOkExists("authorizable_id"); ok {
		authorizable_id := value.(string)
		params.FilterAuthorizableId = &authorizable_id
	}

	if value, ok := d.GetOkExists("authorizable_type"); ok {
		authorizable_type := value.(string)
		params.FilterAuthorizableType = &authorizable_type
	}

	if value, ok := d.GetOkExists("grantee_id"); ok {
		grantee_id := value.(string)
		params.FilterGranteeId = &grantee_id
	}

	if value, ok := d.GetOkExists("grantee_type"); ok {
		grantee_type := value.(string)
		params.FilterGranteeType = &grantee_type
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

	items, err := c.ListAuthorizations(params)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(items) == 0 {
		return diag.Errorf("authorization not found")
	}
	item, _ := items[0].(*client.Authorization)

	d.SetId(item.ID)

	return nil
}
