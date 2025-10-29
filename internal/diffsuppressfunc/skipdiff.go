package diffsuppressfunc

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var _ schema.SchemaDiffSuppressFunc = Skip

func Skip(k, oldValue, newValue string, d *schema.ResourceData) bool {
	return len(oldValue) != 0
}
