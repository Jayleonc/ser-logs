package internal

import (
	"context"
	"errors"
	"github.com/Jayleonc/register/config_center"
	"github.com/spf13/viper"
	"log"
)

type LogServiceClient struct {
	ip string
}

func NewLogServiceClient() (*LogServiceClient, error) {
	// 读取配置文件
	viper.SetConfigFile("dev.yaml")
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
	ip, err := configCenterClient.GetConfig(ctx, "ser-logs/ip")
	if err != nil {
		return nil, err
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
