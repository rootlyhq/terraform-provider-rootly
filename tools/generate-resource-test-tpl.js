const inflect = require('./inflect')

module.exports = (name, resourceSchema, requiredFields, pathIdField) => {
	const namePlural = inflect.pluralize(name)
	const nameCamel = inflect.camelize(name)
	const nameCamelPlural = inflect.camelize(namePlural)

return `package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResource${nameCamel}(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResource${nameCamel},
			},
		},
	})
}

const testAccResource${nameCamel} = \`
resource "rootly_${name}" "test" {
	${testParams(name, resourceSchema, requiredFields || [])}
}
\`
`
}

function testParams(name, schema, required) {
	return (required).map((key) => {
		let val = schema.properties[key].enum ? schema.properties[key].enum[0] : "test"
		switch (schema.properties[key].type) {
			case "boolean":
				return `${key} = false`
			case "string":
				return key == 'url' ? `	url = https://rootly.com/dummy\n` : `${key} = "${val}"`
			case "number":
				return `${key} = 1`
			case "array":
				if (schema.properties[key].items.type === "object" && schema.properties[key].items.properties.id) {
					return `${key} {
						id = "foo"
						name = "bar"
					}`
				}
				return `${key} = ["foo"]`
			case "object":
				return `${key} = {
					id = "foo"
					name = "bar"
				}`
		}
	}).join('\n')
}
