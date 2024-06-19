package serlogs

import (
	"context"
	"errors"
	"fmt"
	"net"
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
		// 解析错误类型并提供详细提示信息
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return fmt.Errorf("initialization failed: unable to reach the log server due to a network timeout")
		case errors.Is(err, context.Canceled):
			return fmt.Errorf("initialization failed: network request to the log server was canceled")
		}

		var opErr *net.OpError
		if errors.As(err, &opErr) {
			return fmt.Errorf("initialization failed: network error occurred while trying to reach the log server: %v", err)
		}

		return fmt.Errorf("initialization failed: failed to ping log server: %v", err)
	}
	return nil
}
