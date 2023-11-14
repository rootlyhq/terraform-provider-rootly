package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
)

func resourceFormFieldPosition() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFormFieldPositionCreate,
		ReadContext:   resourceFormFieldPositionRead,
		UpdateContext: resourceFormFieldPositionUpdate,
		DeleteContext: resourceFormFieldPositionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{

			"form_field_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    true,
				Description: "The ID of the form field.",
			},

			"form": &schema.Schema{
				Type:        schema.TypeString,
				Default:     "web_new_incident_form",
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The form for the position. Value must be one of `web_new_incident_form`, `web_update_incident_form`, `web_incident_post_mortem_form`, `web_incident_mitigation_form`, `web_incident_resolution_form`, `web_incident_cancellation_form`, `web_scheduled_incident_form`, `web_update_scheduled_incident_form`, `incident_post_mortem`, `slack_new_incident_form`, `slack_update_incident_form`, `slack_update_incident_status_form`, `slack_incident_mitigation_form`, `slack_incident_resolution_form`, `slack_incident_cancellation_form`, `slack_scheduled_incident_form`, `slack_update_scheduled_incident_form`.",
			},

			"position": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "The position of the form_field_position",
			},
		},
	}
}

func resourceFormFieldPositionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating FormFieldPosition"))

	s := &client.FormFieldPosition{}

	if value, ok := d.GetOkExists("form_field_id"); ok {
		s.FormFieldId = value.(string)
	}
	if value, ok := d.GetOkExists("form"); ok {
		s.Form = value.(string)
	}
	if value, ok := d.GetOkExists("position"); ok {
		s.Position = value.(int)
	}

	res, err := c.CreateFormFieldPosition(s)
	if err != nil {
		return diag.Errorf("Error creating form_field_position: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a form_field_position resource: %s", d.Id()))

	return resourceFormFieldPositionRead(ctx, d, meta)
}

func resourceFormFieldPositionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading FormFieldPosition: %s", d.Id()))

	item, err := c.GetFormFieldPosition(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("FormFieldPosition (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading form_field_position: %s", d.Id())
	}

	d.Set("form_field_id", item.FormFieldId)
	d.Set("form", item.Form)
	d.Set("position", item.Position)

	return nil
}

func resourceFormFieldPositionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating FormFieldPosition: %s", d.Id()))

	s := &client.FormFieldPosition{}

	if d.HasChange("form_field_id") {
		s.FormFieldId = d.Get("form_field_id").(string)
	}
	if d.HasChange("form") {
		s.Form = d.Get("form").(string)
	}
	if d.HasChange("position") {
		s.Position = d.Get("position").(int)
	}

	_, err := c.UpdateFormFieldPosition(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating form_field_position: %s", err.Error())
	}

	return resourceFormFieldPositionRead(ctx, d, meta)
}

func resourceFormFieldPositionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting FormFieldPosition: %s", d.Id()))

	err := c.DeleteFormFieldPosition(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("FormFieldPosition (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting form_field_position: %s", err.Error())
	}

	d.SetId("")

	return nil
}
