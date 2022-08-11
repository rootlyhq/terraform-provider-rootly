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

func dataSourceServices() *schema.Resource{
	return &schema.Resource{
		ReadContext: dataSourceServicesRead,
		Schema: map[string]*schema.Schema{
			"slug": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
			},
			"services": &schema.Schema{
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
						"public_description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceServicesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListServicesParams)

	slug := d.Get("slug").(string)
	name := d.Get("name").(string)

	params.FilterSlug = &slug
	params.FilterName = &name

	services, err := c.ListServices(params)
	if err != nil {
		return diag.FromErr(err)
	}

	tf_services := make([]interface{}, len(services), len(services))
	for i, service := range services {
		c, _ := service.(*client.Service)
		tf_services[i] = structToLowerFirstMap(*c)
	}

	if err := d.Set("services", tf_services); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}
