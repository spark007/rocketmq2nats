package config

import "github.com/spf13/viper"

func InitConfig() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("rocketmq2nats")
	viper.AddConfigPath("config")

	err := viper.ReadInConfig()
	if err != nil {
		panic("read config failed: " + err.Error())
	}
}
