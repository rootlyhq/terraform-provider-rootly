package diffsuppressfunc

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var _ schema.SchemaDiffSuppressFunc = Time

func Time(k, oldValue, newValue string, d *schema.ResourceData) bool {
	return normalizeTimeString(oldValue) == normalizeTimeString(newValue)
}

// Normalize various time formats to "15:04" (24-hour)
func normalizeTimeString(t string) string {
	t = strings.TrimSpace(strings.ToUpper(t))

	// Attempt 24-hour format first ("08:00")
	if parsed, err := time.Parse("15:04", t); err == nil {
		return parsed.Format("15:04")
	}

	// Attempt common 12-hour formats
	formats := []string{
		"03:04PM",
		"03:04 PM",
		"3:04PM",
		"3:04 PM",
		"03:04AM",
		"03:04 AM",
		"3:04AM",
		"3:04 AM",
	}
	for _, format := range formats {
		if parsed, err := time.Parse(format, t); err == nil {
			return parsed.Format("15:04")
		}
	}

	// Fallback: return original string if no parsing worked
	return t
}
