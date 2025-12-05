import { camelize, humanize, pluralize, singularize } from "inflection";
import { match, P } from "ts-pattern";

type ComputedOptionalRequired =
  | "computed"
  | "optional"
  | "computed_optional"
  | "required";

export type IRType =
  | IRString
  | IRBool
  | IRInt
  | IRObject
  | IRArray
  | IRResource;

interface IRBase {
  kind: string;
  description?: string;
  computedOptionalRequired: ComputedOptionalRequired;
}

export interface IRString extends IRBase {
  kind: "string";
}

export interface IRBool extends IRBase {
  kind: "bool";
}

export interface IRInt extends IRBase {
  kind: "int";
}

export interface IRObject extends IRBase {
  kind: "object";
  fields: Record<string, Exclude<IRType, IRResource>>;
}

export interface IRArray extends IRBase {
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
function toIR({
  schema,
  newSchema,
  updateSchema,
  computedOptionalRequired,
}: {
  schema: any;
  newSchema: any;
  updateSchema: any;
  computedOptionalRequired: ComputedOptionalRequired;
}): Exclude<IRType, IRResource> {
  const common = {
    computedOptionalRequired,
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
      element: toIR({
        schema: schema.items,
        newSchema: newSchema.items,
        updateSchema: updateSchema.items,
        computedOptionalRequired: "required", // TODO: Investigate
      }),
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
                  newSchema: newSchema.properties[propertyName],
                  updateSchema: updateSchema.properties[propertyName],
                  computedOptionalRequired: toComputedOptionalRequired({
                    field: propertyName,
                    schema,
                    newSchema,
                    updateSchema,
                  }),
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

function toComputedOptionalRequired({
  field,
  schema,
  newSchema,
  updateSchema,
}: {
  field: string;
  schema: any;
  newSchema: any;
  updateSchema: any;
}): ComputedOptionalRequired {
  if (field === "dashboard_id") {
    console.log({ schema, newSchema, updateSchema });
  }
  const inRead = field in schema.properties;
  const inCreate = field in newSchema.properties;
  const inUpdate = field in updateSchema.properties;
  const reqCreate = newSchema.required?.includes(field);
  const reqUpdate = updateSchema.required?.includes(field);

  if (reqCreate || reqUpdate) {
    return "required";
  }

  if (inRead && !inCreate && !inUpdate) {
    return "computed";
  }

  if (inCreate || inUpdate) {
    if (inRead) {
      return "computed_optional";
    } else {
      return "optional";
    }
  }

  if (inRead) {
    return "computed";
  }

  throw new Error(
    `Unsupported computedOptionalRequired for field ${field}: ${JSON.stringify({
      schema,
      newSchema,
      updateSchema,
    })}`
  );
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

  const updateResourceSchema = swagger.components.schemas[`update_${name}`];
  if (!updateResourceSchema) {
    throw new Error(`Update resource ${name} not found`);
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
        newSchema: null,
        updateSchema: null,
        computedOptionalRequired: "required",
      })
    : null;

  const getHasQueryParams =
    getSchema?.parameters?.some((param: any) => param.in === "query") ?? false;

  console.log(newResourceSchema);

  // Generate immediate representation of the resource
  const irFields = toIR({
    schema: resourceSchema,
    newSchema: newResourceSchema.properties.data.properties.attributes,
    updateSchema: updateResourceSchema.properties.data.properties.attributes,
    computedOptionalRequired: "required",
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

  console.log(ir);

  return ir;
}
