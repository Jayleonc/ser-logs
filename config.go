package serlogs

// Config holds the configuration for the log client and logger.
type Config struct {
	// Url is the URL of the log server.
	//Url string
	// APIKey is the API key for authenticating with the log server.
	APIKey string
	// AppName is the name of the application sending logs.
	AppName string
	// ServiceName is the name of the service sending logs.
	ServiceName string
	// Host is the host name of the service sending logs.
	Host string
	// Env is the environment in which the service is running.
	Env string
}
