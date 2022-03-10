package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rootly/terraform-provider-rootly/client"
)

func resourceSeverity() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSeverityCreate,
		ReadContext:   resourceSeverityRead,
		UpdateContext: resourceSeverityUpdate,
		DeleteContext: resourceSeverityDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the severity",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the severity",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"color": {
				Description:  "The cikir of the severity",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "#047BF8", // Default value from the API
				ValidateFunc: validCSSHexColor(),
			},
			"slug": {
				Description: "The slug of the severity",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"severity": {
				Description: "The description of the severity",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "medium",
				ValidateFunc: validation.StringInSlice([]string{
					"critical",
					"high",
					"medium",
					"low",
				}, false),
			},
		},
	}
}

func resourceSeverityCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	name := d.Get("name").(string)
	severity := d.Get("severity").(string)

	tflog.Trace(ctx, fmt.Sprintf("Creating Severity: %s", name))

	s := &client.Severity{
		Name:     name,
		Severity: severity,
	}

	if value, ok := d.GetOk("description"); ok {
		s.Description = value.(string)
	}

	if value, ok := d.GetOk("color"); ok {
		s.Color = value.(string)
	}

	if value, ok := d.GetOk("severity"); ok {
		s.Severity = value.(string)
	}

	res, err := c.CreateSeverity(s)
	if err != nil {
		return diag.Errorf("Error creating severity: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a severity resource: %v (%s)", name, d.Id()))

	return resourceSeverityRead(ctx, d, meta)
}

func resourceSeverityRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading Severity: %s", d.Id()))

	severity, err := c.GetSeverity(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Severity (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading severity: %s", d.Id())
	}

	d.Set("name", severity.Name)
	d.Set("description", severity.Description)
	d.Set("color", severity.Color)
	d.Set("slug", severity.Slug)
	d.Set("severity", severity.Severity)

	return nil
}

func resourceSeverityUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating Severity: %s", d.Id()))

	name := d.Get("name").(string)
	severity := d.Get("severity").(string)

	s := &client.Severity{
		Name:     name,
		Severity: severity,
	}

	if d.HasChange("description") {
		s.Description = d.Get("description").(string)
	}

	if d.HasChange("color") {
		s.Color = d.Get("color").(string)
	}

	if d.HasChange("severity") {
		s.Severity = d.Get("severity").(string)
	}

	_, err := c.UpdateSeverity(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating severity: %s", err.Error())
	}

	return resourceSeverityRead(ctx, d, meta)
}

func resourceSeverityDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting Severity: %s", d.Id()))

	err := c.DeleteSeverity(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Severity (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting severity: %s", err.Error())
	}

	d.SetId("")

	return nil
}
