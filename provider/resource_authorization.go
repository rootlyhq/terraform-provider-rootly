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

func resourceAuthorization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAuthorizationCreate,
		ReadContext:   resourceAuthorizationRead,
		UpdateContext: resourceAuthorizationUpdate,
		DeleteContext: resourceAuthorizationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{

			"authorizable_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "The id of the resource being accessed.",
			},

			"authorizable_type": &schema.Schema{
				Type:        schema.TypeString,
				Default:     "Dashboard",
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The type of resource being accessed.. Value must be one of `Dashboard`.",
			},

			"grantee_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "The resource id granted access.",
			},

			"grantee_type": &schema.Schema{
				Type:        schema.TypeString,
				Default:     "User",
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The type of resource granted access.. Value must be one of `User`, `Team`.",
			},

			"permissions": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         false,
				Required:         true,
				Optional:         false,
				Description:      "Value must be one of `read`, `update`, `authorize`, `destroy`.",
			},
		},
	}
}

func resourceAuthorizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating Authorization"))

	s := &client.Authorization{}

	if value, ok := d.GetOkExists("authorizable_id"); ok {
		s.AuthorizableId = value.(string)
	}
	if value, ok := d.GetOkExists("authorizable_type"); ok {
		s.AuthorizableType = value.(string)
	}
	if value, ok := d.GetOkExists("grantee_id"); ok {
		s.GranteeId = value.(string)
	}
	if value, ok := d.GetOkExists("grantee_type"); ok {
		s.GranteeType = value.(string)
	}
	if value, ok := d.GetOkExists("permissions"); ok {
		s.Permissions = value.([]interface{})
	}

	res, err := c.CreateAuthorization(s)
	if err != nil {
		return diag.Errorf("Error creating authorization: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a authorization resource: %s", d.Id()))

	return resourceAuthorizationRead(ctx, d, meta)
}

func resourceAuthorizationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading Authorization: %s", d.Id()))

	item, err := c.GetAuthorization(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Authorization (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading authorization: %s", d.Id())
	}

	d.Set("authorizable_id", item.AuthorizableId)
	d.Set("authorizable_type", item.AuthorizableType)
	d.Set("grantee_id", item.GranteeId)
	d.Set("grantee_type", item.GranteeType)
	d.Set("permissions", item.Permissions)

	return nil
}

func resourceAuthorizationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating Authorization: %s", d.Id()))

	s := &client.Authorization{}

	if d.HasChange("authorizable_id") {
		s.AuthorizableId = d.Get("authorizable_id").(string)
	}
	if d.HasChange("authorizable_type") {
		s.AuthorizableType = d.Get("authorizable_type").(string)
	}
	if d.HasChange("grantee_id") {
		s.GranteeId = d.Get("grantee_id").(string)
	}
	if d.HasChange("grantee_type") {
		s.GranteeType = d.Get("grantee_type").(string)
	}
	if d.HasChange("permissions") {
		s.Permissions = d.Get("permissions").([]interface{})
	}

	_, err := c.UpdateAuthorization(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating authorization: %s", err.Error())
	}

	return resourceAuthorizationRead(ctx, d, meta)
}

func resourceAuthorizationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting Authorization: %s", d.Id()))

	err := c.DeleteAuthorization(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Authorization (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting authorization: %s", err.Error())
	}

	d.SetId("")

	return nil
}
