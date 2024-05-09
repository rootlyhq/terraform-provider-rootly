package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
)

func resourceDashboard() *schema.Resource {
	return &schema.Resource{
		Description: "Manages dashboards.",

		CreateContext: resourceDashboardCreate,
		ReadContext:   resourceDashboardRead,
		UpdateContext: resourceDashboardUpdate,
		DeleteContext: resourceDashboardDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"user_id": {
				Description: "The user_id of the dashboard, for dashboards with owner == user",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"owner": {
				Description: "Whether the dashboard is owned by organization or user.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "user",
				ValidateFunc: validation.StringInSlice([]string{
					"team",
					"user",
				}, false),
			},
			"name": {
				Description: "The name of the dashboard",
				Type:        schema.TypeString,
				Required:    true,
			},
			"slug": {
				Description: "The slug of the dashboard",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"public": {
				Description: "Whether the dashboard is public",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func resourceDashboardCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	name := d.Get("name").(string)
	owner := d.Get("owner").(string)
	public := d.Get("public").(bool)

	tflog.Trace(ctx, fmt.Sprintf("Creating Dashboard: %s", name))

	s := &client.Dashboard{
		Name:   name,
		Owner:  owner,
		Public: &public,
	}

	if value, ok := d.GetOk("user_id"); ok {
		s.UserId = value.(int)
	}

	res, err := c.CreateDashboard(s)
	if err != nil {
		return diag.Errorf("Error creating dashboard: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a dashboard resource: %v (%s)", name, d.Id()))

	return resourceDashboardRead(ctx, d, meta)
}

func resourceDashboardRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading Dashboard: %s", d.Id()))

	dashboard, err := c.GetDashboard(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Dashboard (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading dashboard: %s", d.Id())
	}

	d.Set("user_id", dashboard.UserId)
	d.Set("name", dashboard.Name)
	d.Set("owner", dashboard.Owner)
	d.Set("public", dashboard.Public)

	return nil
}

func resourceDashboardUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating Dashboard: %s", d.Id()))

	name := d.Get("name").(string)

	s := &client.Dashboard{
		Name: name,
	}

	if d.HasChange("public") {
		public := d.Get("public").(bool)
		s.Public = &public
	}

	if d.HasChange("owner") {
		s.Owner = d.Get("owner").(string)
	}

	_, err := c.UpdateDashboard(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating dashboard: %s", err.Error())
	}

	return resourceDashboardRead(ctx, d, meta)
}

func resourceDashboardDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting Dashboard: %s", d.Id()))

	err := c.DeleteDashboard(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Dashboard (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting dashboard: %s", err.Error())
	}

	d.SetId("")

	return nil
}
