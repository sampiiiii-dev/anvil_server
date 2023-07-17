package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Server struct {
		Address       string `mapstructure:"address"`
		BannerWebsite string `mapstructure:"banner_website"`
		Version       string `mapstructure:"version"`
	} `mapstructure:"server"`
	DB struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Pass     string `mapstructure:"pass"`
		Database string `mapstructure:"database"`
	} `mapstructure:"db"`
}

func LoadConfig(scribe *zap.Logger) Config {
	var config Config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		scribe.Warn("Error reading config file", zap.Error(err))
	}

	// Set undefined variables
	viper.SetDefault("Server.Address", "127.0.0.1:1323")

	if err := viper.Unmarshal(&config); err != nil {
		scribe.Warn("Unable to decode into struct", zap.Error(err))
	}

	return config
}
