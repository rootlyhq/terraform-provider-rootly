package client

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/jsonapi"
)

func marshalAttrs(t *testing.T, c *StatusPageComponent) map[string]interface{} {
	t.Helper()
	buf := new(bytes.Buffer)
	if err := jsonapi.MarshalPayload(buf, c); err != nil {
		t.Fatal(err)
	}
	var payload struct {
		Data struct {
			Attributes map[string]interface{} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	return payload.Data.Attributes
}

func TestStatusPageComponentGroupIdMarshal(t *testing.T) {
	// Unchanged update: attribute must be omitted entirely.
	attrs := marshalAttrs(t, &StatusPageComponent{ID: "c1", Name: "renamed"})
	if _, present := attrs["status_page_component_group_id"]; present {
		t.Fatalf("expected group id omitted, got %v", attrs)
	}

	// Clearing: explicit blank must be serialized.
	empty := ""
	attrs = marshalAttrs(t, &StatusPageComponent{ID: "c1", StatusPageComponentGroupId: &empty})
	if v, present := attrs["status_page_component_group_id"]; !present || v != "" {
		t.Fatalf("expected blank group id serialized, got %v (present=%v)", v, present)
	}

	// Setting: value must be serialized.
	gid := "grp-123"
	attrs = marshalAttrs(t, &StatusPageComponent{ID: "c1", StatusPageComponentGroupId: &gid})
	if v := attrs["status_page_component_group_id"]; v != "grp-123" {
		t.Fatalf("expected group id grp-123, got %v", v)
	}
}

func TestStatusPageComponentGroupIdUnmarshalNull(t *testing.T) {
	body := `{"data":{"id":"c1","type":"status_page_components","attributes":{"name":"n","status_page_component_group_id":null}}}`
	c := new(StatusPageComponent)
	if err := jsonapi.UnmarshalPayload(strings.NewReader(body), c); err != nil {
		t.Fatal(err)
	}
	if c.StatusPageComponentGroupId != nil {
		t.Fatalf("expected nil group id, got %q", *c.StatusPageComponentGroupId)
	}
}
