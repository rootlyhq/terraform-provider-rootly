package provider

import (
	"context"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

func dataSourceStatusPage() *schema.Resource{
	return &schema.Resource{
		ReadContext: dataSourceStatusPageRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			
			"title": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: false,
			},
			

			"description": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

			"header_color": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

			"footer_color": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

			"allow_search_engine_index": &schema.Schema{
				Type: schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			

			"show_uptime": &schema.Schema{
				Type: schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			

			"show_uptime_last_days": &schema.Schema{
				Type: schema.TypeMap,
				Elem: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

			"public": &schema.Schema{
				Type: schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			

			"enabled": &schema.Schema{
				Type: schema.TypeBool,
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

func dataSourceStatusPageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListStatusPagesParams)
	page_size := 1
	params.PageSize = &page_size

	

	items, err := c.ListStatusPages(params)
	if err != nil {
		return diag.FromErr(err)
	}

	item, _ := items[0].(*client.StatusPage)

	d.SetId(item.ID)

	return nil
}
