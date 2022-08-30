package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
)

func resourceCustomFieldOption() *schema.Resource{
	return &schema.Resource{
		CreateContext: resourceCustomFieldOptionCreate,
		ReadContext: resourceCustomFieldOptionRead,
		UpdateContext: resourceCustomFieldOptionUpdate,
		DeleteContext: resourceCustomFieldOptionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			
			"custom_field_id": &schema.Schema{
				Type: schema.TypeInt,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: true,
				Description: "The ID of the parent custom field",
			},
			

			"value": &schema.Schema{
				Type: schema.TypeString,
				Computed: false,
				Required: true,
				Optional: false,
				ForceNew: false,
				Description: "The value of the custom_field_option",
			},
			

			"color": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The hex color of the custom_field_option",
			},
			

			"position": &schema.Schema{
				Type: schema.TypeInt,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The position of the custom_field_option",
			},
			
		},
	}
}

func resourceCustomFieldOptionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating CustomFieldOption"))

	s := &client.CustomFieldOption{}

	  if value, ok := d.GetOkExists("custom_field_id"); ok {
		s.CustomFieldId = value.(int)
	}
    if value, ok := d.GetOkExists("value"); ok {
		s.Value = value.(string)
	}
    if value, ok := d.GetOkExists("color"); ok {
		s.Color = value.(string)
	}
    if value, ok := d.GetOkExists("position"); ok {
		s.Position = value.(int)
	}

	res, err := c.CreateCustomFieldOption(s)
	if err != nil {
		return diag.Errorf("Error creating custom_field_option: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a custom_field_option resource: %s", d.Id()))

	return resourceCustomFieldOptionRead(ctx, d, meta)
}

func resourceCustomFieldOptionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading CustomFieldOption: %s", d.Id()))

	item, err := c.GetCustomFieldOption(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("CustomFieldOption (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading custom_field_option: %s", d.Id())
	}

	d.Set("custom_field_id", item.CustomFieldId)
  d.Set("value", item.Value)
  d.Set("color", item.Color)
  d.Set("position", item.Position)

	return nil
}

func resourceCustomFieldOptionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating CustomFieldOption: %s", d.Id()))

	s := &client.CustomFieldOption{}

	  if d.HasChange("custom_field_id") {
		s.CustomFieldId = d.Get("custom_field_id").(int)
	}
    if d.HasChange("value") {
		s.Value = d.Get("value").(string)
	}
    if d.HasChange("color") {
		s.Color = d.Get("color").(string)
	}
    if d.HasChange("position") {
		s.Position = d.Get("position").(int)
	}

	_, err := c.UpdateCustomFieldOption(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating custom_field_option: %s", err.Error())
	}

	return resourceCustomFieldOptionRead(ctx, d, meta)
}

func resourceCustomFieldOptionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting CustomFieldOption: %s", d.Id()))

	err := c.DeleteCustomFieldOption(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("CustomFieldOption (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting custom_field_option: %s", err.Error())
	}

	d.SetId("")

	return nil
}
