package log

import (
	"fmt"
)

type ClientI interface {
	Send(message Entry) error
	Ping() error
}

type Client struct {
	TargetUrl  string
	ApiKey     string
	AppName    string
	HttpClient HTTPClientI
}

func NewLogClient(targetUrl, apiKey, appName string, httpClient HTTPClientI) ClientI {
	return &Client{
		TargetUrl:  targetUrl,
		ApiKey:     apiKey,
		AppName:    appName,
		HttpClient: httpClient,
	}
}

func (c *Client) Send(l Entry) error {
	url := c.TargetUrl
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	err := c.HttpClient.Post(url, headers, l, nil)
	if err != nil {
		return fmt.Errorf("发送日志请求失败: %v", err)
	}

	return nil
}

func (c *Client) Ping() error {
	url := c.TargetUrl + "/ping"
	headers := map[string]string{
		"X-API-KEY":  c.ApiKey,
		"X-APP-NAME": c.AppName,
	}

	var result interface{}
	err := c.HttpClient.Get(url, headers, &result)
	if err != nil {
		return fmt.Errorf("无法连接到日志服务: %v", err)
	}
	return nil
}
