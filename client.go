package cloudant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/mreiferson/go-httpclient"
)

type Client struct {
	httpClient     *http.Client
	rootUri        string
	apiKey         string
	apiPassword    string
	PrintResponses bool
}

// An implementation of 'error' that exposes all the cloudant specific error details.
type CloudantError struct {
	// The status string returned from the HTTP call.
	Status string `json:"-"`

	// The status, as an integer, returned from the HTTP call.
	StatusCode int `json:"-"`

	// NOTE: The following properties are per the JSON API Error Spec
	//       http://jsonapi.org/format/#errors

	// The Cloudant error id.
	Code string `json:"error"`

	// The Cloudant error reason.
	Detail string `json:"reason"`
}

func (e *CloudantError) Error() string {
	return fmt.Sprintf("%s (%d): %s %s", e.Status, e.StatusCode, e.Code, e.Detail)
}

var (
	DefaultTransport *httpclient.Transport = &httpclient.Transport{
		ConnectTimeout:        1 * time.Second,
		RequestTimeout:        5 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
	}
)

func NewClient(rootUri string, apiKey string, apiPassword string) *Client {
	return NewClientWithTransport(rootUri, apiKey, apiPassword, DefaultTransport)
}

// Like NewClient, except that it allows a specific http.Transport to be
// provided for use, rather than DefaultTransport.
func NewClientWithTransport(rootUri string, apiKey string, apiPassword string, transport *httpclient.Transport) *Client {
	return &Client{
		httpClient:     &http.Client{Transport: transport},
		rootUri:        rootUri,
		apiKey:         apiKey,
		apiPassword:    apiPassword,
		PrintResponses: false,
	}
}

// Decodes Cloudant responses.
func (c *Client) decode(r io.Reader, receiver interface{}) error {
	var decoder *json.Decoder

	if !c.PrintResponses {
		decoder = json.NewDecoder(r)
	} else {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r)
		fmt.Printf(buf.String())
		decoder = json.NewDecoder(buf)
	}

	return decoder.Decode(receiver)
}

// Executes an HTTP request.
func (c *Client) doRequest(method, trailing string, headers map[string]string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, c.rootUri+trailing, body)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.apiKey, c.apiPassword)

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	req.Header.Add("Accept", "application/json")
	if method == "POST" || method == "PUT" || method == "PATCH" {
		req.Header.Add("Content-Type", "application/json")
	}

	return c.httpClient.Do(req)
}

func (client *Client) handleResponse(resp *http.Response, err error, successStatusCode int, result interface{}) error {
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != successStatusCode {
		return newCloudantError(resp)
	}

	client.decode(resp.Body, result)

	return err
}

func newCloudantError(resp *http.Response) error {
	ce := &CloudantError{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
	}

	if data, err := ioutil.ReadAll(resp.Body); err != nil {
		ce.Detail = fmt.Sprintf("Can not read HTTP response: %s", err)
		return ce
	} else if err := json.Unmarshal(data, ce); err != nil {
		ce.Detail = fmt.Sprintf("Can not unmarshal JSON response '''%s''': %s", string(data), err)
		return ce
	}

	return ce
}

func NotFoundError() *CloudantError {
	return &CloudantError{
		Status:     "(404) Not Found",
		StatusCode: 404,
		Code:       "not_found",
		Detail:     "missing",
	}
}
