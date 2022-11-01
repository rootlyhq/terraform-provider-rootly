package provider

import (
	"context"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

func dataSourceFormField() *schema.Resource{
	return &schema.Resource{
		ReadContext: dataSourceFormFieldRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			
			"kind": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

			"name": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

			"slug": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

				"enabled": &schema.Schema{
					Type: schema.TypeBool,
					Default: true,
					Optional: true,
				},
				

				"created_at": &schema.Schema{
					Type: schema.TypeMap,
					Description: "Filter by date range using 'lt' and 'gt'.",
					Optional: true,
				},
				
		},
	}
}

func dataSourceFormFieldRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListFormFieldsParams)
	page_size := 1
	params.PageSize = &page_size

	
				if value, ok := d.GetOkExists("slug"); ok {
					slug := value.(string)
					params.FilterSlug = &slug
				}
			

				if value, ok := d.GetOkExists("name"); ok {
					name := value.(string)
					params.FilterName = &name
				}
			

				if value, ok := d.GetOkExists("kind"); ok {
					kind := value.(string)
					params.FilterKind = &kind
				}
			

				if value, ok := d.GetOkExists("enabled"); ok {
					enabled := value.(bool)
					params.FilterEnabled = &enabled
				}
			

				created_at_gt := d.Get("created_at").(map[string]interface{})
				if value, exists := created_at_gt["gt"]; exists {
					v := value.(string)
					params.FilterCreatedAtGt = &v
				}
			

				created_at_lt := d.Get("created_at").(map[string]interface{})
				if value, exists := created_at_lt["lt"]; exists {
					v := value.(string)
					params.FilterCreatedAtLt = &v
				}
			

	items, err := c.ListFormFields(params)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(items) == 0 {
		return diag.Errorf("form_field not found")
	}
	item, _ := items[0].(*client.FormField)

	d.SetId(item.ID)

	return nil
}
