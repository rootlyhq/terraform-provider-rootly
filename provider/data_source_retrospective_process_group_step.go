package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

func dataSourceRetrospectiveProcessGroupStep() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRetrospectiveProcessGroupStepRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"retrospective_step_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func dataSourceRetrospectiveProcessGroupStepRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListRetrospectiveProcessGroupStepsParams)
	page_size := 1
	params.PageSize = &page_size

	if value, ok := d.GetOkExists("retrospective_step_id"); ok {
		retrospective_step_id := value.(string)
		params.FilterRetrospectiveStepId = &retrospective_step_id
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

	retrospective_process_group_id := d.Get("retrospective_process_group_id").(string)
	items, err := c.ListRetrospectiveProcessGroupSteps(retrospective_process_group_id, params)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(items) == 0 {
		return diag.Errorf("retrospective_process_group_step not found")
	}
	item, _ := items[0].(*client.RetrospectiveProcessGroupStep)

	d.SetId(item.ID)

	return nil
}
