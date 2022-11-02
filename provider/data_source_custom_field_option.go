package provider

import (
	"context"
	"strconv"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

func dataSourceCustomFieldOption() *schema.Resource{
	return &schema.Resource{
		Description: "DEPRECATED: Please use `form_field` and `form_field_option` data sources instead.",
		DeprecationMessage: "Please use `form_field` and `form_field_option` data sources instead.",
		ReadContext: dataSourceCustomFieldOptionRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			
			"value": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

			"color": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			
		},
	}
}

func dataSourceCustomFieldOptionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListCustomFieldOptionsParams)
	page_size := 1
	params.PageSize = &page_size

	
				if value, ok := d.GetOkExists("value"); ok {
					value := value.(string)
					params.FilterValue = &value
				}
			

				if value, ok := d.GetOkExists("color"); ok {
					color := value.(string)
					params.FilterColor = &color
				}
			

	custom_field_id := strconv.Itoa(d.Get("custom_field_id").(int))
			items, err := c.ListCustomFieldOptions(custom_field_id, params)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(items) == 0 {
		return diag.Errorf("custom_field_option not found")
	}
	item, _ := items[0].(*client.CustomFieldOption)

	d.SetId(item.ID)

	return nil
}
