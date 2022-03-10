package client

import (
	"github.com/pkg/errors"
	rootlygo "github.com/rootlyhq/rootly-go"
)

type Team struct {
	ID          string `jsonapi:"primary,groups"`
	Name        string `jsonapi:"attr,name,omitempty"`
	Description string `jsonapi:"attr,description,omitempty"`
	Color       string `jsonapi:"attr,color,omitempty"`
}

func (c *Client) CreateTeam(t *Team) (*Team, error) {
	buffer, err := MarshalData(t)
	if err != nil {
		return nil, errors.Errorf("Error marshaling team: %s", err.Error())
	}

	req, err := rootlygo.NewCreateTeamRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to create team: %s", t.ID)
	}

	team, err := UnmarshalData(resp.Body, new(Team))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling team: %s", err.Error())
	}

	return team.(*Team), nil
}

func (c *Client) GetTeam(id string) (*Team, error) {
	req, err := rootlygo.NewGetTeamRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get team: %s", id)
	}

	team, err := UnmarshalData(resp.Body, new(Team))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling team: %s", err.Error())
	}

	return team.(*Team), nil
}

func (c *Client) UpdateTeam(id string, t *Team) (*Team, error) {
	buffer, err := MarshalData(t)
	if err != nil {
		return nil, errors.Errorf("Error marshaling team: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateTeamRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update team: %s", id)
	}

	team, err := UnmarshalData(resp.Body, new(Team))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling team: %s", err.Error())
	}

	return team.(*Team), nil
}

func (c *Client) DeleteTeam(id string) error {
	req, err := rootlygo.NewDeleteTeamRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete team: %s", id)
	}

	return nil
}
