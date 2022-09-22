package provider

import (
	"context"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

func dataSourceFunctionality() *schema.Resource{
	return &schema.Resource{
		ReadContext: dataSourceFunctionalityRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
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
			

				"created_at": &schema.Schema{
					Type: schema.TypeMap,
					Description: "Filter by date range using 'lt' and 'gt'.",
					Optional: true,
				},
				
		},
	}
}

func dataSourceFunctionalityRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListFunctionalitiesParams)
	page_size := 1
	params.PageSize = &page_size

	
				name := d.Get("name").(string)
				params.FilterName = &name
			

				slug := d.Get("slug").(string)
				params.FilterSlug = &slug
			

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
			

	items, err := c.ListFunctionalities(params)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(items) == 0 {
		return diag.Errorf("functionality not found")
	}
	item, _ := items[0].(*client.Functionality)

	d.SetId(item.ID)

	return nil
}
