package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

func dataSourceIncidentPostMortem() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIncidentPostMortemRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"started_at": &schema.Schema{
				Type:        schema.TypeMap,
				Description: "Filter by date range using 'lt' and 'gt'.",
				Optional:    true,
			},

			"mitigated_at": &schema.Schema{
				Type:        schema.TypeMap,
				Description: "Filter by date range using 'lt' and 'gt'.",
				Optional:    true,
			},

			"resolved_at": &schema.Schema{
				Type:        schema.TypeMap,
				Description: "Filter by date range using 'lt' and 'gt'.",
				Optional:    true,
			},

			"created_at": &schema.Schema{
				Type:        schema.TypeMap,
				Description: "Filter by date range using 'lt' and 'gt'.",
				Optional:    true,
			},
		},
	}
}

func dataSourceIncidentPostMortemRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListIncidentPostMortemsParams)
	page_size := 1
	params.PageSize = &page_size

	if value, ok := d.GetOkExists("status"); ok {
		status := value.(string)
		params.FilterStatus = &status
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

	started_at_gt := d.Get("started_at").(map[string]interface{})
	if value, exists := started_at_gt["gt"]; exists {
		v := value.(string)
		params.FilterStartedAtGt = &v
	}

	started_at_lt := d.Get("started_at").(map[string]interface{})
	if value, exists := started_at_lt["lt"]; exists {
		v := value.(string)
		params.FilterStartedAtLt = &v
	}

	mitigated_at_gt := d.Get("mitigated_at").(map[string]interface{})
	if value, exists := mitigated_at_gt["gt"]; exists {
		v := value.(string)
		params.FilterMitigatedAtGt = &v
	}

	mitigated_at_lt := d.Get("mitigated_at").(map[string]interface{})
	if value, exists := mitigated_at_lt["lt"]; exists {
		v := value.(string)
		params.FilterMitigatedAtLt = &v
	}

	resolved_at_gt := d.Get("resolved_at").(map[string]interface{})
	if value, exists := resolved_at_gt["gt"]; exists {
		v := value.(string)
		params.FilterResolvedAtGt = &v
	}

	resolved_at_lt := d.Get("resolved_at").(map[string]interface{})
	if value, exists := resolved_at_lt["lt"]; exists {
		v := value.(string)
		params.FilterResolvedAtLt = &v
	}

	items, err := c.ListIncidentPostMortems(params)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(items) == 0 {
		return diag.Errorf("incident_post_mortem not found")
	}
	item, _ := items[0].(*client.IncidentPostMortem)

	d.SetId(item.ID)

	return nil
}
