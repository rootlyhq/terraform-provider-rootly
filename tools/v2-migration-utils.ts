export const v2Resources = ["dashboard_panel"];

type Schema = any;
type SchemaMod = (schema: Schema) => Schema;

export const SchemaMods: SchemaMod[] = [
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

  function addApiKeyGroupId(schema) {
    // The API requires group_id when kind is "team", but it only exposes the
    // attribute on new_api_key: it is absent from the api_key read model and
    // from update_api_key. Copy it onto the read model so the resource can send
    // it on create, and flag it as create-only so it is neither read back
    // (the API never returns it) nor sent on update (the API rejects it).
    const apiKey = schema.components.schemas["api_key"];
    const newApiKeyAttributes =
      schema.components.schemas["new_api_key"].properties["data"].properties[
        "attributes"
      ];

    if (!apiKey.properties["group_id"]) {
      apiKey.properties["group_id"] = {
        ...newApiKeyAttributes.properties["group_id"],
        tf_create_only: true,
        tf_force_new: true,
      };
    }

    return schema;
  },
];
