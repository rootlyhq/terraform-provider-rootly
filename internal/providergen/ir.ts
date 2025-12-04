import { camelize, humanize, pluralize, singularize } from "inflection";
import { match, P } from "ts-pattern";

interface IRBase {
  kind: string;
  description?: string;
}

interface IRBaseWithComputedOptionalRequired extends IRBase {
  computedOptionalRequired: "computed" | "optional" | "required";
}

export type IRType =
  | IRString
  | IRBool
  | IRInt
  | IRObject
  | IRArray
  | IRResource;

export interface IRString extends IRBaseWithComputedOptionalRequired {
  kind: "string";
}

export interface IRBool extends IRBaseWithComputedOptionalRequired {
  kind: "bool";
}

export interface IRInt extends IRBaseWithComputedOptionalRequired {
  kind: "int";
}

export interface IRObject extends IRBaseWithComputedOptionalRequired {
  kind: "object";
  fields: Record<string, Exclude<IRType, IRResource>>;
}

export interface IRArray extends IRBaseWithComputedOptionalRequired {
  kind: "array";
  element: Exclude<IRType, IRResource>;
}

export interface IRResource extends IRBase {
  kind: "resource";
  resourceType: string;
  listPathIdParam: {
    name: string;
    element: Exclude<IRType, IRResource>;
  } | null;
  getHasQueryParams: boolean;
  idElement: IRString;
  fields: Record<string, Exclude<IRType, IRResource>>;
}

// TODO: Handle computed
export function toIR({
  schema,
  required,
}: {
  schema: any;
  required: boolean | null;
}): Exclude<IRType, IRResource> {
  const common = {
    computedOptionalRequired: required ? "required" : "optional",
    description: schema.description,
  } as const;

  return match(schema)
    .returnType<Exclude<IRType, IRResource>>()
    .with({ type: "string" }, () => ({
      kind: "string",
      ...common,
    }))
    .with({ type: "boolean" }, () => ({
      kind: "bool",
      ...common,
    }))
    .with({ type: "integer" }, { type: "number" }, () => ({
      kind: "int",
      ...common,
    }))
    .with({ type: "array", items: P.record(P.string, P.any) }, (schema) => ({
      kind: "array",
      ...common,
      element: toIR({ schema: schema.items, required: null }),
    }))
    .with(
      {
        type: "object",
        properties: P.record(P.string, P.any),
        required: P.array(P.string).optional(),
      },
      (schema) => {
        return {
          kind: "object",
          ...common,
          fields: Object.fromEntries(
            Object.entries(schema.properties).map(
              ([propertyName, propertySchema]) => [
                propertyName,
                toIR({
                  schema: propertySchema,
                  required: schema.required?.includes(propertyName) ?? null,
                }),
              ]
            )
          ),
        };
      }
    )
    .otherwise(() => {
      throw new Error(
        `Unsupported swagger schema type: ${JSON.stringify(schema)}`
      );
    });
}

export function generateResourceIR({
  swagger,
  name,
}: {
  swagger: any;
  name: string;
}) {
  const resourceSchema = swagger.components.schemas[name];
  if (!resourceSchema) {
    throw new Error(`Resource ${name} not found`);
  }

  const newResourceSchema = swagger.components.schemas[`new_${name}`];
  if (!newResourceSchema) {
    throw new Error(`New resource ${name} not found`);
  }

  const collectionSchema = Object.entries(
    swagger.paths as Record<string, any>
  ).find(
    ([_, pathSchema]) =>
      pathSchema.get &&
      pathSchema.get.operationId === `list${camelize(pluralize(name))}`
  )?.[1];
  if (!collectionSchema) {
    throw new Error(`List path for ${name} not found`);
  }

  const getSchema = Object.entries(swagger.paths as Record<string, any>).find(
    ([_, pathSchema]) =>
      pathSchema.get &&
      pathSchema.get.operationId === `get${camelize(singularize(name))}`
  )?.[1]?.get;
  if (!getSchema) {
    throw new Error(`Get path for ${name} not found`);
  }

  // Get path ID parameter
  const pathIdParameter = collectionSchema?.parameters?.[0]?.name as
    | string
    | undefined;
  const pathIdIR = pathIdParameter
    ? toIR({
        schema: resourceSchema.properties[pathIdParameter],
        required: null,
      })
    : null;

  const getHasQueryParams =
    getSchema?.parameters?.some((param: any) => param.in === "query") ?? false;

  // Generate immediate representation of the resource
  const irFields = toIR({
    schema: resourceSchema,
    required: newResourceSchema.required,
  });
  if (irFields.kind !== "object") {
    throw new Error("Resource root must be an object");
  }

  const ir: IRResource = {
    kind: "resource",
    resourceType: name,
    listPathIdParam:
      pathIdParameter && pathIdIR
        ? { name: pathIdParameter, element: pathIdIR }
        : null,
    getHasQueryParams,
    idElement: {
      kind: "string",
      computedOptionalRequired: "computed",
      description: `The ID of the ${humanize(name, true)}`,
    },
    fields: irFields.fields,
  };

  return ir;
}
