package mapvalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.Map = KeysRequiredValidator{}

type KeysRequiredValidator struct {
	Keys []string
}

// KeysRequired returns a map validator that checks that all of the given keys
// are present in the map.
func KeysRequired(keys ...string) validator.Map {
	return KeysRequiredValidator{
		Keys: keys,
	}
}

func (v KeysRequiredValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("Ensures that the map contains the keys: %s.", strings.Join(v.Keys, ", "))
}

func (v KeysRequiredValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Ensures that the map contains the keys: `%s`.", strings.Join(v.Keys, "`, `"))
}

func (v KeysRequiredValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	elements := req.ConfigValue.Elements()
	for _, key := range v.Keys {
		if _, ok := elements[key]; !ok {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Missing Required Map Key",
				fmt.Sprintf("Key %q is required.", key),
			)
		}
	}
}
