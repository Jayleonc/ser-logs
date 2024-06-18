package serlogs

import "time"

// Logger defines the interface for logging at various levels.
type Logger interface {
	// LogEntry logs a generic log entry.
	LogEntry(LogEntry) error
	// Info logs an informational message.
	Info(module, method, requestID string, fields ...Field) error
	// Warn logs a warning message.
	Warn(module, method, requestID string, fields ...Field) error
	// Error logs an error message.
	Error(module, method, requestID string, fields ...Field) error
}

// logger implements the Logger interface.
type logger struct {
	// client is the log client used to send log entries.
	client LogClient
	// serviceName is the name of the service sending logs.
	serviceName string
	// host is the host name of the service sending logs.
	host string
	// env is the environment in which the service is running.
	env string
}

// NewLogger creates a new logger with the specified configuration.
func NewLogger(targetURL, apiKey, appName, serviceName, host, env string) (Logger, error) {
	c := NewSerLogsClient(Config{
		TargetURL:  targetURL,
		APIKey:     apiKey,
		AppName:    appName,
		HTTPClient: newHTTPClient(time.Second * 10),
	})

	err := c.Ping()
	if err != nil {
		return nil, err
	}

	return &logger{
		client:      c,
		serviceName: serviceName,
		host:        host,
		env:         env,
	}, nil
}

// LogEntry logs a generic log entry.
func (l *logger) LogEntry(entry LogEntry) error {
	entry.ServiceName = l.serviceName
	entry.Host = l.host
	entry.Env = l.env
	return l.client.send(entry)
}

// log is a helper method to create and send a log entry with the specified level and fields.
func (l *logger) log(level, module, method, requestID string, fields ...Field) error {
	logContent := make(map[string]interface{})
	for _, field := range fields {
		logContent[field.Key] = field.Val
	}
	entry := LogEntry{
		ModuleName: module,
		MethodName: method,
		RequestID:  requestID,
		LogLevel:   level,
		LogContent: logContent,
	}
	return l.LogEntry(entry)
}

// Info logs an informational message.
func (l *logger) Info(module, method, requestID string, fields ...Field) error {
	return l.log("INFO", module, method, requestID, fields...)
}

// Warn logs a warning message.
func (l *logger) Warn(module, method, requestID string, fields ...Field) error {
	return l.log("WARN", module, method, requestID, fields...)
}

// Error logs an error message.
func (l *logger) Error(module, method, requestID string, fields ...Field) error {
	return l.log("ERROR", module, method, requestID, fields...)
}
