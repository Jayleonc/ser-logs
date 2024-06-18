package log

type Logger struct {
	Client      *Client
	ServiceName string
	Host        string
	Env         string
}

func NewLogger(targetUrl, apiKey, appName, serviceName, host, env string) (*Logger, error) {
	client := NewLogClient(targetUrl, apiKey, appName)
	err := client.Ping()
	if err != nil {
		return nil, err
	}
	return &Logger{
		Client:      client,
		ServiceName: serviceName,
		Host:        host,
		Env:         env,
	}, nil
}

func (l *Logger) log(level, module, method, requestID string, logContent map[string]interface{}) error {
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
	return l.Client.send(entry)
}

func (l *Logger) Info(module, method, requestID string, logContent map[string]interface{}) error {
	return l.log("INFO", module, method, requestID, logContent)
}

func (l *Logger) Warn(module, method, requestID string, logContent map[string]interface{}) error {
	return l.log("WARN", module, method, requestID, logContent)
}

func (l *Logger) Error(module, method, requestID string, logContent map[string]interface{}) error {
	return l.log("ERROR", module, method, requestID, logContent)
}
