package converter

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func isResource(v any) bool {
	_, ok := v.(*schema.Resource)
	return ok
}
