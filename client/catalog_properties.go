package client

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v5/schema"
)

type CatalogProperty struct {
	ID            string `jsonapi:"primary,catalog_properties"`
	CatalogId     string `jsonapi:"attr,catalog_id,omitempty"`
	Name          string `jsonapi:"attr,name,omitempty"`
	Slug          string `jsonapi:"attr,slug,omitempty"`
	Kind          string `jsonapi:"attr,kind,omitempty"`
	KindCatalogId string `jsonapi:"attr,kind_catalog_id,omitempty"`
	Multiple      *bool  `jsonapi:"attr,multiple,omitempty"`
	Position      int    `jsonapi:"attr,position,omitempty"`
	Required      *bool  `jsonapi:"attr,required,omitempty"`
	CatalogType   string `jsonapi:"attr,catalog_type,omitempty"`
}

// catalogPropertyCollectionPath returns the API path for listing/creating catalog properties
// based on the catalog_type. For native types (service, team, etc.) the path is
// /v1/{plural_type}/properties. For custom catalogs it is /v1/catalogs/{catalog_id}/properties.
func catalogPropertyCollectionPath(catalogType string, catalogId string) (string, error) {
	switch catalogType {
	case "service":
		return "/v1/services/properties", nil
	case "functionality":
		return "/v1/functionalities/properties", nil
	case "team":
		return "/v1/teams/properties", nil
	case "environment":
		return "/v1/environments/properties", nil
	case "cause":
		return "/v1/causes/properties", nil
	case "incident_type":
		return "/v1/incident_types/properties", nil
	case "catalog", "":
		if catalogId == "" {
			return "", fmt.Errorf("catalog_id is required when catalog_type is 'catalog'")
		}
		return fmt.Sprintf("/v1/catalogs/%s/properties", url.PathEscape(catalogId)), nil
	default:
		return "", fmt.Errorf("unsupported catalog_type: %s", catalogType)
	}
}

func (c *Client) ListCatalogProperties(catalogType string, catalogId string) ([]interface{}, error) {
	collectionPath, err := catalogPropertyCollectionPath(catalogType, catalogId)
	if err != nil {
		return nil, fmt.Errorf("Error building request path: %w", err)
	}

	serverURL, err := url.Parse(c.Rootly.Server)
	if err != nil {
		return nil, fmt.Errorf("Error parsing server URL: %w", err)
	}

	operationPath := "." + collectionPath
	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, fmt.Errorf("Error building request URL: %w", err)
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request: %w", err)
	}

	catalog_properties, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(CatalogProperty)))
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling: %w", err)
	}

	return catalog_properties, nil
}

func (c *Client) CreateCatalogProperty(d *CatalogProperty) (*CatalogProperty, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling catalog_property: %w", err)
	}

	collectionPath, err := catalogPropertyCollectionPath(d.CatalogType, d.CatalogId)
	if err != nil {
		return nil, fmt.Errorf("Error building request path: %w", err)
	}

	serverURL, err := url.Parse(c.Rootly.Server)
	if err != nil {
		return nil, fmt.Errorf("Error parsing server URL: %w", err)
	}

	operationPath := "." + collectionPath
	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, fmt.Errorf("Error building request URL: %w", err)
	}

	req, err := http.NewRequest("POST", queryURL.String(), buffer)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}
	req.Header.Add("Content-Type", c.ContentType)

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to perform request to create catalog_property: %s", err)
	}

	data, err := UnmarshalData(resp.Body, new(CatalogProperty))
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling catalog_property: %w", err)
	}

	return data.(*CatalogProperty), nil
}

func (c *Client) GetCatalogProperty(id string) (*CatalogProperty, error) {
	req, err := rootlygo.NewGetCatalogPropertyRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request to get catalog_property: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(CatalogProperty))
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling catalog_property: %w", err)
	}

	return data.(*CatalogProperty), nil
}

func (c *Client) UpdateCatalogProperty(id string, catalog_property *CatalogProperty) (*CatalogProperty, error) {
	buffer, err := MarshalData(catalog_property)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling catalog_property: %w", err)
	}

	req, err := rootlygo.NewUpdateCatalogPropertyRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request to update catalog_property: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(CatalogProperty))
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling catalog_property: %w", err)
	}

	return data.(*CatalogProperty), nil
}

func (c *Client) DeleteCatalogProperty(id string) error {
	req, err := rootlygo.NewDeleteCatalogPropertyRequest(c.Rootly.Server, id)
	if err != nil {
		return fmt.Errorf("Error building request: %w", err)
	}

	_, err = c.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to make request to delete catalog_property: %w", err)
	}

	return nil
}
