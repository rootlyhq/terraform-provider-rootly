import { dereference } from "@apidevtools/json-schema-ref-parser";

type T = any;

export const SWAGGER_MODS: ((swagger: T) => T | Promise<T>)[] = [
  /**
   * Recursively unwraps any `allOf` that contains exactly one item.
   *
   * Example:
   *   { allOf: [ { $ref: "#/something" } ] }
   * becomes:
   *   { $ref: "#/something" }
   */
  function unwrapSingleAllOfDeep<T>(schema: T): T {
    if (Array.isArray(schema)) {
      return schema.map(unwrapSingleAllOfDeep) as T;
    }

    if (schema && typeof schema === "object") {
      const obj: any = schema;

      // If this object *is itself* an allOf with a single entry
      if (obj.allOf && Array.isArray(obj.allOf)) {
        if (obj.allOf.length !== 1) {
          throw new Error(
            `Not implemented allOf length !== 1: ${JSON.stringify(obj)}`
          );
        }
        return unwrapSingleAllOfDeep(obj.allOf[0]);
      }

      // Otherwise recurse into properties
      const out: any = {};
      for (const key of Object.keys(obj)) {
        out[key] = unwrapSingleAllOfDeep(obj[key]);
      }
      return out;
    }

    return schema;
  },
  // Dereference all $refs
  dereference,
  // Fix dashboard_panel
  function fixDashboardPanel(swagger: any) {
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

    const dashboardPanel = swagger.components.schemas["dashboard_panel"];
    const newDashboardPanel = swagger.components.schemas["new_dashboard_panel"];
    const updateDashboardPanel =
      swagger.components.schemas["update_dashboard_panel"];

    const targets = [
      dashboardPanel.properties["params"],
      newDashboardPanel.properties["data"].properties["attributes"].properties[
        "params"
      ],
      updateDashboardPanel.properties["data"].properties["attributes"]
        .properties["params"],
    ];

    for (const target of targets) {
      fixDatasetsGroupByParam(target);
    }

    // Dashboard panel data is not saved to state
    delete dashboardPanel.properties["data"];

    // Dashboard panel requires dashboard_id
    newDashboardPanel.properties["data"].properties["attributes"].required = [
      ...(newDashboardPanel.properties["data"].properties["attributes"]
        .required ?? []),
      "dashboard_id",
    ];

    return swagger;
  },
];
