package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type Client struct {
	Token       string
	ContentType string
	UserAgent   string
	Rootly      rootlygo.Client
}

// Do Intercepts the Request and enriches it with the required information.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("Content-Type", c.ContentType)
	req.Header.Set("User-Agent", c.UserAgent)

	res, err := c.Rootly.Client.Do(req)
	if err != nil {
		return nil, err
	}

	c.logRateLimitHeaders(res)

	if res.StatusCode >= 400 {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		res.Body.Close()

		if res.StatusCode == http.StatusNotFound {
			return nil, NewNotFoundError(string(body))
		}
		return nil, NewRequestError(res.StatusCode, string(body))
	}

	return res, nil
}

// logRateLimitHeaders logs rate limit information from response headers
func (c *Client) logRateLimitHeaders(resp *http.Response) {
	limit := resp.Header.Get("X-RateLimit-Limit")
	remaining := resp.Header.Get("X-RateLimit-Remaining")
	used := resp.Header.Get("X-RateLimit-Used")
	reset := resp.Header.Get("X-RateLimit-Reset")

	if limit != "" || remaining != "" || used != "" || reset != "" {
		ctx := context.Background()
		tflog.Debug(ctx, "Rate limit info", map[string]interface{}{
			"limit":     limit,
			"remaining": remaining,
			"used":      used,
			"reset":     reset,
		})

		if remaining != "" {
			if remainingInt, err := strconv.Atoi(remaining); err == nil && remainingInt < 10 {
				tflog.Warn(ctx, "Rate limit approaching", map[string]interface{}{
					"remaining": remainingInt,
				})
			}
		}
	}
}

// rateLimitRetryPolicy implements a retry policy that respects rate limiting
func rateLimitRetryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if err != nil {
		return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
	}

	if resp.StatusCode == 429 {
		return true, nil
	}

	return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
}

// rateLimitBackoff implements exponential backoff with rate limit awareness
func rateLimitBackoff(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
	if resp != nil && resp.StatusCode == 429 {
		resetHeader := resp.Header.Get("X-RateLimit-Reset")
		if resetHeader != "" {
			if resetTime, err := strconv.ParseInt(resetHeader, 10, 64); err == nil {
				resetDuration := time.Until(time.Unix(resetTime, 0))
				if resetDuration > 0 && resetDuration < max {
					return resetDuration + time.Second
				}
			}
		}

		remainingHeader := resp.Header.Get("X-RateLimit-Remaining")
		if remainingHeader == "0" {
			return time.Duration(attemptNum) * 30 * time.Second
		}
	}

	return retryablehttp.DefaultBackoff(min, max, attemptNum, resp)
}

// NewClient returns a new rootly.Client which can be used to access the API methods.
func NewClient(endpoint, token, userAgent string) (*Client, error) {
	retryableClient := retryablehttp.NewClient()
	retryableClient.RetryMax = 8
	retryableClient.RetryWaitMin = 1 * time.Second
	retryableClient.RetryWaitMax = 5 * time.Minute
	retryableClient.CheckRetry = rateLimitRetryPolicy
	retryableClient.Backoff = rateLimitBackoff

	retryableClient.RequestLogHook = func(logger retryablehttp.Logger, req *http.Request, attempt int) {
		if attempt > 0 {
			logger.Printf("[WARN] Retry attempt %d for %s %s", attempt, req.Method, req.URL)
		}
	}

	retryableClient.ResponseLogHook = func(logger retryablehttp.Logger, resp *http.Response) {
		if resp.StatusCode == 429 {
			remaining := resp.Header.Get("X-RateLimit-Remaining")
			reset := resp.Header.Get("X-RateLimit-Reset")
			logger.Printf("[WARN] Rate limit hit - Remaining: %s, Reset: %s", remaining, reset)
		}
	}

	httpClient := retryableClient.StandardClient()
	httpClient.Transport = logging.NewTransport("Rootly", httpClient.Transport)

	rootlyClient, err := rootlygo.NewClient(
		endpoint,
		rootlygo.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, err
	}

	client := &Client{
		Token:       token,
		ContentType: "application/vnd.api+json",
		UserAgent:   userAgent,
		Rootly:      *rootlyClient,
	}

	return client, nil
}
