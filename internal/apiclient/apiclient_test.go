package apiclient

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rootlyhq/terraform-provider-rootly/v5/client"
)

// The plugin-framework resources remove themselves from state, or treat a
// delete as done, when errors.Is(err, client.NotFoundError{}) holds. This
// checks that a real 404 response reaches them as that error through the
// shared HTTP client.
func TestNotFoundResponseMatchesFrameworkResourceCheck(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"errors":[{"title":"Not found or unauthorized","status":"404"}]}`))
	}))
	t.Cleanup(server.Close)

	legacyClient, apiClient, err := New(server.URL, "token", "test")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx := context.Background()
	calls := map[string]func() error{
		"GetSchedule": func() error {
			_, err := apiClient.GetSchedule(ctx, "missing")
			return err
		},
		"GetScheduleRotation": func() error {
			_, err := apiClient.GetScheduleRotation(ctx, "missing")
			return err
		},
		"GetUser": func() error {
			_, err := apiClient.GetUser(ctx, "missing")
			return err
		},
		"DeleteSchedule": func() error {
			return legacyClient.DeleteSchedule("missing")
		},
		"DeleteScheduleRotation": func() error {
			return legacyClient.DeleteScheduleRotation("missing")
		},
	}

	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := call()
			if !errors.Is(err, client.NotFoundError{}) {
				t.Fatalf("errors.Is(err, client.NotFoundError{}) = false for %v", err)
			}
			if !errors.Is(err, client.NewNotFoundError("")) {
				t.Fatalf("errors.Is(err, client.NewNotFoundError(\"\")) = false for %v", err)
			}
		})
	}
}
