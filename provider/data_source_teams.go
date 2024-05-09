package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
	"strconv"
	"time"
)

func dataSourceTeams() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTeamsRead,
		Schema: map[string]*schema.Schema{
			"slug": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"teams": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"slug": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"color": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTeamsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListTeamsParams)

	slug := d.Get("slug").(string)
	name := d.Get("name").(string)

	params.FilterSlug = &slug
	params.FilterName = &name

	teams, err := c.ListTeams(params)
	if err != nil {
		return diag.FromErr(err)
	}

	tf_teams := make([]interface{}, len(teams), len(teams))
	for i, team := range teams {
		c, _ := team.(*client.Team)
		tf_teams[i] = structToLowerFirstMap(*c)
	}

	if err := d.Set("teams", tf_teams); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}
