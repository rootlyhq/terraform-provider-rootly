package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func TestIANATimeZones_ValidValues(t *testing.T) {
	validValues := []string{
		"America/New_York",
		"America/Los_Angeles",
		"America/Chicago",
		"Europe/London",
		"Asia/Tokyo",
		"Etc/UTC",
		"Pacific/Auckland",
		"Australia/Sydney",
	}

	validateFunc := validation.StringInSlice(IANATimeZones, false)
	for _, v := range validValues {
		_, errors := validateFunc(v, "time_zone")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid IANA time zone: %q", v, errors)
		}
	}
}

func TestIANATimeZones_RejectsRailsNames(t *testing.T) {
	railsOnlyNames := []string{
		"Eastern Time (US & Canada)",
		"Pacific Time (US & Canada)",
		"Central Time (US & Canada)",
		"Mountain Time (US & Canada)",
		"Hawaii",
		"Alaska",
		"London",
		"Tokyo",
		"Sydney",
	}

	validateFunc := validation.StringInSlice(IANATimeZones, false)
	for _, v := range railsOnlyNames {
		_, errors := validateFunc(v, "time_zone")
		if len(errors) == 0 {
			t.Fatalf("%q should NOT be accepted as an IANA time zone", v)
		}
	}
}

func TestRailsTimeZones_AcceptsBothFormats(t *testing.T) {
	validValues := []string{
		"Eastern Time (US & Canada)",
		"America/New_York",
		"Pacific Time (US & Canada)",
		"America/Los_Angeles",
		"UTC",
		"Etc/UTC",
		"London",
		"Europe/London",
	}

	validateFunc := validation.StringInSlice(RailsTimeZones, false)
	for _, v := range validValues {
		_, errors := validateFunc(v, "time_zone")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Rails time zone: %q", v, errors)
		}
	}
}

func TestRailsTimeZones_RejectsInvalidValues(t *testing.T) {
	invalidValues := []string{
		"Not/A/Timezone",
		"EST",
		"PST",
		"GMT-5",
		"",
	}

	validateFunc := validation.StringInSlice(RailsTimeZones, false)
	for _, v := range invalidValues {
		_, errors := validateFunc(v, "time_zone")
		if len(errors) == 0 {
			t.Fatalf("%q should NOT be accepted as a valid time zone", v)
		}
	}
}

func TestIANATimeZones_NotEmpty(t *testing.T) {
	if len(IANATimeZones) == 0 {
		t.Fatal("IANATimeZones should not be empty")
	}
}

func TestRailsTimeZones_NotEmpty(t *testing.T) {
	if len(RailsTimeZones) == 0 {
		t.Fatal("RailsTimeZones should not be empty")
	}
}
