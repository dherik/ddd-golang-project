package config

import "github.com/spf13/viper"

type AppConfig struct {
	HTTPPort  int    `mapstructure:"http_port"`
	JWTSecret string `mapstructure:"jwt_secret"`
	// Add other configuration fields as needed
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
