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
	// targetURL is the URL of the log server.
	targetURL string
	// apiKey is the API key for authenticating with the log server.
	apiKey string
	// appName is the name of the application sending logs.
	appName string
	// httpClient is the HTTP client used to send requests.
	httpClient httpClient
}

// NewSerLogsClient creates a new client for sending log entries.
func NewSerLogsClient(config Config) LogClient {
	return &client{
		targetURL:  config.TargetURL,
		apiKey:     config.APIKey,
		appName:    config.AppName,
		httpClient: config.HTTPClient,
	}
}

// send sends a log entry to the log server.
func (c *client) send(logEntry LogEntry) error {
	url := c.targetURL
	headers := map[string]string{
		"Content-Type": "application/json",
		"X-API-KEY":    c.apiKey,
		"X-APP-NAME":   c.appName,
	}

	err := c.httpClient.Post(url, headers, logEntry, nil)
	if err != nil {
		return fmt.Errorf("failed to send log entry: %v", err)
	}

	return nil
}

// Ping checks the health of the log server.
func (c *client) Ping() error {
	url := c.targetURL + "/ping"
	headers := map[string]string{
		"X-API-KEY":  c.apiKey,
		"X-APP-NAME": c.appName,
	}

	var result interface{}
	err := c.httpClient.Get(url, headers, &result)
	if err != nil {
		return fmt.Errorf("failed to ping log server: %v", err)
	}
	return nil
}
