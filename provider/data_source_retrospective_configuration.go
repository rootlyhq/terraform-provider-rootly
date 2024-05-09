package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

func dataSourceRetrospectiveConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRetrospectiveConfigurationRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"kind": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func dataSourceRetrospectiveConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListRetrospectiveConfigurationsParams)
	page_size := 1
	params.PageSize = &page_size

	if value, ok := d.GetOkExists("kind"); ok {
		kind := value.(string)
		params.FilterKind = &kind
	}

	items, err := c.ListRetrospectiveConfigurations(params)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(items) == 0 {
		return diag.Errorf("retrospective_configuration not found")
	}
	item, _ := items[0].(*client.RetrospectiveConfiguration)

	d.SetId(item.ID)

	return nil
}
