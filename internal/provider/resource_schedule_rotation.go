package provider

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	"github.com/rootlyhq/terraform-provider-rootly/v5/client"
	"github.com/rootlyhq/terraform-provider-rootly/v5/internal/apiclient"
	"github.com/rootlyhq/terraform-provider-rootly/v5/internal/diagutils"
	"github.com/rootlyhq/terraform-provider-rootly/v5/internal/jsonapitypes"
	"github.com/samber/lo"
)

var _ resource.Resource = &ScheduleRotationResource{}
var _ resource.ResourceWithImportState = &ScheduleRotationResource{}

func NewScheduleRotationResource() resource.Resource {
	return &ScheduleRotationResource{}
}

type ScheduleRotationResource struct {
	baseResource
}

func (r *ScheduleRotationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schedule_rotation"
}

func (r *ScheduleRotationResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a schedule rotation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of this resource.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"schedule_id": schema.StringAttribute{
				MarkdownDescription: "The ID of parent schedule.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the schedule rotation.",
				Required:            true,
			},
			"position": schema.Int64Attribute{
				MarkdownDescription: "Position of the schedule rotation.",
				Optional:            true,
				Computed:            true,
			},
			"schedule_rotationable_type": schema.StringAttribute{
				MarkdownDescription: "Schedule rotation type. Value must be one of `ScheduleDailyRotation`, `ScheduleWeeklyRotation`, `ScheduleBiweeklyRotation`, `ScheduleMonthlyRotation`, `ScheduleCustomRotation`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("ScheduleDailyRotation"),
				Validators: []validator.String{
					stringvalidator.OneOf("ScheduleDailyRotation", "ScheduleWeeklyRotation", "ScheduleBiweeklyRotation", "ScheduleMonthlyRotation", "ScheduleCustomRotation"),
				},
			},
			"active_all_week": schema.BoolAttribute{
				MarkdownDescription: "Schedule rotation active all week?",
				Optional:            true,
				Computed:            true,
			},
			"active_days": schema.SetAttribute{
				CustomType:          supertypes.NewSetTypeOf[string](ctx),
				ElementType:         types.StringType,
				MarkdownDescription: "Value must be one of `S`, `M`, `T`, `W`, `R`, `F`, `U`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(stringvalidator.OneOf("S", "M", "T", "W", "R", "F", "U")),
				},
			},
			"active_time_type": schema.StringAttribute{
				MarkdownDescription: "Value must be one of `all_day`, `same_time`, or `custom`. The value chosen will override `active_time_attributes` in any `rootly_schedule_rotation_active_day` resources linked to this `rootly_schedule_rotation`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("all_day", "same_time", "custom"),
				},
			},
			"time_zone": schema.StringAttribute{
				MarkdownDescription: "A valid IANA time zone name.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Etc/UTC"),
			},
			"start_time": schema.StringAttribute{
				MarkdownDescription: "ISO8601 date and time when rotation starts. Shifts will only be created after this time.",
				Optional:            true,
				Computed:            true,
			},
			"end_time": schema.StringAttribute{
				MarkdownDescription: "ISO8601 date and time when rotation ends. Shifts will only be created before this time.",
				Optional:            true,
				Computed:            true,
			},
			"schedule_rotationable_attributes": schema.SingleNestedAttribute{
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[ScheduleRotationResourceScheduleRotationAttributesModel](ctx),
				MarkdownDescription: "handoff_time and/or handoff_day may be required, depending on schedule_rotationable_type. Please see API docs for options based on schedule_rotationable_type: https://docs.rootly.com/api-reference/schedulerotations/creates-a-schedule-rotation#response-data-attributes-schedule-rotationable-attributes",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"handoff_time": schema.StringAttribute{
						MarkdownDescription: "Hand off time. Only applicable for daily, weekly/biweekly, monthly, and custom rotations.",
						Required:            true,
					},
					"handoff_day": schema.StringAttribute{
						MarkdownDescription: "Hand off day. Only applicable for weekly/biweekly, and monthly.",
						Optional:            true,
					},
					"shift_length": schema.Int64Attribute{
						MarkdownDescription: "Shift length for custom rotation.",
						Optional:            true,
					},
					"shift_length_unit": schema.StringAttribute{
						MarkdownDescription: "Shift length unit for custom rotation. Value must be one of `hours`, `days`, `weeks`.",
						Optional:            true,
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"active_time_attributes": schema.SetNestedBlock{
				CustomType:          supertypes.NewSetNestedObjectTypeOf[ScheduleRotationResourceActiveTimeModel](ctx),
				MarkdownDescription: "Schedule rotation's active times.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"start_time": schema.StringAttribute{
							MarkdownDescription: "Start time for schedule rotation active time.",
							Required:            true,
						},
						"end_time": schema.StringAttribute{
							MarkdownDescription: "End time for schedule rotation active time.",
							Required:            true,
						},
					},
				},
			},
			"schedule_rotation_members": schema.SetNestedBlock{
				CustomType:          supertypes.NewSetNestedObjectTypeOf[ScheduleRotationResourceScheduleRotationMemberModel](ctx),
				MarkdownDescription: "Schedule rotation members. You can only add schedule rotation members if your account has schedule nesting feature enabled.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"member_id": schema.StringAttribute{
							MarkdownDescription: "ID of the member.",
							Required:            true,
						},
						"member_type": schema.StringAttribute{
							MarkdownDescription: "Type of member. Value must be one of `Schedule` or `User`.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("Schedule", "User"),
							},
						},
						"position": schema.Int64Attribute{
							MarkdownDescription: "Position of the member in rotation",
							Optional:            true,
							Computed:            true,
							Validators: []validator.Int64{
								int64validator.AtLeast(1),
							},
						},
					},
				},
			},
		},
	}
}

