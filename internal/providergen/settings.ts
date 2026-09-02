import { camelize, humanize, singularize } from "inflection";
import type {
  ClientConfig,
  ComputedOptionalRequired,
  DataSourceConfig,
  ResolvedDataSourceConfig,
  ResolvedResourceConfig,
  ResourceConfig,
} from "./schema";
import type { oas30 } from "openapi3-ts";
import { removeReference } from "./openapi";
import { match, P } from "ts-pattern";

async function importDefaults<T>(pattern: string): Promise<T[]> {
  const files = await Array.fromAsync(
    new Bun.Glob(pattern).scan({ cwd: import.meta.dir }),
  );
  return Promise.all(files.map(async (file) => (await import(file)).default));
}

export const CLIENTS = await importDefaults<ClientConfig>("./clients/*.ts");
export const DATA_SOURCES = await importDefaults<DataSourceConfig>(
  "./data_sources/*.ts",
);
export const RESOURCES =
  await importDefaults<ResourceConfig>("./resources/*.ts");

export function resolveDataSourceConfig({
  doc,
  config,
}: {
  doc: oas30.OpenAPIObject;
  config: DataSourceConfig;
}): ResolvedDataSourceConfig {
  const isSingle = config.strategy === "single";
  const schemaKey = isSingle ? config.name : config.resourceName;

  const read = removeReference(doc.components?.schemas?.[schemaKey]);
  if (!read) {
    throw new Error(`Cannot find schema: ${schemaKey} in doc`);
  }

  const resolved = resolveSchema({
    doc,
    config,
    type: "data_source",
    read,
    create: undefined,
    update: undefined,
    isTopLevel: true,
  });

  return {
    ...config,
    type: "data_source",
    goNames: {
      struct: `${camelize(config.name)}DataSource`,
      model: `${camelize(config.name)}DataSourceModel`,
    },
    schemas: {
      read,
      resolved,
    },
  };
}

export function resolveResourceConfig({
  doc,
  config,
}: {
  doc: oas30.OpenAPIObject;
  config: ResourceConfig;
}): ResolvedResourceConfig {
  const create = removeReference(
    removeReference(
      removeReference(doc.components?.schemas?.[`new_${config.name}`])
        ?.properties?.data,
    )?.properties?.attributes,
  );
  if (!create) {
    throw new Error(
      `Could not find schema for resource ${config.name} creation`,
    );
  }

  const read = removeReference(doc.components?.schemas?.[config.name]);
  if (!read) {
    throw new Error(
      `Could not find schema for resource ${config.name} reading`,
    );
  }

  const update = removeReference(
    removeReference(
      removeReference(doc.components?.schemas?.[`update_${config.name}`])
        ?.properties?.data,
    )?.properties?.attributes,
  );
  if (!update) {
    throw new Error(
      `Could not find schema for resource ${config.name} updating`,
    );
  }

  const resolved = resolveSchema({
    doc,
    config,
    type: "resource",
    read,
    create,
    update,
    isTopLevel: true,
  });

  return {
    ...config,
    type: "resource",
    goNames: {
      struct: `${camelize(config.name)}Resource`,
      model: `${camelize(config.name)}ResourceModel`,
    },
    schemas: { create, read, update, resolved },
  };
}

function resolveSchema({
  doc,
  config,
  type,
  read,
  create,
  update,
  isTopLevel,
}: {
  doc: oas30.OpenAPIObject;
  config: DataSourceConfig | ResourceConfig;
  type: "data_source" | "resource";
  read: oas30.SchemaObject;
  create: oas30.SchemaObject | undefined;
  update: oas30.SchemaObject | undefined;
  isTopLevel: boolean;
}): oas30.SchemaObject {
  // Special handling for top level schemas
  if (isTopLevel) {
    return match([type, config])
      .with(["data_source", { strategy: "list" }], () => {
        const singleSchema = resolveSchema({
          doc,
          config,
          type,
          read,
          create,
          update,
          isTopLevel: false,
        });

        return {
          type: "object",
          properties: {
            [config.name]: {
              type: "array",
              items: {
                ...singleSchema,
                properties: {
                  id: {
                    type: "string",
                    description: `The ID of the ${humanize(singularize(config.name), true)}.`,
                    "x-tf-computed-optional-required": "computed",
                  },
                  ...singleSchema.properties,
                },
                required: ["id", ...(singleSchema.required ?? [])],
              },
              "x-tf-top-level-item-type": true,
              "x-tf-collection-type": "list",
            },
          },
        } satisfies oas30.SchemaObject;
      })
      .with(["data_source", { strategy: "single" }], () =>
        resolveSchema({
          doc,
          config,
          type,
          read: {
            ...read,
            properties: {
              id: {
                type: "string",
                description: `The ID of the ${humanize(config.name, true)}.`,
                "x-tf-computed-optional-required": "required",
              },
              ...read.properties,
            },
          },
          create,
          update,
          isTopLevel: false,
        }),
      )
      .with(["resource", P.any], () =>
        resolveSchema({
          doc,
          config,
          type,
          read: {
            ...read,
            properties: {
              id: {
                type: "string",
                description: `The ID of the ${humanize(config.name, true)}.`,
                "x-tf-computed-optional-required": "computed",
              },
              ...read.properties,
            },
          },
          create,
          update,
          isTopLevel: false,
        }),
      )
      .exhaustive();
  }

  const schema = match(read)
    .with({ type: "array", items: P.record(P.string, P.any) }, (schema) => ({
      ...schema,
      items: resolveSchema({
        doc,
        config,
        type,
        read: removeReference(schema.items),
        create: removeReference(create?.items),
        update: removeReference(update?.items),
        isTopLevel: false,
      }),
      "x-tf-collection-type": type === "data_source" ? "list" : "set",
    }))
    .with(
      { type: "object", properties: P.record(P.string, P.any) },
      (schema) => {
        return {
          ...schema,
          properties: Object.fromEntries(
            Object.entries(schema.properties).map(([key, value]) => {
              return [
                key,
                {
                  ...resolveSchema({
                    doc,
                    config,
                    type,
                    read: removeReference(value),
                    create: removeReference(create?.properties?.[key]),
                    update: removeReference(update?.properties?.[key]),
                    isTopLevel: false,
                  }),
                  "x-tf-computed-optional-required":
                    resolveComputedOptionalRequired({
                      doc,
                      config,
                      name: key,
                      read,
                      create,
                      update,
                    }),
                },
              ];
            }),
          ),
        };
      },
    )
    .with(
      { type: P.union("string", "integer", "number", "boolean") },
      (schema) => schema,
    )
    .exhaustive();

  return {
    ...schema,
    "x-schema-read": read,
    "x-schema-create": create,
    "x-schema-update": update,
  };
}

function resolveComputedOptionalRequired({
  doc,
  config,
  name,
  read,
  create,
  update,
}: {
  doc: oas30.OpenAPIObject;
  config: DataSourceConfig | ResourceConfig;
  name: string;
  read: oas30.SchemaObject;
  create: oas30.SchemaObject | undefined;
  update: oas30.SchemaObject | undefined;
}): ComputedOptionalRequired {
  const inRead = read.properties && name in read.properties;
  const inCreate = create?.properties && name in create.properties;
  const inUpdate = update?.properties && name in update.properties;
  const reqCreate = create?.required && create.required.includes(name);
  const reqUpdate = update?.required && update.required.includes(name);

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
    `Cannot determine computed/optional/required for property ${name}: ${JSON.stringify({ read, create, update })}`,
  );
}
