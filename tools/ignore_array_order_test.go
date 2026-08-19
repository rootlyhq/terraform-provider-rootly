package tools

import (
	"testing"
)

func TestEqualIgnoringOrderResolvesNestedListPath(t *testing.T) {
	if got := listPathFromDiffKey("task_params.0.selected_component_keys.0"); got != "task_params.0.selected_component_keys" {
		t.Fatalf("unexpected nested list path: %q", got)
	}
	if !listsAreEqual(
		[]interface{}{"Service:1", "Service:2"},
		[]interface{}{"Service:2", "Service:1"},
	) {
		t.Fatal("expected an order-only nested list change to be suppressed")
	}
}
