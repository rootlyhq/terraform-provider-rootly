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
  notIn,
}: {
  doc: oas30.OpenAPIObject;
  operationId: string;
  notIn: oas30.ParameterLocation[];
}): oas30.ParameterObject[] | null {
  for (const [_, pathItem] of Object.entries(doc.paths)) {
    for (const method of HTTP_METHODS) {
      const operation = pathItem[method];
      if (operation && operation.operationId === operationId) {
        return (
          [
            ...(pathItem.parameters ?? []),
            ...(operation.parameters ?? []),
          ] as oas30.ParameterObject[]
        ).filter((param) => !notIn.includes(param.in));
      }
    }
  }
  return null;
}
