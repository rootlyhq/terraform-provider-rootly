package provider

import (
	"context"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

func dataSourceCause() *schema.Resource{
	return &schema.Resource{
		ReadContext: dataSourceCauseRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			
			"name": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: false,
			},
			

			"slug": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

			"description": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

			"created_at": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: false,
			},
			

			"updated_at": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: false,
			},
			
		},
	}
}

func dataSourceCauseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListCausesParams)
	page_size := 1
	params.PageSize = &page_size

	
				slug := d.Get("slug").(string)
				params.FilterSlug = &slug
			

				name := d.Get("name").(string)
				params.FilterName = &name
			

	items, err := c.ListCauses(params)
	if err != nil {
		return diag.FromErr(err)
	}

	item, _ := items[0].(*client.Cause)

	d.SetId(item.ID)

	return nil
}
