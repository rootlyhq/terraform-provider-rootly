const fs = require("fs")
const path = require('path')
const swaggerPath = process.argv[2]
const inflect = require('./inflect')
const providerTpl = require('./generate-provider-tpl')
const clientTpl = require('./generate-client-tpl')
const dataSourceTpl = require('./generate-data-source-tpl')
const resourceTpl = require('./generate-resource-tpl')
const resourceTestTpl = require('./generate-resource-test-tpl')
const workflowTpl = require('./generate-workflow-tpl')
const generateTasks = require('./generate-tasks')
const swagger = require(path.resolve(swaggerPath))

const excluded = [
	"dashboard",
	"dashboard_panel",
	"workflow_task",
	"workflow_run",
	"post_mortem_template",
	"incident",
	"incident_post_mortem",
	"incident_action_item",
	"incident_event",
	"incident_feedback",
	"incident_custom_field_selection",
	"pulse",
	"alert",
	"playbook_task",
]

console.log(`Excluding resource from generation:`, excluded)

const resources = Object.keys(swagger.components.schemas).filter((name) => {
	return excluded.indexOf(name) === -1 && collectionPathSchema(name)
})
const dataSources = resources.filter(resourceHasFilters).filter((name) => name !== "workflow")
const taskResources = generateTasks(swagger)

generateProvider(resources, taskResources, dataSources)

resources.forEach((name) => {
	generateClient(name)
	generateResource(name)
	generateResourceTest(name)
})

dataSources.forEach(generateDataSource)

function generateProvider(resources, taskResources, dataSources) {
	const code = providerTpl(resources.filter((name) => name !== 'workflow'), taskResources, dataSources)
	fs.writeFileSync(path.resolve(__dirname, '..', 'provider', 'provider.go'), code)
}

function generateClient(name) {
	const collectionSchema = collectionPathSchema(name)
	const pathIdField = collectionSchema && collectionSchema.parameters && collectionSchema.parameters[0] && collectionSchema.parameters[0].name
	const code = clientTpl(name, resourceSchema(name), pathIdField)
	fs.writeFileSync(path.resolve(__dirname, '..', 'client', `${inflect.pluralize(name)}.go`), code)
}

function resourceHasFilters(name) {
	const collectionSchema = collectionPathSchema(name)
	const filterParameters = collectionSchema.get && collectionSchema.get.parameters
	return filterParameters.filter((filter) => filter.name.match(/filter/i)).length
}

function generateDataSource(name) {
	const collectionSchema = collectionPathSchema(name)
	const filterParameters = collectionSchema.get && collectionSchema.get.parameters
	const pathIdField = collectionSchema && collectionSchema.parameters && collectionSchema.parameters[0] && collectionSchema.parameters[0].name
	const code = dataSourceTpl(name, resourceSchema(name), filterParameters, pathIdField)
	if (code) {
		fs.writeFileSync(path.resolve(__dirname, '..', 'provider', `data_source_${name}.go`), code)
	}
}

function generateResource(name) {
	const collectionSchema = collectionPathSchema(name)
	const pathIdField = collectionSchema && collectionSchema.parameters && collectionSchema.parameters[0] && collectionSchema.parameters[0].name
	let code;
	if (name === "workflow") {
		code = workflowTpl("workflow_incident", resourceSchema(name), requiredFields(name), swagger.components.schemas.incident_trigger_params)
		fs.writeFileSync(path.resolve(__dirname, '..', 'provider', `resource_${name}_incident.go`), code)
		code = workflowTpl("workflow_post_mortem", resourceSchema(name), requiredFields(name), swagger.components.schemas.post_mortem_trigger_params)
		fs.writeFileSync(path.resolve(__dirname, '..', 'provider', `resource_${name}_post_mortem.go`), code)
		code = workflowTpl("workflow_action_item", resourceSchema(name), requiredFields(name), swagger.components.schemas.action_item_trigger_params)
		fs.writeFileSync(path.resolve(__dirname, '..', 'provider', `resource_${name}_action_item.go`), code)
		code = workflowTpl("workflow_alert", resourceSchema(name), requiredFields(name), swagger.components.schemas.alert_trigger_params)
		fs.writeFileSync(path.resolve(__dirname, '..', 'provider', `resource_${name}_alert.go`), code)
		code = workflowTpl("workflow_pulse", resourceSchema(name), requiredFields(name), swagger.components.schemas.pulse_trigger_params)
		fs.writeFileSync(path.resolve(__dirname, '..', 'provider', `resource_${name}_pulse.go`), code)
	} else {
		code = resourceTpl(name, resourceSchema(name), requiredFields(name), pathIdField)
		fs.writeFileSync(path.resolve(__dirname, '..', 'provider', `resource_${name}.go`), code)
	}
}

function generateResourceTest(name) {
	const collectionSchema = collectionPathSchema(name)
	const pathIdField = collectionSchema && collectionSchema.parameters && collectionSchema.parameters[0] && collectionSchema.parameters[0].name
	if (name !== "workflow" && !pathIdField) {
		code = resourceTestTpl(name, resourceSchema(name), requiredFields(name), pathIdField)
		fs.writeFileSync(path.resolve(__dirname, '..', 'provider', `resource_${name}_test.go`), code)
	}
}

function resourceSchema(name) {
	return swagger.components.schemas[name]
}

function requiredFields(name) {
	return swagger.components.schemas[`new_${name}`].properties.data.properties.attributes.required;
}

function collectionPathSchema(name) {
	return Object.keys(swagger.paths).filter((url) => {
		const get = swagger.paths[url].get
		return get && get.operationId.replace(/ /g, '') === `list${inflect.pluralize(inflect.camelize(name))}`
	}).map((url) => swagger.paths[url])[0]
}
