import { oas30 } from "openapi3-ts";

export function assertSchemaObject(
  schema: oas30.SchemaObject | oas30.ReferenceObject,
): asserts schema is oas30.SchemaObject {
  if ("$ref" in schema) {
    throw new Error(`Unsupported: Schema is a reference to ${schema.$ref}`);
  }
}
