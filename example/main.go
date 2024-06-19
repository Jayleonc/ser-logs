package main

import (
	serlogs "github.com/Jayleonc/ser-logs"
	"log"
)

func main() {
	logger, err := serlogs.NewLogger(serlogs.Config{
		Url:         "http://localhost:8080",
		APIKey:      "your-api-key",
		AppName:     "your-app-name",
		ServiceName: "user_service",
		Host:        "localhost",
		Env:         "dev",
	})

	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	err = logger.Info("MyModule", "MyMethod", "12345",
		serlogs.Field{Key: "message", Val: "This is an info log"},
		serlogs.Field{Key: "detail", Val: "This is a detailed message"})
	if err != nil {
		log.Fatalf("Failed to send info log: %v", err)
	}
}
