const inflect = require("./inflect");

module.exports = (name, resourceSchema, pathIdField) => {
  const namePlural = inflect.pluralize(name);
  const nameCamel = inflect.camelize(name);
  const nameCamelPlural = inflect.camelize(namePlural);
  const strconvImport =
    pathIdField && resourceSchema.properties[pathIdField].type === "number"
      ? `"strconv"`
      : "";

  return `package client

import (
	"reflect"
	${strconvImport}
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type ${nameCamel} struct {
	ID string \`jsonapi:"primary,${name === "team" ? "groups" : namePlural}"\`
	${structAttrs(resourceSchema)}
}

func (c *Client) List${nameCamelPlural}(${listFnParams(
    nameCamelPlural,
    pathIdField
  )}) ([]interface{}, error) {
	req, err := rootlygo.NewList${nameCamelPlural}Request(${listClientParams(
    pathIdField
  )})
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	${namePlural}, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(${nameCamel})))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return ${namePlural}, nil
}

`;
};

function listFnParams(nameCamelPlural, nested) {
  if (nested) {
    return `id string, params *rootlygo.List${nameCamelPlural}Params`;
  } else {
    return `params *rootlygo.List${nameCamelPlural}Params`;
  }
}

function listClientParams(nested) {
  if (nested) {
    return `c.Rootly.Server, id, params`;
  } else {
    return `c.Rootly.Server, params`;
  }
}

function structAttr(name, resourceSchema) {
  const schema = resourceSchema.properties[name];
  switch (schema.type) {
    case "string":
      return `${inflect.camelize(
        name
      )} string \`jsonapi:"attr,${name},omitempty"\``;
    case "integer":
    case "number":
      return `${inflect.camelize(
        name
      )} int \`jsonapi:"attr,${name},omitempty"\``;
    case "boolean":
      return `${inflect.camelize(
        name
      )} *bool \`jsonapi:"attr,${name},omitempty"\``;
    case "array":
      return `${inflect.camelize(
        name
      )} []interface{} \`jsonapi:"attr,${name},omitempty"\``;
    case "object":
    default:
      return `${inflect.camelize(
        name
      )} map[string]interface{} \`jsonapi:"attr,${name},omitempty"\``;
  }
}

function structAttrs(resourceSchema) {
  return Object.keys(resourceSchema.properties)
    .filter((name) => {
      return name !== "created_at" && name !== "updated_at";
    })
    .map((name) => {
      return structAttr(name, resourceSchema);
    })
    .join("\n  ");
}
