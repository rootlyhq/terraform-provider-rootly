package provider

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/rootlyhq/terraform-provider-rootly/v5/client"
	"github.com/rootlyhq/terraform-provider-rootly/v5/internal/apiclient"
	"github.com/rootlyhq/terraform-provider-rootly/v5/internal/jsonapitypes"
)

var _ resource.Resource = &ApiKeyRotationResource{}
var _ resource.ResourceWithValidateConfig = &ApiKeyRotationResource{}

const apiKeyRotationDefaultGracePeriodMinutes = 30

func NewApiKeyRotationResource() resource.Resource {
	return &ApiKeyRotationResource{}
}

type ApiKeyRotationResource struct {
	baseResource
}

func (r *ApiKeyRotationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_key_rotation"
}

func (r *ApiKeyRotationResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Rotates the token of the API key that the provider is authenticated with, and exposes the new token as the sensitive ` + "`token`" + ` attribute.

A rotation happens when the resource is created, and again whenever any argument changes (for example ` + "`triggers`" + `). The new token is only returned once, at rotation time, so store it securely (for example in a secret manager) and then update the provider configuration to use it.

Important notes:

- **The API key must be the one the provider is authenticated with.** The Rootly API only lets a key rotate itself. Rotating any other key returns an error.
- **Use one ` + "`rootly_api_key_rotation`" + ` per provider configuration.** Every instance rotates the same key. If two instances share a provider configuration, the second rotation invalidates the first one's token. Terraform can't detect this at plan time, so keep it to one instance (no ` + "`count`" + ` or ` + "`for_each`" + `) per provider configuration.
- **The provider keeps using the old token for the rest of the apply.** It doesn't switch to the new token mid-run. If ` + "`grace_period_minutes`" + ` is ` + "`0`" + `, or the grace period feature isn't enabled for your organization, the old token stops working immediately, and every later API call in the same run fails. That includes refreshes of other resources. Keep the default grace period, and switch the provider to the new token before it ends.
- **Switch the provider to the new token before the next rotation.** The grace period only covers the token that was just replaced. If the provider is still configured with an older token when the key is rotated again, that older token stops working immediately.
- **Destroying this resource doesn't revoke or change any token.** It only removes the resource from the Terraform state.`,
		Attributes: map[string]schema.Attribute{
			"api_key_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the API key to rotate. This must be the API key the provider is authenticated with.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"triggers": schema.MapAttribute{
				MarkdownDescription: "Arbitrary map of values. Changing any of them rotates the API key again. Use it, for example, with `time_rotating` from the `hashicorp/time` provider to rotate on a schedule.",
				ElementType:         types.StringType,
				Optional:            true,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
			"grace_period_minutes": schema.Int64Attribute{
				MarkdownDescription: fmt.Sprintf("How many minutes the previous token stays valid after rotation, from 0 to 1440. Defaults to `%d`. Only applies when the grace period feature is enabled for your organization. Setting it to `0` invalidates the previous token immediately, which breaks the rest of the current apply because the provider is still authenticated with the previous token. Changing this value rotates the API key again.", apiKeyRotationDefaultGracePeriodMinutes),
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(apiKeyRotationDefaultGracePeriodMinutes),
				Validators: []validator.Int64{
					int64validator.Between(0, 1440),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"expires_at": schema.StringAttribute{
				MarkdownDescription: "The expiration date of the new token, in RFC 3339 format (for example `2027-01-01T00:00:00Z`), up to 5 years in the future. Defaults to 90 days after rotation if omitted. Changing this value rotates the API key again. The actual expiration date is exported as `token_expires_at`.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "The new API key token issued by the rotation.",
				Computed:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"token_expires_at": schema.StringAttribute{
				MarkdownDescription: "The expiration date of the new token.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"grace_period_ends_at": schema.StringAttribute{
				MarkdownDescription: "The time when the previous token stops working. Null if the grace period feature isn't enabled for your organization.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *ApiKeyRotationResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data ApiKeyRotationResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.GracePeriodMinutes.IsNull() && !data.GracePeriodMinutes.IsUnknown() && data.GracePeriodMinutes.ValueInt64() == 0 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("grace_period_minutes"),
			"API key rotation without a grace period",
			"With grace_period_minutes = 0, the previous token is invalidated as soon as the API key is rotated. "+
				"The provider is still authenticated with the previous token for the rest of this apply, so every later API call in this run will fail. "+
				"Use a grace period long enough to switch the provider to the new token.",
		)
	}

	if !data.ExpiresAt.IsNull() && !data.ExpiresAt.IsUnknown() {
		if _, err := time.Parse(time.RFC3339, data.ExpiresAt.ValueString()); err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("expires_at"),
				"Invalid expires_at",
				fmt.Sprintf("expires_at must be an RFC 3339 timestamp (for example 2027-01-01T00:00:00Z): %s", err),
			)
		}
	}
}

func (r *ApiKeyRotationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ApiKeyRotationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.RotateApiKey(ctx, data.ApiKeyId.ValueString(), apiclient.RotateApiKey{
		ExpiresAt:          jsonapitypes.NewNullableFromString(data.ExpiresAt),
		GracePeriodMinutes: jsonapitypes.NewNullableFromInt64(data.GracePeriodMinutes),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to rotate API key", err.Error())
		return
	}

	if item.Token == "" {
		resp.Diagnostics.AddError("Unable to rotate API key", "The API did not return a new token.")
		return
	}

	data.Token = types.StringValue(item.Token)
	data.TokenExpiresAt = jsonapitypes.NullableStringValue(item.ExpiresAt)
	data.GracePeriodEndsAt = jsonapitypes.NullableStringValue(item.GracePeriodEndsAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApiKeyRotationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ApiKeyRotationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The token is only returned at rotation time, so there is nothing to refresh. Only check that the API key still
	// exists, so that a revoked key is removed from state.
	if _, err := r.client.GetApiKey(ctx, data.ApiKeyId.ValueString()); err != nil {
		if errors.Is(err, client.NotFoundError{}) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read API key", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApiKeyRotationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Every argument requires replacement, so an update never needs to call the API.
	var data ApiKeyRotationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApiKeyRotationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Rotation cannot be undone. Removing the resource from state leaves the current token in place.
}

type ApiKeyRotationResourceModel struct {
	ApiKeyId           types.String `tfsdk:"api_key_id"`
	Triggers           types.Map    `tfsdk:"triggers"`
	GracePeriodMinutes types.Int64  `tfsdk:"grace_period_minutes"`
	ExpiresAt          types.String `tfsdk:"expires_at"`
	Token              types.String `tfsdk:"token"`
	TokenExpiresAt     types.String `tfsdk:"token_expires_at"`
	GracePeriodEndsAt  types.String `tfsdk:"grace_period_ends_at"`
}
