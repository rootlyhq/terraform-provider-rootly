import type { DataSourceConfig } from "../schema";

export default {
  kind: "data_source",
  name: "services",
  type: "list",
  resourceName: "service",
  description: "Retrieves a list of all services.",
} satisfies DataSourceConfig;
