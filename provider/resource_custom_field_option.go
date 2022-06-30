package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
)

func resourceCustomFieldOption() *schema.Resource {
	return &schema.Resource{
		Description: "Manages custom field options.",

		CreateContext: resourceCustomFieldOptionCreate,
		ReadContext:   resourceCustomFieldOptionRead,
		UpdateContext: resourceCustomFieldOptionUpdate,
		DeleteContext: resourceCustomFieldOptionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"custom_field_id": {
				Description:  "The ID of the parent custom field",
				Type:         schema.TypeInt,
				Required:     true,
			},
			"value": {
				Description: "The value of the custom field option",
				Type:        schema.TypeString,
				Required:    true,
			},
			"color": {
				Description:  "The color of the custom field option",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "#047BF8", // Default value from the API
				ValidateFunc: validCSSHexColor(),
			},
		},
	}
}

func resourceCustomFieldOptionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	value := d.Get("value").(string)
	customFieldId := d.Get("custom_field_id").(int)

	tflog.Trace(ctx, fmt.Sprintf("Creating custom field option: %s", value))

	s := &client.CustomFieldOption{
		Value: value,
		CustomFieldId: customFieldId,
	}

	if value, ok := d.GetOk("color"); ok {
		s.Color = value.(string)
	}

	res, err := c.CreateCustomFieldOption(s)
	if err != nil {
		return diag.Errorf("Error creating custom field option: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created an custom field option resource: %v (%s)", value, d.Id()))

	return resourceCustomFieldOptionRead(ctx, d, meta)
}

func resourceCustomFieldOptionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading custom field option: %s", d.Id()))

	res, err := c.GetCustomFieldOption(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("CustomFieldOption (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading custom field option: %s", d.Id())
	}

	d.Set("custom_field_id", res.CustomFieldId)
	d.Set("value", res.Value)
	d.Set("color", res.Color)

	return nil
}

func resourceCustomFieldOptionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating custom field option: %s", d.Id()))

	value := d.Get("value").(string)

	s := &client.CustomFieldOption{
		Value: value,
	}

	if d.HasChange("color") {
		s.Color = d.Get("color").(string)
	}

	tflog.Debug(ctx, fmt.Sprintf("adding value: %#v", s))
	_, err := c.UpdateCustomFieldOption(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating custom field option: %s", err.Error())
	}

	return resourceCustomFieldOptionRead(ctx, d, meta)
}

func resourceCustomFieldOptionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting custom field option: %s", d.Id()))

	err := c.DeleteCustomFieldOption(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("CustomFieldOption (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting custom field option: %s", err.Error())
	}

	d.SetId("")

	return nil
}
