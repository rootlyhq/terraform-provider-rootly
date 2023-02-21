package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type PostmortemTemplate struct {
	ID string `jsonapi:"primary,post_mortem_templates"`
	Name string `jsonapi:"attr,name,omitempty"`
  Default *bool `jsonapi:"attr,default,omitempty"`
  Content string `jsonapi:"attr,content,omitempty"`
  Format string `jsonapi:"attr,format,omitempty"`
}

func (c *Client) ListPostmortemTemplates(params *rootlygo.ListPostmortemTemplatesParams) ([]interface{}, error) {
	req, err := rootlygo.NewListPostmortemTemplatesRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	post_mortem_templates, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(PostmortemTemplate)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return post_mortem_templates, nil
}

func (c *Client) CreatePostmortemTemplate(d *PostmortemTemplate) (*PostmortemTemplate, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling post_mortem_template: %s", err.Error())
	}

	req, err := rootlygo.NewCreatePostmortemTemplateRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create post_mortem_template: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(PostmortemTemplate))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling post_mortem_template: %s", err.Error())
	}

	return data.(*PostmortemTemplate), nil
}

func (c *Client) GetPostmortemTemplate(id string) (*PostmortemTemplate, error) {
	req, err := rootlygo.NewGetPostmortemTemplateRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get post_mortem_template: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(PostmortemTemplate))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling post_mortem_template: %s", err.Error())
	}

	return data.(*PostmortemTemplate), nil
}

func (c *Client) UpdatePostmortemTemplate(id string, post_mortem_template *PostmortemTemplate) (*PostmortemTemplate, error) {
	buffer, err := MarshalData(post_mortem_template)
	if err != nil {
		return nil, errors.Errorf("Error marshaling post_mortem_template: %s", err.Error())
	}

	req, err := rootlygo.NewUpdatePostmortemTemplateRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update post_mortem_template: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(PostmortemTemplate))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling post_mortem_template: %s", err.Error())
	}

	return data.(*PostmortemTemplate), nil
}

func (c *Client) DeletePostmortemTemplate(id string) error {
	req, err := rootlygo.NewDeletePostmortemTemplateRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete post_mortem_template: %s", err.Error())
	}

	return nil
}
