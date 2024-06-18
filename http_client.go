package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type HTTPClientI interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string, headers map[string]string, result interface{}) error
	Post(url string, headers map[string]string, body, result interface{}) error
	Put(url string, headers map[string]string, body, result interface{}) error
	Delete(url string, headers map[string]string, body, result interface{}) error
}

type HTTPClient struct {
	client *http.Client
}

func NewHTTPClient(timeout time.Duration) HTTPClientI {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

func (c *HTTPClient) doRequest(method, url string, headers map[string]string, body, result interface{}) error {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("请求体序列化失败: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s 请求失败，状态码: %d", method, resp.StatusCode)
	}

	if result != nil {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("读取响应体失败: %v", err)
		}
		return json.Unmarshal(respBody, result)
	}

	return nil
}

func (c *HTTPClient) Get(url string, headers map[string]string, result interface{}) error {
	return c.doRequest("GET", url, headers, nil, result)
}

func (c *HTTPClient) Post(url string, headers map[string]string, body, result interface{}) error {
	return c.doRequest("POST", url, headers, body, result)
}

func (c *HTTPClient) Put(url string, headers map[string]string, body, result interface{}) error {
	return c.doRequest("PUT", url, headers, body, result)
}

func (c *HTTPClient) Delete(url string, headers map[string]string, body, result interface{}) error {
	return c.doRequest("DELETE", url, headers, body, result)
}
