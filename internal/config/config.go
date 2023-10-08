package config

import "github.com/spf13/viper"

type AppConfig struct {
	DatabaseURL string `mapstructure:"database_url"`
	RabbitMQURL string `mapstructure:"rabbitmq_url"`
	HTTPPort    int    `mapstructure:"http_port"`
	// Add other configuration fields as needed
}

func LoadConfig(configFilePath string) (*AppConfig, error) {
	viper.SetConfigFile(configFilePath)

	// Optionally, set default values for configuration fields
	// viper.SetDefault("database_url", "default_database_url")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config AppConfig
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
