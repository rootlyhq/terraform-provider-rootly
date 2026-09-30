package apiclient

import (
	"bytes"
	"context"
	"io"

	"github.com/rootlyhq/jsonapi"
)

type Schedule struct {
	ID                                 string                                       `jsonapi:"primary,schedules"`
	Name                               string                                       `jsonapi:"attr,name"`
	Description                        jsonapi.NullableAttr[string]                 `jsonapi:"attr,description"`
	AllTimeCoverage                    jsonapi.NullableAttr[bool]                   `jsonapi:"attr,all_time_coverage"`
	SlackUserGroup                     jsonapi.NullableAttr[ScheduleSlackUserGroup] `jsonapi:"attr,slack_user_group"`
	SlackChannel                       jsonapi.NullableAttr[ScheduleSlackChannel]   `jsonapi:"attr,slack_channel"`
	OwnerGroupIds                      jsonapi.NullableAttr[[]string]               `jsonapi:"attr,owner_group_ids"`
	OwnerUserId                        jsonapi.NullableAttr[int64]                  `jsonapi:"attr,owner_user_id"`
	SyncLinearEnabled                  jsonapi.NullableAttr[bool]                   `jsonapi:"attr,sync_linear_enabled"`
	IncludeShadowsInSlackNotifications jsonapi.NullableAttr[bool]                   `jsonapi:"attr,include_shadows_in_slack_notifications"`
	ShiftStartNotificationsEnabled     jsonapi.NullableAttr[bool]                   `jsonapi:"attr,shift_start_notifications_enabled"`
	ShiftUpdateNotificationsEnabled    jsonapi.NullableAttr[bool]                   `jsonapi:"attr,shift_update_notifications_enabled"`
	ShiftReportEnabled                 jsonapi.NullableAttr[bool]                   `jsonapi:"attr,shift_report_enabled"`
	ShiftReportDayOfWeek               jsonapi.NullableAttr[string]                 `jsonapi:"attr,shift_report_day_of_week"`
	ShiftReportTimeOfDay               jsonapi.NullableAttr[string]                 `jsonapi:"attr,shift_report_time_of_day"`
	ShiftReportTimeZone                jsonapi.NullableAttr[string]                 `jsonapi:"attr,shift_report_time_zone"`
	TimeZone                           jsonapi.NullableAttr[string]                 `jsonapi:"attr,time_zone"`
}

type ScheduleSlackUserGroup struct {
	Id   jsonapi.NullableAttr[string] `jsonapi:"attr,id"`
	Name jsonapi.NullableAttr[string] `jsonapi:"attr,name"`
}

type ScheduleSlackChannel struct {
	Id   jsonapi.NullableAttr[string] `jsonapi:"attr,id"`
	Name jsonapi.NullableAttr[string] `jsonapi:"attr,name"`
}

func (c *Client) GetSchedule(ctx context.Context, id string) (*Schedule, error) {
	httpResp, err := c.ClientWithResponses.GetSchedule(ctx, id)
	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	var resp Schedule
	if err := jsonapi.UnmarshalPayload(bytes.NewReader(body), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) CreateSchedule(ctx context.Context, req Schedule) (*Schedule, error) {
	var buf bytes.Buffer
	if err := jsonapi.MarshalPayload(&buf, &req); err != nil {
		return nil, err
	}

	httpResp, err := c.ClientWithResponses.CreateScheduleWithBody(ctx, "application/vnd.api+json", bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	var resp Schedule
	if err := jsonapi.UnmarshalPayload(bytes.NewReader(body), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) UpdateSchedule(ctx context.Context, req Schedule) (*Schedule, error) {
	var buf bytes.Buffer
	if err := jsonapi.MarshalPayload(&buf, &req); err != nil {
		return nil, err
	}

	httpResp, err := c.ClientWithResponses.UpdateScheduleWithBody(ctx, req.ID, "application/vnd.api+json", bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	var resp Schedule
	if err := jsonapi.UnmarshalPayload(bytes.NewReader(body), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