func (r *ScheduleRotationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ScheduleRotationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item := diagutils.MergeDiagnostics(data.ToApi(ctx))(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.CreateScheduleRotation(ctx, *item)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create schedule rotation", err.Error())
		return
	}

	data.Id = types.StringValue(res.ID)

	item, err = r.client.GetScheduleRotation(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read schedule rotation", err.Error())
		return
	}

	resp.Diagnostics.Append(data.FromApi(ctx, *item)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ScheduleRotationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ScheduleRotationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.GetScheduleRotation(ctx, data.Id.ValueString())
	if err != nil {
		if errors.Is(err, client.NotFoundError{}) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read schedule rotation", err.Error())
		return
	}

	resp.Diagnostics.Append(data.FromApi(ctx, *item)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ScheduleRotationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ScheduleRotationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item := diagutils.MergeDiagnostics(data.ToApi(ctx))(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateScheduleRotation(ctx, *item)
	if err != nil {
		resp.Diagnostics.AddError("Unable to update schedule rotation", err.Error())
		return
	}

	item, err = r.client.GetScheduleRotation(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read schedule rotation", err.Error())
		return
	}

	resp.Diagnostics.Append(data.FromApi(ctx, *item)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ScheduleRotationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ScheduleRotationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.legacyClient.DeleteScheduleRotation(data.Id.ValueString())
	if err != nil {
		if errors.Is(err, client.NotFoundError{}) {
			return
		}
		resp.Diagnostics.AddError("Unable to delete schedule rotation", err.Error())
	}
}

func (r *ScheduleRotationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type ScheduleRotationResourceModel struct {
	Id                             types.String                                                                                  `tfsdk:"id"`
	ScheduleId                     types.String                                                                                  `tfsdk:"schedule_id"`
	Name                           types.String                                                                                  `tfsdk:"name"`
	Position                       types.Int64                                                                                   `tfsdk:"position"`
	ScheduleRotationableType       types.String                                                                                  `tfsdk:"schedule_rotationable_type"`
	ActiveAllWeek                  types.Bool                                                                                    `tfsdk:"active_all_week"`
	ActiveDays                     supertypes.SetValueOf[string]                                                                 `tfsdk:"active_days"`
	ActiveTimeType                 types.String                                                                                  `tfsdk:"active_time_type"`
	TimeZone                       types.String                                                                                  `tfsdk:"time_zone"`
	StartTime                      types.String                                                                                  `tfsdk:"start_time"`
	EndTime                        types.String                                                                                  `tfsdk:"end_time"`
	ScheduleRotationableAttributes supertypes.SingleNestedObjectValueOf[ScheduleRotationResourceScheduleRotationAttributesModel] `tfsdk:"schedule_rotationable_attributes"`
	ActiveTimeAttributes           supertypes.SetNestedObjectValueOf[ScheduleRotationResourceActiveTimeModel]                    `tfsdk:"active_time_attributes"`
	ScheduleRotationMembers        supertypes.SetNestedObjectValueOf[ScheduleRotationResourceScheduleRotationMemberModel]        `tfsdk:"schedule_rotation_members"`
}

func (m *ScheduleRotationResourceModel) FromApi(ctx context.Context, data apiclient.ScheduleRotation) diag.Diagnostics {
	m.Id = types.StringValue(data.ID)
	m.ScheduleId = types.StringValue(data.ScheduleId)
	m.Name = types.StringValue(data.Name)
	m.Position = types.Int64Value(data.Position)
	m.ScheduleRotationableType = types.StringValue(data.ScheduleRotationableType)
	m.ActiveAllWeek = jsonapitypes.NullableBoolValue(data.ActiveAllWeek)
	m.ActiveDays = supertypes.NewSetValueOfSlice(ctx, data.ActiveDays)
	m.ActiveTimeType = jsonapitypes.NullableStringValue(data.ActiveTimeType)
	m.TimeZone = types.StringValue(data.TimeZone)
	m.StartTime = jsonapitypes.NullableStringValue(data.StartTime)
	m.EndTime = jsonapitypes.NullableStringValue(data.EndTime)
	m.ScheduleRotationableAttributes = supertypes.NewSingleNestedObjectValueOf(ctx, &ScheduleRotationResourceScheduleRotationAttributesModel{
		HandoffTime:     jsonapitypes.NullableStringValue(data.ScheduleRotationableAttributes.HandoffTime),
		HandoffDay:      jsonapitypes.NullableStringValue(data.ScheduleRotationableAttributes.HandoffDay),
		ShiftLength:     jsonapitypes.NullableInt64Value(data.ScheduleRotationableAttributes.ShiftLength),
		ShiftLengthUnit: jsonapitypes.NullableStringValue(data.ScheduleRotationableAttributes.ShiftLengthUnit),
	})

	if m.ActiveTimeAttributes.IsKnown() {
		m.ActiveTimeAttributes = supertypes.NewSetNestedObjectValueOfValueSlice(ctx, lo.Map(data.ActiveTimeAttributes, func(v apiclient.ScheduleRotationActiveTimeAttributes, _ int) ScheduleRotationResourceActiveTimeModel {
			return ScheduleRotationResourceActiveTimeModel{
				StartTime: types.StringValue(v.StartTime),
				EndTime:   types.StringValue(v.EndTime),
			}
		}))
	}

	if m.ScheduleRotationMembers.IsKnown() {
		m.ScheduleRotationMembers = supertypes.NewSetNestedObjectValueOfValueSlice(ctx, lo.Map(data.ScheduleRotationMembers, func(v apiclient.ScheduleRotationMember, _ int) ScheduleRotationResourceScheduleRotationMemberModel {
			return ScheduleRotationResourceScheduleRotationMemberModel{
				MemberId:   types.StringValue(v.MemberID),
				MemberType: types.StringValue(v.MemberType),
				Position:   types.Int64Value(v.Position),
			}
		}))
	}

	return nil
}

func (m *ScheduleRotationResourceModel) ToApi(ctx context.Context) (*apiclient.ScheduleRotation, diag.Diagnostics) {
	var diags diag.Diagnostics
	var data apiclient.ScheduleRotation

	if !m.Id.IsNull() && !m.Id.IsUnknown() {
		data.ID = m.Id.ValueString()
	}

	if !m.ScheduleId.IsNull() && !m.ScheduleId.IsUnknown() {
		data.ScheduleId = m.ScheduleId.ValueString()
	}

	if !m.Name.IsNull() && !m.Name.IsUnknown() {
		data.Name = m.Name.ValueString()
	}

	if !m.Position.IsNull() && !m.Position.IsUnknown() {
		data.Position = m.Position.ValueInt64()
	}

	if !m.ScheduleRotationableType.IsNull() && !m.ScheduleRotationableType.IsUnknown() {
		data.ScheduleRotationableType = m.ScheduleRotationableType.ValueString()
	}

	if !m.ActiveAllWeek.IsNull() && !m.ActiveAllWeek.IsUnknown() {
		data.ActiveAllWeek.Set(m.ActiveAllWeek.ValueBool())
	} else {
		data.ActiveAllWeek.SetNull()
	}

	if !m.ActiveDays.IsNull() && !m.ActiveDays.IsUnknown() {
		vv, diagss := m.ActiveDays.Get(ctx)
		diags.Append(diagss...)
		if diags.HasError() {
			return nil, diags
		}

		data.ActiveDays = vv
	} else {
		data.ActiveDays = []string{}
	}

	if !m.ActiveTimeType.IsNull() && !m.ActiveTimeType.IsUnknown() {
		data.ActiveTimeType.Set(m.ActiveTimeType.ValueString())
	} else {
		data.ActiveTimeType.SetNull()
	}

	if !m.TimeZone.IsNull() && !m.TimeZone.IsUnknown() {
		data.TimeZone = m.TimeZone.ValueString()
	}

	if !m.StartTime.IsNull() && !m.StartTime.IsUnknown() {
		data.StartTime.Set(m.StartTime.ValueString())
	} else {
		data.StartTime.SetNull()
	}

	if !m.EndTime.IsNull() && !m.EndTime.IsUnknown() {
		data.EndTime.Set(m.EndTime.ValueString())
	} else {
		data.EndTime.SetNull()
	}

	if !m.ScheduleRotationableAttributes.IsNull() && !m.ScheduleRotationableAttributes.IsUnknown() {
		vv, diagss := m.ScheduleRotationableAttributes.Get(ctx)
		diags.Append(diagss...)
		if diags.HasError() {
			return nil, diags
		}

		vvv, diagss := vv.ToApi(ctx)
		diags.Append(diagss...)
		if diags.HasError() {
			return nil, diags
		}

		data.ScheduleRotationableAttributes = *vvv
	}

	data.ActiveTimeAttributes = []apiclient.ScheduleRotationActiveTimeAttributes{}
	if !m.ActiveTimeAttributes.IsNull() && !m.ActiveTimeAttributes.IsUnknown() {
		vv, diagss := m.ActiveTimeAttributes.Get(ctx)
		diags.Append(diagss...)
		if diags.HasError() {
			return nil, diags
		}

		for _, v := range vv {
			vvv, diagss := v.ToApi(ctx)
			diags.Append(diagss...)
			if diags.HasError() {
				return nil, diags
			}

			data.ActiveTimeAttributes = append(data.ActiveTimeAttributes, *vvv)
		}
	}

	data.ScheduleRotationMembers = []apiclient.ScheduleRotationMember{}
	if !m.ScheduleRotationMembers.IsNull() && !m.ScheduleRotationMembers.IsUnknown() {
		vv, diagss := m.ScheduleRotationMembers.Get(ctx)
		diags.Append(diagss...)
		if diags.HasError() {
			return nil, diags
		}

		for _, v := range vv {
			vvv, diagss := v.ToApi(ctx)
			diags.Append(diagss...)
			if diags.HasError() {
				return nil, diags
			}

			data.ScheduleRotationMembers = append(data.ScheduleRotationMembers, *vvv)
		}
	}

	return &data, diags
}

type ScheduleRotationResourceScheduleRotationAttributesModel struct {
	HandoffTime     types.String `tfsdk:"handoff_time"`
	HandoffDay      types.String `tfsdk:"handoff_day"`
	ShiftLength     types.Int64  `tfsdk:"shift_length"`
	ShiftLengthUnit types.String `tfsdk:"shift_length_unit"`
}

func (m *ScheduleRotationResourceScheduleRotationAttributesModel) ToApi(ctx context.Context) (*apiclient.ScheduleRotationScheduleRotationableAttributes, diag.Diagnostics) {
	var data apiclient.ScheduleRotationScheduleRotationableAttributes

	if !m.HandoffTime.IsNull() && !m.HandoffTime.IsUnknown() {
		data.HandoffTime.Set(m.HandoffTime.ValueString())
	} else {
		data.HandoffTime.SetNull()
	}

	if !m.HandoffDay.IsNull() && !m.HandoffDay.IsUnknown() {
		data.HandoffDay.Set(m.HandoffDay.ValueString())
	} else {
		data.HandoffDay.SetNull()
	}

	if !m.ShiftLength.IsNull() && !m.ShiftLength.IsUnknown() {
		data.ShiftLength.Set(m.ShiftLength.ValueInt64())
	} else {
		data.ShiftLength.SetNull()
	}

	if !m.ShiftLengthUnit.IsNull() && !m.ShiftLengthUnit.IsUnknown() {
		data.ShiftLengthUnit.Set(m.ShiftLengthUnit.ValueString())
	} else {
		data.ShiftLengthUnit.SetNull()
	}

	return &data, nil
}

type ScheduleRotationResourceActiveTimeModel struct {
	StartTime types.String `tfsdk:"start_time"`
	EndTime   types.String `tfsdk:"end_time"`
}

func (m *ScheduleRotationResourceActiveTimeModel) ToApi(ctx context.Context) (*apiclient.ScheduleRotationActiveTimeAttributes, diag.Diagnostics) {
	return &apiclient.ScheduleRotationActiveTimeAttributes{
		StartTime: m.StartTime.ValueString(),
		EndTime:   m.EndTime.ValueString(),
	}, nil
}

type ScheduleRotationResourceScheduleRotationMemberModel struct {
	MemberId   types.String `tfsdk:"member_id"`
	MemberType types.String `tfsdk:"member_type"`
	Position   types.Int64  `tfsdk:"position"`
}

func (m *ScheduleRotationResourceScheduleRotationMemberModel) ToApi(ctx context.Context) (*apiclient.ScheduleRotationMember, diag.Diagnostics) {
	return &apiclient.ScheduleRotationMember{
		MemberID:   m.MemberId.ValueString(),
		MemberType: m.MemberType.ValueString(),
		Position:   m.Position.ValueInt64(),
	}, nil
}
