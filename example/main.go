package main

import (
	serlogs "github.com/Jayleonc/ser-logs"
	"log"
)

func main() {
	// 初始化 Logger
	logger, err := serlogs.NewLogger("http://localhost:8080/logs", "your-api-key", "your-app-name", "user_service", "localhost", "dev")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	logger.Info("MyModule", "MyMethod", "12345",
		serlogs.Field{Key: "message", Val: "This is an info log"},
		serlogs.Field{Key: "detail", Val: "This is a detailed message"})

}
