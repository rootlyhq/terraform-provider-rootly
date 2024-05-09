package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v2/tools"
)

func resourceRetrospectiveConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRetrospectiveConfigurationCreate,
		ReadContext:   resourceRetrospectiveConfigurationRead,
		UpdateContext: resourceRetrospectiveConfigurationUpdate,
		DeleteContext: resourceRetrospectiveConfigurationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"kind": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "The kind of the retrospective configuration.",
			},
			"severity_ids": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Severities associated with the retrospective configuration.",
			},
			"group_ids": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Teams associated with the retrospective configuration.",
			},
			"incident_type_ids": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Incident types associated with the retrospective configuration.",
			},
		},
	}
}

func resourceRetrospectiveConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	item, present := items[0].(*client.RetrospectiveConfiguration)

	if !present {
		return diag.Errorf("Error creating retrospective configuration: %s", err.Error())
	}

	d.SetId(item.ID)

	if value, ok := d.GetOkExists("severity_ids"); ok {
		item.SeverityIds = value.([]interface{})
	}

	if value, ok := d.GetOkExists("group_ids"); ok {
		item.GroupIds = value.([]interface{})
	}

	if value, ok := d.GetOkExists("incident_type_ids"); ok {
		item.IncidentTypeIds = value.([]interface{})
	}

	_, err = c.UpdateRetrospectiveConfiguration(d.Id(), item)
	if err != nil {
		return diag.Errorf("Error creating retrospective configuration: %s", err.Error())
	}

	return resourceRetrospectiveConfigurationRead(ctx, d, meta)
}

func resourceRetrospectiveConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading RetrospectiveConfiguration: %s", d.Id()))

	item, err := c.GetRetrospectiveConfiguration(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("RetrospectiveConfiguration (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading retrospective configuration: %s", d.Id())
	}

	d.Set("kind", item.Kind)
	d.Set("severity_ids", item.SeverityIds)
	d.Set("group_ids", item.GroupIds)
	d.Set("incident_type_ids", item.IncidentTypeIds)

	return nil
}

func resourceRetrospectiveConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating RetrospectiveConfiguration: %s", d.Id()))

	s := &client.RetrospectiveConfiguration{}

	if d.HasChange("kind") {
		s.Kind = d.Get("kind").(string)
	}
	if d.HasChange("severity_ids") {
		s.SeverityIds = d.Get("severity_ids").([]interface{})
	}
	if d.HasChange("group_ids") {
		s.GroupIds = d.Get("group_ids").([]interface{})
	}
	if d.HasChange("incident_type_ids") {
		s.IncidentTypeIds = d.Get("incident_type_ids").([]interface{})
	}

	_, err := c.UpdateRetrospectiveConfiguration(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating retrospective configuration: %s", err.Error())
	}

	return resourceRetrospectiveConfigurationRead(ctx, d, meta)
}

func resourceRetrospectiveConfigurationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("Deleting RetrospectiveConfiguration: %s", d.Id()))

	d.SetId("")

	return nil
}
