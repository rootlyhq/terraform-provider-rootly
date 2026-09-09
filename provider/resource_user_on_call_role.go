// Hand-maintained: `user` is excluded from codegen (see tools/generate.js).

package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v5/client"
)

func resourceUserOnCallRole() *schema.Resource {
	return &schema.Resource{
		Description: "Manages the on-call role assigned to a Rootly user. The resource id is the user id. " +
			"Destroying it stops managing the assignment but leaves the user's on-call role unchanged.",
		CreateContext: resourceUserOnCallRoleCreate,
		ReadContext:   resourceUserOnCallRoleRead,
		UpdateContext: resourceUserOnCallRoleUpdate,
		DeleteContext: resourceUserOnCallRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"user_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the user to assign the on-call role to.",
			},

			"on_call_role_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the on-call role to assign.",
			},
		},
	}
}

func resourceUserOnCallRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	userID := d.Get("user_id").(string)
	tflog.Trace(ctx, fmt.Sprintf("Assigning on-call role to user: %s", userID))

	// Note: the API returns 404 (not 422) when on_call_role_id is invalid or
	// belongs to another team.
	if _, err := c.UpdateUserOnCallRole(userID, d.Get("on_call_role_id").(string)); err != nil {
		return diag.Errorf("Error assigning on-call role to user %s: %s", userID, err.Error())
	}

	d.SetId(userID)

	return resourceUserOnCallRoleRead(ctx, d, meta)
}

func resourceUserOnCallRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading on-call role for user: %s", d.Id()))

	user, err := c.GetUser(d.Id())
	if err != nil {
		if errors.Is(err, client.NewNotFoundError("")) && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("User (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading user %s: %s", d.Id(), err.Error())
	}

	d.Set("user_id", user.ID)
	d.Set("on_call_role_id", user.OnCallRoleId())

	return nil
}

func resourceUserOnCallRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating on-call role for user: %s", d.Id()))

	if _, err := c.UpdateUserOnCallRole(d.Id(), d.Get("on_call_role_id").(string)); err != nil {
		return diag.Errorf("Error updating on-call role for user %s: %s", d.Id(), err.Error())
	}

	return resourceUserOnCallRoleRead(ctx, d, meta)
}

// Delete is intentionally a no-op: we only stop managing the assignment and
// leave the user's current on-call role in place.
func resourceUserOnCallRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("Unmanaging on-call role for user: %s (no-op)", d.Id()))
	d.SetId("")
	return nil
}
