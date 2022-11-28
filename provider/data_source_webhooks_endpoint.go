package provider

import (
	"context"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

func dataSourceWebhooksEndpoint() *schema.Resource{
	return &schema.Resource{
		ReadContext: dataSourceWebhooksEndpointRead,
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

func dataSourceWebhooksEndpointRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListWebhooksEndpointsParams)
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
			

	items, err := c.ListWebhooksEndpoints(params)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(items) == 0 {
		return diag.Errorf("webhooks_endpoint not found")
	}
	item, _ := items[0].(*client.WebhooksEndpoint)

	d.SetId(item.ID)

	return nil
}
