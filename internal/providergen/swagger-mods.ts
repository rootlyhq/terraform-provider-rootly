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
    swagger.components.schemas.dashboard_panel.properties.params.properties.datasets.items.properties.group_by =
      {
        type: "string",
        nullable: true,
      };
    delete swagger.components.schemas.dashboard_panel.properties.data;
    return swagger;
  },
];
