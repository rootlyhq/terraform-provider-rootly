import { oas30 } from "openapi3-ts";

const HTTP_METHODS = [
  "get",
  "post",
  "put",
  "delete",
  "patch",
  "head",
  "options",
  "trace",
] as const;

export function getParametersByOperationId({
  doc,
  operationId,
  onlyLocations,
  excludeLocations,
}: {
  doc: oas30.OpenAPIObject;
  operationId: string;
  onlyLocations?: oas30.ParameterLocation[];
  excludeLocations?: oas30.ParameterLocation[];
}): oas30.ParameterObject[] {
  for (const [_, pathItem] of Object.entries(doc.paths)) {
    for (const method of HTTP_METHODS) {
      const operation = pathItem[method];
      if (operation && operation.operationId === operationId) {
        return (
          [
            ...(pathItem.parameters ?? []),
            ...(operation.parameters ?? []),
          ] as oas30.ParameterObject[]
        ).filter((param) => {
          if (onlyLocations && !onlyLocations.includes(param.in)) {
            return false;
          }
          if (excludeLocations && excludeLocations.includes(param.in)) {
            return false;
          }
          return true;
        });
      }
    }
  }

  return [];
}

export function removeReference<T>(value: undefined): undefined;
export function removeReference<T>(value: T | oas30.ReferenceObject): T;
export function removeReference<T>(
  value: T | oas30.ReferenceObject | undefined,
): T | undefined;
export function removeReference<T>(
  value: T | oas30.ReferenceObject | undefined,
): T | undefined {
  if (typeof value === "object" && value !== null && !("$ref" in value)) {
    return value as T;
  }
  return undefined;
}
