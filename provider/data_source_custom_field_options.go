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

func dataSourceCustomFieldOptions() *schema.Resource{
	return &schema.Resource{
		ReadContext: dataSourceCustomFieldOptionsRead,
		Schema: map[string]*schema.Schema{
			"custom_field_id": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
			},
			"value": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
			},
			"color": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
			},
			"custom_field_options": &schema.Schema{
				Type: schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type: schema.TypeString,
							Computed: true,
						},
						"custom_field_id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
						},
						"color": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
						},
						"position": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCustomFieldOptionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListCustomFieldOptionsParams)

	customFieldId := strconv.Itoa(d.Get("custom_field_id").(int))
	value := d.Get("value").(string)
	color := d.Get("color").(string)

	params.FilterValue = &value
	params.FilterColor = &color

	custom_field_options, err := c.ListCustomFieldOptions(customFieldId, params)
	if err != nil {
		return diag.FromErr(err)
	}

	tf_custom_field_options := make([]interface{}, len(custom_field_options), len(custom_field_options))
	for i, custom_field_option := range custom_field_options {
		c, _ := custom_field_option.(*client.CustomFieldOption)
		tf_custom_field_options[i] = structToLowerFirstMap(*c)
	}

	if err := d.Set("custom_field_options", tf_custom_field_options); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}
