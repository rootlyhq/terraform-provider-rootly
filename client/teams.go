package client

import (
	"bytes"
	"github.com/google/jsonapi"
	"github.com/pkg/errors"
	rootlygo "github.com/rootlyhq/rootly-go"
)

type Team struct {
	ID          string `jsonapi:"primary,groups"`
	Name        string `jsonapi:"attr,name,omitempty"`
	Description string `jsonapi:"attr,description,omitempty"`
}

func (t Team) Marshal() (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	if err := jsonapi.MarshalPayload(buffer, t); err != nil {
		return nil, errors.Errorf("Error marshaling team (creation): %s", err.Error())
	}

	return buffer, nil
}

func (c *Client) CreateTeam(t *Team) (*Team, error) {
	buffer := new(bytes.Buffer)
	if err := jsonapi.MarshalPayload(buffer, t); err != nil {
		return nil, errors.Errorf("Error marshaling team (creation): %s", err.Error())
	}

	req, err := rootlygo.NewCreateTeamRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling team (creation): %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to create team" + "\n\n" + buffer.String())
	}

	team := new(Team)
	if err := jsonapi.UnmarshalPayload(resp.Body, team); err != nil {
		return nil, errors.Errorf("Error unmarshaling team (creation): %s", err.Error())
	}

	return team, nil
}

func (c *Client) GetTeam(id string) (*Team, error) {
	req, err := rootlygo.NewGetTeamRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get team: %s", id)
	}

	team := new(Team)
	if err := jsonapi.UnmarshalPayload(resp.Body, team); err != nil {
		return nil, errors.Errorf("Error unmarshaling team (read): %s", err.Error())
	}

	return team, nil
}

func (c *Client) UpdateTeam(id string, t *Team) (*Team, error) {
	buffer := new(bytes.Buffer)
	if err := jsonapi.MarshalPayload(buffer, t); err != nil {
		return nil, errors.Errorf("Error marshaling team (update): %s", err.Error())
	}

	req, err := rootlygo.NewUpdateTeamRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling team (update): %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update team: %s", id)
	}

	team := new(Team)
	if err := jsonapi.UnmarshalPayload(resp.Body, team); err != nil {
		return nil, errors.Errorf("Error unmarshaling team (update): %s", err.Error())
	}

	return team, nil
}

func (c *Client) DeleteTeam(id string) error {
	req, err := rootlygo.NewDeleteTeamRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error unmarshaling team (delete): %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete team: %s", id)
	}

	return nil
}
