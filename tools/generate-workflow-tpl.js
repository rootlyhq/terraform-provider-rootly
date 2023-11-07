const inflect = require("inflect");

function includeTools(resourceSchema) {
  for (var key in resourceSchema.properties) {
    if (resourceSchema.properties[key].type === "boolean") {
      return true;
    }
    if (resourceSchema.properties[key].type === "array" && resourceSchema.properties[key].items?.type !== "object") {
      return true;
    }
  }
  return false;
}

module.exports = (name, resourceSchema, requiredFields, taskParamsSchema) => {
  // convert {anyOf: [{enum: [null]}, {type: "string"}]} to {type: "string"}
  Object.keys(taskParamsSchema.properties).forEach((key) => {
    const propSchema = taskParamsSchema.properties[key];
    if (propSchema.anyOf) {
      taskParamsSchema.properties[key] = propSchema.anyOf
        .filter((x) => x.type)
        .pop();
      delete taskParamsSchema.properties[key].enum;
    }
  });
  const namePlural = inflect.pluralize(name);
  const nameCamel = inflect.camelize(name);
  const nameCamelPlural = inflect.camelize(namePlural);
  const tools = includeTools(resourceSchema)
    ? `"github.com/rootlyhq/terraform-provider-rootly/tools"`
    : "";

  return `package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	${tools}
)

func resource${nameCamel}() *schema.Resource {
	return &schema.Resource{
		CreateContext: resource${nameCamel}Create,
		ReadContext: resource${nameCamel}Read,
		UpdateContext: resource${nameCamel}Update,
		DeleteContext: resource${nameCamel}Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema {
			${schemaFields(resourceSchema, requiredFields, taskParamsSchema)}
		},
	}
}

func resource${nameCamel}Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating ${nameCamel}"))

	s := &client.Workflow{}

	${createResourceFields(resourceSchema)}

	res, err := c.CreateWorkflow(s)
	if err != nil {
		return diag.Errorf("Error creating ${name}: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a ${name} resource: %s", d.Id()))

	return resource${nameCamel}Read(ctx, d, meta)
}

func resource${nameCamel}Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading ${nameCamel}: %s", d.Id()))

	item, err := c.GetWorkflow(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("${nameCamel} (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading ${name}: %s", d.Id())
	}

	${setResourceFields(resourceSchema)}

	return nil
}

func resource${nameCamel}Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating ${nameCamel}: %s", d.Id()))

	s := &client.Workflow{}

	${updateResourceFields(resourceSchema)}

	_, err := c.UpdateWorkflow(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating ${name}: %s", err.Error())
	}

	return resource${nameCamel}Read(ctx, d, meta)
}

func resource${nameCamel}Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting ${nameCamel}: %s", d.Id()))

	err := c.DeleteWorkflow(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("${nameCamel} (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting ${name}: %s", err.Error())
	}

	d.SetId("")

	return nil
}
`;
};

function excludeDateFields(field) {
  return field !== "created_at" && field !== "updated_at";
}

function setResourceFields(resourceSchema) {
  return Object.keys(resourceSchema.properties)
    .filter(excludeDateFields)
    .map((field) => {
      if (field === "trigger_params") {
        return `
        tps := make([]interface{}, 1, 1)
        tps[0] = item.TriggerParams
        d.Set("trigger_params", tps)
			`;
      } else {
        return `d.Set("${field}", item.${inflect.camelize(field)})`;
      }
    })
    .join("\n  ");
}

function createResourceFields(resourceSchema) {
  return Object.keys(resourceSchema.properties)
    .filter(excludeDateFields)
    .map((field) => {
      const schema = resourceSchema.properties[field];
      if (field === "trigger_params") {
        return `  if value, ok := d.GetOkExists("${field}"); ok {
			s.TriggerParams = value.([]interface{})[0].(map[string]interface{})
	}`;
      } else if (schema.type === "boolean") {
        return `  if value, ok := d.GetOkExists("${field}"); ok {
				s.${inflect.camelize(field)} = tools.Bool(value.(${jsonapiToGoType(
          schema.type
        )}))
			}`;
      } else {
        return `  if value, ok := d.GetOkExists("${field}"); ok {
			s.${inflect.camelize(field)} = value.(${jsonapiToGoType(schema.type)})
	}`;
      }
    })
    .join("\n  ");
}

function updateResourceFields(resourceSchema) {
  return Object.keys(resourceSchema.properties)
    .filter(excludeDateFields)
    .map((field) => {
      const schema = resourceSchema.properties[field];
      if (field === "trigger_params") {
        return `
				if d.HasChange("trigger_params") {
					tps := d.Get("trigger_params").([]interface{})
					for _, tpsi := range tps {
						s.TriggerParams = tpsi.(map[string]interface{})
					}
				}`;
      } else if (schema.type === "boolean") {
        return `  if d.HasChange("${field}") {
				s.${inflect.camelize(field)} = tools.Bool(d.Get("${field}").(${jsonapiToGoType(
          schema.type
        )}))
			}`;
      } else {
        return `  if d.HasChange("${field}") {
				s.${inflect.camelize(field)} = d.Get("${field}").(${jsonapiToGoType(
          schema.type
        )})
			}`;
      }
    })
    .join("\n  ");
}

