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

func dataSourceFunctionalities() *schema.Resource{
	return &schema.Resource{
		ReadContext: dataSourceFunctionalitiesRead,
		Schema: map[string]*schema.Schema{
			"slug": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
			},
			"functionalities": &schema.Schema{
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

func dataSourceFunctionalitiesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListFunctionalitiesParams)

	slug := d.Get("slug").(string)
	name := d.Get("name").(string)

	params.FilterSlug = &slug
	params.FilterName = &name

	functionalities, err := c.ListFunctionalities(params)
	if err != nil {
		return diag.FromErr(err)
	}

	tf_functionalities := make([]interface{}, len(functionalities), len(functionalities))
	for i, functionality := range functionalities {
		c, _ := functionality.(*client.Functionality)
		tf_functionalities[i] = structToLowerFirstMap(*c)
	}

	if err := d.Set("functionalities", tf_functionalities); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}
