package config

import "github.com/spf13/viper"

type NatsConfig struct {
	Host    string `yaml:"Host"`
	Port    int    `yaml:"port"`
	Subject string `yaml:"subject"`
}

type RocketMQConfig struct {
	Host      string `yaml:"Host"`
	Port      int    `yaml:"port"`
	GroupName string `yaml:"groupName"`
	TopicName string `yaml:"topicName"`
	MaxRetry  int    `yaml:"maxRetry"`
	Offset    int    `yaml:"offset"`
	Timeout   int    `yaml:"timeout"`
}

type Config struct {
	NatsCfg     *NatsConfig
	RocketMqCfg *RocketMQConfig
}

func NewGlobalConfig() *Config {
	natsCfg := &NatsConfig{
		Host:    viper.GetString("nats.host"),
		Port:    viper.GetInt("nats.port"),
		Subject: viper.GetString("nats.subject"),
	}
	rocketMqCfg := &RocketMQConfig{
		Host:      viper.GetString("rocketmq.host"),
		Port:      viper.GetInt("rocketmq.port"),
		GroupName: viper.GetString("rocketmq.groupName"),
		TopicName: viper.GetString("rocketmq.topicName"),
		MaxRetry:  viper.GetInt("rocketmq.maxRetry"),
		Offset:    viper.GetInt("rocketmq.offset"),
		Timeout:   viper.GetInt("rocketmq.timeout"),
	}

	return &Config{
		NatsCfg:     natsCfg,
		RocketMqCfg: rocketMqCfg,
	}
}
