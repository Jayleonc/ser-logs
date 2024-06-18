package serlogs

import "time"

// Field represents a key-value pair for log content.
type Field struct {
	Key string
	Val any
}

// Logger defines the interface for logging at various levels.
type Logger interface {
	LogEntry(LogEntry) error
	Info(module, method, requestID string, fields ...Field) error
	Warn(module, method, requestID string, fields ...Field) error
	Error(module, method, requestID string, fields ...Field) error
}

// logger implements the Logger interface.
type logger struct {
	Client      LogClient
	ServiceName string
	Host        string
	Env         string
}

// NewLogger creates a new logger with the specified configuration.
func NewLogger(targetURL, apiKey, appName, serviceName, host, env string) (Logger, error) {
	httpClient := NewHTTPClient(10 * time.Second)
	client := NewLogClient(targetURL, apiKey, appName, httpClient)
	err := client.Ping()
	if err != nil {
		return nil, err
	}
	return &logger{
		Client:      client,
		ServiceName: serviceName,
		Host:        host,
		Env:         env,
	}, nil
}

// LogEntry logs a generic log entry.
func (l *logger) LogEntry(entry LogEntry) error {
	entry.ServiceName = l.ServiceName
	entry.Host = l.Host
	entry.Env = l.Env
	return l.Client.Send(entry)
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
