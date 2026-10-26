// Hand-maintained: api_key is excluded from codegen (see tools/generate.js) because
// it exposes the create/rotate-time token and a rotation_trigger value. Changing
// rotation_trigger calls the rotate endpoint, which can only rotate the key used to
// authenticate the request, so rotation only works when the provider's api_token
// is the managed key itself.

package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rootlyhq/terraform-provider-rootly/v5/client"
)

func resourceApiKey() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages a Rootly API key. The `token` is only returned by the API when the key is created or rotated, so it is not available after import until the key is rotated. To rotate a key, change `rotation_trigger`: the provider calls the rotate endpoint and stores the new token. The API only allows a key to rotate itself, so the provider's `api_token` must be this key's token, and you must update the provider's token after rotating.",
		CreateContext: resourceApiKeyCreate,
		ReadContext:   resourceApiKeyRead,
		UpdateContext: resourceApiKeyUpdate,
		DeleteContext: resourceApiKeyDelete,
		CustomizeDiff: resourceApiKeyCustomizeDiff,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				Sensitive:   false,
				ForceNew:    false,
				WriteOnly:   false,
				Description: "The name of the API key",
			},

			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
				ForceNew:    false,
				WriteOnly:   false,
				Description: "A description of the API key",
			},

			"kind": &schema.Schema{
				Type:         schema.TypeString,
				Default:      "personal",
				Required:     false,
				Optional:     true,
				Sensitive:    false,
				ForceNew:     false,
				WriteOnly:    false,
				Description:  "The kind of the API key. Value must be one of `personal`, `team`, `organization`.",
				ValidateFunc: validation.StringInSlice([]string{"personal", "team", "organization"}, false),
			},

			"group_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
				ForceNew:    true,
				WriteOnly:   false,
				Description: "The group (team) ID for team API keys. Keys created via the API have exactly one group; for keys with multiple groups this is the first group.",
			},

			"role_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
				ForceNew:    false,
				WriteOnly:   false,
				Description: "The role ID",
			},

			"on_call_role_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
				ForceNew:    false,
				WriteOnly:   false,
				Description: "The on-call role ID",
			},

			"expires_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				Sensitive:   false,
				ForceNew:    false,
				WriteOnly:   false,
				Description: "Expiration date",
			},

			"last_used_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
				ForceNew:    false,
				WriteOnly:   false,
				Description: "Date of last use",
			},

			"grace_period_ends_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
				ForceNew:    false,
				WriteOnly:   false,
				Description: "Grace period end date",
			},

			"token": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The API key token. Only returned when the key is created or rotated, so it is empty for imported keys until they are rotated. Stored in state as a sensitive value.",
			},

			"rotation_trigger": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Arbitrary value that, when changed, rotates the API key's token. Not sent to the API.",
			},

			"rotation_grace_period_minutes": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 1440),
				Description:  "How many minutes the old token stays valid after rotation, if the grace period is enabled for your organization. Set to 0 to invalidate the old token immediately. The API defaults to 30.",
			},
		},
	}
}

func resourceApiKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating ApiKey"))

	s := &client.ApiKey{}

	if value, ok := d.GetOkExists("name"); ok {
		s.Name = value.(string)
	}
	if value, ok := d.GetOkExists("description"); ok {
		s.Description = value.(string)
	}
	if value, ok := d.GetOkExists("kind"); ok {
		s.Kind = value.(string)
	}
	if value, ok := d.GetOkExists("group_id"); ok {
		s.GroupId = value.(string)
	}
	if value, ok := d.GetOkExists("role_id"); ok {
		s.RoleId = value.(string)
	}
	if value, ok := d.GetOkExists("on_call_role_id"); ok {
		s.OnCallRoleId = value.(string)
	}
	if value, ok := d.GetOkExists("expires_at"); ok {
		s.ExpiresAt = value.(string)
	}
	if value, ok := d.GetOkExists("last_used_at"); ok {
		s.LastUsedAt = value.(string)
	}
	if value, ok := d.GetOkExists("grace_period_ends_at"); ok {
		s.GracePeriodEndsAt = value.(string)
	}

	res, err := c.CreateApiKey(s)
	if err != nil {
		return diag.Errorf("Error creating api_key: %s", err.Error())
	}

	d.SetId(res.ID)
	d.Set("token", res.Token)
	tflog.Trace(ctx, fmt.Sprintf("created a api_key resource: %s", d.Id()))

	return resourceApiKeyRead(ctx, d, meta)
}

func resourceApiKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading ApiKey: %s", d.Id()))

	item, err := c.GetApiKey(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if errors.Is(err, client.NewNotFoundError("")) && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("ApiKey (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading api_key: %s", d.Id())
	}

	setApiKeyAttributes(d, item)

	return nil
}

func setApiKeyAttributes(d *schema.ResourceData, item *client.ApiKey) {
	d.Set("name", item.Name)
	d.Set("description", item.Description)
	d.Set("kind", item.Kind)
	d.Set("group_id", item.GroupId)
	d.Set("role_id", item.RoleId)
	d.Set("on_call_role_id", item.OnCallRoleId)
	d.Set("expires_at", item.ExpiresAt)
	d.Set("last_used_at", item.LastUsedAt)
	d.Set("grace_period_ends_at", item.GracePeriodEndsAt)
}

func resourceApiKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating ApiKey: %s", d.Id()))

	s := &client.ApiKey{}

	if d.HasChange("name") {
		s.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		s.Description = d.Get("description").(string)
	}
	if d.HasChange("kind") {
		s.Kind = d.Get("kind").(string)
	}
	if d.HasChange("group_id") {
		s.GroupId = d.Get("group_id").(string)
	}
	if d.HasChange("role_id") {
		s.RoleId = d.Get("role_id").(string)
	}
	if d.HasChange("on_call_role_id") {
		s.OnCallRoleId = d.Get("on_call_role_id").(string)
	}
	if d.HasChange("expires_at") {
		s.ExpiresAt = d.Get("expires_at").(string)
	}
	if d.HasChange("last_used_at") {
		s.LastUsedAt = d.Get("last_used_at").(string)
	}
	if d.HasChange("grace_period_ends_at") {
		s.GracePeriodEndsAt = d.Get("grace_period_ends_at").(string)
	}

	if d.HasChangesExcept("rotation_trigger", "rotation_grace_period_minutes", "token", "grace_period_ends_at") {
		_, err := c.UpdateApiKey(d.Id(), s)
		if err != nil {
			return diag.Errorf("Error updating api_key: %s", err.Error())
		}
	}

	if d.HasChange("rotation_trigger") {
		tflog.Trace(ctx, fmt.Sprintf("Rotating ApiKey: %s", d.Id()))

		// Pass expires_at so the API doesn't reset it to its 90 day default.
		r := &client.ApiKeyRotation{ExpiresAt: d.Get("expires_at").(string)}
		if value, ok := d.GetOkExists("rotation_grace_period_minutes"); ok {
			minutes := value.(int)
			r.GracePeriodMinutes = &minutes
		}

		res, err := c.RotateApiKey(d.Id(), r)
		if err != nil {
			return diag.Errorf("Error rotating api_key: %s", err.Error())
		}

		// The provider authenticated with the old token, which may already be invalid,
		// so take state from the rotate response instead of reading the key back.
		d.Set("token", res.Token)
		setApiKeyAttributes(d, res)
		return nil
	}

	return resourceApiKeyRead(ctx, d, meta)
}

// The token is only known after rotating, so mark it unknown in the plan.
func resourceApiKeyCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	if d.Id() != "" && d.HasChange("rotation_trigger") {
		if err := d.SetNewComputed("token"); err != nil {
			return err
		}
		return d.SetNewComputed("grace_period_ends_at")
	}
	return nil
}

func resourceApiKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting ApiKey: %s", d.Id()))

	err := c.DeleteApiKey(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if errors.Is(err, client.NewNotFoundError("")) && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("ApiKey (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting api_key: %s", err.Error())
	}

	d.SetId("")

	return nil
}
