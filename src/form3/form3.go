package form3

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
)

const (
	defaultBaseURL = "http://localhost:8080/v1/"
	contentType    = "application/vnd.api+json"
)

// A Client manages communication with the form3 API.
type Client struct {
	client  *http.Client // HTTP client used to commincate with the API.
	BaseURL *url.URL     // Base URL for API requests.
}

// BaseURL returns the base URL for the API.
func BaseURL() string {
	envBaseURL := os.Getenv("BASE_URL")
	if len(envBaseURL) > 0 {
		return envBaseURL
	}
	return defaultBaseURL
}

// NewClient returns a new form3 API client.
func NewClient() (client *Client) {
	baseURL, _ := url.Parse(BaseURL())
	c := &Client{client: &http.Client{}, BaseURL: baseURL}
	return c
}

// NewRequest creates an API request.
// If body is present it is JSON encoded and set as the request body.
func (client *Client) NewRequest(method string, urlStr string, body interface{}) (*http.Request, error) {
	u, err := client.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", contentType)
	return req, nil
}

// Do sends the request and returns the API response. The ctx input parameter must not be nil.
// The response body is JSON decoded into the provided model value.
// If the ctx has been cancelled then ctx.Err() is returned.
func (client *Client) Do(ctx context.Context, req *http.Request, model interface{}) (*http.Response, error) {
	if ctx == nil {
		return nil, errors.New("context.Context must not be nil")
	}
	req = req.WithContext(ctx)
	resp, err := client.client.Do(req)
	if err != nil {
		// Return context error if it has been cancelled.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}
	defer resp.Body.Close()

	decErr := json.NewDecoder(resp.Body).Decode(model)
	if decErr != nil && decErr != io.EOF {
		err = decErr
	}
	return resp, err
}
