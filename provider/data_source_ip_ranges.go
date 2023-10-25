package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
)

func dataSourceIpRanges() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIpRangesRead,
		Schema: map[string]*schema.Schema{
			"ipv4": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "IPv4 addresses used by rootly.com services",
			},
			"ipv6": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "IPv6 addresses used by rootly.com services",
			},
		},
	}
}

func dataSourceIpRangesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	ip_range, err := c.GetIpRanges()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("ip_ranges")
	d.Set("ipv4", ip_range.Ipv4)
	d.Set("ipv6", ip_range.Ipv6)

	return nil
}
