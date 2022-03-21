package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootly/terraform-provider-rootly/client"
)

func resourceFunctionality() *schema.Resource {
	return &schema.Resource{
		Description: "Manages functionalities (e.g Logging In, Search, Adds items to Cart).",

		CreateContext: resourceFunctionalityCreate,
		ReadContext:   resourceFunctionalityRead,
		UpdateContext: resourceFunctionalityUpdate,
		DeleteContext: resourceFunctionalityDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the functionality",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the functionality",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"color": {
				Description:  "The color of the severity",
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
		},
	}
}

func resourceFunctionalityCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	name := d.Get("name").(string)

	tflog.Trace(ctx, fmt.Sprintf("Creating Functionality: %s", name))

	s := &client.Functionality{
		Name: name,
	}

	if value, ok := d.GetOk("description"); ok {
		s.Description = value.(string)
	}

	if value, ok := d.GetOk("color"); ok {
		s.Color = value.(string)
	}

	res, err := c.CreateFunctionality(s)
	if err != nil {
		return diag.Errorf("Error creating functionality: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a functionality resource: %v (%s)", name, d.Id()))

	return resourceFunctionalityRead(ctx, d, meta)
}

func resourceFunctionalityRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading Functionality: %s", d.Id()))

	functionality, err := c.GetFunctionality(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Functionality (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading functionality: %s", d.Id())
	}

	d.Set("name", functionality.Name)
	d.Set("description", functionality.Description)
	d.Set("color", functionality.Color)
	d.Set("slug", functionality.Slug)

	return nil
}

func resourceFunctionalityUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating Functionality: %s", d.Id()))

	name := d.Get("name").(string)

	s := &client.Functionality{
		Name: name,
	}

	if d.HasChange("description") {
		s.Description = d.Get("description").(string)
	}

	if d.HasChange("color") {
		s.Color = d.Get("color").(string)
	}

	_, err := c.UpdateFunctionality(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating functionality: %s", err.Error())
	}

	return resourceFunctionalityRead(ctx, d, meta)
}

func resourceFunctionalityDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting Functionality: %s", d.Id()))

	err := c.DeleteFunctionality(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Functionality (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting functionality: %s", err.Error())
	}

	d.SetId("")

	return nil
}
