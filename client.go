package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ClientI interface {
	Send(entry Entry) error
	Ping() error
}

type Client struct {
	TargetUrl  string
	ApiKey     string
	AppName    string
	HttpClient *http.Client
}

func NewLogClient(targetUrl, apiKey, appName string) *Client {
	return &Client{
		TargetUrl:  targetUrl,
		ApiKey:     apiKey,
		AppName:    appName,
		HttpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) Send(entry Entry) error {
	url := c.TargetUrl
	headers := map[string]string{
		"Content-Type": "application/json",
		"X-API-KEY":    c.ApiKey,
		"X-APP-NAME":   c.AppName,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("请求体序列化失败: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) Ping() error {
	url := c.TargetUrl + "/ping"
	headers := map[string]string{
		"X-API-KEY":  c.ApiKey,
		"X-APP-NAME": c.AppName,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	return nil
}
