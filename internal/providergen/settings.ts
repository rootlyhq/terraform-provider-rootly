import type { ClientConfig, DataSourceConfig, ResourceConfig } from "./schema";

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
