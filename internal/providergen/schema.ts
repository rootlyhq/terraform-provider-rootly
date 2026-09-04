import type { oas30 } from "openapi3-ts";
import { match, P } from "ts-pattern";
import { getParametersByOperationId, removeReference } from "./openapi";
import assert from "node:assert";
import { camelize } from "inflection";

export interface ClientConfig {
  name: string;
}

interface DataSourceListConfig {
  strategy: "list";
  resourceName: string;
}

interface DataSourceSingleConfig {
  strategy: "single";
}

export type DataSourceConfig = {
  name: string;
  description?: string;
  modifyDef?: (def: DataSourceDef) => DataSourceDef;
} & (DataSourceListConfig | DataSourceSingleConfig);

export type ResourceConfig = {
  name: string;
  description?: string;
  modifyDef?: (def: ResourceDef) => ResourceDef;
};

declare module "openapi3-ts/oas30" {
  interface SchemaObject {
    /** Defines the jsonapi attribute tag. */
    "x-go-jsonapi-tag"?: "primary" | "attr";
    /** Defines the jsonapi type. */
    "x-go-jsonapi-type"?: string;
    /** Overrides the nested object type name. Only used for nested objects. */
    "x-go-nested-type"?: string;
    /** Defines if the property is computed, optional or required. */
    "x-tf-computed-optional-required"?: ComputedOptionalRequired;
    /** Indicates this is the top level item type for the data source. Only used for plural data sources. */
    "x-tf-top-level-item-type"?: boolean;
    /** The type of collection to use for arrays. */
    "x-tf-collection-type"?: "list" | "set";
    /** The schema to use for reading the property. */
    "x-schema-read"?: oas30.SchemaObject;
    /** The schema to use for creating the property. */
    "x-schema-create"?: oas30.SchemaObject;
    /** The schema to use for updating the property. */
    "x-schema-update"?: oas30.SchemaObject;
    /** Indicates that this property should be ignored. */
    tf_ignore?: boolean;
  }
}

interface AttributeBase {
  name: string;
  type: string;
  description?: string;
  deprecationMessage?: string;
  sensitive?: string;
  nullable?: boolean;
  computedOptionalRequired: ComputedOptionalRequired;
  validators?: string[];
  planModifiers?: string[];
  schemas?: {
    read: oas30.SchemaObject;
    create?: oas30.SchemaObject;
    update?: oas30.SchemaObject;
  };
  paramSchemas?: {
    read?: oas30.ParameterObject;
    create?: oas30.ParameterObject;
    update?: oas30.ParameterObject;
    delete?: oas30.ParameterObject;
  };
}

export interface AttributeString extends AttributeBase {
  type: "string";
  enum?: string[];
  default?: string;
}

export interface AttributeBool extends AttributeBase {
  type: "bool";
  default?: boolean;
}

export interface AttributeInt64 extends AttributeBase {
  type: "int64";
  default?: number;
}

export interface AttributeObject extends AttributeBase {
  type: "object";
  attributes: AttributeType[];
  blocks: AttributeBlockType[];
}

export interface AttributeList extends AttributeBase {
  type: "list";
  elementType: "string" | "int64";
}

export interface AttributeSet extends AttributeBase {
  type: "set";
  elementType: "string" | "int64";
}

export interface AttributeListNested extends AttributeBase {
  type: "list_nested";
  attributes: AttributeType[];
  blocks: AttributeBlockType[];
  hacks?: {
    /** The behavior for unknown values. Default is "omit" */
    unknownBehavior?: "omit" | "empty" | "null";
    /** The behavior for null values. Default is "null" */
    nullBehavior?: "omit" | "empty" | "null";
    /** The behavior for empty values. Default is "empty" */
    emptyBehavior?: "omit" | "empty" | "null";
  };
}

export interface AttributeSetNested extends AttributeBase {
  type: "set_nested";
  attributes: AttributeType[];
  blocks: AttributeBlockType[];
  hacks?: {
    /** The behavior for unknown values. Default is "omit" */
    unknownBehavior?: "omit" | "empty" | "null";
    /** The behavior for null values. Default is "null" */
    nullBehavior?: "omit" | "empty" | "null";
    /** The behavior for empty values. Default is "empty" */
    emptyBehavior?: "omit" | "empty" | "null";
  };
}

