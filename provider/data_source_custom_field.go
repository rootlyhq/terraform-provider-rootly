package provider

import (
	"context"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

func dataSourceCustomField() *schema.Resource{
	return &schema.Resource{
		ReadContext: dataSourceCustomFieldRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			
			"label": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

			"kind": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

				"enabled": &schema.Schema{
					Type: schema.TypeBool,
					Default: true,
					Optional: true,
				},
				

			"slug": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			
		},
	}
}

func dataSourceCustomFieldRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListCustomFieldsParams)
	page_size := 1
	params.PageSize = &page_size

	
				slug := d.Get("slug").(string)
				params.FilterSlug = &slug
			

				label := d.Get("label").(string)
				params.FilterLabel = &label
			

				kind := d.Get("kind").(string)
				params.FilterKind = &kind
			

				enabled := d.Get("enabled").(bool)
				params.FilterEnabled = &enabled
			

	items, err := c.ListCustomFields(params)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(items) == 0 {
		return diag.Errorf("custom_field not found")
	}
	item, _ := items[0].(*client.CustomField)

	d.SetId(item.ID)

	return nil
}
