package polling

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const defaultListTimeout = 60 * time.Second

// ListFunc is a function that performs a list/search API call and returns
// the results. Data source Read functions should wrap their client.List*
// call in a ListFunc.
type ListFunc func() ([]interface{}, error)

// WaitForList retries listFn until it returns a non-empty result or the
// timeout elapses. This handles the eventual consistency window introduced
// by asynchronous Searchkick/Elasticsearch indexing — a resource that was
// just created may not be immediately searchable via filtered list queries.
func WaitForList(ctx context.Context, resourceName string, listFn ListFunc) ([]interface{}, error) {
	var items []interface{}
	err := retry.RetryContext(ctx, defaultListTimeout, func() *retry.RetryError {
		var listErr error
		items, listErr = listFn()
		if listErr != nil {
			return retry.NonRetryableError(listErr)
		}
		if len(items) == 0 {
			return retry.RetryableError(fmt.Errorf("%s not found, retrying...", resourceName))
		}
		return nil
	})
	return items, err
}
