export function renameTfKeysDeep<T>(value: T): T {
  const prefix = "tf_";
  const replacement = "x-rootly-";

  if (Array.isArray(value)) {
    return value.map(renameTfKeysDeep) as T;
  }

  if (value && typeof value === "object") {
    const obj = value as Record<string, any>;
    const out: Record<string, any> = {};

    for (const [key, val] of Object.entries(obj)) {
      const newKey = key.startsWith(prefix)
        ? replacement + key.slice(prefix.length)
        : key;

      out[newKey] = renameTfKeysDeep(val);
    }

    return out as T;
  }

  return value;
}

/**
 * Recursively unwraps any `allOf` that contains exactly one item.
 *
 * Example:
 *   { allOf: [ { $ref: "#/something" } ] }
 * becomes:
 *   { $ref: "#/something" }
 */
export function unwrapSingleAllOfDeep<T>(schema: T): T {
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
}
type AnyObject = Record<string, any>;

/**
 * Recursively unwraps JSON:API-style { data: { id, type, attributes } } objects.
 * Works on plain JSON objects, arrays, or schema-like structures.
 */
export function unwrapJsonApi<T extends AnyObject>(obj: T): T {
  if (obj === null || typeof obj !== "object") return obj;

  // Handle arrays
  if (Array.isArray(obj)) {
    return obj.map(unwrapJsonApi) as any;
  }

  // Detect JSON:API resource: object with 'data' containing id, type, attributes
  const data = obj.data;
  const props = data && typeof data === "object" ? data : null;

  const looksLikeJsonApi =
    props &&
    ("attributes" in props || "id" in props || "type" in props) &&
    typeof props.attributes === "object";

  if (looksLikeJsonApi) {
    // Replace the whole object with the contents of attributes
    return unwrapJsonApi(props.attributes) as T;
  }

  // Otherwise, recurse into all keys
  const result: AnyObject = {};
  for (const [key, value] of Object.entries(obj)) {
    result[key] = unwrapJsonApi(value);
  }
  return result as T;
}
