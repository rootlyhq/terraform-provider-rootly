package provider

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/rootlyhq/terraform-provider-rootly/v5/client"
	"github.com/rootlyhq/terraform-provider-rootly/v5/internal/apiclient"
	"github.com/rootlyhq/terraform-provider-rootly/v5/internal/diagutils"
	"github.com/rootlyhq/terraform-provider-rootly/v5/internal/jsonapitypes"
)

var _ resource.Resource = &UserIncidentResponseRoleResource{}
var _ resource.ResourceWithImportState = &UserIncidentResponseRoleResource{}

func NewUserIncidentResponseRoleResource() resource.Resource {
	return &UserIncidentResponseRoleResource{}
}

type UserIncidentResponseRoleResource struct {
	baseResource
}

func (r *UserIncidentResponseRoleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_incident_response_role"
}

func (r *UserIncidentResponseRoleResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a user's incident response role assignment.",
		Attributes: map[string]schema.Attribute{
			"user_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the user.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the incident response role to assign to the user.",
				Required:            true,
			},
		},
	}
}

func (r *UserIncidentResponseRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data UserIncidentResponseRoleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item := diagutils.MergeDiagnostics(data.ToApi(ctx))(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.UpdateUser(ctx, *item)
	if err != nil {
		resp.Diagnostics.AddError("Unable to assign incident response role to user", err.Error())
		return
	}

	resp.Diagnostics.Append(data.FromApi(ctx, *item)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserIncidentResponseRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data UserIncidentResponseRoleResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.GetUser(ctx, data.UserId.ValueString())
	if err != nil {
		if errors.Is(err, client.NotFoundError{}) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read user", err.Error())
		return
	}

	resp.Diagnostics.Append(data.FromApi(ctx, *item)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserIncidentResponseRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data UserIncidentResponseRoleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item := diagutils.MergeDiagnostics(data.ToApi(ctx))(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.UpdateUser(ctx, *item)
	if err != nil {
		resp.Diagnostics.AddError("Unable to update user incident response role", err.Error())
		return
	}

	resp.Diagnostics.Append(data.FromApi(ctx, *item)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserIncidentResponseRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning("Deleting user incident response role has no effect", "Please remove the user incident response role from the Terraform configuration instead.")
}

func (r *UserIncidentResponseRoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("user_id"), req, resp)
}

type UserIncidentResponseRoleResourceModel struct {
	UserId types.String `tfsdk:"user_id"`
	RoleId types.String `tfsdk:"role_id"`
}

func (m *UserIncidentResponseRoleResourceModel) FromApi(ctx context.Context, data apiclient.User) (diags diag.Diagnostics) {
	m.UserId = types.StringValue(data.ID)

	if data.Role == nil {
		diags.AddError("Unable to update user incident response role", "User role not found")
		return
	}
	m.RoleId = types.StringValue(data.Role.ID)

	return
}

func (m *UserIncidentResponseRoleResourceModel) ToApi(ctx context.Context) (*apiclient.User, diag.Diagnostics) {
	var data apiclient.User

	data.ID = m.UserId.ValueString()
	data.RoleId = jsonapitypes.NewNullableFromString(m.RoleId)

	return &data, nil
}
