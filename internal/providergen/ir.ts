import { camelize, humanize, pluralize, singularize } from "inflection";
import { match, P } from "ts-pattern";

export interface DataSource {
  name: string;
}

export interface DataSourceIR extends DataSource {}

export interface Resource {
  name: string;
  description?: string;
  mode?: "legacy" | "modern";
}

export interface ResourceIR extends Resource {
  mode: "legacy" | "modern";
  attributes: IRType[];
  listPathIdAttribute?: IRString;
  idAttribute: IRString;
  getHasQueryParams: boolean;
  original: {
    resourceSchema: any;
    newResourceSchema: any;
    updateResourceSchema: any;
    collectionSchema: any;
    getSchema: any;
  };
}

export interface BaseAttribute {
  name: string;
  description: string;
}

type ComputedOptionalRequired =
  | "computed"
  | "optional"
  | "computed_optional"
  | "required";

export type IRType =
  | IRString
  | IRBool
  | IRInt
  | IRList
  | IRSet
  | IRObject
  | IRListNested
  | IRSetNested;

interface IRBase {
  name: string;
  type: string;
  description?: string;
  deprecationMessage?: string;
  sensitive?: boolean;
  nullable?: boolean;
  computedOptionalRequired: ComputedOptionalRequired;
  validators?: string[];
  planModifiers?: string[];
  original?: {
    schema: any;
    newSchema: any;
    updateSchema: any;
  };
}

export interface IRString extends IRBase {
  type: "string";
  enum?: string[];
}

export interface IRBool extends IRBase {
  type: "bool";
}

export interface IRInt extends IRBase {
  type: "int";
}

export interface IRObject extends IRBase {
  type: "object";
  attributes: IRType[];
}

export interface IRList extends IRBase {
  type: "list";
  elementType: "string";
}

export interface IRSet extends IRBase {
  type: "set";
  elementType: "string";
}

export interface IRListNested extends IRBase {
  type: "list_nested";
  attributes: IRType[];
}

export interface IRSetNested extends IRBase {
  type: "set_nested";
  attributes: IRType[];
}

export interface IRResource extends IRBase {
  type: "resource";
  resourceType: string;
  listPathIdParam: {
    name: string;
    element: IRType;
  } | null;
  getHasQueryParams: boolean;
  idElement: IRString;
  fields: Record<string, IRType>;
}

function toIR({
  name,
  computedOptionalRequired,
  schema,
  newSchema,
  updateSchema,
}: {
  name: string;
  computedOptionalRequired: ComputedOptionalRequired;
  schema: any;
  newSchema: any;
  updateSchema: any;
}): IRType {
  const common: Pick<
    IRBase,
    | "name"
    | "computedOptionalRequired"
    | "description"
    | "nullable"
    | "original"
  > = {
    name,
    computedOptionalRequired,
    description: schema.description,
    nullable: schema.nullable ?? false,
    original: {
      schema,
      newSchema,
      updateSchema,
    },
  };

  return match(schema)
    .returnType<IRType>()
    .with({ type: "string" }, () => ({
      type: "string",
      enum: schema.enum,
      ...common,
    }))
    .with({ type: "integer" }, () => ({
      type: "int",
      ...common,
    }))
    .with({ type: "array", items: { type: "string" } }, () => ({
      type: "set",
      elementType: "string",
      ...common,
    }))
    .with(
      {
        type: "object",
        properties: P.record(P.string, P.any),
        required: P.array(P.string).optional(),
      },
      (schema) => ({
        type: "object",
        attributes: Object.entries(schema.properties)
          .filter(([_, schema]: [string, any]) => !schema.tf_ignore)
          .map(([propertyName, propertySchema]: [string, any]) => {
            return toIR({
              name: propertyName,
              computedOptionalRequired: toComputedOptionalRequired({
                name: propertyName,
                schema,
                newSchema,
                updateSchema,
              }),
              schema: propertySchema,
              newSchema: newSchema.properties[propertyName],
              updateSchema: updateSchema.properties[propertyName],
            });
          }),
        ...common,
      })
    )
    .otherwise(() => {
      throw new Error(`Unsupported type: ${JSON.stringify(schema)}`);
    });
}

function toComputedOptionalRequired({
  name,
  schema,
  newSchema,
  updateSchema,
}: {
  name: string;
  schema: any;
  newSchema: any;
  updateSchema: any;
}): ComputedOptionalRequired {
  const inRead = name in schema.properties;
  const inCreate = name in newSchema.properties;
  const inUpdate = name in updateSchema.properties;
  const reqCreate = newSchema.required?.includes(name);
  const reqUpdate = updateSchema.required?.includes(name);

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
    `Unsupported computedOptionalRequired for field ${name}: ${JSON.stringify({
      schema,
      newSchema,
      updateSchema,
    })}`
  );
}

