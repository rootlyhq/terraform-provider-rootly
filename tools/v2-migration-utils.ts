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
];
