package serlogs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// httpClient defines the interface for HTTP client operations.
type httpClient interface {
	// Do sends an HTTP request and returns an HTTP response.
	Do(req *http.Request) (*http.Response, error)
	// Get sends a GET request to the specified URL.
	Get(url string, headers map[string]string, result interface{}) error
	// Post sends a POST request to the specified URL with the given body.
	Post(url string, headers map[string]string, body, result interface{}) error
	// Put sends a PUT request to the specified URL with the given body.
	Put(url string, headers map[string]string, body, result interface{}) error
	// Delete sends a DELETE request to the specified URL with the given body.
	Delete(url string, headers map[string]string, body, result interface{}) error
}

// httpclient implements the httpClient interface.
type httpclient struct {
	client *http.Client
}

// newHTTPClient creates a new httpclient with the specified timeout.
func newHTTPClient(timeout time.Duration) httpClient {
	return &httpclient{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Do sends an HTTP request and returns an HTTP response.
func (c *httpclient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

// doRequest is a helper method for sending HTTP requests.
func (c *httpclient) doRequest(method, url string, headers map[string]string, body, result interface{}) error {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to serialize request body: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &httpClientError{
			StatusCode: resp.StatusCode,
			Err:        fmt.Errorf("%s request failed, status code: %d", method, resp.StatusCode),
		}
	}

	if result != nil {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %v", err)
		}
		return json.Unmarshal(respBody, result)
	}

	return nil
}

// Get sends a GET request to the specified URL.
func (c *httpclient) Get(url string, headers map[string]string, result interface{}) error {
	return c.doRequest("GET", url, headers, nil, result)
}

// Post sends a POST request to the specified URL with the given body.
func (c *httpclient) Post(url string, headers map[string]string, body, result interface{}) error {
	return c.doRequest("POST", url, headers, body, result)
}

// Put sends a PUT request to the specified URL with the given body.
func (c *httpclient) Put(url string, headers map[string]string, body, result interface{}) error {
	return c.doRequest("PUT", url, headers, body, result)
}

// Delete sends a DELETE request to the specified URL with the given body.
func (c *httpclient) Delete(url string, headers map[string]string, body, result interface{}) error {
	return c.doRequest("DELETE", url, headers, body, result)
}

type httpClientError struct {
	StatusCode int
	Err        error
}

func (e *httpClientError) Error() string {
	return fmt.Sprintf("%s request failed, status code: %d, error: %v", e.Err, e.StatusCode, e.Err)
}

func (e *httpClientError) Unwrap() error {
	return e.Err
}
