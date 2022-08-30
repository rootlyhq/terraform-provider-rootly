const inflect = require('inflect');

module.exports = (name, resourceSchema, collectionSchema, pathIdField) => {
	const namePlural = inflect.pluralize(name)
	const nameCamel = inflect.camelize(name)
	const nameCamelPlural = inflect.camelize(namePlural)
	const strconvImport = pathIdField && resourceSchema.properties[pathIdField].type === 'number' ? `"strconv"` : ''

return `package provider

import (
	"context"
	${strconvImport}
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

func dataSource${nameCamel}() *schema.Resource{
	return &schema.Resource{
		ReadContext: dataSource${nameCamel}Read,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			${schemaFields(resourceSchema)}
		},
	}
}

func dataSource${nameCamel}Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.List${nameCamelPlural}Params)
	page_size := 1
	params.PageSize = &page_size

	${setFilterFields(name, resourceSchema, collectionSchema)}

	${listFn(nameCamelPlural, resourceSchema, pathIdField)}
	if err != nil {
		return diag.FromErr(err)
	}

	item, _ := items[0].(*client.${nameCamel})

	d.SetId(item.ID)

	return nil
}
`}

function listFn(nameCamelPlural, resourceSchema, pathIdField) {
	if (pathIdField) {
		if (resourceSchema.properties[pathIdField].type === 'number') {
			return `${pathIdField} := strconv.Itoa(d.Get("${pathIdField}").(int))
			items, err := c.List${nameCamelPlural}(${pathIdField}, params)`
		} else {
			return `${pathIdField} := d.Get("${pathIdField}").(string)
			items, err := c.List${nameCamelPlural}(${pathIdField}, params)`
		}
	} else {
		return `items, err := c.List${nameCamelPlural}(params)`
	}
}

function setFilterFields(name, resourceSchema, collectionSchema) {
	return (collectionSchema.get.parameters || []).filter((paramSchema) => {
		return paramSchema.name.match(/^filter/)
	}).map((paramSchema) => {
		const filterField = inflect.underscore(paramSchema.name.replace("filter[", "").replace("]", ""))
		const fieldSchema = resourceSchema.properties[filterField]
		// TODO remove
		if (fieldSchema) {
			return `
				${filterField} := d.Get("${filterField}").(${jsonapiToGoType(fieldSchema.type)})
				params.Filter${inflect.camelize(filterField)} = &${filterField}
			`
		}
	}).filter((x) => x).join('\n')
}

function jsonapiToGoType(type) {
	switch (type) {
		case 'string':
			return 'string'
		case 'number':
			return 'int'
		case 'boolean':
			return 'bool'
		default:
			return 'string'
	}
}

function schemaFields(resourceSchema) {
	return Object.keys(resourceSchema.properties).map((field) => {
		return schemaField(field, resourceSchema)
	}).join('\n')
}

function schemaField(name, resourceSchema) {
	const schema = resourceSchema.properties[name]
	const optional = (resourceSchema.required || []).indexOf(name) === -1 ? "true" : "false"
	switch (schema.type) {
		case 'string':
			return `
			"${name}": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: ${optional},
			},
			`
		case 'number':
			return `
			"${name}": &schema.Schema{
				Type: schema.TypeInt,
				Computed: true,
				Optional: ${optional},
			},
			`
		case 'boolean':
			return `
			"${name}": &schema.Schema{
				Type: schema.TypeBool,
				Computed: true,
				Optional: ${optional},
			},
			`
		case 'array':
			if (schema.items && schema.items.type === "object") {
				return `
				"${name}": &schema.Schema{
					Type: schema.TypeList,
					Computed: true,
					Optional: ${optional},
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"id": &schema.Schema{
								Type: schema.TypeString,
								Required: true,
							},
							"name": &schema.Schema{
								Type: schema.TypeString,
								Required: true,
							},
						},
					},
				},
				`
			} else {
				return `
				"${name}": &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Computed: true,
					Optional: ${optional},
				},
				`
			}
		case 'object':
		default:
			return `
			"${name}": &schema.Schema{
				Type: schema.TypeMap,
				Elem: schema.TypeString,
				Computed: true,
				Optional: ${optional},
			},
			`
	}
}
