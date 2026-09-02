import type { DataSourceConfig } from "../schema";

export default {
  name: "services",
  strategy: "list",
  resourceName: "service",
  description: "Retrieves a list of all services.",
} satisfies DataSourceConfig;