export type AttributeType =
  | AttributeString
  | AttributeBool
  | AttributeInt64
  | AttributeObject
  | AttributeList
  | AttributeSet
  | AttributeListNested
  | AttributeSetNested;

export type AttributeBlockType = AttributeListNested | AttributeSetNested;

export type ComputedOptionalRequired =
  | "computed"
  | "optional"
  | "computed_optional"
  | "required";

export interface DataSourceDef {
  name: string;
  description?: string;
  attributes: AttributeType[];
  blocks: AttributeBlockType[];
  goNames: {
    /** Name of the struct that represents the base client. */
    clientBase: string;
    /** Name of the struct that represents the resource. */
    struct: `${string}DataSource`;
    /** Name of the struct that represents the model of the resource. */
    model: `${string}DataSourceModel`;
  };
}

export interface ResourceDef {
  name: string;
  description?: string;
  attributes: AttributeType[];
  blocks: AttributeBlockType[];
  goNames: {
    /** Name of the struct that represents the base client. */
    clientBase: string;
    /** Name of the struct that represents the resource. */
    struct: `${string}Resource`;
    /** Name of the struct that represents the model of the resource. */
    model: `${string}ResourceModel`;
  };
}

function openapiSchemaToAttribute({
  name,
  computedOptionalRequired,
  schemas,
  options,
}: {
  name: string;
  computedOptionalRequired: ComputedOptionalRequired;
  schemas: {
    read: oas30.SchemaObject;
    create?: oas30.SchemaObject;
    update?: oas30.SchemaObject;
  };
  options: {
    defaultCollectionType: "list" | "set";
    defaultNestedCollectionType: "list_nested" | "set_nested";
    collectionsAsBlocks: boolean;
    readOnly?: boolean;
  };
}): AttributeType {
  const common: Pick<
    AttributeBase,
    "name" | "description" | "nullable" | "computedOptionalRequired" | "schemas"
  > = {
    name,
    description: cleanDescription(schemas.read.description),
    nullable: schemas.read.nullable ?? undefined,
    computedOptionalRequired: options.readOnly
      ? "computed"
      : computedOptionalRequired,
    schemas,
  };

  return match(schemas.read)
    .returnType<AttributeType>()
    .with(
      {
        type: "string",
        enum: P.array(P.string).optional(),
        default: P.string.optional(),
      },
      (schema) => ({
        ...common,
        type: "string",
        enum: schema.enum,
        default: schema.default,
      }),
    )
    .with(
      {
        type: "integer",
        enum: P.array(P.number).optional(),
        default: P.number.optional(),
      },
      (schema) => ({
        ...common,
        type: "int64",
        enum: schema.enum,
        default: schema.default,
        validators: schema.enum
          ? [
              `int64validator.OneOf(${schema.enum.map((value) => JSON.stringify(value)).join(", ")})`,
            ]
          : undefined,
      }),
    )
    .with({ type: "boolean", default: P.boolean.optional() }, (schema) => ({
      ...common,
      type: "bool",
      default: schema.default,
    }))
    .with({ type: "array", items: { type: "string" } }, () => ({
      ...common,
      type: options.defaultCollectionType,
      elementType: "string",
    }))
    .with({ type: "array", items: { type: "integer" } }, () => ({
      ...common,
      type: options.defaultCollectionType,
      elementType: "int64",
    }))
    .with({ type: "array", items: { type: "object" } }, (schema) => {
      const tmpAttr = openapiSchemaToAttribute({
        name: "item",
        computedOptionalRequired,
        schemas: {
          read: schema.items,
          create: schema.items,
          update: schema.items,
        },
        options,
      });
      assert(tmpAttr.type === "object");
      return {
        ...common,
        type: options.defaultNestedCollectionType,
        attributes: tmpAttr.attributes,
        blocks: tmpAttr.blocks,
      };
    })
    .with(
      {
        type: "object",
        properties: P.record(P.string, P.any),
        required: P.array(P.string).optional(),
      },
      (schema) => {
        const allAttributes = Object.entries(schema.properties)
          .filter(([, property]) => !removeReference(property).tf_ignore)
          .map(([name, property]) =>
            openapiSchemaToAttribute({
              name,
              computedOptionalRequired: toComputedOptionalRequired({
                name,
                schemas,
              }),
              schemas: {
                read: removeReference(property),
                create: removeReference(schemas.create?.properties?.[name]),
                update: removeReference(schemas.update?.properties?.[name]),
              },
              options,
            }),
          );

        const attributes = options.collectionsAsBlocks
          ? allAttributes.filter((attr) => !isCollectionAttribute(attr))
          : allAttributes;
        const blocks = options.collectionsAsBlocks
          ? allAttributes.filter((attr) => isCollectionAttribute(attr))
          : [];

        return {
          ...common,
          type: "object",
          attributes,
          blocks,
        };
      },
    )
    .otherwise((schema) => {
      throw new Error(`Unsupported schema: ${JSON.stringify(schema)}`);
    });
}

