package log

import (
	"time"
)

type Logger interface {
	Info(module, method, requestID string, logContent map[string]interface{})
	Warn(module, method, requestID string, logContent map[string]interface{})
	Error(module, method, requestID string, logContent map[string]interface{})
}

type logger struct {
	ServiceName string
	Host        string
	Env         string
	sender      SenderI
}

func NewLogger(sender SenderI, serviceName, host, env string) Logger {
	return &logger{
		sender:      sender,
		ServiceName: serviceName,
		Host:        host,
		Env:         env,
	}
}

func (l *logger) log(level, module, method, requestID string, logContent map[string]interface{}) {
	entry := Entry{
		ServiceName: l.ServiceName,
		ModuleName:  module,
		MethodName:  method,
		RequestID:   requestID,
		LogLevel:    level,
		LogContent:  logContent,
		Host:        l.Host,
		Env:         l.Env,
	}
	l.sender.LogChan() <- entry
}

func (l *logger) Info(module, method, requestID string, logContent map[string]interface{}) {
	l.log("INFO", module, method, requestID, logContent)
}

func (l *logger) Warn(module, method, requestID string, logContent map[string]interface{}) {
	l.log("WARN", module, method, requestID, logContent)
}

func (l *logger) Error(module, method, requestID string, logContent map[string]interface{}) {
	l.log("ERROR", module, method, requestID, logContent)
}

// NewSimpleLogger simplifies the initialization process for the logger.
func NewSimpleLogger(targetUrl, apiKey, appName, serviceName, host, env string) Logger {
	httpClient := NewHTTPClient(10 * time.Second)
	client := NewLogClient(targetUrl, apiKey, appName, httpClient)
	sender := NewLogSender(client)
	logger := NewLogger(sender, serviceName, host, env)
	sender.Start()
	return logger
}
