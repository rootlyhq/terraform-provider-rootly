export type ValidIdentifier = string;
export type SchemaComputedOptionalRequired =
  | "computed"
  | "computed_optional"
  | "optional"
  | "required";
export type CodeImports = [CodeImport, ...CodeImport[]];
export type SchemaBoolValidators = SchemaBoolValidator[];
export type SchemaDynamicValidators = SchemaDynamicValidator[];
export type SchemaFloat64Validators = SchemaFloat64Validator[];
export type SchemaInt64Validators = SchemaInt64Validator[];
export type SchemaElementType =
  | {
      [k: string]: unknown;
    }
  | {
      [k: string]: unknown;
    }
  | {
      [k: string]: unknown;
    }
  | {
      [k: string]: unknown;
    }
  | {
      [k: string]: unknown;
    }
  | {
      [k: string]: unknown;
    }
  | {
      [k: string]: unknown;
    }
  | {
      [k: string]: unknown;
    }
  | {
      [k: string]: unknown;
    };
export type SchemaListValidators = SchemaListValidator[];
export type SchemaObjectValidators = SchemaObjectValidator[];
export type SchemaMapValidators = SchemaMapValidator[];
export type SchemaNumberValidators = SchemaNumberValidator[];
export type SchemaObjectAttributeTypes = [
  SchemaObjectAttributeType,
  ...SchemaObjectAttributeType[]
];
export type SchemaObjectAttributeType =
  | {
      [k: string]: unknown;
    }
  | {
      [k: string]: unknown;
    }
  | {
      [k: string]: unknown;
    }
  | {
      [k: string]: unknown;
    }
  | {
      [k: string]: unknown;
    }
  | {
      [k: string]: unknown;
    }
  | {
      [k: string]: unknown;
    }
  | {
      [k: string]: unknown;
    }
  | {
      [k: string]: unknown;
    }
  | {
      [k: string]: unknown;
    };
export type SchemaSetValidators = SchemaSetValidator[];
export type SchemaStringValidators = SchemaStringValidator[];
export type DatasourceAttributes = (
  | DatasourceBoolAttribute
  | DatasourceDynamicAttribute
  | DatasourceFloat64Attribute
  | DatasourceInt64Attribute
  | DatasourceListAttribute
  | DatasourceListNestedAttribute
  | DatasourceMapAttribute
  | DatasourceMapNestedAttribute
  | DatasourceNumberAttribute
  | DatasourceObjectAttribute
  | DatasourceSetAttribute
  | DatasourceSetNestedAttribute
  | DatasourceSingleNestedAttribute
  | DatasourceStringAttribute
)[];
export type DatasourceBlocks = (
  | DatasourceListNestedBlock
  | DatasourceSetNestedBlock
  | DatasourceSingleNestedBlock
)[];
export type SchemaOptionalRequired = "optional" | "required";
export type ProviderAttributes = (
  | ProviderBoolAttribute
  | ProviderDynamicAttribute
  | ProviderFloat64Attribute
  | ProviderInt64Attribute
  | ProviderListAttribute
  | ProviderListNestedAttribute
  | ProviderMapAttribute
  | ProviderMapNestedAttribute
  | ProviderNumberAttribute
  | ProviderObjectAttribute
  | ProviderSetAttribute
  | ProviderSetNestedAttribute
  | ProviderSingleNestedAttribute
  | ProviderStringAttribute
)[];
export type ProviderBlocks = (
  | ProviderListNestedBlock
  | ProviderSetNestedBlock
  | ProviderSingleNestedBlock
)[];
export type SchemaBoolDefault =
  | {
      [k: string]: unknown;
    }
  | {
      [k: string]: unknown;
    };
export type SchemaBoolPlanModifiers = SchemaBoolPlanModifier[];
export type SchemaDynamicDefault = {
  [k: string]: unknown;
};
export type SchemaDynamicPlanModifiers = SchemaDynamicPlanModifier[];
export type SchemaFloat64Default =
  | {
      [k: string]: unknown;
    }
  | {
      [k: string]: unknown;
    };
export type SchemaFloat64PlanModifiers = SchemaFloat64PlanModifier[];
export type SchemaInt64Default =
  | {
      [k: string]: unknown;
    }
  | {
      [k: string]: unknown;
    };
