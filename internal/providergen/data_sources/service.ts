import type { DataSourceConfig } from "../schema";

export default {
  name: "service",
  strategy: "single",
  description: "Retrieves a service.",
} satisfies DataSourceConfig;