export function generateDataSourceIR({
  swagger,
  dataSource,
}: {
  swagger: any;
  dataSource: DataSource;
}) {
  const resourceSchema = swagger.components.schemas[dataSource.name];
  if (!resourceSchema) {
    throw new Error(`Resource ${dataSource.name} not found`);
  }

  const getSchema = Object.entries(swagger.paths as Record<string, any>).find(
    ([_, pathSchema]) =>
      pathSchema.get &&
      pathSchema.get.operationId ===
        `get${camelize(singularize(dataSource.name))}`
  )?.[1]?.get;
  if (!getSchema) {
    throw new Error(`Get path for ${dataSource.name} not found`);
  }

  const attributes = Object.entries(resourceSchema.properties)
    .filter(([name]) => !["created_at", "updated_at"].includes(name))
    .map(([name, schema]: [string, any]) => ({
      name,
      type: "string",
      description: schema.description,
      nullable: schema.nullable || undefined,
      original: schema,
    }));

  return {
    name: dataSource.name,
    attributes,
    original: {
      resourceSchema,
      getSchema,
    },
  };
}

export function generateResourceIR({
  swagger,
  resource,
}: {
  swagger: any;
  resource: Resource;
}): ResourceIR {
  const resourceSchema = swagger.components.schemas[resource.name];
  if (!resourceSchema) {
    throw new Error(`Resource ${resource.name} not found`);
  }

  const newResourceSchema = swagger.components.schemas[`new_${resource.name}`];
  if (!newResourceSchema) {
    throw new Error(`New resource ${resource.name} not found`);
  }

  const updateResourceSchema =
    swagger.components.schemas[`update_${resource.name}`];
  if (!updateResourceSchema) {
    throw new Error(`Update resource ${resource.name} not found`);
  }

  const collectionSchema = Object.entries(
    swagger.paths as Record<string, any>
  ).find(
    ([_, pathSchema]) =>
      pathSchema.get &&
      pathSchema.get.operationId === `list${camelize(pluralize(resource.name))}`
  )?.[1];
  if (!collectionSchema) {
    throw new Error(`List path for ${resource.name} not found`);
  }

  const getSchema = Object.entries(swagger.paths as Record<string, any>).find(
    ([_, pathSchema]) =>
      pathSchema.get &&
      pathSchema.get.operationId ===
        `get${camelize(singularize(resource.name))}`
  )?.[1]?.get;
  if (!getSchema) {
    throw new Error(`Get path for ${resource.name} not found`);
  }

  // Get ID parameter
  const listPathIdParameterName = collectionSchema?.parameters?.[0]?.name as
    | string
    | undefined;
  const listPathIdAttribute = listPathIdParameterName
    ? toIR({
        name: listPathIdParameterName,
        schema: resourceSchema.properties[listPathIdParameterName],
        newSchema: null,
        updateSchema: null,
        computedOptionalRequired: "required",
      })
    : undefined;
  if (listPathIdAttribute && listPathIdAttribute.type !== "string") {
    throw new Error(`ID attribute for ${resource.name} must be a string`);
  }

  const getHasQueryParams =
    getSchema?.parameters?.some((param: any) => param.in === "query") ?? false;

  // Generate immediate representation of the resource
  const tempIR = toIR({
    name: resource.name,
    computedOptionalRequired: "required",
    schema: resourceSchema,
    newSchema: newResourceSchema.properties.data.properties.attributes,
    updateSchema: updateResourceSchema.properties.data.properties.attributes,
  });
  if (tempIR.type !== "object") {
    throw new Error("Resource root must be an object");
  }
  const attributes = tempIR.attributes.filter(
    (attribute) => !["created_at", "updated_at"].includes(attribute.name)
  );

  return {
    ...resource,
    mode: resource.mode ?? "modern",
    description: resourceSchema.description,
    listPathIdAttribute: listPathIdAttribute,
    idAttribute: {
      name: "id",
      type: "string",
      description: `The ID of the ${humanize(resource.name, true)}`,
      computedOptionalRequired: "computed",
    },
    getHasQueryParams,
    attributes,
    original: {
      resourceSchema,
      newResourceSchema,
      updateResourceSchema,
      collectionSchema,
      getSchema,
    },
  };
}
