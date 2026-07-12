package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v5/client"
)

func dataSourcePostmortemTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePostmortemTemplateRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"content": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"format": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePostmortemTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	slug := d.Get("slug").(string)

	item, err := c.GetPostmortemTemplate(slug)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(item.ID)
	d.Set("name", item.Name)
	d.Set("slug", slug)
	d.Set("default", item.Default)
	d.Set("content", item.Content)
	d.Set("format", item.Format)

	return nil
}
