export interface ClientConfig {
  name: string;
  actions?: {
    list?: {
      enabled: true;
    };
    get?: {
      enabled: true;
    };
  };
}

interface DataSourceListConfig {
  type: "list";
}

interface DataSourceSingleConfig {
  type: "single";
}

export type DataSourceConfig = {
  name: string;
  resource: string;
  description?: string;
} & (DataSourceListConfig | DataSourceSingleConfig);

// TODO
export type ResourceConfig = {
  name: string;
  description?: string;
};

export type ComputedOptionalRequired =
  | "computed"
  | "optional"
  | "computed_optional"
  | "required";

declare module "openapi3-ts/oas30" {
  interface SchemaObject {
    // Defines the jsonapi attribute tag
    "x-tf-jsonapi-tag"?: "primary" | "attr";
    // Defines the jsonapi type
    "x-tf-jsonapi-type"?: string;
    // Defines if the property is computed, optional or required.
    "x-tf-computed-optional-required"?: ComputedOptionalRequired;
    // Indicates this is the top level item type for the data source. Only used for plural data sources.
    "x-tf-top-level-item-type"?: boolean;
  }
}
