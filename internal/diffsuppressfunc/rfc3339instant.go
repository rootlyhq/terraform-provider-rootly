package diffsuppressfunc

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// RFC3339Instant suppresses diff for time values
func RFC3339Instant(k, oldValue, newValue string, d *schema.ResourceData) bool {
	newRFC3339time, errNew := time.Parse(time.RFC3339, newValue)
	oldRFC3339time, errOld := time.Parse(time.RFC3339, oldValue)
	if errNew != nil || errOld != nil {
		return false
	}
	return newRFC3339time.Equal(oldRFC3339time)
}
