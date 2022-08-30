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
				Optional: false,
			},
			

			"description": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

				"notify_emails": &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Computed: true,
					Optional: true,
				},
				

			"color": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

				"slack_channels": &schema.Schema{
					Type: schema.TypeList,
					Computed: true,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"id": &schema.Schema{
								Type: schema.TypeString,
								Required: true,
							},
							"name": &schema.Schema{
								Type: schema.TypeString,
								Required: true,
							},
						},
					},
				},
				

				"slack_aliases": &schema.Schema{
					Type: schema.TypeList,
					Computed: true,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"id": &schema.Schema{
								Type: schema.TypeString,
								Required: true,
							},
							"name": &schema.Schema{
								Type: schema.TypeString,
								Required: true,
							},
						},
					},
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

func dataSourceTeamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListTeamsParams)
	page_size := 1
	params.PageSize = &page_size

	
				name := d.Get("name").(string)
				params.FilterName = &name
			

				color := d.Get("color").(string)
				params.FilterColor = &color
			

	items, err := c.ListTeams(params)
	if err != nil {
		return diag.FromErr(err)
	}

	item, _ := items[0].(*client.Team)

	d.SetId(item.ID)

	return nil
}