function jsonapiToGoType(type) {
  switch (type) {
    case "string":
      return "string";
    case "integer":
    case "number":
      return "int";
    case "boolean":
      return "bool";
    case "array":
      return "[]interface{}";
    case "object":
      return "map[string]interface{}";
    default:
      return "map[string]interface{}";
  }
}

function schemaFields(resourceSchema, requiredFields, taskParamsSchema) {
  return Object.keys(resourceSchema.properties)
    .filter(excludeDateFields)
    .map((field) => {
      return schemaField(
        field,
        resourceSchema,
        requiredFields,
        taskParamsSchema
      );
    })
    .join("\n");
}

function annotatedDescription(schema) {
  const description = (schema.description || (schema.items && schema.items.description) || "").replace(/"/g, '\\"');
  if (schema.enum) {
    return `${
      !!description ? `${description}. ` : ""
    }Value must be one off ${schema.enum
      .map((val) => `\`${val}\``)
      .join(", ")}.`;
  }
  if (
    schema.type === "object" &&
    schema.properties.id &&
    schema.properties.name
  ) {
    return `Map must contain two fields, \`id\` and \`name\`. ${description}`;
  }
  if (schema.type === "array" && schema.items && schema.items.enum) {
    return `${
      !!description ? `${description}. ` : ""
    }Value must be one of ${schema.items.enum
      .map((val) => `\`${val}\``)
      .join(", ")}.`;
  }
  return description;
}

function schemaField(name, resourceSchema, requiredFields, taskParamsSchema) {
  const schema = resourceSchema.properties[name];
  const optional =
    (requiredFields || []).indexOf(name) === -1 || schema.enum
      ? "true"
      : "false";
  const required =
    (requiredFields || []).indexOf(name) === -1 || schema.enum
      ? "false"
      : "true";
  let defaultValue;
  if (schema.default) {
    defaultValue = `Default: "${schema.default}"`;
  } else if (schema.enum && schema.enum.length > 0) {
    defaultValue = `Default: "${schema.enum[0]}"`;
  } else {
    defaultValue = `Computed: ${optional}`;
  }
  const description = annotatedDescription(schema);
  switch (schema.type) {
    case "string":
      return `
			"${name}": &schema.Schema {
				Type: schema.TypeString,
				${defaultValue},
				Required: ${required},
				Optional: ${optional},
				Description: "${description}",
			},
			`;
    case "number":
      return `
			"${name}": &schema.Schema {
				Type: schema.TypeInt,
				Computed: ${optional},
				Required: ${required},
				Optional: ${optional},
				Description: "${description}",
			},
			`;
    case "boolean":
      if (name === "enabled") {
        return `
				"${name}": &schema.Schema {
					Type: schema.TypeBool,
					Default: true,
					Optional: true,
				},
				`;
      }
      return `
			"${name}": &schema.Schema {
				Type: schema.TypeBool,
				Computed: ${optional},
				Required: ${required},
				Optional: ${optional},
				Description: "${description}",
			},
			`;
    case "array":
      if (schema.items && schema.items.type === "object") {
        return `
				"${name}": &schema.Schema {
					Type: schema.TypeList,
					Computed: ${optional},
					Required: ${required},
					Optional: ${optional},
					Description: "${description}",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema {
							"id": &schema.Schema {
								Type: schema.TypeString,
								Required: true,
							},
							"name": &schema.Schema {
								Type: schema.TypeString,
								Required: true,
							},
						},
					},
				},
				`;
      } else {
        return `
				"${name}": &schema.Schema {
					Type: schema.TypeList,
					Elem: &schema.Schema {
						Type: schema.TypeString,
					},
					DiffSuppressFunc: tools.EqualIgnoringOrder,
					Computed: ${optional},
					Required: ${required},
					Optional: ${optional},
					Description: "${description}",
				},
				`;
      }
    case "object":
    default:
      if (name === "trigger_params") {
        return `
				"${name}": &schema.Schema {
					Type: schema.TypeList,
					MinItems: 1,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema {
							${schemaFields(taskParamsSchema, taskParamsSchema.required)}
						},
					},
					Computed: true,
					Optional: true,
				},
				`;
      }
      return `
			"${name}": &schema.Schema {
				Type: schema.TypeMap,
				Elem: &schema.Schema {
					Type: schema.TypeString,
				},
				Computed: ${optional},
				Required: ${required},
				Optional: ${optional},
				Description: "${description}",
			},
			`;
  }
}
