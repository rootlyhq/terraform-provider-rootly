package provider

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	"github.com/rootlyhq/terraform-provider-rootly/v5/client"
	"github.com/rootlyhq/terraform-provider-rootly/v5/internal/apiclient"
	"github.com/rootlyhq/terraform-provider-rootly/v5/internal/diagutils"
	"github.com/rootlyhq/terraform-provider-rootly/v5/internal/fwtypes"
	"github.com/rootlyhq/terraform-provider-rootly/v5/internal/jsonapitypes"
)

var _ resource.Resource = &ScheduleResource{}
var _ resource.ResourceWithImportState = &ScheduleResource{}
var _ resource.ResourceWithUpgradeState = &ScheduleResource{}

func NewScheduleResource() resource.Resource {
	return &ScheduleResource{}
}

type ScheduleResource struct {
	baseResource
}

func (r *ScheduleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schedule"
}

func (r *ScheduleResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a schedule.",
		Version:             1,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of this resource.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the schedule.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the schedule.",
				Optional:            true,
				Computed:            true,
			},
			"all_time_coverage": schema.BoolAttribute{
				MarkdownDescription: "24/7 coverage of the schedule. Value must be one of true or false.",
				Optional:            true,
				Computed:            true,
			},
			"owner_group_ids": schema.SetAttribute{
				CustomType:          supertypes.NewSetTypeOf[string](ctx),
				ElementType:         types.StringType,
				MarkdownDescription: "The owning teams for this schedules.",
				Optional:            true,
				Computed:            true,
			},
			"owner_user_id": schema.Int64Attribute{
				MarkdownDescription: "ID of user assigned as owner of the schedule. Defaults to the API token's user if not specified.",
				Optional:            true,
				Computed:            true,
			},
			"sync_linear_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the schedule is synced with Linear. Value must be one of true or false.",
				Optional:            true,
				Computed:            true,
			},
			"include_shadows_in_slack_notifications": schema.BoolAttribute{
				MarkdownDescription: "Whether shadow users are included in Slack notifications and user group syncing. Value must be one of true or false.",
				Optional:            true,
				Computed:            true,
			},
			"shift_start_notifications_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether to send a Slack message every time a new shift begins. Requires `slack_channel` to be set. Value must be one of true or false.",
				Optional:            true,
				Computed:            true,
			},
			"shift_update_notifications_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether to send a Slack message whenever a shift is updated (overrides, removed users, rotation changes, etc.). Requires `slack_channel` to be set. Value must be one of true or false.",
				Optional:            true,
				Computed:            true,
			},
			"shift_report_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the weekly shift summary report is sent. Requires `slack_channel` to be set. Value must be one of true or false.",
				Optional:            true,
				Computed:            true,
			},
			"shift_report_day_of_week": schema.StringAttribute{
				MarkdownDescription: "Day of week the weekly shift summary is sent. Value must be one of `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, `saturday`, `sunday`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"),
				},
			},
			"shift_report_time_of_day": schema.StringAttribute{
				MarkdownDescription: "Time of day the weekly shift summary is sent, in `HH:MM` 24-hour format.",
				Optional:            true,
				Computed:            true,
			},
			"shift_report_time_zone": schema.StringAttribute{
				MarkdownDescription: "IANA time zone used for the weekly shift summary (e.g. `Australia/Sydney`).",
				Optional:            true,
				Computed:            true,
			},
			"time_zone": schema.StringAttribute{
				MarkdownDescription: "A valid IANA time zone name. Only applicable when config_one_timezone_per_schedule_enabled is true for the organization.",
				Optional:            true,
				Computed:            true,
			},
			"slack_user_group": schema.SingleNestedAttribute{
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[ScheduleResourceSlackUserGroupModel](ctx),
				MarkdownDescription: "Synced slack group of the schedule. To set, specify all nested attributes. To remove, set the attribute to an empty object `{}`.",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "Slack user group ID.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.AlsoRequires(path.MatchRoot("slack_user_group").AtName("name")),
						},
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "Slack user group name.",
						Optional:            true,
					},
				},
			},
			"slack_channel": schema.SingleNestedAttribute{
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[ScheduleResourceSlackChannelModel](ctx),
				MarkdownDescription: "Synced slack channel of the schedule. To set, specify all nested attributes. To remove, set the attribute to an empty object `{}`.",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "Slack channel ID.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.AlsoRequires(path.MatchRoot("slack_channel").AtName("name")),
						},
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "Slack channel name.",
						Optional:            true,
					},
				},
			},
			"business_hours": schema.SingleNestedAttribute{
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[ScheduleResourceBusinessHoursModel](ctx),
				MarkdownDescription: "Controls shadow paging on the schedule. To enable, specify all nested attributes. To remove, set the attribute to an empty object `{}`. `start_time` and `end_time` are HH:MM 24-hour format strings.",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"start_time": schema.StringAttribute{
						MarkdownDescription: "Start time in HH:MM 24-hour format.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.AlsoRequires(path.MatchRoot("business_hours").AtName("end_time")),
							stringvalidator.AlsoRequires(path.MatchRoot("business_hours").AtName("include_weekends")),
						},
					},
					"end_time": schema.StringAttribute{
						MarkdownDescription: "End time in HH:MM 24-hour format.",
						Optional:            true,
					},
					"include_weekends": schema.BoolAttribute{
						MarkdownDescription: "Whether to include weekends.",
						Optional:            true,
					},
				},
			},
		},
	}
}

