package serlogs

import (
	"fmt"
)

// LogClient defines the interface for sending log entries and performing a ping check.
type LogClient interface {
	// Send sends a log entry to the log server.
	send(LogEntry) error
	// Ping checks the health of the log server.
	Ping() error
}

// client implements the LogClient interface.
type client struct {
	// url is the URL of the log server.
	url string
	// apiKey is the API key for authenticating with the log server.
	apiKey string
	// appName is the name of the application sending logs.
	appName string
	// httpClient is the HTTP client used to send requests.
	httpClient httpClient
}

// NewSerLogsClient creates a new client for sending log entries.
func NewSerLogsClient(url, apiKey, appName string, httpClient httpClient) LogClient {
	return &client{
		url:        url,
		apiKey:     apiKey,
		appName:    appName,
		httpClient: httpClient,
	}
}

// send sends a log entry to the log server.
func (c *client) send(logEntry LogEntry) error {
	url := c.url + "/logs"
	headers := map[string]string{
		"Content-Type": "application/json",
		"X-API-KEY":    c.apiKey,  // todo 待确定
		"X-APP-NAME":   c.appName, // todo 待确定
	}

	err := c.httpClient.Post(url, headers, logEntry, nil)
	if err != nil {
		return fmt.Errorf("failed to send log entry: %v", err)
	}

	return nil
}

// Ping checks the health of the log server.
func (c *client) Ping() error {
	url := c.url + "/ping"
	headers := map[string]string{
		"X-API-KEY":  c.apiKey,  // todo 待确定
		"X-APP-NAME": c.appName, // todo 待确定
	}

	var result interface{}
	err := c.httpClient.Get(url, headers, &result)
	if err != nil {
		return fmt.Errorf("failed to ping log server: %v", err)
	}
	return nil
}
