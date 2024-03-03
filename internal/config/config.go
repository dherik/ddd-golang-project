package config

import "github.com/spf13/viper"

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

	// Optionally, set default values for configuration fields
	// viper.SetDefault("database_url", "default_database_url")
	viper.SetDefault("http_port", 3333)

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