export type SchemaInt64PlanModifiers = SchemaInt64PlanModifier[];
export type SchemaListDefault = {
  [k: string]: unknown;
};
export type SchemaListPlanModifiers = SchemaListPlanModifier[];
export type SchemaObjectPlanModifiers = SchemaObjectPlanModifier[];
export type SchemaMapDefault = {
  [k: string]: unknown;
};
export type SchemaMapPlanModifiers = SchemaMapPlanModifier[];
export type SchemaNumberDefault = {
  [k: string]: unknown;
};
export type SchemaNumberPlanModifiers = SchemaNumberPlanModifier[];
export type SchemaObjectDefault = {
  [k: string]: unknown;
};
export type SchemaSetDefault = {
  [k: string]: unknown;
};
export type SchemaSetPlanModifiers = SchemaSetPlanModifier[];
export type SchemaStringDefault =
  | {
      [k: string]: unknown;
    }
  | {
      [k: string]: unknown;
    };
export type SchemaStringPlanModifiers = SchemaStringPlanModifier[];
export type ResourceAttributes = (
  | ResourceBoolAttribute
  | ResourceDynamicAttribute
  | ResourceFloat64Attribute
  | ResourceInt64Attribute
  | ResourceListAttribute
  | ResourceListNestedAttribute
  | ResourceMapAttribute
  | ResourceMapNestedAttribute
  | ResourceNumberAttribute
  | ResourceObjectAttribute
  | ResourceSetAttribute
  | ResourceSetNestedAttribute
  | ResourceSingleNestedAttribute
  | ResourceStringAttribute
)[];
export type ResourceBlocks = (
  | ResourceListNestedBlock
  | ResourceSetNestedBlock
  | ResourceSingleNestedBlock
)[];

