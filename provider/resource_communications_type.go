// DO NOT MODIFY: This file is generated by tools/generate.js. Any changes will be overwritten during the next build.

package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	
)

func resourceCommunicationsType() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCommunicationsTypeCreate,
		ReadContext: resourceCommunicationsTypeRead,
		UpdateContext: resourceCommunicationsTypeUpdate,
		DeleteContext: resourceCommunicationsTypeDelete,
		Importer: &schema.ResourceImporter {
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema {
			
			"name": &schema.Schema {
				Type: schema.TypeString,
				Computed: false,
				Required: true,
				Optional: false,
				ForceNew: false,
				Description: "The name of the communications type",
				
			},
			

			"slug": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The slug of the communications type",
				
			},
			

			"description": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The description of the communications type",
				
			},
			

			"color": &schema.Schema {
				Type: schema.TypeString,
				Computed: false,
				Required: true,
				Optional: false,
				ForceNew: false,
				Description: "The color of the communications type",
				
			},
			

		"position": &schema.Schema {
			Type: schema.TypeInt,
			Computed: true,
			Required: false,
			Optional: true,
			ForceNew: false,
			Description: "Position of the communications type",
			
		},
		
		},
	}
}

func resourceCommunicationsTypeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating CommunicationsType"))

	s := &client.CommunicationsType{}

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
    if value, ok := d.GetOkExists("position"); ok {
				s.Position = value.(int)
			}

	res, err := c.CreateCommunicationsType(s)
	if err != nil {
		return diag.Errorf("Error creating communications_type: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a communications_type resource: %s", d.Id()))

	return resourceCommunicationsTypeRead(ctx, d, meta)
}

func resourceCommunicationsTypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading CommunicationsType: %s", d.Id()))

	item, err := c.GetCommunicationsType(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if errors.Is(err, client.NewNotFoundError("")) && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("CommunicationsType (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading communications_type: %s", d.Id())
	}

	d.Set("name", item.Name)
  d.Set("slug", item.Slug)
  d.Set("description", item.Description)
  d.Set("color", item.Color)
  d.Set("position", item.Position)

	return nil
}

func resourceCommunicationsTypeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating CommunicationsType: %s", d.Id()))

	s := &client.CommunicationsType{}

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
    if d.HasChange("position") {
				s.Position = d.Get("position").(int)
			}

	_, err := c.UpdateCommunicationsType(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating communications_type: %s", err.Error())
	}

	return resourceCommunicationsTypeRead(ctx, d, meta)
}

func resourceCommunicationsTypeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting CommunicationsType: %s", d.Id()))

	err := c.DeleteCommunicationsType(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if errors.Is(err, client.NewNotFoundError("")) && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("CommunicationsType (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting communications_type: %s", err.Error())
	}

	d.SetId("")

	return nil
}
