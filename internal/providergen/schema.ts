import type { oas30 } from "openapi3-ts";
import { match, P } from "ts-pattern";
import { removeReference } from "./openapi";
import assert from "node:assert";

export interface ClientConfig {
  name: string;
  actions?: {
    list?: {
      enabled: true;
    };
    get?: {
      enabled: true;
    };
    create?: {
      enabled: true;
    };
    update?: {
      enabled: true;
    };
  };
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
} & (DataSourceListConfig | DataSourceSingleConfig);

export type ResolvedDataSourceConfig = {
  type: "data_source";
  name: string;
  description?: string;
  goNames: {
    /** Name of the struct that represents the data source. */
    struct: `${string}DataSource`;
    /** Name of the struct that represents the model of the data source. */
    model: `${string}DataSourceModel`;
  };
  schemas: {
    read: oas30.SchemaObject;
    resolved: oas30.SchemaObject;
  };
} & (DataSourceListConfig | DataSourceSingleConfig);

// TODO
export type ResourceConfig = {
  name: string;
  description?: string;
};

export type ResolvedResourceConfig = {
  type: "resource";
  name: string;
  description?: string;
  goNames: {
    /** Name of the struct that represents the resource. */
    struct: `${string}Resource`;
    /** Name of the struct that represents the model of the resource. */
    model: `${string}ResourceModel`;
  };
  schemas: {
    create: oas30.SchemaObject;
    read: oas30.SchemaObject;
    update: oas30.SchemaObject;
    resolved: oas30.SchemaObject;
  };
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
  nullable?: string;
  computedOptionalRequired: ComputedOptionalRequired;
  validators?: string[];
  planModifiers?: string[];
  schemas?: {
    read: oas30.SchemaObject;
    create?: oas30.SchemaObject;
    update?: oas30.SchemaObject;
  };
}

export interface AttributeString extends AttributeBase {
  type: "string";
  enum?: string[];
}

export interface AttributeBool extends AttributeBase {
  type: "bool";
}

export interface AttributeInt64 extends AttributeBase {
  type: "int64";
}

export interface AttributeObject extends AttributeBase {
  type: "object";
  attributes: AttributeType[];
  blocks: (AttributeListNested | AttributeSetNested)[];
}

export interface AttributeList extends AttributeBase {
  type: "list";
  elementType: "string";
}

export interface AttributeSet extends AttributeBase {
  type: "set";
  elementType: "string";
}

export interface AttributeListNested extends AttributeBase {
  type: "list_nested";
  attributes: AttributeType[];
  blocks: (AttributeListNested | AttributeSetNested)[];
}

export interface AttributeSetNested extends AttributeBase {
  type: "set_nested";
  attributes: AttributeType[];
  blocks: (AttributeListNested | AttributeSetNested)[];
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

export type ComputedOptionalRequired =
  | "computed"
  | "optional"
  | "computed_optional"
  | "required";

export interface ResourceDef {
  name: string;
  description?: string;
  attributes: AttributeType[];
  blocks: (AttributeListNested | AttributeSetNested)[];
}

function openapiToAttribute({
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
  };
}): AttributeType {
  const common: Pick<
    AttributeBase,
    "name" | "description" | "computedOptionalRequired"
  > = {
    name,
    description: schemas.read.description,
    computedOptionalRequired,
  };

  return match(schemas.read)
    .returnType<AttributeType>()
    .with(
      {
        type: "string",
        enum: P.array(P.string).optional(),
      },
      (schema) => ({
        ...common,
        type: "string",
        enum: schema.enum,
      }),
    )
    .with({ type: "integer" }, () => ({
      ...common,
      type: "int64",
    }))
    .with({ type: "boolean" }, () => ({
      ...common,
      type: "bool",
    }))
    .with({ type: "array", items: { type: "string" } }, () => ({
      ...common,
      type: options.defaultCollectionType,
      elementType: "string",
    }))
    .with({ type: "array", items: { type: "object" } }, (schema) => {
      const tmpAttr = openapiToAttribute({
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
            openapiToAttribute({
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

        const collectionTypes = ["list_nested", "set_nested"];
        const attributes = options.collectionsAsBlocks
          ? allAttributes.filter((attr) => !collectionTypes.includes(attr.type))
          : allAttributes;
        const blocks = options.collectionsAsBlocks
          ? allAttributes.filter((attr) => collectionTypes.includes(attr.type))
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

  const tmpAttr = openapiToAttribute({
    name: config.name,
    computedOptionalRequired: "required",
    schemas: {
      read: readSchema,
      create: createSchema,
      update: updateSchema,
    },
    options: {
      defaultCollectionType: "set",
      defaultNestedCollectionType: "set_nested",
      collectionsAsBlocks: true,
    },
  });
  assert(tmpAttr.type === "object");

  return {
    name: config.name,
    description: config.description,
    attributes: tmpAttr.attributes,
    blocks: tmpAttr.blocks,
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
