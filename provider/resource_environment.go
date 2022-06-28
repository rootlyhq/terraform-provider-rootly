package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
)

func resourceEnvironment() *schema.Resource {
	return &schema.Resource{
		Description: "Manages incident environments (e.g production, development).",

		CreateContext: resourceEnvironmentCreate,
		ReadContext:   resourceEnvironmentRead,
		UpdateContext: resourceEnvironmentUpdate,
		DeleteContext: resourceEnvironmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the environment",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the environment",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"color": {
				Description:  "The color of the environment",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "#047BF8", // Default value from the API
				ValidateFunc: validCSSHexColor(),
			},
			"slug": {
				Description: "The slug of the environment",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceEnvironmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	name := d.Get("name").(string)

	tflog.Trace(ctx, fmt.Sprintf("Creating Environment: %s", name))

	s := &client.Environment{
		Name: name,
	}

	if value, ok := d.GetOk("description"); ok {
		s.Description = value.(string)
	}

	if value, ok := d.GetOk("color"); ok {
		s.Color = value.(string)
	}

	res, err := c.CreateEnvironment(s)
	if err != nil {
		return diag.Errorf("Error creating environment: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a environment resource: %v (%s)", name, d.Id()))

	return resourceEnvironmentRead(ctx, d, meta)
}

func resourceEnvironmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading Environment: %s", d.Id()))

	environment, err := c.GetEnvironment(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Environment (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading environment: %s", d.Id())
	}

	d.Set("name", environment.Name)
	d.Set("description", environment.Description)
	d.Set("color", environment.Color)
	d.Set("slug", environment.Slug)

	return nil
}

func resourceEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating Environment: %s", d.Id()))

	name := d.Get("name").(string)

	s := &client.Environment{
		Name: name,
	}

	if d.HasChange("description") {
		s.Description = d.Get("description").(string)
	}

	if d.HasChange("color") {
		s.Color = d.Get("color").(string)
	}

	_, err := c.UpdateEnvironment(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating environment: %s", err.Error())
	}

	return resourceEnvironmentRead(ctx, d, meta)
}

func resourceEnvironmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting Environment: %s", d.Id()))

	err := c.DeleteEnvironment(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Environment (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting environment: %s", err.Error())
	}

	d.SetId("")

	return nil
}