export interface MySchema {
  datasources?: Datasource[];
  provider: Provider;
  resources?: Resource[];
  version: string;
  [k: string]: unknown;
}
export interface Datasource {
  name: ValidIdentifier;
  schema: {
    attributes?: DatasourceAttributes;
    blocks?: DatasourceBlocks;
    description?: string;
    markdown_description?: string;
    deprecation_message?: string;
    [k: string]: unknown;
  };
  [k: string]: unknown;
}
export interface DatasourceBoolAttribute {
  name: ValidIdentifier;
  bool: {
    associated_external_type?: SchemaAssociatedExternalType;
    computed_optional_required: SchemaComputedOptionalRequired;
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    sensitive?: boolean;
    validators?: SchemaBoolValidators;
  };
}
export interface SchemaAssociatedExternalType {
  import?: CodeImport;
  type: string;
}
export interface CodeImport {
  alias?: string;
  path: string;
  [k: string]: unknown;
}
export interface SchemaCustomType {
  import?: CodeImport;
  type: string;
  value_type: string;
  [k: string]: unknown;
}
export interface SchemaBoolValidator {
  custom: SchemaCustomValidator;
  [k: string]: unknown;
}
export interface SchemaCustomValidator {
  imports?: CodeImports;
  schema_definition: string;
  [k: string]: unknown;
}
export interface DatasourceDynamicAttribute {
  name: ValidIdentifier;
  dynamic: {
    associated_external_type?: SchemaAssociatedExternalType;
    computed_optional_required: SchemaComputedOptionalRequired;
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    sensitive?: boolean;
    validators?: SchemaDynamicValidators;
  };
}
export interface SchemaDynamicValidator {
  custom: SchemaCustomValidator;
  [k: string]: unknown;
}
export interface DatasourceFloat64Attribute {
  name: ValidIdentifier;
  float64: {
    associated_external_type?: SchemaAssociatedExternalType;
    computed_optional_required: SchemaComputedOptionalRequired;
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    sensitive?: boolean;
    validators?: SchemaFloat64Validators;
  };
}
export interface SchemaFloat64Validator {
  custom: SchemaCustomValidator;
  [k: string]: unknown;
}
export interface DatasourceInt64Attribute {
  name: ValidIdentifier;
  int64: {
    associated_external_type?: SchemaAssociatedExternalType;
    computed_optional_required: SchemaComputedOptionalRequired;
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    sensitive?: boolean;
    validators?: SchemaInt64Validators;
  };
}
export interface SchemaInt64Validator {
  custom: SchemaCustomValidator;
  [k: string]: unknown;
}
export interface DatasourceListAttribute {
  name: ValidIdentifier;
  list: {
    associated_external_type?: SchemaAssociatedExternalType;
    computed_optional_required: SchemaComputedOptionalRequired;
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    element_type: SchemaElementType;
    sensitive?: boolean;
    validators?: SchemaListValidators;
  };
}
export interface SchemaListValidator {
  custom: SchemaCustomValidator;
  [k: string]: unknown;
}
export interface DatasourceListNestedAttribute {
  name: ValidIdentifier;
  list_nested: {
    computed_optional_required: SchemaComputedOptionalRequired;
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    nested_object: DatasourceNestedAttributeObject;
    sensitive?: boolean;
    validators?: SchemaListValidators;
  };
}
export interface DatasourceNestedAttributeObject {
  associated_external_type?: SchemaAssociatedExternalType;
  attributes?: DatasourceAttributes;
  custom_type?: SchemaCustomType;
  validators?: SchemaObjectValidators;
  [k: string]: unknown;
}
export interface SchemaObjectValidator {
  custom: SchemaCustomValidator;
  [k: string]: unknown;
}
export interface DatasourceMapAttribute {
  name: ValidIdentifier;
  map: {
    associated_external_type?: SchemaAssociatedExternalType;
    computed_optional_required: SchemaComputedOptionalRequired;
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    element_type: SchemaElementType;
    sensitive?: boolean;
    validators?: SchemaMapValidators;
  };
}
export interface SchemaMapValidator {
  custom: SchemaCustomValidator;
  [k: string]: unknown;
}
export interface DatasourceMapNestedAttribute {
  name: ValidIdentifier;
  map_nested: {
    computed_optional_required: SchemaComputedOptionalRequired;
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    nested_object: DatasourceNestedAttributeObject;
    sensitive?: boolean;
    validators?: SchemaMapValidators;
  };
}
export interface DatasourceNumberAttribute {
  name: ValidIdentifier;
  number: {
    associated_external_type?: SchemaAssociatedExternalType;
    computed_optional_required: SchemaComputedOptionalRequired;
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    sensitive?: boolean;
    validators?: SchemaNumberValidators;
  };
}
export interface SchemaNumberValidator {
  custom: SchemaCustomValidator;
  [k: string]: unknown;
}
export interface DatasourceObjectAttribute {
  name: ValidIdentifier;
  object: {
    associated_external_type?: SchemaAssociatedExternalType;
    attribute_types: SchemaObjectAttributeTypes;
    computed_optional_required: SchemaComputedOptionalRequired;
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    sensitive?: boolean;
    validators?: SchemaObjectValidators;
  };
}
export interface DatasourceSetAttribute {
  name: ValidIdentifier;
  set: {
    associated_external_type?: SchemaAssociatedExternalType;
    computed_optional_required: SchemaComputedOptionalRequired;
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    element_type: SchemaElementType;
    sensitive?: boolean;
    validators?: SchemaSetValidators;
  };
}
export interface SchemaSetValidator {
  custom: SchemaCustomValidator;
  [k: string]: unknown;
}
export interface DatasourceSetNestedAttribute {
  name: ValidIdentifier;
  set_nested: {
    computed_optional_required: SchemaComputedOptionalRequired;
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    nested_object: DatasourceNestedAttributeObject;
    sensitive?: boolean;
    validators?: SchemaSetValidators;
  };
}
export interface DatasourceSingleNestedAttribute {
  name: ValidIdentifier;
  single_nested: {
    associated_external_type?: SchemaAssociatedExternalType;
    attributes?: DatasourceAttributes;
    computed_optional_required: SchemaComputedOptionalRequired;
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    sensitive?: boolean;
    validators?: SchemaObjectValidators;
  };
}
export interface DatasourceStringAttribute {
  name: ValidIdentifier;
  string: {
    associated_external_type?: SchemaAssociatedExternalType;
    computed_optional_required: SchemaComputedOptionalRequired;
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    sensitive?: boolean;
    validators?: SchemaStringValidators;
  };
}
export interface SchemaStringValidator {
  custom: SchemaCustomValidator;
  [k: string]: unknown;
}
export interface DatasourceListNestedBlock {
  name: ValidIdentifier;
  list_nested: {
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    nested_object: DatasourceNestedBlockObject;
    validators?: SchemaListValidators;
  };
}
export interface DatasourceNestedBlockObject {
  associated_external_type?: SchemaAssociatedExternalType;
  attributes?: DatasourceAttributes;
  blocks?: DatasourceBlocks;
  custom_type?: SchemaCustomType;
  validators?: SchemaObjectValidators;
  [k: string]: unknown;
}
export interface DatasourceSetNestedBlock {
  name: ValidIdentifier;
  set_nested: {
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    nested_object: DatasourceNestedBlockObject;
    validators?: SchemaSetValidators;
  };
}
export interface DatasourceSingleNestedBlock {
  name: ValidIdentifier;
  single_nested: {
    associated_external_type?: SchemaAssociatedExternalType;
    attributes?: DatasourceAttributes;
    blocks?: DatasourceBlocks;
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    validators?: SchemaObjectValidators;
  };
}
export interface Provider {
  name: ValidIdentifier;
  schema?: {
    attributes?: ProviderAttributes;
    blocks?: ProviderBlocks;
    description?: string;
    markdown_description?: string;
    deprecation_message?: string;
    [k: string]: unknown;
  };
  [k: string]: unknown;
}
export interface ProviderBoolAttribute {
  name: ValidIdentifier;
  bool: {
    associated_external_type?: SchemaAssociatedExternalType;
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    optional_required: SchemaOptionalRequired;
    sensitive?: boolean;
    validators?: SchemaBoolValidators;
  };
}
export interface ProviderDynamicAttribute {
  name: ValidIdentifier;
  dynamic: {
    associated_external_type?: SchemaAssociatedExternalType;
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    optional_required: SchemaOptionalRequired;
    sensitive?: boolean;
    validators?: SchemaDynamicValidators;
  };
}
export interface ProviderFloat64Attribute {
  name: ValidIdentifier;
  float64: {
    associated_external_type?: SchemaAssociatedExternalType;
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    optional_required: SchemaOptionalRequired;
    sensitive?: boolean;
    validators?: SchemaFloat64Validators;
  };
}
export interface ProviderInt64Attribute {
  name: ValidIdentifier;
  int64: {
    associated_external_type?: SchemaAssociatedExternalType;
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    optional_required: SchemaOptionalRequired;
    sensitive?: boolean;
    validators?: SchemaInt64Validators;
  };
}
export interface ProviderListAttribute {
  name: ValidIdentifier;
  list: {
    associated_external_type?: SchemaAssociatedExternalType;
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    element_type: SchemaElementType;
    optional_required: SchemaOptionalRequired;
    sensitive?: boolean;
    validators?: SchemaListValidators;
  };
}
export interface ProviderListNestedAttribute {
  name: ValidIdentifier;
  list_nested: {
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    nested_object: ProviderNestedAttributeObject;
    optional_required: SchemaOptionalRequired;
    sensitive?: boolean;
    validators?: SchemaListValidators;
  };
}
export interface ProviderNestedAttributeObject {
  associated_external_type?: SchemaAssociatedExternalType;
  attributes?: ProviderAttributes;
  custom_type?: SchemaCustomType;
  validators?: SchemaObjectValidators;
  [k: string]: unknown;
}
export interface ProviderMapAttribute {
  name: ValidIdentifier;
  map: {
    associated_external_type?: SchemaAssociatedExternalType;
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    element_type: SchemaElementType;
    optional_required: SchemaOptionalRequired;
    sensitive?: boolean;
    validators?: SchemaMapValidators;
  };
}
export interface ProviderMapNestedAttribute {
  name: ValidIdentifier;
  map_nested: {
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    nested_object: ProviderNestedAttributeObject;
    optional_required: SchemaOptionalRequired;
    sensitive?: boolean;
    validators?: SchemaMapValidators;
  };
}
export interface ProviderNumberAttribute {
  name: ValidIdentifier;
  number: {
    associated_external_type?: SchemaAssociatedExternalType;
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    optional_required: SchemaOptionalRequired;
    sensitive?: boolean;
    validators?: SchemaNumberValidators;
  };
}
export interface ProviderObjectAttribute {
  name: ValidIdentifier;
  object: {
    associated_external_type?: SchemaAssociatedExternalType;
    attribute_types: SchemaObjectAttributeTypes;
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    optional_required: SchemaOptionalRequired;
    sensitive?: boolean;
    validators?: SchemaObjectValidators;
  };
}
export interface ProviderSetAttribute {
  name: ValidIdentifier;
  set: {
    associated_external_type?: SchemaAssociatedExternalType;
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    element_type: SchemaElementType;
    optional_required: SchemaOptionalRequired;
    sensitive?: boolean;
    validators?: SchemaSetValidators;
  };
}
export interface ProviderSetNestedAttribute {
  name: ValidIdentifier;
  set_nested: {
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    nested_object: ProviderNestedAttributeObject;
    optional_required: SchemaOptionalRequired;
    sensitive?: boolean;
    validators?: SchemaSetValidators;
  };
}
export interface ProviderSingleNestedAttribute {
  name: ValidIdentifier;
  single_nested: {
    associated_external_type?: SchemaAssociatedExternalType;
    attributes?: ProviderAttributes;
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    optional_required: SchemaOptionalRequired;
    sensitive?: boolean;
    validators?: SchemaObjectValidators;
  };
}
export interface ProviderStringAttribute {
  name: ValidIdentifier;
  string: {
    associated_external_type?: SchemaAssociatedExternalType;
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    optional_required: SchemaOptionalRequired;
    sensitive?: boolean;
    validators?: SchemaStringValidators;
  };
}
export interface ProviderListNestedBlock {
  name: ValidIdentifier;
  list_nested: {
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    nested_object: ProviderNestedBlockObject;
    validators?: SchemaListValidators;
  };
}
export interface ProviderNestedBlockObject {
  associated_external_type?: SchemaAssociatedExternalType;
  attributes?: ProviderAttributes;
  blocks?: ProviderBlocks;
  custom_type?: SchemaCustomType;
  validators?: SchemaObjectValidators;
  [k: string]: unknown;
}
export interface ProviderSetNestedBlock {
  name: ValidIdentifier;
  set_nested: {
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    nested_object: ProviderNestedBlockObject;
    validators?: SchemaSetValidators;
  };
}
export interface ProviderSingleNestedBlock {
  name: ValidIdentifier;
  single_nested: {
    associated_external_type?: SchemaAssociatedExternalType;
    attributes?: ProviderAttributes;
    blocks?: ProviderBlocks;
    custom_type?: SchemaCustomType;
    deprecation_message?: string;
    description?: string;
    validators?: SchemaObjectValidators;
  };
}
export interface Resource {
  name: ValidIdentifier;
  schema: {
    attributes?: ResourceAttributes;
    blocks?: ResourceBlocks;
    description?: string;
    markdown_description?: string;
    deprecation_message?: string;
    [k: string]: unknown;
  };
  [k: string]: unknown;
}
export interface ResourceBoolAttribute {
  name: ValidIdentifier;
  bool: {
    associated_external_type?: SchemaAssociatedExternalType;
    computed_optional_required: SchemaComputedOptionalRequired;
    custom_type?: SchemaCustomType;
    default?: SchemaBoolDefault;
    deprecation_message?: string;
    description?: string;
    plan_modifiers?: SchemaBoolPlanModifiers;
    sensitive?: boolean;
    validators?: SchemaBoolValidators;
  };
}
export interface SchemaBoolPlanModifier {
  custom: SchemaCustomPlanModifier;
  [k: string]: unknown;
}
export interface SchemaCustomPlanModifier {
  imports?: CodeImports;
  schema_definition: string;
  [k: string]: unknown;
}
export interface ResourceDynamicAttribute {
  name: ValidIdentifier;
  dynamic: {
    associated_external_type?: SchemaAssociatedExternalType;
    computed_optional_required: SchemaComputedOptionalRequired;
    custom_type?: SchemaCustomType;
    default?: SchemaDynamicDefault;
    deprecation_message?: string;
    description?: string;
    plan_modifiers?: SchemaDynamicPlanModifiers;
    sensitive?: boolean;
    validators?: SchemaDynamicValidators;
  };
}
export interface SchemaDynamicPlanModifier {
  custom: SchemaCustomPlanModifier;
  [k: string]: unknown;
}
export interface ResourceFloat64Attribute {
  name: ValidIdentifier;
  float64: {
    associated_external_type?: SchemaAssociatedExternalType;
    computed_optional_required: SchemaComputedOptionalRequired;
    custom_type?: SchemaCustomType;
    default?: SchemaFloat64Default;
    deprecation_message?: string;
    description?: string;
    plan_modifiers?: SchemaFloat64PlanModifiers;
    sensitive?: boolean;
    validators?: SchemaFloat64Validators;
  };
}
export interface SchemaFloat64PlanModifier {
  custom: SchemaCustomPlanModifier;
  [k: string]: unknown;
}
export interface ResourceInt64Attribute {
  name: ValidIdentifier;
  int64: {
    associated_external_type?: SchemaAssociatedExternalType;
    computed_optional_required: SchemaComputedOptionalRequired;
    custom_type?: SchemaCustomType;
    default?: SchemaInt64Default;
    deprecation_message?: string;
    description?: string;
    plan_modifiers?: SchemaInt64PlanModifiers;
    sensitive?: boolean;
    validators?: SchemaInt64Validators;
  };
}
export interface SchemaInt64PlanModifier {
  custom: SchemaCustomPlanModifier;
  [k: string]: unknown;
}
export interface ResourceListAttribute {
  name: ValidIdentifier;
  list: {
    associated_external_type?: SchemaAssociatedExternalType;
    computed_optional_required: SchemaComputedOptionalRequired;
    custom_type?: SchemaCustomType;
    default?: SchemaListDefault;
    deprecation_message?: string;
    description?: string;
    element_type: SchemaElementType;
    plan_modifiers?: SchemaListPlanModifiers;
    sensitive?: boolean;
    validators?: SchemaListValidators;
  };
}
export interface SchemaListPlanModifier {
  custom: SchemaCustomPlanModifier;
  [k: string]: unknown;
}
export interface ResourceListNestedAttribute {
  name: ValidIdentifier;
  list_nested: {
    computed_optional_required: SchemaComputedOptionalRequired;
    custom_type?: SchemaCustomType;
    default?: SchemaListDefault;
    deprecation_message?: string;
    description?: string;
    nested_object: ResourceNestedAttributeObject;
    plan_modifiers?: SchemaListPlanModifiers;
    sensitive?: boolean;
    validators?: SchemaListValidators;
  };
}
export interface ResourceNestedAttributeObject {
  associated_external_type?: SchemaAssociatedExternalType;
  attributes?: ResourceAttributes;
  custom_type?: SchemaCustomType;
  plan_modifiers?: SchemaObjectPlanModifiers;
  validators?: SchemaObjectValidators;
  [k: string]: unknown;
}
export interface SchemaObjectPlanModifier {
  custom: SchemaCustomPlanModifier;
  [k: string]: unknown;
}
export interface ResourceMapAttribute {
  name: ValidIdentifier;
  map: {
    associated_external_type?: SchemaAssociatedExternalType;
    computed_optional_required: SchemaComputedOptionalRequired;
    custom_type?: SchemaCustomType;
    default?: SchemaMapDefault;
    deprecation_message?: string;
    description?: string;
    element_type: SchemaElementType;
    plan_modifiers?: SchemaMapPlanModifiers;
    sensitive?: boolean;
    validators?: SchemaMapValidators;
  };
}
export interface SchemaMapPlanModifier {
  custom: SchemaCustomPlanModifier;
  [k: string]: unknown;
}
export interface ResourceMapNestedAttribute {
  name: ValidIdentifier;
  map_nested: {
    computed_optional_required: SchemaComputedOptionalRequired;
    custom_type?: SchemaCustomType;
    default?: SchemaMapDefault;
    deprecation_message?: string;
    description?: string;
    nested_object: ResourceNestedAttributeObject;
    plan_modifiers?: SchemaMapPlanModifiers;
    sensitive?: boolean;
    validators?: SchemaMapValidators;
  };
}
export interface ResourceNumberAttribute {
  name: ValidIdentifier;
  number: {
    associated_external_type?: SchemaAssociatedExternalType;
    computed_optional_required: SchemaComputedOptionalRequired;
    custom_type?: SchemaCustomType;
    default?: SchemaNumberDefault;
    deprecation_message?: string;
    description?: string;
    plan_modifiers?: SchemaNumberPlanModifiers;
    sensitive?: boolean;
    validators?: SchemaNumberValidators;
  };
}
export interface SchemaNumberPlanModifier {
  custom: SchemaCustomPlanModifier;
  [k: string]: unknown;
}
export interface ResourceObjectAttribute {
  name: ValidIdentifier;
  object: {
    associated_external_type?: SchemaAssociatedExternalType;
    attribute_types: SchemaObjectAttributeTypes;
    computed_optional_required: SchemaComputedOptionalRequired;
    custom_type?: SchemaCustomType;
    default?: SchemaObjectDefault;
    deprecation_message?: string;
    description?: string;
    plan_modifiers?: SchemaObjectPlanModifiers;
    sensitive?: boolean;
    validators?: SchemaObjectValidators;
  };
}
export interface ResourceSetAttribute {
  name: ValidIdentifier;
  set: {
    associated_external_type?: SchemaAssociatedExternalType;
    computed_optional_required: SchemaComputedOptionalRequired;
    custom_type?: SchemaCustomType;
    default?: SchemaSetDefault;
    deprecation_message?: string;
    description?: string;
    element_type: SchemaElementType;
    plan_modifiers?: SchemaSetPlanModifiers;
    sensitive?: boolean;
    validators?: SchemaSetValidators;
  };
}
export interface SchemaSetPlanModifier {
  custom: SchemaCustomPlanModifier;
  [k: string]: unknown;
}
export interface ResourceSetNestedAttribute {
  name: ValidIdentifier;
  set_nested: {
    computed_optional_required: SchemaComputedOptionalRequired;
    custom_type?: SchemaCustomType;
    default?: SchemaSetDefault;
    deprecation_message?: string;
    description?: string;
    nested_object: ResourceNestedAttributeObject;
    plan_modifiers?: SchemaSetPlanModifiers;
    sensitive?: boolean;
    validators?: SchemaSetValidators;
  };
}
export interface ResourceSingleNestedAttribute {
  name: ValidIdentifier;
  single_nested: {
    associated_external_type?: SchemaAssociatedExternalType;
    attributes?: ResourceAttributes;
    computed_optional_required: SchemaComputedOptionalRequired;
    custom_type?: SchemaCustomType;
    default?: SchemaObjectDefault;
    deprecation_message?: string;
    description?: string;
    plan_modifiers?: SchemaObjectPlanModifiers;
    sensitive?: boolean;
    validators?: SchemaObjectValidators;
  };
}
export interface ResourceStringAttribute {
  name: ValidIdentifier;
  string: {
    associated_external_type?: SchemaAssociatedExternalType;
    computed_optional_required: SchemaComputedOptionalRequired;
    custom_type?: SchemaCustomType;
    default?: SchemaStringDefault;
    deprecation_message?: string;
    description?: string;
    plan_modifiers?: SchemaStringPlanModifiers;
    sensitive?: boolean;
    validators?: SchemaStringValidators;
  };
}
export interface SchemaStringPlanModifier {
  custom: SchemaCustomPlanModifier;
  [k: string]: unknown;
}
export interface ResourceListNestedBlock {
  name: ValidIdentifier;
  list_nested: {
    custom_type?: SchemaCustomType;
    default?: SchemaListDefault;
    deprecation_message?: string;
    description?: string;
    nested_object: ResourceNestedBlockObject;
    plan_modifiers?: SchemaListPlanModifiers;
    validators?: SchemaListValidators;
  };
}
export interface ResourceNestedBlockObject {
  associated_external_type?: SchemaAssociatedExternalType;
  attributes?: ResourceAttributes;
  blocks?: ResourceBlocks;
  custom_type?: SchemaCustomType;
  plan_modifiers?: SchemaObjectPlanModifiers;
  validators?: SchemaObjectValidators;
  [k: string]: unknown;
}
export interface ResourceSetNestedBlock {
  name: ValidIdentifier;
  set_nested: {
    custom_type?: SchemaCustomType;
    default?: SchemaSetDefault;
    deprecation_message?: string;
    description?: string;
    nested_object: ResourceNestedBlockObject;
    plan_modifiers?: SchemaSetPlanModifiers;
    validators?: SchemaSetValidators;
  };
}
export interface ResourceSingleNestedBlock {
  name: ValidIdentifier;
  single_nested: {
    associated_external_type?: SchemaAssociatedExternalType;
    attributes?: ResourceAttributes;
    blocks?: ResourceBlocks;
    custom_type?: SchemaCustomType;
    default?: SchemaObjectDefault;
    deprecation_message?: string;
    description?: string;
    plan_modifiers?: SchemaObjectPlanModifiers;
    validators?: SchemaObjectValidators;
  };
}
