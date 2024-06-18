

# Ser-Logs SDK 

`ser-logs` SDK 是一个用于将服务消息记录到集中日志服务器的 Go 库。

它提供了易于使用的方法，用于在各种级别（Info、Warn、Error）记录日志，并允许您在日志中包含额外的上下文字段。 

## 安装 

要安装 SDK，请运行：

```sh
go get github.com/Jayleonc/ser-logs
```

## 使用方法

### 初始化

首先，使用您的配置创建一个 `Config` 结构体：

```go
package main

import (
	"github.com/Jayleonc/ser-logs"
	"log"
)

func main() {
	config := serlogs.Config{
		TargetURL:   "http://localhost:8080/logs",
		APIKey:      "your-api-key",
		AppName:     "your-app-name",
		ServiceName: "user_service",
		Host:        "localhost",
		Env:         "dev",
	}

	logger, err := serlogs.NewLogger(config)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
}
```

### 记录日志

您可以使用各种级别（Info、Warn、Error）记录日志，并添加额外的上下文字段：

```go
logger.Info("MyModule", "MyMethod", "12345",
	serlogs.String("message", "This is an info log"),
	serlogs.String("detail", "This is a detailed message"))

logger.Warn("MyModule", "MyMethod", "12345",
	serlogs.String("message", "This is a warning log"),
	serlogs.String("detail", "This is a detailed message"))

logger.Error("MyModule", "MyMethod", "12345",
	serlogs.String("message", "This is an error log"),
	serlogs.String("detail", "This is a detailed message"))

```

### 并发记录日志

要并发记录日志，您可以使用 Go 协程：

```go
var wg sync.WaitGroup
for i := 0; i < 10; i++ {
	wg.Add(1)
	go func(i int) {
		defer wg.Done()
		err := logger.Info("MyModule", "MyMethod", fmt.Sprintf("request-%d", i),
			serlogs.String("message", fmt.Sprintf("This is log number %d", i)),
			serlogs.String("detail", fmt.Sprintf("This is log number %d", i)))
		if err != nil {
			log.Printf("Failed to send log %d: %v", i, err)
		} else {
			log.Printf("Log %d sent successfully", i)
		}
	}(i)
}
wg.Wait()

fmt.Println("All logs have been sent")
```

- ### API 参考

  #### Logger

  `Logger` 接口提供了用于在不同级别记录消息的方法。

  - `LogEntry(entry LogEntry) error`：记录一个通用的日志条目。
  - `Info(module, method, requestID string, fields ...Field) error`：记录一条信息日志。
  - `Warn(module, method, requestID string, fields ...Field) error`：记录一条警告日志。
  - `Error(module, method, requestID string, fields ...Field) error`：记录一条错误日志。

#### LogEntry

`LogEntry` 结构体表示要发送到日志服务器的日志条目。

```go
type LogEntry struct {
	ServiceName string                 `json:"service_name" binding:"required"`
	ModuleName  string                 `json:"module_name" binding:"required"`
	MethodName  string                 `json:"method_name" binding:"required"`
	RequestID   string                 `json:"request_id" binding:"required"`
	LogLevel    string                 `json:"log_level" binding:"required"`
	LogContent  map[string]interface{} `json:"log_content" binding:"required"`
	Host        string                 `json:"host"`
	Env         string                 `json:"env"`
}
```

#### Field

`Field` 结构体表示日志内容的键值对。提供了几个辅助函数来创建不同类型的 `Field` 实例：

```go
func Error(err error) Field {
	return Field{"error", err}
}

func Int64(key string, val int64) Field {
	return Field{key, val}
}

func Uint8(key string, val uint8) Field {
	return Field{key, val}
}

func String(key, val string) Field {
	return Field{key, val}
}

func Slice(key string, val ...any) Field {
	return Field{key, val}
}

func Bool(key string, val bool) Field {
	return Field{key, val}
}
```

### 错误处理

如果日志条目未能发送，所有日志记录方法都会返回一个错误。建议在应用程序中适当地处理这些错误。

```go
err := logger.Info("MyModule", "MyMethod", "12345",
	serlogs.String("message", "This is an info log"))
if err != nil {
	log.Printf("Failed to send info log: %v", err)
}
```

### HTTP 客户端优化

`Client` 结构体使用一个可配置超时的优化 HTTP 客户端。当前默认10秒，不可配置。它处理发送日志条目和执行健康检查。

```go
type Client struct {
	TargetURL  string
	APIKey     string
	AppName    string
	HTTPClient HTTPClientI
}

func NewLogClient(targetURL, apiKey, appName string) *Client {
	return &Client{
		TargetURL:  targetURL,
		APIKey:     apiKey,
		AppName:    appName,
		HTTPClient: NewHTTPClient(10 * time.Second),
	}
}
```

### 完整代码结构

```plaintext
ser-logs/
├── client.go
├── entity.go
├── fields.go
├── http_client.go
├── logger.go
└── README.md
```
