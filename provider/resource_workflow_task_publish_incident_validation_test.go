package provider

import (
	"strings"
	"testing"
)

func TestValidatePublishIncidentComponentSelection(t *testing.T) {
	t.Run("accepts statuses for selected keys", func(t *testing.T) {
		params := map[string]interface{}{
			"selected_component_keys": []interface{}{"Service:1"},
			"selected_component_statuses": map[string]interface{}{
				"Service:1": "major_outage",
			},
		}

		if err := validatePublishIncidentComponentSelection(params); err != nil {
			t.Fatalf("expected valid component selection: %v", err)
		}
	})

	t.Run("rejects statuses for unselected keys", func(t *testing.T) {
		params := map[string]interface{}{
			"selected_component_keys": []interface{}{"Service:1"},
			"selected_component_statuses": map[string]interface{}{
				"Service:2": "major_outage",
			},
		}

		err := validatePublishIncidentComponentSelection(params)
		if err == nil || !strings.Contains(err.Error(), "Service:2") {
			t.Fatalf("expected an unknown-key validation error, got: %v", err)
		}
	})
}