func (r *ScheduleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ScheduleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item := diagutils.MergeDiagnostics(data.ToApi(ctx))(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.CreateSchedule(ctx, *item)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create schedule", err.Error())
		return
	}

	data.Id = types.StringValue(res.ID)

	item, err = r.client.GetSchedule(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read schedule", err.Error())
		return
	}

	resp.Diagnostics.Append(data.FromApi(ctx, *item)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ScheduleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ScheduleResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.GetSchedule(ctx, data.Id.ValueString())
	if err != nil {
		if errors.Is(err, client.NotFoundError{}) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read schedule", err.Error())
		return
	}

	resp.Diagnostics.Append(data.FromApi(ctx, *item)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ScheduleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ScheduleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item := diagutils.MergeDiagnostics(data.ToApi(ctx))(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateSchedule(ctx, *item)
	if err != nil {
		resp.Diagnostics.AddError("Unable to update schedule", err.Error())
		return
	}

	item, err = r.client.GetSchedule(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read schedule", err.Error())
		return
	}

	resp.Diagnostics.Append(data.FromApi(ctx, *item)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ScheduleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ScheduleResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.legacyClient.DeleteSchedule(data.Id.ValueString())
	if err != nil {
		if errors.Is(err, client.NotFoundError{}) {
			return
		}
		resp.Diagnostics.AddError("Unable to delete schedule", err.Error())
	}
}

func (r *ScheduleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ScheduleResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"slack_user_group": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				// The old slack_user_group string attribute is dropped and replaced with a nested object.
				// The new nested object will be populated from the API on the next read operation.
			},
		},
	}
}

type ScheduleResourceModel struct {
	Id                                 types.String                                                              `tfsdk:"id"`
	Name                               types.String                                                              `tfsdk:"name"`
	Description                        types.String                                                              `tfsdk:"description"`
	AllTimeCoverage                    types.Bool                                                                `tfsdk:"all_time_coverage"`
	SlackUserGroup                     supertypes.SingleNestedObjectValueOf[ScheduleResourceSlackUserGroupModel] `tfsdk:"slack_user_group"`
	SlackChannel                       supertypes.SingleNestedObjectValueOf[ScheduleResourceSlackChannelModel]   `tfsdk:"slack_channel"`
	OwnerGroupIds                      supertypes.SetValueOf[string]                                             `tfsdk:"owner_group_ids"`
	OwnerUserId                        types.Int64                                                               `tfsdk:"owner_user_id"`
	SyncLinearEnabled                  types.Bool                                                                `tfsdk:"sync_linear_enabled"`
	IncludeShadowsInSlackNotifications types.Bool                                                                `tfsdk:"include_shadows_in_slack_notifications"`
	ShiftStartNotificationsEnabled     types.Bool                                                                `tfsdk:"shift_start_notifications_enabled"`
	ShiftUpdateNotificationsEnabled    types.Bool                                                                `tfsdk:"shift_update_notifications_enabled"`
	ShiftReportEnabled                 types.Bool                                                                `tfsdk:"shift_report_enabled"`
	ShiftReportDayOfWeek               types.String                                                              `tfsdk:"shift_report_day_of_week"`
	ShiftReportTimeOfDay               types.String                                                              `tfsdk:"shift_report_time_of_day"`
	ShiftReportTimeZone                types.String                                                              `tfsdk:"shift_report_time_zone"`
	TimeZone                           types.String                                                              `tfsdk:"time_zone"`
	BusinessHours                      supertypes.SingleNestedObjectValueOf[ScheduleResourceBusinessHoursModel]  `tfsdk:"business_hours"`
}

func (m *ScheduleResourceModel) FromApi(ctx context.Context, data apiclient.Schedule) (diags diag.Diagnostics) {
	m.Id = types.StringValue(data.ID)
	m.Name = types.StringValue(data.Name)
	m.Description = jsonapitypes.NullableStringValue(data.Description)
	m.AllTimeCoverage = jsonapitypes.NullableBoolValue(data.AllTimeCoverage)
	m.OwnerUserId = jsonapitypes.NullableInt64Value(data.OwnerUserId)
	m.SyncLinearEnabled = jsonapitypes.NullableBoolValue(data.SyncLinearEnabled)
	m.IncludeShadowsInSlackNotifications = jsonapitypes.NullableBoolValue(data.IncludeShadowsInSlackNotifications)
	m.ShiftStartNotificationsEnabled = jsonapitypes.NullableBoolValue(data.ShiftStartNotificationsEnabled)
	m.ShiftUpdateNotificationsEnabled = jsonapitypes.NullableBoolValue(data.ShiftUpdateNotificationsEnabled)
	m.ShiftReportEnabled = jsonapitypes.NullableBoolValue(data.ShiftReportEnabled)
	m.ShiftReportDayOfWeek = jsonapitypes.NullableStringValue(data.ShiftReportDayOfWeek)
	m.ShiftReportTimeOfDay = jsonapitypes.NullableStringValue(data.ShiftReportTimeOfDay)
	m.ShiftReportTimeZone = jsonapitypes.NullableStringValue(data.ShiftReportTimeZone)
	m.TimeZone = jsonapitypes.NullableStringValue(data.TimeZone)
	m.OwnerGroupIds = jsonapitypes.NullableSetValueOfSlice(ctx, data.OwnerGroupIds)

	if v, err := data.SlackUserGroup.Get(); err == nil {
		var mm ScheduleResourceSlackUserGroupModel
		diags.Append(mm.FromApi(ctx, v)...)
		if diags.HasError() {
			return
		}
		m.SlackUserGroup = supertypes.NewSingleNestedObjectValueOf(ctx, &mm)
	} else {
		m.SlackUserGroup = supertypes.NewSingleNestedObjectValueOfNull[ScheduleResourceSlackUserGroupModel](ctx)
	}

	if v, err := data.SlackChannel.Get(); err == nil {
		var mm ScheduleResourceSlackChannelModel
		diags.Append(mm.FromApi(ctx, v)...)
		if diags.HasError() {
			return
		}
		m.SlackChannel = supertypes.NewSingleNestedObjectValueOf(ctx, &mm)
	} else {
		m.SlackChannel = supertypes.NewSingleNestedObjectValueOf(ctx, &ScheduleResourceSlackChannelModel{
			Id:   types.StringNull(),
			Name: types.StringNull(),
		})
	}

	if v, err := data.BusinessHours.Get(); err == nil {
		var mm ScheduleResourceBusinessHoursModel
		diags.Append(mm.FromApi(ctx, v)...)
		if diags.HasError() {
			return
		}
		m.BusinessHours = supertypes.NewSingleNestedObjectValueOf(ctx, &mm)
	} else {
		m.BusinessHours = supertypes.NewSingleNestedObjectValueOf(ctx, &ScheduleResourceBusinessHoursModel{
			StartTime:       types.StringNull(),
			EndTime:         types.StringNull(),
			IncludeWeekends: types.BoolNull(),
		})
	}

	return
}

func (m *ScheduleResourceModel) ToApi(ctx context.Context) (*apiclient.Schedule, diag.Diagnostics) {
	var diags diag.Diagnostics
	var data apiclient.Schedule

	if fwtypes.IsKnown(m.Id) {
		data.ID = m.Id.ValueString()
	}

	if fwtypes.IsKnown(m.Name) {
		data.Name = m.Name.ValueString()
	}

	if fwtypes.IsKnown(m.Description) {
		data.Description.Set(m.Description.ValueString())
	}

	if fwtypes.IsKnown(m.AllTimeCoverage) {
		data.AllTimeCoverage.Set(m.AllTimeCoverage.ValueBool())
	}

	if fwtypes.IsKnown(m.OwnerUserId) {
		data.OwnerUserId.Set(m.OwnerUserId.ValueInt64())
	}

	if fwtypes.IsKnown(m.SyncLinearEnabled) {
		data.SyncLinearEnabled.Set(m.SyncLinearEnabled.ValueBool())
	}

	if fwtypes.IsKnown(m.IncludeShadowsInSlackNotifications) {
		data.IncludeShadowsInSlackNotifications.Set(m.IncludeShadowsInSlackNotifications.ValueBool())
	}

	if fwtypes.IsKnown(m.ShiftStartNotificationsEnabled) {
		data.ShiftStartNotificationsEnabled.Set(m.ShiftStartNotificationsEnabled.ValueBool())
	}

	if fwtypes.IsKnown(m.ShiftUpdateNotificationsEnabled) {
		data.ShiftUpdateNotificationsEnabled.Set(m.ShiftUpdateNotificationsEnabled.ValueBool())
	}

	if fwtypes.IsKnown(m.ShiftReportEnabled) {
		data.ShiftReportEnabled.Set(m.ShiftReportEnabled.ValueBool())
	}

	if fwtypes.IsKnown(m.ShiftReportDayOfWeek) {
		data.ShiftReportDayOfWeek.Set(m.ShiftReportDayOfWeek.ValueString())
	}

	if fwtypes.IsKnown(m.ShiftReportTimeOfDay) {
		data.ShiftReportTimeOfDay.Set(m.ShiftReportTimeOfDay.ValueString())
	}

	if fwtypes.IsKnown(m.ShiftReportTimeZone) {
		data.ShiftReportTimeZone.Set(m.ShiftReportTimeZone.ValueString())
	}

	if fwtypes.IsKnown(m.TimeZone) {
		data.TimeZone.Set(m.TimeZone.ValueString())
	}

	if fwtypes.IsKnown(m.OwnerGroupIds) {
		vv, diagss := m.OwnerGroupIds.Get(ctx)
		diags.Append(diagss...)
		if diags.HasError() {
			return nil, diags
		}
		data.OwnerGroupIds.Set(vv)
	}

	if fwtypes.IsKnown(m.SlackUserGroup) {
		vv, diagss := m.SlackUserGroup.Get(ctx)
		diags.Append(diagss...)
		if diags.HasError() {
			return nil, diags
		}

		vvv, diagss := vv.ToApi(ctx)
		diags.Append(diagss...)
		if diags.HasError() {
			return nil, diags
		}

		data.SlackUserGroup.Set(*vvv)
	}

	if fwtypes.IsKnown(m.SlackChannel) {
		vv, diagss := m.SlackChannel.Get(ctx)
		diags.Append(diagss...)
		if diags.HasError() {
			return nil, diags
		}

		vvv, ok, diagss := vv.ToApi(ctx)
		diags.Append(diagss...)
		if diags.HasError() {
			return nil, diags
		}

		if ok {
			data.SlackChannel.Set(*vvv)
		} else {
			data.SlackChannel.SetNull()
		}
	}

	if fwtypes.IsKnown(m.BusinessHours) {
		vv, diagss := m.BusinessHours.Get(ctx)
		diags.Append(diagss...)
		if diags.HasError() {
			return nil, diags
		}

		vvv, ok, diagss := vv.ToApi(ctx)
		diags.Append(diagss...)
		if diags.HasError() {
			return nil, diags
		}

		if ok {
			data.BusinessHours.Set(*vvv)
		} else {
			data.BusinessHours.SetNull()
		}
	}

	return &data, diags
}

type ScheduleResourceSlackUserGroupModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (m *ScheduleResourceSlackUserGroupModel) ToApi(ctx context.Context) (*apiclient.ScheduleSlackUserGroup, diag.Diagnostics) {
	var data apiclient.ScheduleSlackUserGroup

	if fwtypes.IsKnown(m.Id) {
		data.Id.Set(m.Id.ValueString())
	}

	if fwtypes.IsKnown(m.Name) {
		data.Name.Set(m.Name.ValueString())
	}

	return &data, nil
}

func (m *ScheduleResourceSlackUserGroupModel) FromApi(ctx context.Context, data apiclient.ScheduleSlackUserGroup) diag.Diagnostics {
	m.Id = jsonapitypes.NullableStringValue(data.Id)
	m.Name = jsonapitypes.NullableStringValue(data.Name)

	return nil
}

type ScheduleResourceSlackChannelModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (m *ScheduleResourceSlackChannelModel) ToApi(ctx context.Context) (*apiclient.ScheduleSlackChannel, bool, diag.Diagnostics) {
	var data apiclient.ScheduleSlackChannel
	var ok bool

	if fwtypes.IsKnown(m.Id) {
		ok = true
		data.Id.Set(m.Id.ValueString())
	} else {
		data.Id.SetNull()
	}

	if fwtypes.IsKnown(m.Name) {
		ok = true
		data.Name.Set(m.Name.ValueString())
	} else {
		data.Name.SetNull()
	}

	return &data, ok, nil
}

func (m *ScheduleResourceSlackChannelModel) FromApi(ctx context.Context, data apiclient.ScheduleSlackChannel) diag.Diagnostics {
	m.Id = jsonapitypes.NullableStringValue(data.Id)
	m.Name = jsonapitypes.NullableStringValue(data.Name)

	return nil
}

type ScheduleResourceBusinessHoursModel struct {
	StartTime       types.String `tfsdk:"start_time"`
	EndTime         types.String `tfsdk:"end_time"`
	IncludeWeekends types.Bool   `tfsdk:"include_weekends"`
}

func (m *ScheduleResourceBusinessHoursModel) ToApi(ctx context.Context) (*apiclient.ScheduleBusinessHours, bool, diag.Diagnostics) {
	var data apiclient.ScheduleBusinessHours
	var ok bool

	if fwtypes.IsKnown(m.StartTime) {
		ok = true
		data.StartTime.Set(m.StartTime.ValueString())
	}

	if fwtypes.IsKnown(m.EndTime) {
		ok = true
		data.EndTime.Set(m.EndTime.ValueString())
	}

	if fwtypes.IsKnown(m.IncludeWeekends) {
		ok = true
		data.IncludeWeekends.Set(m.IncludeWeekends.ValueBool())
	}

	return &data, ok, nil
}

func (m *ScheduleResourceBusinessHoursModel) FromApi(ctx context.Context, data apiclient.ScheduleBusinessHours) diag.Diagnostics {
	m.StartTime = jsonapitypes.NullableStringValue(data.StartTime)
	m.EndTime = jsonapitypes.NullableStringValue(data.EndTime)
	m.IncludeWeekends = jsonapitypes.NullableBoolValue(data.IncludeWeekends)

	return nil
}
