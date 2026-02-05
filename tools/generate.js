const fs = require("fs");
const path = require("path");
const swaggerPath = process.argv[2];
const filterResource = process.argv[3] || null;
const inflect = require("./inflect");
const providerTpl = require("./generate-provider-tpl");
const clientTpl = require("./generate-client-tpl");
const clientReadOnlyTpl = require("./generate-read-only-client-tpl");
const dataSourceTpl = require("./generate-data-source-tpl");
const resourceTpl = require("./generate-resource-tpl");
const workflowTpl = require("./generate-workflow-tpl");
const generateWorkflowTaskResources = require("./generate-tasks");
const v2MigrationUtils = require("./v2-migration-utils");
const swagger = v2MigrationUtils.SchemaMods.reduce((schema, mod) => mod(schema), require(path.resolve(swaggerPath)));

const excluded = {
  dataSources: [
    "alert",
    "alert_event",
    "alerts_source", // cannot auto-generate because of plural/singular filter mismatch (filter[source_types] vs source_type)
    "audit",
    "catalog",
    "catalog_field",
    "catalog_entity",
    "catalog_entity_property",
    "communications_group",
    "custom_field_option",
    "custom_field",
    "dashboard",
    "incident_action_item",
    "incident_custom_field_selection",
    "incident_event_functionality",
    "incident_event_service",
    "incident_event",
    "incident_feedback",
    "incident_form_field_selection",
    "ip_ranges",
    "post_mortem_template",
    "pulse",
    "retrospective_configuration",
    "retrospective_process",
    "retrospective_step",
    "schedule", // cannot auto-generate because of schema upgrade logic
    "schedule_rotation",
    "secret",
    "shift",
    "user_notification_rule",
    "webhooks_delivery",
    "workflow_run",
    "workflow_task",
  ],
  resources: [
    "alert",
    "alert_event",
    "alert_route",
    "alert_group",
    "alerts_source",
    "audit",
    "catalog",
    "catalog_field",
    "catalog_entity",
    "catalog_entity_property",
    "communications_group",
    "communications_template", // cannot auto-generate because of custom nested JSON:API format handling (IR-3529)
    "custom_field_option",
    "custom_field",
    "dashboard",
    "escalation_path",
    "escalation_policy",
    "incident_action_item",
    "incident_custom_field_selection",
    "incident_event_functionality",
    "incident_event_service",
    "incident_event",
    "incident_feedback",
    "incident_form_field_selection",
    "incident_post_mortem",
    "incident",
    "ip_ranges",
    "live_call_router",
    "on_call_role",
    "override_shift",
    "post_mortem_template",
    "pulse",
    "retrospective_configuration",
    "retrospective_process",
    "retrospective_step",
    "secret",
    "schedule", // cannot auto-generate because of schema upgrade logic
    "schedule_rotation",
    "shift",
    "team",
    "user",
    "user_notification_rule",
    "webhooks_delivery",
    "workflow_alert", // cannot auto-generate because codegen doesn't handle nested objects in trigger_params (alert_payload_conditions requires complex nested schema)
    "workflow_run",
    "workflow_task",
  ]
}

const readOnlyCollections = [
  "incident_post_mortem",
  "incident",
  "user",
]

function main() {
  if (filterResource) {
    console.log(`Generating code for resource: ${filterResource}`);
    if (resources().includes(filterResource)) {
      if (readOnlyCollections.includes(filterResource)) {
        generateReadOnlyClient(filterResource);
      } else {
        generateClient(filterResource);
      }
      generateResource(filterResource);
    } else if (dataSources().includes(filterResource)) {
      if (readOnlyCollections.includes(filterResource)) {
        generateReadOnlyClient(filterResource);
      } else {
        generateClient(filterResource);
      }
      generateDataSource(filterResource);
    } else {
      console.error(`Error: Resource '${filterResource}' not found in resources or data sources`);
      console.error(`Available resources: ${resources().slice(0, 10).join(', ')}...`);
      process.exit(1);
    }
  } else {
    // Generate everything
    generateProvider(resources(), workflowTaskResources(), dataSources())
    generateClients()
    generateResources()
    generateWorkflowTaskResources(workflowTaskResources(), swagger)
    generateDataSources()
  }
}

main()

function resources() {
  return Object.keys(swagger.components.schemas).filter((name) => {
    return !excluded.resources.includes(name) && collectionPathSchema(name);
  });
}

function dataSources() {
  return Object.keys(swagger.components.schemas).filter((name) => {
    return !excluded.dataSources.includes(name) && collectionPathSchema(name) && resourceHasFilters(name);
  });
}

function workflowTaskResources() {
  return Object.keys(swagger.components.schemas)
    .filter((key) => key.match(/_task_params/))
    .map((key) => key.replace("_task_params", ""))
}

function generateProvider(resources, taskResources, dataSources) {
  const code = providerTpl(
    resources.filter((name) => name !== "workflow"),
    taskResources,
    dataSources
  );
  fs.writeFileSync(
    path.resolve(__dirname, "..", "provider", "provider.go"),
    code
  );
}

function generateClients() {
  new Set([...resources(), ...dataSources()]).forEach((name) => {
    if (readOnlyCollections.includes(name)) {
      generateReadOnlyClient(name)
    } else {
      generateClient(name)
    }
  })
}

function generateResources() {
  resources().forEach(generateResource)
}

function generateDataSources() {
  dataSources().forEach(generateDataSource)
}

