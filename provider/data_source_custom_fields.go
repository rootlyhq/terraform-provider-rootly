package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

func dataSourceCustomFields() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCustomFieldsRead,
		Schema: map[string]*schema.Schema{
			"slug": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"label": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"kind": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"custom_fields": &schema.Schema{
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
						"label": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"kind": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"shown": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"required": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceCustomFieldsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListCustomFieldsParams)

	slug := d.Get("slug").(string)
	label := d.Get("label").(string)
	kind := d.Get("kind").(string)
	enabled := d.Get("enabled").(bool)

	params.FilterSlug = &slug
	params.FilterLabel = &label
	params.FilterKind = &kind
	params.FilterEnabled = &enabled

	custom_fields, err := c.ListCustomFields(params)
	if err != nil {
		return diag.FromErr(err)
	}

	tf_custom_fields := make([]interface{}, len(custom_fields), len(custom_fields))
	for i, custom_field := range custom_fields {
		c, _ := custom_field.(*client.CustomField)
		tf_custom_fields[i] = structToLowerFirstMap(*c)
	}

	if err := d.Set("custom_fields", tf_custom_fields); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}
