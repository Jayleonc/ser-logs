package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/Jayleonc/register/config_center"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type LogServiceClient struct {
	ip string
}

const defaultConfigFilePath = "internal/dev.yaml"

func NewLogServiceClient() (*LogServiceClient, error) {
	getwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	fmt.Println(getwd)
	// 读取配置文件

	// 获取当前文件的绝对路径
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, errors.New("failed to get current file path")
	}
	configFilePath := filepath.Join(filepath.Dir(filename), "dev.yaml")

	fmt.Println("ser-logs sdk using config file:", configFilePath)

	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file: %v", err)
		return nil, err
	}

	// 获取etcd地址
	etcdAddresses := viper.GetStringSlice("etcd.addresses")
	if len(etcdAddresses) == 0 {
		log.Fatalf("Etcd addresses not provided in config file")
		return nil, errors.New("etcd addresses not provided in config file")
	}

	// 初始化配置中心客户端
	configCenterClient, err := config_center.NewClient(
		config_center.WithEtcdAddresses(etcdAddresses),
		// 添加其他必要的配置
	)
	if err != nil {
		log.Fatalf("Failed to create config center client: %v", err)
		return nil, err
	}

	ctx := context.Background()
	ip, err := configCenterClient.GetConfig(ctx, "ser-logs/host")
	if err != nil {
		return nil, err
	}

	// 检查 ip 是否包含协议前缀
	if !strings.HasPrefix(ip, "http://") {
		ip = "http://" + ip
	}

	return &LogServiceClient{
		ip: ip,
	}, nil
}

func (c *LogServiceClient) GetIP() string {
	return c.ip
}

func (c *LogServiceClient) SendLog(logData string) {
	// 使用 c.ip 发送日志
	log.Printf("Sending log to log service at %s: %s", c.ip, logData)
	// 这里可以添加具体的发送日志逻辑，例如HTTP请求等
}