function generateReadOnlyClient(name) {
  const collectionSchema = collectionPathSchema(name);
  const pathIdField =
    collectionSchema &&
    collectionSchema.parameters &&
    collectionSchema.parameters[0] &&
    collectionSchema.parameters[0].name;
  const code = clientReadOnlyTpl(name, resourceSchema(name), pathIdField, hasQueryParam(name));
  fs.writeFileSync(
    path.resolve(__dirname, "..", "client", `${inflect.pluralize(name)}.go`),
    code
  );
}

function generateClient(name) {
  const collectionSchema = collectionPathSchema(name);
  const pathIdField =
    collectionSchema &&
    collectionSchema.parameters &&
    collectionSchema.parameters[0] &&
    collectionSchema.parameters[0].name;
  const code = clientTpl(name, resourceSchema(name), pathIdField, hasQueryParam(name));
  fs.writeFileSync(
    path.resolve(__dirname, "..", "client", `${inflect.pluralize(name)}.go`),
    code
  );
}

function resourceHasFilters(name) {
  const collectionSchema = collectionPathSchema(name);
  const filterParameters =
    collectionSchema.get && collectionSchema.get.parameters;
  return filterParameters.filter((filter) => filter.name.match(/filter/i))
    .length;
}

function generateDataSource(name) {
  const collectionSchema = collectionPathSchema(name);
  const filterParameters =
    collectionSchema.get && collectionSchema.get.parameters;
  const pathIdField =
    collectionSchema &&
    collectionSchema.parameters &&
    collectionSchema.parameters[0] &&
    collectionSchema.parameters[0].name;
  const code = dataSourceTpl(
    name,
    resourceSchema(name),
    filterParameters,
    pathIdField
  );
  if (code) {
    fs.writeFileSync(
      path.resolve(__dirname, "..", "provider", `data_source_${name}.go`),
      code
    );
  }
}

function generateResource(name) {
  const collectionSchema = collectionPathSchema(name);
  const pathIdField =
    collectionSchema &&
    collectionSchema.parameters &&
    collectionSchema.parameters[0] &&
    collectionSchema.parameters[0].name;
  const schema = resourceSchema(name);
  let code;
  if (name === "workflow") {
    code = workflowTpl(
      "workflow_incident",
      resourceSchema(name),
      requiredFields(name),
      swagger.components.schemas.incident_trigger_params
    );
    fs.writeFileSync(
      path.resolve(__dirname, "..", "provider", `resource_${name}_incident.go`),
      code
    );
    code = workflowTpl(
      "workflow_post_mortem",
      resourceSchema(name),
      requiredFields(name),
      swagger.components.schemas.post_mortem_trigger_params
    );
    fs.writeFileSync(
      path.resolve(
        __dirname,
        "..",
        "provider",
        `resource_${name}_post_mortem.go`
      ),
      code
    );
    code = workflowTpl(
      "workflow_action_item",
      resourceSchema(name),
      requiredFields(name),
      swagger.components.schemas.action_item_trigger_params
    );
    fs.writeFileSync(
      path.resolve(
        __dirname,
        "..",
        "provider",
        `resource_${name}_action_item.go`
      ),
      code
    );
    code = workflowTpl(
      "workflow_alert",
      resourceSchema(name),
      requiredFields(name),
      swagger.components.schemas.alert_trigger_params
    );
    fs.writeFileSync(
      path.resolve(__dirname, "..", "provider", `resource_${name}_alert.go`),
      code
    );
    code = workflowTpl(
      "workflow_pulse",
      resourceSchema(name),
      requiredFields(name),
      swagger.components.schemas.pulse_trigger_params
    );
    fs.writeFileSync(
      path.resolve(__dirname, "..", "provider", `resource_${name}_pulse.go`),
      code
    );
    code = workflowTpl(
      "workflow_simple",
      resourceSchema(name),
      requiredFields(name),
      swagger.components.schemas.simple_trigger_params
    );
    fs.writeFileSync(
      path.resolve(__dirname, "..", "provider", `resource_${name}_simple.go`),
      code
    );
  } else {
    code = resourceTpl(
      name,
      schema,
      requiredFields(name),
      pathIdField
    );
    fs.writeFileSync(
      path.resolve(__dirname, "..", "provider", `resource_${name}.go`),
      code
    );
  }
}

function resourceSchema(name) {
  return swagger.components.schemas[name];
}

function requiredFields(name) {
  const schemaName = `new_${name}`;
  const schema = swagger.components.schemas[schemaName];
  if (!schema) {
    console.warn(`Schema '${schemaName}' not found for resource '${name}'. Skipping required fields check.`);
    return [];
  }
  if (!schema.properties || !schema.properties.data || !schema.properties.data.properties || !schema.properties.data.properties.attributes) {
    console.warn(`Schema '${schemaName}' exists but doesn't have the expected structure for resource '${name}'. Skipping required fields check.`);
    return [];
  }
  return schema.properties.data.properties.attributes.required || [];
}

function collectionPathSchema(name) {
  return Object.keys(swagger.paths)
    .filter((url) => {
      const get = swagger.paths[url].get;
      return (
        get &&
        get.operationId &&
        get.operationId.replace(/ /g, "") ===
          `list${inflect.pluralize(inflect.camelize(name))}`
      );
    })
    .map((url) => swagger.paths[url])[0];
}

function hasQueryParam(name) {
  const paramsSchema = Object.keys(swagger.paths)
    .filter((url) => {
      const get = swagger.paths[url].get;
      return (
        get &&
        get.operationId &&
        get.operationId.replace(/ /g, "") ===
          `get${inflect.singularize(inflect.camelize(name))}`
      );
    })
    .map((url) => swagger.paths[url])[0]?.get;

  return paramsSchema &&
    paramsSchema.parameters &&
    paramsSchema.parameters.some((param) => param.in === "query");
}
