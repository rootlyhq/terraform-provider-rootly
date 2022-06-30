package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	"github.com/rootlyhq/terraform-provider-rootly/tools"
)

func resourceCustomField() *schema.Resource {
	return &schema.Resource{
		Description: "Manages custom fields.",

		CreateContext: resourceCustomFieldCreate,
		ReadContext:   resourceCustomFieldRead,
		UpdateContext: resourceCustomFieldUpdate,
		DeleteContext: resourceCustomFieldDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"label": {
				Description: "The name of the custom field",
				Type:        schema.TypeString,
				Required:    true,
			},
			"kind": {
				Description: "The kind of the custom field",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"text",
					"select",
					"multi_select",
				}, false),
			},
			"description": {
				Description: "The description of the custom field",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"enabled": {
				Description: "Whether the custom field is enabled or not",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"shown": {
				Description: "Where the custom field is shown.",
				Type: schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"incident_form",
						"incident_mitigation_form",
						"incident_resolution_form",
						"incident_post_mortem_form",
						"incident_slack_form",
						"incident_mitigation_slack_form",
						"incident_resolution_slack_form",
						"incident_post_mortem",
					}, false),
				},
			},
			"required": {
				Description: "Where the custom field is required.",
				Type: schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"incident_form",
						"incident_mitigation_form",
						"incident_resolution_form",
						"incident_post_mortem_form",
						"incident_slack_form",
						"incident_mitigation_slack_form",
						"incident_resolution_slack_form",
					}, false),
				},
			},
		},
	}
}

func resourceCustomFieldCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	label := d.Get("label").(string)
	kind := d.Get("kind").(string)

	tflog.Trace(ctx, fmt.Sprintf("Creating custom field: %s", label))

	s := &client.CustomField{
		Label: label,
		Kind: kind,
	}

	if v, ok := d.GetOk("shown"); ok {
		s.Shown = v.([]interface{})
	}

	if v, ok := d.GetOk("required"); ok {
		s.Required = v.([]interface{})
	}

	if value, ok := d.GetOk("description"); ok {
		s.Description = value.(string)
	}

	if v, ok := d.GetOk("enabled"); ok {
		s.Enabled = tools.Bool(v.(bool))
	}

	res, err := c.CreateCustomField(s)
	if err != nil {
		return diag.Errorf("Error creating custom field: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created an custom field resource: %v (%s)", label, d.Id()))

	return resourceCustomFieldRead(ctx, d, meta)
}

func resourceCustomFieldRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading custom field: %s", d.Id()))

	res, err := c.GetCustomField(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("CustomField (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading custom field: %s", d.Id())
	}

	d.Set("label", res.Label)
	d.Set("kind", res.Kind)
	d.Set("description", res.Description)
	d.Set("enabled", res.Enabled)
	d.Set("shown", res.Shown)
	d.Set("required", res.Required)

	return nil
}

func resourceCustomFieldUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating custom field: %s", d.Id()))

	label := d.Get("label").(string)

	s := &client.CustomField{
		Label: label,
	}

	if d.HasChange("kind") {
		s.Kind = d.Get("kind").(string)
	}

	if d.HasChange("description") {
		s.Description = d.Get("description").(string)
	}

	if d.HasChange("enabled") {
		s.Enabled = tools.Bool(d.Get("enabled").(bool))
	}

	if d.HasChange("shown") {
		s.Shown = d.Get("shown").([]interface{})
	}

	if d.HasChange("required") {
		s.Required = d.Get("required").([]interface{})
	}

	tflog.Debug(ctx, fmt.Sprintf("adding value: %#v", s))
	_, err := c.UpdateCustomField(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating custom field: %s", err.Error())
	}

	return resourceCustomFieldRead(ctx, d, meta)
}

func resourceCustomFieldDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting custom field: %s", d.Id()))

	err := c.DeleteCustomField(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("CustomField (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting custom field: %s", err.Error())
	}

	d.SetId("")

	return nil
}
