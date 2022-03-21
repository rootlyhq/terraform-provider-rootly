package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootly/terraform-provider-rootly/client"
)

func resourceService() *schema.Resource {
	return &schema.Resource{
		Description: "Manages Services (e.g elasticsearch-prod, redis-preprod, customer-postgresql-prod).",

		CreateContext: resourceServiceCreate,
		ReadContext:   resourceServiceRead,
		UpdateContext: resourceServiceUpdate,
		DeleteContext: resourceServiceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the service",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "For internal use only",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"public_description": {
				Description: "This will be displayed on your status pages to explain to your customer the use of this service.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"color": {
				Description:  "The color chosen for the service",
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

	if value, ok := d.GetOk("public_description"); ok {
		s.PublicDescription = value.(string)
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
	d.Set("public_description", service.PublicDescription)
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

	if d.HasChange("public_description") {
		s.PublicDescription = d.Get("public_description").(string)
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
