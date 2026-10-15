package provider

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
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
	"github.com/rootlyhq/terraform-provider-rootly/v5/internal/fwtypes"
	"github.com/rootlyhq/terraform-provider-rootly/v5/internal/jsonapitypes"
	intmapvalidator "github.com/rootlyhq/terraform-provider-rootly/v5/internal/validators/mapvalidator"
	intsetvalidator "github.com/rootlyhq/terraform-provider-rootly/v5/internal/validators/setvalidator"
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
				MarkdownDescription: "Value must be one of `all_day`, `same_time`, or `custom`. The value chosen will override `active_time_attributes` in any `rootly_schedule_rotation_active_day` resources linked to this `rootly_schedule_rotation`. Defaults to `all_day`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("all_day"),
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
			"schedule_rotationable_attributes": schema.MapAttribute{
				CustomType:  supertypes.NewMapTypeOf[string](ctx),
				ElementType: types.StringType,
				MarkdownDescription: "Schedule rotation type specific attributes. `handoff_time` and/or `handoff_day` may be required, depending on `schedule_rotationable_type`. The following keys are supported:\n" +
					"  - `handoff_time` - (Required) Hand off time, e.g. `09:00`. Applicable for daily, weekly/biweekly, monthly, and custom rotations.\n" +
					"  - `handoff_day` - (Optional) Hand off day. Applicable for weekly/biweekly and monthly rotations.\n" +
					"  - `shift_length` - (Optional) Shift length, as an integer (e.g. `7` or `\"7\"`). Applicable for custom rotations.\n" +
					"  - `shift_length_unit` - (Optional) Shift length unit. Applicable for custom rotations. Value must be one of `hours`, `days`, `weeks`.\n\n" +
					"Please see API docs for options based on `schedule_rotationable_type`: https://docs.rootly.com/api-reference/schedulerotations/creates-a-schedule-rotation#response-data-attributes-schedule-rotationable-attributes",
				Required: true,
				Validators: []validator.Map{
					mapvalidator.KeysAre(stringvalidator.OneOf("handoff_time", "handoff_day", "shift_length", "shift_length_unit")),
					intmapvalidator.KeysRequired("handoff_time"),
					intmapvalidator.ValueStringAtKey("shift_length", stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9]+$`), "must be an integer")),
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
				Validators: []validator.Set{
					intsetvalidator.UniqueByAttribute("position"),
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
	Id                             types.String                                                                           `tfsdk:"id"`
	ScheduleId                     types.String                                                                           `tfsdk:"schedule_id"`
	Name                           types.String                                                                           `tfsdk:"name"`
	Position                       types.Int64                                                                            `tfsdk:"position"`
	ScheduleRotationableType       types.String                                                                           `tfsdk:"schedule_rotationable_type"`
	ActiveAllWeek                  types.Bool                                                                             `tfsdk:"active_all_week"`
	ActiveDays                     supertypes.SetValueOf[string]                                                          `tfsdk:"active_days"`
	ActiveTimeType                 types.String                                                                           `tfsdk:"active_time_type"`
	TimeZone                       types.String                                                                           `tfsdk:"time_zone"`
	StartTime                      types.String                                                                           `tfsdk:"start_time"`
	EndTime                        types.String                                                                           `tfsdk:"end_time"`
	ScheduleRotationableAttributes supertypes.MapValueOf[string]                                                          `tfsdk:"schedule_rotationable_attributes"`
	ActiveTimeAttributes           supertypes.SetNestedObjectValueOf[ScheduleRotationResourceActiveTimeModel]             `tfsdk:"active_time_attributes"`
	ScheduleRotationMembers        supertypes.SetNestedObjectValueOf[ScheduleRotationResourceScheduleRotationMemberModel] `tfsdk:"schedule_rotation_members"`
}

func (m *ScheduleRotationResourceModel) FromApi(ctx context.Context, data apiclient.ScheduleRotation) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Id = types.StringValue(data.ID)
	m.ScheduleId = types.StringValue(data.ScheduleId)
	m.Name = types.StringValue(data.Name)
	m.Position = types.Int64Value(data.Position)
	m.ScheduleRotationableType = types.StringValue(data.ScheduleRotationableType)
	m.ActiveAllWeek = jsonapitypes.NullableBoolValue(data.ActiveAllWeek)
	m.ActiveDays = supertypes.NewSetValueOfSlice(ctx, data.ActiveDays)
	m.ActiveTimeType = types.StringValue(data.ActiveTimeType)
	m.TimeZone = types.StringValue(data.TimeZone)
	m.StartTime = jsonapitypes.NullableStringValue(data.StartTime)
	m.EndTime = jsonapitypes.NullableStringValue(data.EndTime)

	scheduleRotationableAttributes := map[string]string{}
	if v, err := data.ScheduleRotationableAttributes.HandoffTime.Get(); err == nil {
		scheduleRotationableAttributes["handoff_time"] = v
	}
	if v, err := data.ScheduleRotationableAttributes.HandoffDay.Get(); err == nil {
		scheduleRotationableAttributes["handoff_day"] = v
	}
	if v, err := data.ScheduleRotationableAttributes.ShiftLength.Get(); err == nil {
		scheduleRotationableAttributes["shift_length"] = strconv.FormatInt(v, 10)
	}
	if v, err := data.ScheduleRotationableAttributes.ShiftLengthUnit.Get(); err == nil {
		scheduleRotationableAttributes["shift_length_unit"] = v
	}
	var diagss diag.Diagnostics
	m.ScheduleRotationableAttributes, diagss = supertypes.NewMapValueOfMap(ctx, scheduleRotationableAttributes)
	diags.Append(diagss...)
	if diags.HasError() {
		return diags
	}

	if fwtypes.IsKnown(m.ActiveTimeAttributes) {
		m.ActiveTimeAttributes = supertypes.NewSetNestedObjectValueOfValueSlice(ctx, lo.Map(data.ActiveTimeAttributes, func(v apiclient.ScheduleRotationActiveTimeAttributes, _ int) ScheduleRotationResourceActiveTimeModel {
			return ScheduleRotationResourceActiveTimeModel{
				StartTime: types.StringValue(v.StartTime),
				EndTime:   types.StringValue(v.EndTime),
			}
		}))
	}

	if fwtypes.IsKnown(m.ScheduleRotationMembers) {
		m.ScheduleRotationMembers = supertypes.NewSetNestedObjectValueOfValueSlice(ctx, lo.Map(data.ScheduleRotationMembers, func(v apiclient.ScheduleRotationMember, _ int) ScheduleRotationResourceScheduleRotationMemberModel {
			return ScheduleRotationResourceScheduleRotationMemberModel{
				MemberId:   types.StringValue(v.MemberID),
				MemberType: types.StringValue(v.MemberType),
				Position:   types.Int64Value(v.Position),
			}
		}))
	}

	return diags
}

func (m *ScheduleRotationResourceModel) ToApi(ctx context.Context) (*apiclient.ScheduleRotation, diag.Diagnostics) {
	var diags diag.Diagnostics
	var data apiclient.ScheduleRotation

	if fwtypes.IsKnown(m.Id) {
		data.ID = m.Id.ValueString()
	}

	if fwtypes.IsKnown(m.ScheduleId) {
		data.ScheduleId = m.ScheduleId.ValueString()
	}

	if fwtypes.IsKnown(m.Name) {
		data.Name = m.Name.ValueString()
	}

	if fwtypes.IsKnown(m.Position) {
		data.Position = m.Position.ValueInt64()
	}

	if fwtypes.IsKnown(m.ScheduleRotationableType) {
		data.ScheduleRotationableType = m.ScheduleRotationableType.ValueString()
	}

	if fwtypes.IsKnown(m.ActiveAllWeek) {
		data.ActiveAllWeek.Set(m.ActiveAllWeek.ValueBool())
	}

	if fwtypes.IsKnown(m.ActiveDays) {
		vv, diagss := m.ActiveDays.Get(ctx)
		diags.Append(diagss...)
		if diags.HasError() {
			return nil, diags
		}

		data.ActiveDays = vv
	} else {
		data.ActiveDays = []string{}
	}

	if fwtypes.IsKnown(m.ActiveTimeType) {
		data.ActiveTimeType = m.ActiveTimeType.ValueString()
	}

	if fwtypes.IsKnown(m.TimeZone) {
		data.TimeZone = m.TimeZone.ValueString()
	}

	if fwtypes.IsKnown(m.StartTime) {
		data.StartTime.Set(m.StartTime.ValueString())
	} else {
		data.StartTime.SetNull()
	}

	if fwtypes.IsKnown(m.EndTime) {
		data.EndTime.Set(m.EndTime.ValueString())
	} else {
		data.EndTime.SetNull()
	}

	if fwtypes.IsKnown(m.ScheduleRotationableAttributes) {
		vv, diagss := m.ScheduleRotationableAttributes.Get(ctx)
		diags.Append(diagss...)
		if diags.HasError() {
			return nil, diags
		}

		vvv, diagss := scheduleRotationableAttributesToApi(vv)
		diags.Append(diagss...)
		if diags.HasError() {
			return nil, diags
		}

		data.ScheduleRotationableAttributes = *vvv
	}

	data.ActiveTimeAttributes = []apiclient.ScheduleRotationActiveTimeAttributes{}
	if fwtypes.IsKnown(m.ActiveTimeAttributes) {
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
	if fwtypes.IsKnown(m.ScheduleRotationMembers) {
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

func scheduleRotationableAttributesToApi(m map[string]string) (*apiclient.ScheduleRotationScheduleRotationableAttributes, diag.Diagnostics) {
	var diags diag.Diagnostics
	var data apiclient.ScheduleRotationScheduleRotationableAttributes

	if v, ok := m["handoff_time"]; ok {
		data.HandoffTime.Set(v)
	} else {
		data.HandoffTime.SetNull()
	}

	if v, ok := m["handoff_day"]; ok {
		data.HandoffDay.Set(v)
	} else {
		data.HandoffDay.SetNull()
	}

	if v, ok := m["shift_length"]; ok {
		vv, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			diags.AddAttributeError(
				path.Root("schedule_rotationable_attributes").AtMapKey("shift_length"),
				"Invalid shift_length",
				fmt.Sprintf("Unable to parse %q as an integer: %s", v, err),
			)
			return nil, diags
		}
		data.ShiftLength.Set(vv)
	} else {
		data.ShiftLength.SetNull()
	}

	if v, ok := m["shift_length_unit"]; ok {
		data.ShiftLengthUnit.Set(v)
	} else {
		data.ShiftLengthUnit.SetNull()
	}

	return &data, diags
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
