package serlogs

import (
	"fmt"
)

// LogClient defines the interface for sending log entries and performing a ping check.
type LogClient interface {
	Send(LogEntry) error
	Ping() error
}

// Client implements the LogClient interface.
type Client struct {
	TargetURL  string
	APIKey     string
	AppName    string
	HTTPClient HTTPClientI
}

// NewLogClient creates a new Client for sending log entries.
func NewLogClient(targetURL, apiKey, appName string, httpClient HTTPClientI) *Client {
	return &Client{
		TargetURL:  targetURL,
		APIKey:     apiKey,
		AppName:    appName,
		HTTPClient: httpClient,
	}
}

// Send sends a log entry to the log server.
func (c *Client) Send(logEntry LogEntry) error {
	url := c.TargetURL
	headers := map[string]string{
		"Content-Type": "application/json",
		"X-API-KEY":    c.APIKey,
		"X-APP-NAME":   c.AppName,
	}

	err := c.HTTPClient.Post(url, headers, logEntry, nil)
	if err != nil {
		return fmt.Errorf("failed to send log entry: %v", err)
	}

	return nil
}

// Ping checks the health of the log server.
func (c *Client) Ping() error {
	url := c.TargetURL + "/ping"
	headers := map[string]string{
		"X-API-KEY":  c.APIKey,
		"X-APP-NAME": c.AppName,
	}

	var result interface{}
	err := c.HTTPClient.Get(url, headers, &result)
	if err != nil {
		return fmt.Errorf("failed to ping log server: %v", err)
	}
	return nil
}
