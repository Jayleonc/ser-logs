package serlogs

// LogEntry represents a log entry to be sent to the log server.
type LogEntry struct {
	// ServiceName is the name of the service sending the log.
	ServiceName string `json:"service_name" binding:"required"`
	// ModuleName is the name of the module sending the log.
	ModuleName string `json:"module_name" binding:"required"`
	// MethodName is the name of the method sending the log.
	MethodName string `json:"method_name" binding:"required"`
	// RequestID is the unique identifier for the request.
	RequestID string `json:"request_id" binding:"required"`
	// LogLevel is the level of the log (e.g., INFO, WARN, ERROR).
	LogLevel string `json:"log_level" binding:"required"`
	// LogContent is the content of the log entry.
	LogContent map[string]interface{} `json:"log_content" binding:"required"`
	// Host is the host name of the service sending the log.
	Host string `json:"host"`
	// Env is the environment in which the service is running.
	Env string `json:"env"`
}
