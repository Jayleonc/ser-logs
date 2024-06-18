package log

type LoggerI interface {
	log(entry Entry) error
	Info(module, method, requestID string, logContent map[string]interface{}) error
	Warn(module, method, requestID string, logContent map[string]interface{}) error
	Error(module, method, requestID string, logContent map[string]interface{}) error
}

type Logger struct {
	Client      ClientI
	ServiceName string
	Host        string
	Env         string
}

func NewLogger(targetUrl, apiKey, appName, serviceName, host, env string) (LoggerI, error) {
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

func (l *Logger) log(logRequest Entry) error {
	logRequest.ServiceName = l.ServiceName
	logRequest.Host = l.Host
	logRequest.Env = l.Env
	return l.Client.Send(logRequest)
}

func (l *Logger) Info(module, method, requestID string, logContent map[string]interface{}) error {
	logRequest := Entry{
		ModuleName: module,
		MethodName: method,
		RequestID:  requestID,
		LogLevel:   "INFO",
		LogContent: logContent,
	}
	return l.log(logRequest)
}

func (l *Logger) Warn(module, method, requestID string, logContent map[string]interface{}) error {
	logRequest := Entry{
		ModuleName: module,
		MethodName: method,
		RequestID:  requestID,
		LogLevel:   "WARN",
		LogContent: logContent,
	}
	return l.log(logRequest)
}

func (l *Logger) Error(module, method, requestID string, logContent map[string]interface{}) error {
	logRequest := Entry{
		ModuleName: module,
		MethodName: method,
		RequestID:  requestID,
		LogLevel:   "ERROR",
		LogContent: logContent,
	}
	return l.log(logRequest)
}
