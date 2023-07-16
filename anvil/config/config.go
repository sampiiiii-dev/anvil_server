package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress string
}

func LoadConfig(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")

	var config Config
	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	err := viper.Unmarshal(&config)
	return config, err
}
