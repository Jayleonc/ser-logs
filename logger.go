package serlogs

import "time"

// Logger defines the interface for logging at various levels.
type Logger interface {
	// LogEntry logs a generic log entry.
	LogEntry(LogEntry) error
	// Info logs an informational message.
	Info(module, method, requestId string, fields ...Field) error
	// Warn logs a warning message.
	Warn(module, method, requestId string, fields ...Field) error
	// Error logs an error message.
	Error(module, method, requestId string, fields ...Field) error
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
func NewLogger(cfg Config) (Logger, error) {
	c, err := NewSerLogsClient(cfg.APIKey, cfg.AppName, newHTTPClient(time.Second*10))
	if err != nil {
		return nil, err
	}

	err = c.Ping()
	if err != nil {
		return nil, err
	}

	return &logger{
		client:      c,
		serviceName: cfg.ServiceName,
		host:        cfg.Host,
		env:         cfg.Env,
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
func (l *logger) log(level, module, method, requestId string, fields ...Field) error {
	logContent := make(map[string]interface{})
	for _, field := range fields {
		logContent[field.Key] = field.Val
	}
	entry := LogEntry{
		ModuleName: module,
		MethodName: method,
		RequestID:  requestId,
		LogLevel:   level,
		LogContent: logContent,
	}
	return l.LogEntry(entry)
}

// Info logs an informational message.
func (l *logger) Info(module, method, requestId string, fields ...Field) error {
	return l.log("INFO", module, method, requestId, fields...)
}

// Warn logs a warning message.
func (l *logger) Warn(module, method, requestId string, fields ...Field) error {
	return l.log("WARN", module, method, requestId, fields...)
}

// Error logs an error message.
func (l *logger) Error(module, method, requestId string, fields ...Field) error {
	return l.log("ERROR", module, method, requestId, fields...)
}