function openapiParametersToAttributes({
  params,
}: {
  params: {
    read: oas30.ParameterObject[];
    list?: oas30.ParameterObject[];
    create: oas30.ParameterObject[];
    update: oas30.ParameterObject[];
    delete: oas30.ParameterObject[];
  };
}): AttributeType[] {
  const allParamNames = Array.from(
    new Set([
      ...params.read.map((param) => param.name),
      ...(params.list?.map((param) => param.name) ?? []),
      ...params.create.map((param) => param.name),
      ...params.update.map((param) => param.name),
      ...params.delete.map((param) => param.name),
    ]),
  );

  return allParamNames.map((name) => {
    const readParam = params.read.find((param) => param.name === name);
    const listParam = params.list?.find((param) => param.name === name);
    const createParam = params.create.find((param) => param.name === name);
    const updateParam = params.update.find((param) => param.name === name);
    const deleteParam = params.delete.find((param) => param.name === name);

    const baseParam =
      readParam ?? listParam ?? createParam ?? updateParam ?? deleteParam!;
    const baseParamSchema = removeReference(baseParam.schema);
    assert(baseParamSchema, `Parameter ${name} has no schema`);
    assert(
      baseParamSchema.type === "string",
      `Parameter ${name} has no string type`,
    );

    const {
      computedOptionalRequired,
      planModifiers,
    }: {
      computedOptionalRequired: ComputedOptionalRequired;
      planModifiers?: string[];
    } = (() => {
      // TODO: Handle listParam

      if (createParam && !readParam && !updateParam && !deleteParam) {
        return {
          computedOptionalRequired: "required",
          planModifiers: ["stringplanmodifier.RequiresReplace()"],
        };
      }

      if (!createParam && readParam && updateParam && deleteParam) {
        return {
          computedOptionalRequired: "computed",
          planModifiers: ["stringplanmodifier.UseStateForUnknown()"],
        };
      }

      throw new Error(
        `Unsupported parameter ${name} computed optional required: ${JSON.stringify({ readParam, createParam, updateParam, deleteParam })}`,
      );
    })();

    return {
      type: "string",
      name,
      computedOptionalRequired,
      description: baseParam.description ?? baseParamSchema.description,
      enum: baseParamSchema.enum,
      planModifiers,
      paramSchemas: {
        read: readParam,
        list: listParam,
        create: createParam,
        update: updateParam,
        delete: deleteParam,
      },
    };
  });
}

