package client

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
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
	req.Header.Set("User-Agent", "terraform-provider-rootly/v2.15.0")

	res, err := c.Rootly.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		body, err := ioutil.ReadAll(res.Body)
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

// NewClient returns a new rootly.Client which can be used to access the API methods.
func NewClient(endpoint, token, userAgent string) (*Client, error) {
	retryableClient := retryablehttp.NewClient()
	retryableClient.RetryMax = 5
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
