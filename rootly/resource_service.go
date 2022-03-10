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

func resourceService() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceCreate,
		ReadContext:   resourceServiceRead,
		UpdateContext: resourceServiceUpdate,
		DeleteContext: resourceServiceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the service",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the service",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"color": {
				Description:  "The cikir of the service",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "#047BF8", // Default value from the API
				ValidateFunc: validCSSHexColor(),
			},
			"slug": {
				Description: "The slug of the service",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"service": {
				Description: "The description of the service",
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

func resourceServiceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	name := d.Get("name").(string)

	tflog.Trace(ctx, fmt.Sprintf("Creating Service: %s", name))

	s := &client.Service{
		Name: name,
	}

	if value, ok := d.GetOk("description"); ok {
		s.Description = value.(string)
	}

	if value, ok := d.GetOk("color"); ok {
		s.Color = value.(string)
	}

	res, err := c.CreateService(s)
	if err != nil {
		return diag.Errorf("Error creating service: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a service resource: %v (%s)", name, d.Id()))

	return resourceServiceRead(ctx, d, meta)
}

func resourceServiceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading Service: %s", d.Id()))

	service, err := c.GetService(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Service (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading service: %s", d.Id())
	}

	d.Set("name", service.Name)
	d.Set("description", service.Description)
	d.Set("color", service.Color)
	d.Set("slug", service.Slug)

	return nil
}

func resourceServiceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating Service: %s", d.Id()))

	name := d.Get("name").(string)

	s := &client.Service{
		Name: name,
	}

	if d.HasChange("description") {
		s.Description = d.Get("description").(string)
	}

	if d.HasChange("color") {
		s.Color = d.Get("color").(string)
	}

	_, err := c.UpdateService(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating service: %s", err.Error())
	}

	return resourceServiceRead(ctx, d, meta)
}

func resourceServiceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting Service: %s", d.Id()))

	err := c.DeleteService(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Service (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting service: %s", err.Error())
	}

	d.SetId("")

	return nil
}
