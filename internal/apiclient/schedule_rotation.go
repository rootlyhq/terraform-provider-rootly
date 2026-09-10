package apiclient

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/rootlyhq/jsonapi"
)

type ScheduleRotation struct {
	ID                             string                                         `jsonapi:"primary,schedule_rotations"`
	ScheduleId                     string                                         `jsonapi:"attr,schedule_id,omitempty"`
	Name                           string                                         `jsonapi:"attr,name"`
	Position                       int64                                          `jsonapi:"attr,position,omitempty"`
	ScheduleRotationableType       string                                         `jsonapi:"attr,schedule_rotationable_type"`
	ActiveAllWeek                  jsonapi.NullableAttr[bool]                     `jsonapi:"attr,active_all_week,omitempty"`
	ActiveDays                     []string                                       `jsonapi:"attr,active_days,omitempty"`
	ActiveTimeType                 jsonapi.NullableAttr[string]                   `jsonapi:"attr,active_time_type,omitempty"`
	TimeZone                       string                                         `jsonapi:"attr,time_zone"`
	StartTime                      jsonapi.NullableAttr[string]                   `jsonapi:"attr,start_time,omitempty"`
	EndTime                        jsonapi.NullableAttr[string]                   `jsonapi:"attr,end_time,omitempty"`
	ScheduleRotationableAttributes ScheduleRotationScheduleRotationableAttributes `jsonapi:"attr,schedule_rotationable_attributes"`
	ActiveTimeAttributes           []ScheduleRotationActiveTimeAttributes         `jsonapi:"attr,active_time_attributes,omitempty"`
	ScheduleRotationMembers        []ScheduleRotationMember                       `jsonapi:"attr,schedule_rotation_members,omitempty"`
}

// HACK: Request expects schedule_rotation_members as an attribute, but response returns it as a relationship.
type scheduleRotationInternal struct {
	ID                      string                   `jsonapi:"primary,schedule_rotations"`
	ScheduleRotationMembers []ScheduleRotationMember `jsonapi:"relation,schedule_rotation_members,omitempty"`
}

type ScheduleRotationScheduleRotationableAttributes struct {
	HandoffTime     jsonapi.NullableAttr[string] `jsonapi:"attr,handoff_time,omitempty"`
	HandoffDay      jsonapi.NullableAttr[string] `jsonapi:"attr,handoff_day,omitempty"`
	ShiftLength     jsonapi.NullableAttr[int64]  `jsonapi:"attr,shift_length,omitempty"`
	ShiftLengthUnit jsonapi.NullableAttr[string] `jsonapi:"attr,shift_length_unit,omitempty"`
}

type ScheduleRotationActiveTimeAttributes struct {
	StartTime string `jsonapi:"attr,start_time"`
	EndTime   string `jsonapi:"attr,end_time"`
}

type ScheduleRotationMember struct {
	ID         string `jsonapi:"primary,schedule_rotation_members"`
	MemberType string `jsonapi:"attr,member_type"`
	MemberID   string `jsonapi:"attr,member_id"`
	Position   int64  `jsonapi:"attr,position"`
}

func (c *Client) GetScheduleRotation(ctx context.Context, id string) (*ScheduleRotation, error) {
	httpResp, err := c.ClientWithResponses.GetScheduleRotation(ctx, id, func(ctx context.Context, req *http.Request) error {
		q := req.URL.Query()
		q.Add("include", "schedule_rotation_members")
		req.URL.RawQuery = q.Encode()
		return nil
	})
	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	var resp ScheduleRotation
	if err := jsonapi.UnmarshalPayload(bytes.NewReader(body), &resp); err != nil {
		return nil, err
	}

	// HACK
	var intResp scheduleRotationInternal
	if err := jsonapi.UnmarshalPayload(bytes.NewReader(body), &intResp); err != nil {
		return nil, err
	}
	resp.ScheduleRotationMembers = intResp.ScheduleRotationMembers

	return &resp, nil
}

func (c *Client) CreateScheduleRotation(ctx context.Context, req ScheduleRotation) (*ScheduleRotation, error) {
	var buf bytes.Buffer
	if err := jsonapi.MarshalPayload(&buf, &req); err != nil {
		return nil, err
	}

	httpResp, err := c.ClientWithResponses.CreateScheduleRotationWithBody(ctx, req.ScheduleId, "application/vnd.api+json", bytes.NewReader(buf.Bytes()), func(ctx context.Context, req *http.Request) error {
		q := req.URL.Query()
		q.Add("include", "schedule_rotation_members")
		req.URL.RawQuery = q.Encode()
		return nil
	})
	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	var resp ScheduleRotation
	if err := jsonapi.UnmarshalPayload(bytes.NewReader(body), &resp); err != nil {
		return nil, err
	}

	// HACK
	var intResp scheduleRotationInternal
	if err := jsonapi.UnmarshalPayload(bytes.NewReader(body), &intResp); err != nil {
		return nil, err
	}
	resp.ScheduleRotationMembers = intResp.ScheduleRotationMembers

	return &resp, nil
}

func (c *Client) UpdateScheduleRotation(ctx context.Context, req ScheduleRotation) (*ScheduleRotation, error) {
	// schedule_id has RequiresReplace in the provider schema — it can never
	// change on update. Sending it in the PUT body triggers a server-side 500
	// when schedule_rotation_members is also included. Strip it before marshaling.
	req.ScheduleId = ""

	var buf bytes.Buffer
	if err := jsonapi.MarshalPayload(&buf, &req); err != nil {
		return nil, err
	}

	httpResp, err := c.ClientWithResponses.UpdateScheduleRotationWithBody(ctx, req.ID, "application/vnd.api+json", bytes.NewReader(buf.Bytes()), func(ctx context.Context, req *http.Request) error {
		q := req.URL.Query()
		q.Add("include", "schedule_rotation_members")
		req.URL.RawQuery = q.Encode()
		return nil
	})
	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	var resp ScheduleRotation
	if err := jsonapi.UnmarshalPayload(bytes.NewBuffer(body), &resp); err != nil {
		return nil, err
	}

	// HACK
	var intResp scheduleRotationInternal
	if err := jsonapi.UnmarshalPayload(bytes.NewBuffer(body), &intResp); err != nil {
		return nil, err
	}
	resp.ScheduleRotationMembers = intResp.ScheduleRotationMembers

	return &resp, nil
}
