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
	"github.com/rootlyhq/terraform-provider-rootly/v2/tools"
)

func resourceIncidentRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIncidentRoleCreate,
		ReadContext: resourceIncidentRoleRead,
		UpdateContext: resourceIncidentRoleUpdate,
		DeleteContext: resourceIncidentRoleDelete,
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
				Description: "The name of the incident role",
				
			},
			

			"slug": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The slug of the incident role",
				
			},
			

			"summary": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The summary of the incident role",
				
			},
			

			"description": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The description of the incident role",
				
			},
			

		"position": &schema.Schema {
			Type: schema.TypeInt,
			Computed: true,
			Required: false,
			Optional: true,
			ForceNew: false,
			Description: "Position of the incident role",
			
		},
		

			"optional": &schema.Schema {
				Type: schema.TypeBool,
				Computed: true,
				Required: false,
				Optional: true,
				Description: "Value must be one of true or false",
				
			},
			

				"enabled": &schema.Schema {
					Type: schema.TypeBool,
					Default: true,
					Optional: true,
					
				},
				

			"allow_multi_user_assignment": &schema.Schema {
				Type: schema.TypeBool,
				Computed: true,
				Required: false,
				Optional: true,
				Description: "Value must be one of true or false",
				
			},
			
		},
	}
}

func resourceIncidentRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating IncidentRole"))

	s := &client.IncidentRole{}

	  if value, ok := d.GetOkExists("name"); ok {
				s.Name = value.(string)
			}
    if value, ok := d.GetOkExists("slug"); ok {
				s.Slug = value.(string)
			}
    if value, ok := d.GetOkExists("summary"); ok {
				s.Summary = value.(string)
			}
    if value, ok := d.GetOkExists("description"); ok {
				s.Description = value.(string)
			}
    if value, ok := d.GetOkExists("position"); ok {
				s.Position = value.(int)
			}
    if value, ok := d.GetOkExists("optional"); ok {
				s.Optional = tools.Bool(value.(bool))
			}
    if value, ok := d.GetOkExists("enabled"); ok {
				s.Enabled = tools.Bool(value.(bool))
			}
    if value, ok := d.GetOkExists("allow_multi_user_assignment"); ok {
				s.AllowMultiUserAssignment = tools.Bool(value.(bool))
			}

	res, err := c.CreateIncidentRole(s)
	if err != nil {
		return diag.Errorf("Error creating incident_role: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a incident_role resource: %s", d.Id()))

	return resourceIncidentRoleRead(ctx, d, meta)
}

func resourceIncidentRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading IncidentRole: %s", d.Id()))

	item, err := c.GetIncidentRole(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if errors.Is(err, client.NewNotFoundError("")) && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("IncidentRole (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading incident_role: %s", d.Id())
	}

	d.Set("name", item.Name)
  d.Set("slug", item.Slug)
  d.Set("summary", item.Summary)
  d.Set("description", item.Description)
  d.Set("position", item.Position)
  d.Set("optional", item.Optional)
  d.Set("enabled", item.Enabled)
  d.Set("allow_multi_user_assignment", item.AllowMultiUserAssignment)

	return nil
}

func resourceIncidentRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating IncidentRole: %s", d.Id()))

	s := &client.IncidentRole{}

	  if d.HasChange("name") {
				s.Name = d.Get("name").(string)
			}
    if d.HasChange("slug") {
				s.Slug = d.Get("slug").(string)
			}
    if d.HasChange("summary") {
				s.Summary = d.Get("summary").(string)
			}
    if d.HasChange("description") {
				s.Description = d.Get("description").(string)
			}
    if d.HasChange("position") {
				s.Position = d.Get("position").(int)
			}
    if d.HasChange("optional") {
				s.Optional = tools.Bool(d.Get("optional").(bool))
			}
    if d.HasChange("enabled") {
				s.Enabled = tools.Bool(d.Get("enabled").(bool))
			}
    if d.HasChange("allow_multi_user_assignment") {
				s.AllowMultiUserAssignment = tools.Bool(d.Get("allow_multi_user_assignment").(bool))
			}

	_, err := c.UpdateIncidentRole(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating incident_role: %s", err.Error())
	}

	return resourceIncidentRoleRead(ctx, d, meta)
}

func resourceIncidentRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting IncidentRole: %s", d.Id()))

	err := c.DeleteIncidentRole(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if errors.Is(err, client.NewNotFoundError("")) && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("IncidentRole (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting incident_role: %s", err.Error())
	}

	d.SetId("")

	return nil
}
