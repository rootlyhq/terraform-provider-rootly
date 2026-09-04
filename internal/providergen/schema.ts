import type { oas30 } from "openapi3-ts";
import { match, P } from "ts-pattern";
import { getParametersByOperationId, removeReference } from "./openapi";
import assert from "node:assert";
import { camelize, pluralize, singularize } from "inflection";

export interface ClientConfig {
  name: string;
}

export type DataSourceConfig = {
  name: string;
  description?: string;
  strategy: "list" | "single";
  modifyDef?: (def: DataSourceDef) => DataSourceDef;
};

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
  hints?: {
    /** Indicates that this attribute is an ID from OpenAPI. */
    isOpenApiId?: boolean;
  };
}

export interface AttributeBool extends AttributeBase {
  type: "bool";
  default?: boolean;
}

export interface AttributeInt64 extends AttributeBase {
  type: "int64";
  default?: number;
}

export interface AttributeSingleNested extends AttributeBase {
  type: "single_nested";
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
  hints?: {
    /** Only used for data sources. Indicates that this is the list of the top level item type. */
    isTopLevelCollection?: boolean;
  };
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
  hints?: {
    /** Only used for data sources. Indicates that this is the list of the top level item type. */
    isTopLevelCollection?: boolean;
  };
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
  | AttributeSingleNested
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
  strategy: "single" | "list";
  attributes: AttributeType[];
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
      assert(tmpAttr.type === "single_nested");
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
          type: "single_nested",
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
    read?: oas30.ParameterObject[];
    list?: oas30.ParameterObject[];
    create?: oas30.ParameterObject[];
    update?: oas30.ParameterObject[];
    delete?: oas30.ParameterObject[];
  };
}): AttributeType[] {
  const allParamNames = Array.from(
    new Set([
      ...(params.read?.map((param) => param.name) ?? []),
      ...(params.list?.map((param) => param.name) ?? []),
      ...(params.create?.map((param) => param.name) ?? []),
      ...(params.update?.map((param) => param.name) ?? []),
      ...(params.delete?.map((param) => param.name) ?? []),
    ]),
  );

  return allParamNames.map((name) => {
    const readParam = params.read?.find((param) => param.name === name);
    const listParam = params.list?.find((param) => param.name === name);
    const createParam = params.create?.find((param) => param.name === name);
    const updateParam = params.update?.find((param) => param.name === name);
    const deleteParam = params.delete?.find((param) => param.name === name);

    const baseParam =
      readParam ?? listParam ?? createParam ?? updateParam ?? deleteParam!;
    const baseParamSchema = removeReference(baseParam.schema);
    assert(baseParamSchema, `Parameter ${name} has no schema`);
    assert(
      baseParamSchema.type === "string",
      `Parameter ${name} has no string type`,
    );

    const computedParam: Pick<
      AttributeType,
      "computedOptionalRequired" | "planModifiers"
    > = (() => {
      if (
        !createParam &&
        (readParam || listParam) &&
        !updateParam &&
        !deleteParam
      ) {
        return {
          computedOptionalRequired: "required",
        };
      }

      if (
        createParam &&
        (!readParam || !listParam) &&
        !updateParam &&
        !deleteParam
      ) {
        return {
          computedOptionalRequired: "required",
          planModifiers: ["stringplanmodifier.RequiresReplace()"],
        };
      }

      if (
        !createParam &&
        (readParam || listParam) &&
        updateParam &&
        deleteParam
      ) {
        return {
          computedOptionalRequired: "computed",
          planModifiers: ["stringplanmodifier.UseStateForUnknown()"],
        };
      }

      throw new Error(
        `Unsupported parameter ${name} computed optional required: ${JSON.stringify({ readParam, createParam, updateParam, deleteParam })}`,
      );
    })();

    const defaultDescription = match(name)
      .with("id", () => "The ID of the resource.")
      .otherwise(() => undefined);

    return {
      type: "string",
      name,
      computedOptionalRequired: computedParam.computedOptionalRequired,
      description:
        baseParam.description ??
        baseParamSchema.description ??
        defaultDescription,
      enum: baseParamSchema.enum,
      planModifiers: computedParam.planModifiers,
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
  return match(config)
    .returnType<DataSourceDef>()
    .with({ strategy: "single" }, (config) =>
      generateDataSourceDefSingle({
        doc,
        config,
      }),
    )
    .with({ strategy: "list" }, (config) =>
      generateDataSourceDefList({
        doc,
        config,
      }),
    )
    .exhaustive();
}

export function generateDataSourceDefSingle({
  doc,
  config,
}: {
  doc: oas30.OpenAPIObject;
  config: DataSourceConfig & { strategy: "single" };
}): DataSourceDef {
  // Get read schema
  const readSchema = removeReference(doc.components?.schemas?.[config.name]);
  if (!readSchema) {
    throw new Error(`Read schema "${config.name}" not found`);
  }

  // Get read path params
  const readPathParams = getParametersByOperationId({
    doc,
    operationId: `get${camelize(config.name)}`,
    onlyLocations: ["path"],
  });

  const tmpPathAttrs = openapiParametersToAttributes({
    params: {
      read: readPathParams,
    },
  });

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
  assert(tmpAttr.type === "single_nested");
  assert(tmpAttr.blocks.length === 0);

  return {
    name: config.name,
    description: cleanDescription(config.description ?? tmpAttr.description),
    strategy: config.strategy,
    attributes: mergeAttributeTypes(tmpAttr.attributes, tmpPathAttrs),
    goNames: {
      clientBase: camelize(config.name),
      struct: `${camelize(config.name)}DataSource`,
      model: `${camelize(config.name)}DataSourceModel`,
    },
  };
}

export function generateDataSourceDefList({
  doc,
  config,
}: {
  doc: oas30.OpenAPIObject;
  config: DataSourceConfig & { strategy: "list" };
}): DataSourceDef {
  const singularName = singularize(config.name);
  const pluralName = pluralize(config.name);

  // Get read schema
  const readSchema = removeReference(doc.components?.schemas?.[singularName]);
  if (!readSchema) {
    throw new Error(`Read schema "${singularName}" not found`);
  }

  // Get list path params
  const listPathParams = getParametersByOperationId({
    doc,
    operationId: `list${camelize(pluralName)}`,
    onlyLocations: ["path"],
  });

  const tmpPathAttrs = openapiParametersToAttributes({
    params: {
      list: listPathParams,
    },
  });

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
  assert(tmpAttr.type === "single_nested");
  assert(tmpAttr.blocks.length === 0);

  // TODO: Generalize
  const itemIdAttribute = {
    type: "string",
    name: "id",
    computedOptionalRequired: "computed",
    description: "The ID of the resource.",
    hints: { isOpenApiId: true },
  } satisfies AttributeString;

  const collectionAttribute = {
    ...tmpAttr,
    type: "list_nested",
    name: config.name,
    attributes: mergeAttributeTypes([itemIdAttribute], tmpAttr.attributes),
    hints: {
      isTopLevelCollection: true,
    },
  } satisfies AttributeListNested;

  return {
    name: config.name,
    description: cleanDescription(config.description ?? tmpAttr.description),
    strategy: config.strategy,
    attributes: mergeAttributeTypes(tmpPathAttrs, [collectionAttribute]),
    goNames: {
      clientBase: camelize(singularName),
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

  // Get list path params
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
  assert(tmpAttr.type === "single_nested");

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
): attr is AttributeSingleNested | AttributeListNested | AttributeSetNested {
  return (
    attr.type === "single_nested" ||
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

export function resolveDescription(
  attribute: AttributeType,
): string | undefined {
  const parts: string[] = [];

  // Add ForceNew description
  const hasRequiresReplace =
    attribute.planModifiers &&
    attribute.planModifiers.some((planModifier) =>
      planModifier.includes("RequiresReplace"),
    );
  if (hasRequiresReplace) {
    parts.push(`<i style="color:red;font-weight: bold">(ForceNew)</i>`);
  }

  // Add description
  if (attribute.description) {
    parts.push(attribute.description);
  }

  // Add enum description
  if ("enum" in attribute && attribute.enum && attribute.enum.length > 0) {
    parts.push(
      `Value must be one of ${attribute.enum.map((v) => `\`${v}\``).join(", ")}.`,
    );
  }

  return parts.length > 0 ? parts.join(" ") : undefined;
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
