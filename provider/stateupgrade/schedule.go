package stateupgrade

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ScheduleV0 is the pre-v2.x rootly_schedule schema used to decode state
// containing the removed slack_user_group field. Frozen — do not update.
func ScheduleV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"slack_user_group": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

// UpgradeScheduleV0ToV1 drops the legacy slack_user_group field.
func UpgradeScheduleV0ToV1(_ context.Context, rawState map[string]any, _ any) (map[string]any, error) {
	delete(rawState, "slack_user_group")
	return rawState, nil
}
