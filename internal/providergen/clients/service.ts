import type { ClientConfig } from "../schema";

export default {
  name: "service",
  actions: {
    list: {
      enabled: true,
    },
    get: {
      enabled: true,
    },
  },
} satisfies ClientConfig;
