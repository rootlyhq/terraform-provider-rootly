import { camelize } from "inflection";
import type {
  ClientConfig,
  DataSourceConfig,
  ResolvedDataSourceConfig,
  ResolvedResourceConfig,
  ResourceConfig,
} from "./schema";

async function importDefaults<T>(pattern: string): Promise<T[]> {
  const files = await Array.fromAsync(
    new Bun.Glob(pattern).scan({ cwd: import.meta.dir }),
  );
  return Promise.all(files.map(async (file) => (await import(file)).default));
}

export const CLIENTS = await importDefaults<ClientConfig>("./clients/*.ts");
export const DATA_SOURCES = (
  await importDefaults<DataSourceConfig>("./data_sources/*.ts")
).map(resolveDataSourceConfig);
export const RESOURCES = (
  await importDefaults<ResourceConfig>("./resources/*.ts")
).map(resolveResourceConfig);

function resolveDataSourceConfig(
  config: DataSourceConfig,
): ResolvedDataSourceConfig {
  return {
    ...config,
    kind: "data_source",
    goNames: {
      struct: `${camelize(config.name)}DataSource`,
      model: `${camelize(config.name)}DataSourceModel`,
    },
  };
}

function resolveResourceConfig(config: ResourceConfig): ResolvedResourceConfig {
  return {
    ...config,
    kind: "resource",
    goNames: {
      struct: `${camelize(config.name)}Resource`,
      model: `${camelize(config.name)}ResourceModel`,
    },
  };
}
