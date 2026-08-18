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
    "x-tf-computed-optional-required"?: ComputedOptionalRequired;
    // Inidcates this is the top level item type for the data source. Only used for plural data sources.
    "x-tf-top-level-item-type"?: boolean;
  }
}
