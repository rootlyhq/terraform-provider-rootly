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
  resourceName: string;
}

interface DataSourceSingleConfig {
  type: "single";
}

export type DataSourceConfig = {
  name: string;
  description?: string;
} & (DataSourceListConfig | DataSourceSingleConfig);

export type ResolvedDataSourceConfig = {
  kind: "data_source";
  name: string;
  description?: string;
  goNames: {
    /** Name of the struct that represents the data source. */
    struct: `${string}DataSource`;
    /** Name of the struct that represents the model of the data source. */
    model: `${string}DataSourceModel`;
  };
} & (DataSourceListConfig | DataSourceSingleConfig);

// TODO
export type ResourceConfig = {
  name: string;
  description?: string;
};

export type ResolvedResourceConfig = {
  kind: "resource";
  name: string;
  description?: string;
  goNames: {
    /** Name of the struct that represents the resource. */
    struct: `${string}Resource`;
    /** Name of the struct that represents the model of the resource. */
    model: `${string}ResourceModel`;
  };
};

export type ComputedOptionalRequired =
  | "computed"
  | "optional"
  | "computed_optional"
  | "required";

declare module "openapi3-ts/oas30" {
  interface SchemaObject {
    /** Defines the jsonapi attribute tag. */
    "x-go-jsonapi-tag"?: "primary" | "attr";
    /** Defines the jsonapi type. */
    "x-go-jsonapi-type"?: string;
    /** Overrides the nested object type name. Only used for nested objects. */
    "x-go-nested-type"?: string;
    /** Defines if the property is computed, optional or required. */
    "x-tf-computed-optional-required"?: ComputedOptionalRequired;
    /** Indicates this is the top level item type for the data source. Only used for plural data sources. */
    "x-tf-top-level-item-type"?: boolean;
  }
}
