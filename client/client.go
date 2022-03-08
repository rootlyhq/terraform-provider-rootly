package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
)

type Client struct {
	Endpoint       string
	DefaultHeaders map[string]string
	Client         http.Client
	Timeout        int
}

type Data struct {
	Type       string
	Attributes map[string]interface{}
}

// NewClient returns a new rootly.Client which can be used to access the API methods.
func NewClient(endpoint, token, userAgent string) (*Client, error) {
	httpClient := cleanhttp.DefaultClient()
	httpClient.Transport = logging.NewTransport("Rootly", httpClient.Transport)
	httpClient.Timeout = 60 * time.Second

	client := &Client{
		Endpoint: endpoint,
		DefaultHeaders: map[string]string{
			"Content-Type": "application/vnd.api+json",
			"User-Agent":   userAgent,
		},
		Client: *httpClient,
	}

	if token != "" {
		client.DefaultHeaders["Authorization"] = fmt.Sprintf("Bearer %s", token)
	}

	return client, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	// Handle headers
	for k, v := range c.DefaultHeaders {
		req.Header.Set(k, v)
	}

	res, err := c.Client.Do(req)
	defer c.Client.CloseIdleConnections()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		if res.StatusCode == http.StatusNotFound {
			return nil, NewNotFoundError(string(body))
		}
		return nil, NewRequestError(res.StatusCode, string(body))
	}

	return body, err
}

func (c *Client) makeUrl(path string) string {
	return fmt.Sprintf("%s/%s", c.Endpoint, path)
}
