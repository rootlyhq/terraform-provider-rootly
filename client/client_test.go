package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/go-cleanhttp"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type requestTesterFunc func(r *http.Request) error

type testRoundTripper struct {
	requestTester requestTesterFunc
	response      string
	statusCode    int
}

func (rt testRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	if rt.requestTester != nil {
		err := rt.requestTester(r)
		if err != nil {
			return nil, fmt.Errorf("request tester failed %w", err)
		}
	}

	responseBody := ioutil.NopCloser(strings.NewReader(rt.response))
	statusCode := 200
	if rt.statusCode != 0 {
		statusCode = rt.statusCode
	}
	return &http.Response{
		Body:       responseBody,
		StatusCode: statusCode,
	}, nil
}

func buildClient(rt testRoundTripper) Client {
	client := cleanhttp.DefaultClient()
	client.Transport = rt
	rootlyClient, _ := rootlygo.NewClient(
		"",
		rootlygo.WithHTTPClient(client),
	)
	jsonClient := Client{
		Token:       "",
		UserAgent:   "",
		ContentType: "",
		Rootly:      *rootlyClient,
	}

	return jsonClient
}

func expectHeader(expectedHeader string, expectedValue string) requestTesterFunc {
	return func(r *http.Request) error {
		actualValue := r.Header.Get(expectedHeader)
		if actualValue != expectedValue {
			return fmt.Errorf(
				"expected header %s to have value %s but it had value %s",
				expectedHeader,
				expectedValue,
				actualValue,
			)
		}
		return nil
	}
}

func TestHandlesNilOutputParam(t *testing.T) {
	rt := testRoundTripper{
		response: `{"responseParam": 1}`,
	}
	jsonClient := buildClient(rt)
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Errorf("err building request: %s", err.Error())
	}
	_, err = jsonClient.Do(req)
	if err != nil {
		t.Errorf("err returned from Get: %s", err.Error())
	}
}

func Test4XXStatusCodeIsError(t *testing.T) {
	rt := testRoundTripper{
		statusCode: 404,
	}
	jsonClient := buildClient(rt)
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Errorf("err building request: %s", err.Error())
	}
	_, err = jsonClient.Do(req)
	if reqErr, ok := err.(NotFoundError); !ok {
		t.Errorf("err was not of type client.NotFoundError, it was of type %T\n%s", err, reqErr)
	} else if reqErr.StatusCode != 404 {
		t.Errorf("err.StatusCode was not 404, it was %d", reqErr.StatusCode)
	}
}

func Test5XXStatusCodeIsError(t *testing.T) {
	rt := testRoundTripper{
		statusCode: 502,
	}
	jsonClient := buildClient(rt)
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Errorf("err building request: %s", err.Error())
	}
	_, err = jsonClient.Do(req)
	if reqErr, ok := err.(RequestError); !ok {
		t.Errorf("err was not of type client.RequestError, it was of type %T", err)
	} else if reqErr.StatusCode != 502 {
		t.Errorf("err.StatusCode was not 502, it was %d", reqErr.StatusCode)
	}
}
