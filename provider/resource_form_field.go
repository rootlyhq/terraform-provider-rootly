package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	"github.com/rootlyhq/terraform-provider-rootly/tools"
)

func resourceFormField() *schema.Resource{
	return &schema.Resource{
		CreateContext: resourceFormFieldCreate,
		ReadContext: resourceFormFieldRead,
		UpdateContext: resourceFormFieldUpdate,
		DeleteContext: resourceFormFieldDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			
			"kind": &schema.Schema{
				Type: schema.TypeString,
				Default: "custom",
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The kind of the form field. Value must be one of `custom`, `title`, `summary`, `severity`, `environments`, `types`, `services`, `functionalities`, `teams`, `visibility`, `mark_as_test`, `mark_as_backfilled`, `labels`, `notify_emails`, `trigger_manual_workflows`, `show_ongoing_incidents`, `attach_alerts`, `manual_starting_datetime_field`.",
				
			},
			

			"input_kind": &schema.Schema{
				Type: schema.TypeString,
				Default: "text",
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The input kind of the form field. Value must be one of `text`, `textarea`, `select`, `multi_select`, `date`, `datetime`, `users`, `number`, `checkbox`, `tags`.",
				
			},
			

			"name": &schema.Schema{
				Type: schema.TypeString,
				Computed: false,
				Required: true,
				Optional: false,
				ForceNew: false,
				Description: "The name of the form field",
				
			},
			

			"slug": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The slug of the form field",
				
			},
			

			"description": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The description of the form field",
				
			},
			

				"shown": &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Computed: true,
					Required: false,
					Optional: true,
					Description: "Value must be one of `web_new_incident_form`, `web_update_incident_form`, `web_incident_post_mortem_form`, `web_incident_mitigation_form`, `web_incident_resolution_form`, `web_scheduled_incident_form`, `web_update_scheduled_incident_form`, `incident_post_mortem`, `slack_new_incident_form`, `slack_update_incident_form`, `slack_update_incident_status_form`, `slack_incident_mitigation_form`, `slack_incident_resolution_form`, `slack_scheduled_incident_form`, `slack_update_scheduled_incident_form`.",
				},
				

				"required": &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Computed: true,
					Required: false,
					Optional: true,
					Description: "Value must be one of `web_new_incident_form`, `web_update_incident_form`, `web_incident_post_mortem_form`, `web_incident_mitigation_form`, `web_incident_resolution_form`, `web_scheduled_incident_form`, `web_update_scheduled_incident_form`, `slack_new_incident_form`, `slack_update_incident_form`, `slack_update_incident_status_form`, `slack_incident_mitigation_form`, `slack_incident_resolution_form`, `slack_scheduled_incident_form`, `slack_update_scheduled_incident_form`.",
				},
				

				"enabled": &schema.Schema{
					Type: schema.TypeBool,
					Default: true,
					Optional: true,
					
				},
				

				"default_values": &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Computed: true,
					Required: false,
					Optional: true,
					Description: "",
				},
				
		},
	}
}

func resourceFormFieldCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating FormField"))

	s := &client.FormField{}

	  if value, ok := d.GetOkExists("kind"); ok {
				s.Kind = value.(string)
			}
    if value, ok := d.GetOkExists("input_kind"); ok {
				s.InputKind = value.(string)
			}
    if value, ok := d.GetOkExists("name"); ok {
				s.Name = value.(string)
			}
    if value, ok := d.GetOkExists("slug"); ok {
				s.Slug = value.(string)
			}
    if value, ok := d.GetOkExists("description"); ok {
				s.Description = value.(string)
			}
    if value, ok := d.GetOkExists("shown"); ok {
				s.Shown = value.([]interface{})
			}
    if value, ok := d.GetOkExists("required"); ok {
				s.Required = value.([]interface{})
			}
    if value, ok := d.GetOkExists("enabled"); ok {
				s.Enabled = tools.Bool(value.(bool))
			}
    if value, ok := d.GetOkExists("default_values"); ok {
				s.DefaultValues = value.([]interface{})
			}

	res, err := c.CreateFormField(s)
	if err != nil {
		return diag.Errorf("Error creating form_field: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a form_field resource: %s", d.Id()))

	return resourceFormFieldRead(ctx, d, meta)
}

func resourceFormFieldRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading FormField: %s", d.Id()))

	item, err := c.GetFormField(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("FormField (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading form_field: %s", d.Id())
	}

	d.Set("kind", item.Kind)
  d.Set("input_kind", item.InputKind)
  d.Set("name", item.Name)
  d.Set("slug", item.Slug)
  d.Set("description", item.Description)
  d.Set("shown", item.Shown)
  d.Set("required", item.Required)
  d.Set("enabled", item.Enabled)
  d.Set("default_values", item.DefaultValues)

	return nil
}

func resourceFormFieldUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating FormField: %s", d.Id()))

	s := &client.FormField{}

	  if d.HasChange("kind") {
				s.Kind = d.Get("kind").(string)
			}
    if d.HasChange("input_kind") {
				s.InputKind = d.Get("input_kind").(string)
			}
    if d.HasChange("name") {
				s.Name = d.Get("name").(string)
			}
    if d.HasChange("slug") {
				s.Slug = d.Get("slug").(string)
			}
    if d.HasChange("description") {
				s.Description = d.Get("description").(string)
			}
    if d.HasChange("shown") {
				s.Shown = d.Get("shown").([]interface{})
			}
    if d.HasChange("required") {
				s.Required = d.Get("required").([]interface{})
			}
    if d.HasChange("enabled") {
				s.Enabled = tools.Bool(d.Get("enabled").(bool))
			}
    if d.HasChange("default_values") {
				s.DefaultValues = d.Get("default_values").([]interface{})
			}

	_, err := c.UpdateFormField(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating form_field: %s", err.Error())
	}

	return resourceFormFieldRead(ctx, d, meta)
}

func resourceFormFieldDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting FormField: %s", d.Id()))

	err := c.DeleteFormField(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("FormField (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting form_field: %s", err.Error())
	}

	d.SetId("")

	return nil
}
