package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/pkg/errors"
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
			return nil, errors.Wrap(err, "request tester failed")
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

func expectBody(expectedBody string) requestTesterFunc {
	return func(r *http.Request) error {
		var body = ""
		if r.Body != nil {
			bodyBytes, err := ioutil.ReadAll(r.Body)
			defer func() {
				_ = r.Body.Close()
			}()
			if err != nil {
				return errors.Wrap(err, "error reading request body")
			}
			body = string(bodyBytes)
		}
		if body != expectedBody {
			return errors.New("test request body != expected request body")
		}
		return nil
	}
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
	client := http.Client{
		Transport: rt,
	}
	jsonClient := Client{
		Endpoint: "",
		Client:   client,
	}
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Errorf("err building request: %s", err.Error())
	}
	_, err = jsonClient.doRequest(req)
	if err != nil {
		t.Errorf("err returned from Get: %s", err.Error())
	}
}

func Test4XXStatusCodeIsError(t *testing.T) {
	rt := testRoundTripper{
		statusCode: 404,
	}
	jsonClient := Client{
		Endpoint: "",
		Client: http.Client{
			Transport: rt,
		},
	}
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Errorf("err building request: %s", err.Error())
	}
	_, err = jsonClient.doRequest(req)
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
	client := http.Client{
		Transport: rt,
	}
	jsonClient := Client{
		Endpoint: "",
		Client:   client,
	}
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Errorf("err building request: %s", err.Error())
	}
	_, err = jsonClient.doRequest(req)
	if reqErr, ok := err.(RequestError); !ok {
		t.Errorf("err was not of type client.RequestError, it was of type %T", err)
	} else if reqErr.StatusCode != 502 {
		t.Errorf("err.StatusCode was not 502, it was %d", reqErr.StatusCode)
	}
}

func TestAddsDefaultHeaders(t *testing.T) {
	rt := testRoundTripper{
		requestTester: expectHeader("access-token", "abc"),
	}
	client := http.Client{
		Transport: rt,
	}
	jsonClient := Client{
		Endpoint: "",
		DefaultHeaders: map[string]string{
			"access-token": "abc",
		},
		Client: client,
	}
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Errorf("err building request: %s", err.Error())
	}
	_, err = jsonClient.doRequest(req)
	if err != nil {
		t.Errorf("error in test get: %s", err.Error())
	}
}

func TestAddsAuthorizationHeader(t *testing.T) {
	client, err := NewClient("endpoint", "test", "userAgent")

	if err != nil {
		t.Errorf("error creating client: %s", err.Error())
	}

	if _, ok := client.DefaultHeaders["Authorization"]; !ok {
		t.Errorf("Authorization header not set")
	}
}
