package provider

import (
	"context"
	"strconv"
	"time"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

func dataSourceEnvironments() *schema.Resource{
	return &schema.Resource{
		ReadContext: dataSourceEnvironmentsRead,
		Schema: map[string]*schema.Schema{
			"slug": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
			},
			"environments": &schema.Schema{
				Type: schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type: schema.TypeString,
							Computed: true,
						},
						"slug": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
						},
						"color": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceEnvironmentsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListEnvironmentsParams)

	slug := d.Get("slug").(string)
	name := d.Get("name").(string)

	params.FilterSlug = &slug
	params.FilterName = &name

	environments, err := c.ListEnvironments(params)
	if err != nil {
		return diag.FromErr(err)
	}

	tf_environments := make([]interface{}, len(environments), len(environments))
	for i, environment := range environments {
		c, _ := environment.(*client.Environment)
		tf_environments[i] = structToLowerFirstMap(*c)
	}

	if err := d.Set("environments", tf_environments); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}