export function generateDataSourceDef({
  doc,
  config,
}: {
  doc: oas30.OpenAPIObject;
  config: DataSourceConfig;
}): DataSourceDef {
  const isSingle = config.strategy === "single";
  const operationId = `${isSingle ? "get" : "list"}${camelize(config.name)}`;

  // Get read schema
  const readSchemaKey = isSingle ? config.name : config.resourceName;
  const readSchema = removeReference(doc.components?.schemas?.[readSchemaKey]);
  if (!readSchema) {
    throw new Error(`Read schema "${readSchemaKey}" not found`);
  }

  // Get path params
  const pathParams = getParametersByOperationId({
    doc,
    operationId,
    onlyLocations: ["path"],
  });

  // TODO: Add support for path params

  const tmpAttr = openapiSchemaToAttribute({
    name: config.name,
    computedOptionalRequired: "computed",
    schemas: {
      read: readSchema,
    },
    options: {
      defaultCollectionType: "list",
      defaultNestedCollectionType: "list_nested",
      collectionsAsBlocks: false,
      readOnly: true,
    },
  });
  assert(tmpAttr.type === "object");

  return {
    name: config.name,
    description: config.description ?? tmpAttr.description,
    attributes: tmpAttr.attributes,
    blocks: tmpAttr.blocks,
    goNames: {
      clientBase: camelize(readSchemaKey),
      struct: `${camelize(config.name)}DataSource`,
      model: `${camelize(config.name)}DataSourceModel`,
    },
  };
}

export function generateResourceDef({
  doc,
  config,
}: {
  doc: oas30.OpenAPIObject;
  config: ResourceConfig;
}): ResourceDef {
  // Get read schema
  const readSchema = removeReference(doc.components?.schemas?.[config.name]);
  if (!readSchema) {
    throw new Error(`Read schema "${config.name}" not found`);
  }

  // Get create schema
  const createSchema = removeReference(
    removeReference(
      removeReference(doc.components?.schemas?.[`new_${config.name}`])
        ?.properties?.data,
    )?.properties?.attributes,
  );
  if (!createSchema) {
    throw new Error(`Create schema "new_${config.name}" not found`);
  }

  // Get update schema
  const updateSchema = removeReference(
    removeReference(
      removeReference(doc.components?.schemas?.[`update_${config.name}`])
        ?.properties?.data,
    )?.properties?.attributes,
  );
  if (!updateSchema) {
    throw new Error(`Update schema "update_${config.name}" not found`);
  }

  // Get read path params
  const readPathParams = getParametersByOperationId({
    doc,
    operationId: `get${camelize(config.name)}`,
    onlyLocations: ["path"],
  });

  // Get read path params
  const listPathParams = getParametersByOperationId({
    doc,
    operationId: `list${camelize(config.name)}`,
    onlyLocations: ["path"],
  });

  // Get create path params
  const createPathParams = getParametersByOperationId({
    doc,
    operationId: `create${camelize(config.name)}`,
    onlyLocations: ["path"],
  });

  // Get update path params
  const updatePathParams = getParametersByOperationId({
    doc,
    operationId: `update${camelize(config.name)}`,
    onlyLocations: ["path"],
  });

  // Get Delete path params
  const deletePathParams = getParametersByOperationId({
    doc,
    operationId: `delete${camelize(config.name)}`,
    onlyLocations: ["path"],
  });

  const tmpPathAttrs = openapiParametersToAttributes({
    params: {
      read: readPathParams,
      list: listPathParams,
      create: createPathParams,
      update: updatePathParams,
      delete: deletePathParams,
    },
  });

  const tmpAttr = openapiSchemaToAttribute({
    name: config.name,
    computedOptionalRequired: "required",
    schemas: {
      read: readSchema,
      create: createSchema,
      update: updateSchema,
    },
    options: {
      defaultCollectionType: "set",
      defaultNestedCollectionType: "list_nested",
      collectionsAsBlocks: true,
    },
  });
  assert(tmpAttr.type === "object");

  return {
    name: config.name,
    description: cleanDescription(config.description ?? tmpAttr.description),
    attributes: mergeAttributeTypes(tmpAttr.attributes, tmpPathAttrs),
    blocks: tmpAttr.blocks,
    goNames: {
      clientBase: camelize(config.name),
      struct: `${camelize(config.name)}Resource`,
      model: `${camelize(config.name)}ResourceModel`,
    },
  };
}

