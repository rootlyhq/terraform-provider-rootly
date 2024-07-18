package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	"github.com/rootlyhq/terraform-provider-rootly/v2/tools"
)

func resourcePostmortemTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePostmortemTemplateCreate,
		ReadContext:   resourcePostmortemTemplateRead,
		UpdateContext: resourcePostmortemTemplateUpdate,
		DeleteContext: resourcePostmortemTemplateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "The name of the postmortem template",
			},

			"default": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Default selected template when editing a postmortem. Value must be one of true or false",
			},

			"content": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "The postmortem template. Liquid syntax and markdown are supported",
			},

			"format": &schema.Schema{
				Type:        schema.TypeString,
				Default:     "html",
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The format of the input. Value must be one of `html`, `markdown`.",
			},
		},
	}
}

func resourcePostmortemTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating PostmortemTemplate"))

	s := &client.PostmortemTemplate{}

	if value, ok := d.GetOkExists("name"); ok {
		s.Name = value.(string)
	}
	if value, ok := d.GetOkExists("default"); ok {
		s.Default = tools.Bool(value.(bool))
	}
	if value, ok := d.GetOkExists("content"); ok {
		s.Content = value.(string)
	}
	if value, ok := d.GetOkExists("format"); ok {
		s.Format = value.(string)
	}

	res, err := c.CreatePostmortemTemplate(s)
	if err != nil {
		return diag.Errorf("Error creating post_mortem_template: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a post_mortem_template resource: %s", d.Id()))

	return resourcePostmortemTemplateRead(ctx, d, meta)
}

func resourcePostmortemTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading PostmortemTemplate: %s", d.Id()))

	item, err := c.GetPostmortemTemplate(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("PostmortemTemplate (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading post_mortem_template: %s", d.Id())
	}

	d.Set("name", item.Name)
	d.Set("default", item.Default)
	d.Set("content", item.Content)
	d.Set("format", item.Format)

	return nil
}

func resourcePostmortemTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating PostmortemTemplate: %s", d.Id()))

	s := &client.PostmortemTemplate{}

	if d.HasChange("name") {
		s.Name = d.Get("name").(string)
	}
	if d.HasChange("default") {
		s.Default = tools.Bool(d.Get("default").(bool))
	}
	if d.HasChange("content") {
		s.Content = d.Get("content").(string)
	}
	if d.HasChange("format") {
		s.Format = d.Get("format").(string)
	}

	_, err := c.UpdatePostmortemTemplate(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating post_mortem_template: %s", err.Error())
	}

	return resourcePostmortemTemplateRead(ctx, d, meta)
}

func resourcePostmortemTemplateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting PostmortemTemplate: %s", d.Id()))

	err := c.DeletePostmortemTemplate(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("PostmortemTemplate (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting post_mortem_template: %s", err.Error())
	}

	d.SetId("")

	return nil
}
