package provider

import (
	"context"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

func dataSourceTeam() *schema.Resource{
	return &schema.Resource{
		ReadContext: dataSourceTeamRead,
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
			

			"color": &schema.Schema{
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

func dataSourceTeamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListTeamsParams)
	page_size := 1
	params.PageSize = &page_size

	
				slug := d.Get("slug").(string)
				params.FilterSlug = &slug
			

				name := d.Get("name").(string)
				params.FilterName = &name
			

				color := d.Get("color").(string)
				params.FilterColor = &color
			

	items, err := c.ListTeams(params)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(items) == 0 {
		return diag.Errorf("team not found")
	}
	item, _ := items[0].(*client.Team)

	d.SetId(item.ID)

	return nil
}