function toComputedOptionalRequired({
  name,
  schemas,
}: {
  name: string;
  schemas: {
    read: oas30.SchemaObject;
    create?: oas30.SchemaObject;
    update?: oas30.SchemaObject;
  };
}): ComputedOptionalRequired {
  const inRead = schemas.read.properties && name in schemas.read.properties;
  const inCreate =
    schemas.create?.properties && name in schemas.create.properties;
  const inUpdate =
    schemas.update?.properties && name in schemas.update.properties;

  const reqCreate = schemas.create?.required?.includes(name);
  const reqUpdate = schemas.update?.required?.includes(name);

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
    `Unsupported computedOptionalRequired for field ${name}: ${JSON.stringify({ schemas, inRead, inCreate, inUpdate, reqCreate, reqUpdate })}`,
  );
}

export function mergeAttributeTypes(
  ...lists: AttributeType[][]
): AttributeType[] {
  const mergedMap = new Map<string, AttributeType>();

  for (const list of lists) {
    for (const attr of list) {
      const existing = mergedMap.get(attr.name);

      if (!existing) {
        mergedMap.set(attr.name, structuredClone(attr));
        continue;
      }

      // If types match and are complex nested types, recursively merge children
      if (
        existing.type === attr.type &&
        isNestedAttribute(attr) &&
        isNestedAttribute(existing)
      ) {
        const mergedChildren = mergeAttributeTypes(
          existing.attributes ?? [],
          attr.attributes ?? [],
        );

        const mergedBlocks = mergeNestedBlocks(
          existing.blocks ?? [],
          attr.blocks ?? [],
        );

        mergedMap.set(attr.name, {
          ...existing,
          ...attr, // Latter attribute shallow-overrides parent properties
          attributes: mergedChildren,
          blocks: mergedBlocks,
        } as AttributeType);
      } else {
        // Simple override (or type changed): latter completely replaces former
        mergedMap.set(attr.name, structuredClone(attr));
      }
    }
  }

  return Array.from(mergedMap.values());
}

function isNestedAttribute(
  attr: AttributeType,
): attr is AttributeObject | AttributeListNested | AttributeSetNested {
  return (
    attr.type === "object" ||
    attr.type === "list_nested" ||
    attr.type === "set_nested"
  );
}

export function isCollectionAttribute(
  attr: AttributeType,
): attr is AttributeListNested | AttributeSetNested {
  return attr.type === "list_nested" || attr.type === "set_nested";
}

function mergeNestedBlocks(
  existingBlocks: AttributeBlockType[],
  newBlocks: AttributeBlockType[],
): AttributeBlockType[] {
  const blockMap = new Map<string, AttributeBlockType>();

  for (const block of [...existingBlocks, ...newBlocks]) {
    const existing = blockMap.get(block.name);
    if (!existing) {
      blockMap.set(block.name, structuredClone(block));
      continue;
    }

    if (existing.type === block.type) {
      blockMap.set(block.name, {
        ...existing,
        ...block,
        attributes: mergeAttributeTypes(
          existing.attributes ?? [],
          block.attributes ?? [],
        ),
        blocks: mergeNestedBlocks(existing.blocks ?? [], block.blocks ?? []),
      });
    } else {
      blockMap.set(block.name, structuredClone(block));
    }
  }

  return Array.from(blockMap.values());
}

function cleanDescription(description: string | undefined): string | undefined {
  if (!description) {
    return undefined;
  }
  description = description.trim();
  if (!description.endsWith(".")) {
    description += ".";
  }
  return description;
}

export function withEnumDescription(
  description: string | undefined,
  enumValues: (string | number)[] | undefined,
): string | undefined {
  if (!enumValues || enumValues.length === 0) {
    return description;
  }
  if (description) {
    description += " ";
  } else {
    description = "";
  }
  description += `Value must be one of ${enumValues.map((v) => `\`${v}\``).join(", ")}.`;
  return description;
}

export function buildValidators(attribute: AttributeType): string[] {
  if (
    !("enum" in attribute) ||
    !attribute.enum ||
    attribute.enum.length === 0
  ) {
    return [];
  }

  return match(attribute)
    .returnType<string[]>()
    .with({ type: "string", enum: P.array(P.string) }, (schema) => [
      `stringvalidator.OneOf(${schema.enum.map((value) => JSON.stringify(value)).join(", ")})`,
      ...(schema.validators ?? []),
    ])
    .exhaustive();
}
