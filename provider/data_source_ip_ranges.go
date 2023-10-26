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
			"integrations_ipv4": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "IPv4 addresses used by Rootly integrations.",
			},
			"integrations_ipv6": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "IPv6 addresses used by Rootly integrations.",
			},
			"webhooks_ipv4": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "IPv4 addresses used by Rootly webhooks.",
			},
			"webhooks_ipv6": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "IPv6 addresses used by Rootly webhooks.",
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
	d.Set("integrations_ipv4", ip_range.IntegrationsIpv4)
	d.Set("integrations_ipv6", ip_range.IntegrationsIpv6)
	d.Set("webhooks_ipv4", ip_range.WebhooksIpv4)
	d.Set("webhooks_ipv6", ip_range.WebhooksIpv6)

	return nil
}
