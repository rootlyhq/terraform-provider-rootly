package apiclient

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/DataDog/jsonapi"
	"github.com/oapi-codegen/nullable"
)

type ScheduleRotation struct {
	ID                             string                                         `jsonapi:"primary,schedule_rotations"`
	ScheduleId                     string                                         `jsonapi:"attr" json:"schedule_id"`
	Name                           string                                         `jsonapi:"attr" json:"name"`
	Position                       int64                                          `jsonapi:"attr" json:"position,omitempty"`
	ScheduleRotationableType       string                                         `jsonapi:"attr" json:"schedule_rotationable_type"`
	ActiveAllWeek                  nullable.Nullable[bool]                        `jsonapi:"attr" json:"active_all_week,omitempty"`
	ActiveDays                     []string                                       `jsonapi:"attr" json:"active_days,omitempty"`
	ActiveTimeType                 nullable.Nullable[string]                      `jsonapi:"attr" json:"active_time_type,omitempty"`
	TimeZone                       string                                         `jsonapi:"attr" json:"time_zone"`
	StartTime                      nullable.Nullable[string]                      `jsonapi:"attr" json:"start_time,omitempty"`
	EndTime                        nullable.Nullable[string]                      `jsonapi:"attr" json:"end_time,omitempty"`
	ScheduleRotationableAttributes ScheduleRotationScheduleRotationableAttributes `jsonapi:"attr" json:"schedule_rotationable_attributes"`
	ActiveTimeAttributes           []ScheduleRotationActiveTimeAttributes         `jsonapi:"attr" json:"active_time_attributes,omitempty"`
	ScheduleRotationMembers        []ScheduleRotationMember                       `jsonapi:"attr" json:"schedule_rotation_members,omitempty"`
}

// HACK: Request expects schedule_rotation_members as an attribute, but response returns it as a relationship.
type scheduleRotationInternal struct {
	ID                      string                   `jsonapi:"primary,schedule_rotations"`
	ScheduleRotationMembers []ScheduleRotationMember `jsonapi:"rel" json:"schedule_rotation_members,omitempty"`
}

type ScheduleRotationScheduleRotationableAttributes struct {
	HandoffTime     nullable.Nullable[string] `jsonapi:"attr" json:"handoff_time,omitempty"`
	HandoffDay      nullable.Nullable[string] `jsonapi:"attr" json:"handoff_day,omitempty"`
	ShiftLength     nullable.Nullable[int64]  `jsonapi:"attr" json:"shift_length,omitempty"`
	ShiftLengthUnit nullable.Nullable[string] `jsonapi:"attr" json:"shift_length_unit,omitempty"`
}

type ScheduleRotationActiveTimeAttributes struct {
	StartTime string `jsonapi:"attr" json:"start_time"`
	EndTime   string `jsonapi:"attr" json:"end_time"`
}

type ScheduleRotationMember struct {
	ID         string `jsonapi:"primary,schedule_rotation_members"`
	MemberType string `jsonapi:"attr" json:"member_type"`
	MemberID   string `jsonapi:"attr" json:"member_id"`
	Position   int64  `jsonapi:"attr" json:"position"`
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

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	var resp ScheduleRotation
	if err := jsonapi.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	// HACK
	var intResp scheduleRotationInternal
	if err := jsonapi.Unmarshal(body, &intResp); err != nil {
		return nil, err
	}
	resp.ScheduleRotationMembers = intResp.ScheduleRotationMembers

	return &resp, nil
}

func (c *Client) CreateScheduleRotation(ctx context.Context, req ScheduleRotation) (*ScheduleRotation, error) {
	body, err := jsonapi.Marshal(&req, jsonapi.MarshalClientMode())
	if err != nil {
		return nil, err
	}

	httpResp, err := c.ClientWithResponses.CreateScheduleRotationWithBody(ctx, req.ScheduleId, "application/vnd.api+json", bytes.NewReader(body), func(ctx context.Context, req *http.Request) error {
		q := req.URL.Query()
		q.Add("include", "schedule_rotation_members")
		req.URL.RawQuery = q.Encode()
		return nil
	})
	if err != nil {
		return nil, err
	}

	body, err = io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	var resp ScheduleRotation
	if err := jsonapi.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	// HACK
	var intResp scheduleRotationInternal
	if err := jsonapi.Unmarshal(body, &intResp); err != nil {
		return nil, err
	}
	resp.ScheduleRotationMembers = intResp.ScheduleRotationMembers

	return &resp, nil
}

func (c *Client) UpdateScheduleRotation(ctx context.Context, req ScheduleRotation) (*ScheduleRotation, error) {
	body, err := jsonapi.Marshal(&req)
	if err != nil {
		return nil, err
	}

	httpResp, err := c.ClientWithResponses.UpdateScheduleRotationWithBody(ctx, req.ID, "application/vnd.api+json", bytes.NewReader(body), func(ctx context.Context, req *http.Request) error {
		q := req.URL.Query()
		q.Add("include", "schedule_rotation_members")
		req.URL.RawQuery = q.Encode()
		return nil
	})
	if err != nil {
		return nil, err
	}

	body, err = io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	var resp ScheduleRotation
	if err := jsonapi.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	// HACK
	var intResp scheduleRotationInternal
	if err := jsonapi.Unmarshal(body, &intResp); err != nil {
		return nil, err
	}
	resp.ScheduleRotationMembers = intResp.ScheduleRotationMembers

	return &resp, nil
}
