package client

import (
	"fmt"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	rootlygo "github.com/rootlyhq/rootly-go"
	"net/http"
)

type Client struct {
	Token       string
	ContentType string
	UserAgent   string
	Rootly      rootlygo.Client
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("Content-Type", c.ContentType)
	req.Header.Set("User-Agent", c.UserAgent)

	return c.Rootly.Client.Do(req)
}

// NewClient returns a new rootly.Client which can be used to access the API methods.
func NewClient(endpoint, token, userAgent string) (*Client, error) {
	httpClient := cleanhttp.DefaultClient()
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
