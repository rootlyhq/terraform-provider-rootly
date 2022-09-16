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

func resourceCustomField() *schema.Resource{
	return &schema.Resource{
		CreateContext: resourceCustomFieldCreate,
		ReadContext: resourceCustomFieldRead,
		UpdateContext: resourceCustomFieldUpdate,
		DeleteContext: resourceCustomFieldDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			
			"label": &schema.Schema{
				Type: schema.TypeString,
				Computed: false,
				Required: true,
				Optional: false,
				ForceNew: false,
				Description: "The name of the custom_field",
			},
			

			"kind": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The kind of the custom_field",
			},
			

				"enabled": &schema.Schema{
					Type: schema.TypeBool,
					Default: true,
					Optional: true,
				},
				

			"slug": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The slug of the custom_field",
			},
			

			"description": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The description of the custom_field",
			},
			

				"shown": &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Computed: true,
					Required: false,
					Optional: true,
					Description: "",
				},
				

				"required": &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Computed: true,
					Required: false,
					Optional: true,
					Description: "",
				},
				

			"default": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The default value for text field kinds.",
			},
			

			"position": &schema.Schema{
				Type: schema.TypeInt,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The position of the custom_field",
			},
			
		},
	}
}

func resourceCustomFieldCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating CustomField"))

	s := &client.CustomField{}

	  if value, ok := d.GetOkExists("label"); ok {
				s.Label = value.(string)
			}
    if value, ok := d.GetOkExists("kind"); ok {
				s.Kind = value.(string)
			}
    if value, ok := d.GetOkExists("enabled"); ok {
				s.Enabled = tools.Bool(value.(bool))
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
    if value, ok := d.GetOkExists("default"); ok {
				s.Default = value.(string)
			}
    if value, ok := d.GetOkExists("position"); ok {
				s.Position = value.(int)
			}

	res, err := c.CreateCustomField(s)
	if err != nil {
		return diag.Errorf("Error creating custom_field: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a custom_field resource: %s", d.Id()))

	return resourceCustomFieldRead(ctx, d, meta)
}

func resourceCustomFieldRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading CustomField: %s", d.Id()))

	item, err := c.GetCustomField(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("CustomField (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading custom_field: %s", d.Id())
	}

	d.Set("label", item.Label)
  d.Set("kind", item.Kind)
  d.Set("enabled", item.Enabled)
  d.Set("slug", item.Slug)
  d.Set("description", item.Description)
  d.Set("shown", item.Shown)
  d.Set("required", item.Required)
  d.Set("default", item.Default)
  d.Set("position", item.Position)

	return nil
}

func resourceCustomFieldUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating CustomField: %s", d.Id()))

	s := &client.CustomField{}

	  if d.HasChange("label") {
				s.Label = d.Get("label").(string)
			}
    if d.HasChange("kind") {
				s.Kind = d.Get("kind").(string)
			}
    if d.HasChange("enabled") {
				s.Enabled = tools.Bool(d.Get("enabled").(bool))
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
    if d.HasChange("default") {
				s.Default = d.Get("default").(string)
			}
    if d.HasChange("position") {
				s.Position = d.Get("position").(int)
			}

	_, err := c.UpdateCustomField(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating custom_field: %s", err.Error())
	}

	return resourceCustomFieldRead(ctx, d, meta)
}

func resourceCustomFieldDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting CustomField: %s", d.Id()))

	err := c.DeleteCustomField(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("CustomField (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting custom_field: %s", err.Error())
	}

	d.SetId("")

	return nil
}
