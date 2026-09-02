import type { ClientConfig } from "../schema";

export default {
  name: "escalation_path",
  actions: {
    list: {
      enabled: true,
    },
    get: {
      enabled: true,
    },
    create: {
      enabled: true,
    },
    update: {
      enabled: true,
    },
  },
} satisfies ClientConfig;
