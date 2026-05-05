package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rootlygo_ "github.com/rootlyhq/rootly-go"
	"github.com/rootlyhq/terraform-provider-rootly/v5/client"
	"github.com/rootlyhq/terraform-provider-rootly/v5/internal/polling"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v5/schema"
)

func dataSourceCatalogEntity() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCatalogEntityRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"catalog_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"created_at": &schema.Schema{
				Type:        schema.TypeMap,
				Description: "Filter by date range using 'lt' and 'gt'.",
				Optional:    true,
			},
		},
	}
}

func dataSourceCatalogEntityRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListCatalogEntitiesParams)
	page_size := 1
	params.PageSize = &page_size

	if value, ok := d.GetOkExists("name"); ok {
		name := value.(string)
		params.FilterName = &name
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

	catalog_id := d.Get("catalog_id").(string)

	items, err := polling.WaitForList(ctx, "catalog_entity", func() ([]interface{}, error) {
		return c.ListCatalogEntities(rootlygo_.ID(catalog_id), params)
	})
	if err != nil {
		return diag.FromErr(err)
	}

	item, _ := items[0].(*client.CatalogEntity)

	d.SetId(item.ID)
	d.Set("name", item.Name)
	return nil
}
