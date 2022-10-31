package provider

import (
	"context"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

func dataSourceFormFieldOption() *schema.Resource{
	return &schema.Resource{
		ReadContext: dataSourceFormFieldOptionRead,
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

func dataSourceFormFieldOptionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListFormFieldOptionsParams)
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
			

	form_field_id := d.Get("form_field_id").(string)
			items, err := c.ListFormFieldOptions(form_field_id, params)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(items) == 0 {
		return diag.Errorf("form_field_option not found")
	}
	item, _ := items[0].(*client.FormFieldOption)

	d.SetId(item.ID)

	return nil
}
