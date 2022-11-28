package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	
)

func resourceEnvironment() *schema.Resource{
	return &schema.Resource{
		CreateContext: resourceEnvironmentCreate,
		ReadContext: resourceEnvironmentRead,
		UpdateContext: resourceEnvironmentUpdate,
		DeleteContext: resourceEnvironmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			
			"name": &schema.Schema{
				Type: schema.TypeString,
				Computed: false,
				Required: true,
				Optional: false,
				ForceNew: false,
				Description: "The name of the environment",
			},
			

			"slug": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The slug of the environment",
			},
			

			"description": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The description of the environment",
			},
			

			"color": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "",
			},
			
		},
	}
}


func resourceEnvironmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating Environment"))

	s := &client.Environment{}

	  if value, ok := d.GetOkExists("name"); ok {
				s.Name = value.(string)
			}
    if value, ok := d.GetOkExists("slug"); ok {
				s.Slug = value.(string)
			}
    if value, ok := d.GetOkExists("description"); ok {
				s.Description = value.(string)
			}
    if value, ok := d.GetOkExists("color"); ok {
				s.Color = value.(string)
			}

	res, err := c.CreateEnvironment(s)
	if err != nil {
		return diag.Errorf("Error creating environment: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a environment resource: %s", d.Id()))

	return resourceEnvironmentRead(ctx, d, meta)
}


func resourceEnvironmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading Environment: %s", d.Id()))

	item, err := c.GetEnvironment(d.Id())
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

	d.Set("name", item.Name)
  d.Set("slug", item.Slug)
  d.Set("description", item.Description)
  d.Set("color", item.Color)

	return nil
}


func resourceEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating Environment: %s", d.Id()))

	s := &client.Environment{}

	  if d.HasChange("name") {
				s.Name = d.Get("name").(string)
			}
    if d.HasChange("slug") {
				s.Slug = d.Get("slug").(string)
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

