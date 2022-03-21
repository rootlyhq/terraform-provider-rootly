package provider

import (
	"testing"
)

func TestValidPermissionAction(t *testing.T) {
	validValues := []string{
		"#112233",
		"#123",
		"#000233",
		"#023",
	}
	for _, v := range validValues {
		_, errors := validCSSHexColor()(v, "action")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid CSS Hex color: %q", v, errors)
		}
	}

	invalidNames := []string{
		"invalid",
		"#abcd",
		"#-12",
	}
	for _, v := range invalidNames {
		_, errors := validCSSHexColor()(v, "action")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid CSS Hex color", v)
		}
	}
}
