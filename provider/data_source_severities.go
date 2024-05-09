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

func dataSourceSeverities() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSeveritiesRead,
		Schema: map[string]*schema.Schema{
			"slug": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"severities": &schema.Schema{
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
						"severity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"color": {
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

func dataSourceSeveritiesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListSeveritiesParams)

	slug := d.Get("slug").(string)
	name := d.Get("name").(string)
	severity := d.Get("severity").(string)
	color := d.Get("color").(string)

	params.FilterSlug = &slug
	params.FilterName = &name
	params.FilterSeverity = &severity
	params.FilterColor = &color

	severities, err := c.ListSeverities(params)
	if err != nil {
		return diag.FromErr(err)
	}

	tf_severities := make([]interface{}, len(severities), len(severities))
	for i, severity := range severities {
		x, _ := severity.(*client.Severity)
		tf_severities[i] = structToLowerFirstMap(*x)
	}

	if err := d.Set("severities", tf_severities); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}
