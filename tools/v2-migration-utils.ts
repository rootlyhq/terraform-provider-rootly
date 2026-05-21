export const v2Resources = ["dashboard_panel"];

type Schema = any;
type SchemaMod = (schema: Schema) => Schema;

export const SchemaMods: SchemaMod[] = [
  function fixIncidentPermissionSet(schema) {
    const props = schema.components.schemas["incident_permission_set"].properties;
    props["private_incident_permissions"].tf_computed = true;
    props["public_incident_permissions"].tf_computed = true;
    return schema;
  },

  function fixRolePermissions(schema) {
    const props = schema.components.schemas["role"].properties
    const fields = [
      "alerts_permissions",
      "api_keys_permissions",
      "audits_permissions",
      "billing_permissions",
      "catalogs_permissions",
      "communication_permissions",
      "edge_connector_permissions",
      "environments_permissions",
      "form_fields_permissions",
      "functionalities_permissions",
      "groups_permissions",
      "incident_causes_permissions",
      "incident_communication_permissions",
      "incident_feedbacks_permissions",
      "incident_roles_permissions",
      "incident_types_permissions",
      "incidents_permissions",
      "integrations_permissions",
      "invitations_permissions",
      "paging_permissions",
      "playbooks_permissions",
      "private_incidents_permissions",
      "pulses_permissions",
      "retrospective_permissions",
      "roles_permissions",
      "secrets_permissions",
      "services_permissions",
      "severities_permissions",
      "slas_permissions",
      "status_pages_permissions",
      "sub_statuses_permissions",
      "webhooks_permissions",
      "workflows_permissions",
    ]
    for (const field of fields) {
      if (props[field] !== undefined) {
        props[field].tf_computed = true;
      }
    }

    return schema;
  },

  function fixDashboardPanel(schema) {
    // Converts {"oneOf": [{"type": "string"}, {"type": "object"}]} to {"type": "string"}
    function fixDatasetsGroupByParam(paramsSchema: any) {
      const datasetsItemsProperties =
        paramsSchema.properties["datasets"].items.properties;
      const groupBySchema = datasetsItemsProperties["group_by"];
      const { oneOf, ...rest } = groupBySchema;
      datasetsItemsProperties["group_by"] = {
        ...oneOf[0],
        ...rest,
      };
    }

    const targets = [
      schema.components.schemas["dashboard_panel"].properties["params"],
      schema.components.schemas["new_dashboard_panel"].properties["data"]
        .properties["attributes"].properties["params"],
      schema.components.schemas["update_dashboard_panel"].properties["data"]
        .properties["attributes"].properties["params"],
    ];

    for (const target of targets) {
      fixDatasetsGroupByParam(target);
    }

    // Dashboard panel data is not saved to state
    delete schema.components.schemas["dashboard_panel"].properties["data"];

    // Dashboard panel requires dashboard_id
    schema.components.schemas["new_dashboard_panel"].properties[
      "data"
    ].properties["attributes"].required.push("dashboard_id");

    // Dashboard panel params description is not computed
    schema.components.schemas["dashboard_panel"].properties[
      "params"
    ].properties["description"].tf_computed = false;

    return schema;
  },
];
