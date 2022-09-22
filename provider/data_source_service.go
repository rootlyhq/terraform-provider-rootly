package provider

import (
	"context"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

func dataSourceService() *schema.Resource{
	return &schema.Resource{
		ReadContext: dataSourceServiceRead,
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
			

			"backstage_id": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

			"created_at": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			
		},
	}
}

func dataSourceServiceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListServicesParams)
	page_size := 1
	params.PageSize = &page_size

	
				name := d.Get("name").(string)
				params.FilterName = &name
			

				slug := d.Get("slug").(string)
				params.FilterSlug = &slug
			

				backstage_id := d.Get("backstage_id").(string)
				params.FilterBackstageId = &backstage_id
			

	items, err := c.ListServices(params)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(items) == 0 {
		return diag.Errorf("service not found")
	}
	item, _ := items[0].(*client.Service)

	d.SetId(item.ID)

	return nil
}
