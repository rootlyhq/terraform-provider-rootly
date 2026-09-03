import { produce } from "immer";
import type {
  AttributeSetNested,
  AttributeString,
  ResourceConfig,
} from "../schema";

export default {
  name: "escalation_path",
  description: "Manages an escalation path.",
  modifyDef: produce((def) => {
    const rules = def.blocks.find(
      (attr) => attr.name === "rules",
    ) as AttributeSetNested;
    rules.hacks = { nullBehavior: "empty" };

    const timeRestrictions = def.blocks.find(
      (attr) => attr.name === "time_restrictions",
    ) as AttributeSetNested;
    timeRestrictions.hacks = { nullBehavior: "empty" };

    const notificationTypeRules = def.blocks.find(
      (attr) => attr.name === "notification_type_rules",
    ) as AttributeSetNested;
    notificationTypeRules.hacks = {
      unknownBehavior: "omit",
      nullBehavior: "omit",
      emptyBehavior: "omit",
    };

    const notificationTypeFallback = def.attributes.find(
      (attr) => attr.name === "notification_type_fallback",
    ) as AttributeString;
    notificationTypeFallback.default = undefined;
  }),
} satisfies ResourceConfig;
