import type { DataSourceConfig } from "../schema";

export default {
  kind: "data_source",
  name: "service",
  type: "single",
  description: "Retrieves a service.",
} satisfies DataSourceConfig;
