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
    const isObj = (v: any) => v && typeof v === "object" && !Array.isArray(v);

    const merge = (a: any, b: any): any => {
      if (isObj(a) && isObj(b)) {
        const o: any = {};
        for (const k of new Set([...Object.keys(a), ...Object.keys(b)])) {
          if (k in a && k in b) o[k] = merge(a[k], b[k]);
          else if (k in b) o[k] = b[k];
          else o[k] = a[k];
        }
        return o;
      }
      return b !== undefined ? b : a;
    };

    const inner = (node: any): any => {
      if (Array.isArray(node)) return node.map(inner);
      if (!isObj(node)) return node;

      // unwrap: allOf with exactly 1 element
      if (Array.isArray(node.allOf) && node.allOf.length === 1) {
        const { allOf, ...rest } = node;
        const unwrapped = inner(allOf[0]);
        return inner(merge(unwrapped, rest));
      }

      // normal object: recurse into keys
      const out: any = {};
      for (const [k, v] of Object.entries(node)) out[k] = inner(v);
      return out;
    };

    return inner(schema);
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
