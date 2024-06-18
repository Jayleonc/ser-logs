package serlogs

type Config struct {
	// TargetURL is the URL of the log server.
	TargetURL string
	// APIKey is the API key for authenticating with the log server.
	APIKey string
	// AppName is the name of the application sending logs.
	AppName string
	// HTTPClient is the HTTP client used to send requests.
	HTTPClient httpClient
}
