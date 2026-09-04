import { produce } from "immer";
import type {
  AttributeListNested,
  AttributeSetNested,
  AttributeString,
  ResourceConfig,
} from "../schema";

export default {
  name: "escalation_path",
  description: "Manages an escalation path.",
  modifyDef: produce((def) => {
    const rules = def.blocks.find(
      (v) => v.name === "rules",
    ) as AttributeListNested;
    rules.hacks = { nullBehavior: "empty" };

    const rulesOperator = rules.attributes.find(
      (v) => v.name === "operator",
    ) as AttributeString;
    rulesOperator.enum = [
      "is",
      "is_not",
      "contains",
      "does_not_contain",
      "is_one_of",
      "is_not_one_of",
      "is_empty",
      "is_not_empty",
      "contains_key",
      "does_not_contain_key",
      "starts_with",
      "does_not_start_with",
      "matches",
      "does_not_match",
      "is_set",
      "is_not_set",
    ];

    const rulesRuleType = rules.attributes.find(
      (v) => v.name === "rule_type",
    ) as AttributeString;
    rulesRuleType.enum = undefined;

    const timeRestrictions = def.blocks.find(
      (v) => v.name === "time_restrictions",
    ) as AttributeListNested;
    timeRestrictions.hacks = { nullBehavior: "empty" };

    const notificationTypeRules = def.blocks.find(
      (v) => v.name === "notification_type_rules",
    ) as AttributeListNested;
    notificationTypeRules.hacks = {
      unknownBehavior: "omit",
      nullBehavior: "omit",
      emptyBehavior: "omit",
    };

    const notificationTypeRulesConditions = notificationTypeRules.blocks.find(
      (v) => v.name === "conditions",
    ) as AttributeListNested;
    const notificationTypeRulesConditionsOperator =
      notificationTypeRulesConditions.attributes.find(
        (v) => v.name === "operator",
      ) as AttributeString;
    notificationTypeRulesConditionsOperator.enum = [
      "is",
      "is_not",
      "contains",
      "does_not_contain",
      "is_one_of",
      "is_not_one_of",
      "is_empty",
      "is_not_empty",
      "contains_key",
      "does_not_contain_key",
      "starts_with",
      "does_not_start_with",
      "matches",
      "does_not_match",
      "is_set",
      "is_not_set",
    ];

    const notificationTypeFallback = def.attributes.find(
      (v) => v.name === "notification_type_fallback",
    ) as AttributeString;
    notificationTypeFallback.default = undefined;
  }),
} satisfies ResourceConfig;
