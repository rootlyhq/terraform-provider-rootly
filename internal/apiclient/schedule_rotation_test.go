package apiclient

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/rootlyhq/jsonapi"
)

func TestScheduleRotationScheduleRotationableAttributesRoundTrip(t *testing.T) {
	want := map[string]any{
		"handoff_time": "09:00",
		"end_of_day":   "17:00",
		"shift_length": "7",
	}
	input := ScheduleRotation{
		ID:                             "rotation-id",
		ScheduleRotationableAttributes: want,
	}

	var payload bytes.Buffer
	if err := jsonapi.MarshalPayload(&payload, &input); err != nil {
		t.Fatalf("marshalling schedule rotation: %v", err)
	}

	var got ScheduleRotation
	if err := jsonapi.UnmarshalPayload(&payload, &got); err != nil {
		t.Fatalf("unmarshalling schedule rotation: %v", err)
	}
	if !reflect.DeepEqual(got.ScheduleRotationableAttributes, want) {
		t.Errorf("round-tripped attributes = %#v, want %#v", got.ScheduleRotationableAttributes, want)
	}
}
