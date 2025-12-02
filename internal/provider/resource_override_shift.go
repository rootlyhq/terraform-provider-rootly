package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/rootlyhq/terraform-provider-rootly/v2/internal/provider/resource_override_shift"
	"github.com/rootlyhq/terraform-provider-rootly/v2/internal/rootlytypes"
	rootly "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
	"github.com/sanity-io/litter"
)

var _ resource.Resource = (*overrideShiftResource)(nil)

func NewOverrideShiftResource() resource.Resource {
	return &overrideShiftResource{}
}

type overrideShiftResource struct {
	baseResource
}

func (r *overrideShiftResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_override_shift"
}

func (r *overrideShiftResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_override_shift.OverrideShiftResourceSchema(ctx)
}

func (r *overrideShiftResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_override_shift.OverrideShiftModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	var body rootly.CreateOverrideShiftApplicationVndAPIPlusJSONRequestBody
	body.Data.Attributes.EndsAt, _ = data.EndsAt.ValueRFC3339Time()
	body.Data.Attributes.StartsAt, _ = data.StartsAt.ValueRFC3339Time()
	body.Data.Attributes.UserId = int(data.UserId.ValueInt64())
	httpResp, err := r.client.CreateOverrideShiftWithApplicationVndAPIPlusJSONBodyWithResponse(ctx, data.ScheduleId.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create override shift, got error: %s", err))
		return
	} else if httpResp.StatusCode() != http.StatusCreated {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create override shift, got status code %d, body: %s", httpResp.StatusCode(), string(httpResp.Body)))
		return
	} else if httpResp.ApplicationvndApiJSON201 == nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create override shift, got status code %d, empty body", httpResp.StatusCode()))
		return
	}

	litter.Dump(httpResp.ApplicationvndApiJSON201)
	attributes := httpResp.ApplicationvndApiJSON201.Data.Attributes
	data.EndsAt = rootlytypes.NewRFC3339ValueMust(attributes.EndsAt)
	data.Id = types.StringValue(httpResp.ApplicationvndApiJSON201.Data.Id)
	data.IsOverride = types.BoolValue(attributes.IsOverride)
	data.RotationId = types.StringPointerValue(attributes.RotationId)
	data.ScheduleId = types.StringValue(attributes.ScheduleId)
	// data.ShiftOverride
	data.StartsAt = rootlytypes.NewRFC3339ValueMust(attributes.StartsAt)
	data.UserId = rootlytypes.IntPointerValue(attributes.UserId)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *overrideShiftResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_override_shift.OverrideShiftModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	httpResp, err := r.client.GetOverrideShiftWithResponse(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got error: %s", err))
		return
	} else if httpResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got status code %d, body: %s", httpResp.StatusCode(), string(httpResp.Body)))
		return
	} else if httpResp.ApplicationvndApiJSON200 == nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got status code %d, empty body", httpResp.StatusCode()))
		return
	}

	attributes := httpResp.ApplicationvndApiJSON200.Data.Attributes
	data.EndsAt = rootlytypes.NewRFC3339ValueMust(attributes.EndsAt)
	data.Id = types.StringValue(httpResp.ApplicationvndApiJSON200.Data.Id)
	data.IsOverride = types.BoolValue(attributes.IsOverride)
	data.RotationId = types.StringPointerValue(attributes.RotationId)
	data.ScheduleId = types.StringValue(attributes.ScheduleId)
	// data.ShiftOverride
	data.StartsAt = rootlytypes.NewRFC3339ValueMust(attributes.StartsAt)
	data.UserId = rootlytypes.IntPointerValue(attributes.UserId)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *overrideShiftResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_override_shift.OverrideShiftModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *overrideShiftResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_override_shift.OverrideShiftModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	httpResp, err := r.client.DeleteOverrideShiftWithResponse(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete, got error: %s", err))
	} else if httpResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete, got status code %d, body: %s", httpResp.StatusCode(), string(httpResp.Body)))
	}
}
