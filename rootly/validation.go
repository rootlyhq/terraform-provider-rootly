package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"regexp"
)

func validCSSHexColor() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`),
		"must be a valid action (usually starts with lambda:)",
	)
}
