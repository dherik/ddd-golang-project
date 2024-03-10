package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type HTTP struct {
	Port int `mapstructure:"port"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

type JWT struct {
	Secret string `mapstructure:"secret"`
}

type AppConfig struct {
	HTTP     HTTP     `mapstructure:"http"`
	Database Database `mapstructure:"database"`
	JWT      JWT      `mapstructure:"jwt"`
}

func LoadConfig(configFilePath string) (*AppConfig, error) {
	viper.SetConfigFile(configFilePath)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config AppConfig
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
