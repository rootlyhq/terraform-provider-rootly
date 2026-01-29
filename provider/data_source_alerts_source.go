package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

func dataSourceAlertsSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAlertsSourceRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"source_type": &schema.Schema{
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"email", "app_dynamics", "catchpoint", "datadog", "alertmanager", "google_cloud", "grafana", "sentry", "generic_webhook", "cloud_watch", "checkly", "azure", "new_relic", "splunk", "chronosphere", "app_optics", "bug_snag", "honeycomb", "monte_carlo", "nagios", "prtg"}, false),
			},

			"status": &schema.Schema{
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"connected", "setup_complete", "setup_incomplete"}, false),
			},
		},
	}
}

func dataSourceAlertsSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListAlertsSourcesParams)
	page_size := 1
	params.PageSize = &page_size

	// Manually set filters because the API uses plural filter names (filter[source_types],
	// filter[statuses]) but the schema properties are singular (source_type, status).
	// The generator doesn't handle this plural/singular mismatch, so we set filters manually.
	if value, ok := d.GetOkExists("name"); ok {
		name := value.(string)
		params.FilterName = &name
	}

	if value, ok := d.GetOkExists("source_type"); ok {
		source_type := value.(string)
		params.FilterSourceTypes = &source_type
	}

	if value, ok := d.GetOkExists("status"); ok {
		status := value.(string)
		params.FilterStatuses = &status
	}

	items, err := c.ListAlertsSources(params)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(items) == 0 {
		return diag.Errorf("alerts_source not found")
	}
	item, _ := items[0].(*client.AlertsSource)

	d.SetId(item.ID)
	d.Set("name", item.Name)
	d.Set("source_type", item.SourceType)
	d.Set("status", item.Status)

	return nil
}
