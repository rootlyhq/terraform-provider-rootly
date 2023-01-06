package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
	"strconv"
	"time"
)

func dataSourceCauses() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCausesRead,
		Schema: map[string]*schema.Schema{
			"slug": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"causes": &schema.Schema{
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

func dataSourceCausesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListCausesParams)

	slug := d.Get("slug").(string)
	name := d.Get("name").(string)

	params.FilterSlug = &slug
	params.FilterName = &name

	causes, err := c.ListCauses(params)
	if err != nil {
		return diag.FromErr(err)
	}

	tf_causes := make([]interface{}, len(causes), len(causes))
	for i, cause := range causes {
		c, _ := cause.(*client.Cause)
		tf_causes[i] = structToLowerFirstMap(*c)
	}

	if err := d.Set("causes", tf_causes); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}
