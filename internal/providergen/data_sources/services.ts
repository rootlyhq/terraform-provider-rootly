import type { DataSourceConfig } from "../schema";

export default {
  name: "services",
  strategy: "list",
  description: "Retrieves a list of all services.",
} satisfies DataSourceConfig;
