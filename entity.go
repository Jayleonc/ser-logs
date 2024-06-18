package log

type Entry struct {
	ServiceName string                 `json:"service_name" binding:"required"`
	ModuleName  string                 `json:"module_name" binding:"required"`
	MethodName  string                 `json:"method_name" binding:"required"`
	RequestID   string                 `json:"request_id" binding:"required"`
	LogLevel    string                 `json:"log_level" binding:"required"`
	LogContent  map[string]interface{} `json:"log_content" binding:"required"`
	Host        string                 `json:"host"`
	Env         string                 `json:"env"`
}
